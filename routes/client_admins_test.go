package routes

import "os"
import "fmt"
import "bytes"
import "testing"
import "net/http"

import "github.com/joho/godotenv"
import "github.com/labstack/gommon/log"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/activity"

func after(database *db.Connection) {
	var u models.User
	var c models.Client

	database.Where("client_id = ?", "test1-id").Find(&c)
	database.Where("email = ?", "test-1@client-admin-test.com").Find(&u)

	database.Unscoped().Delete(&models.ClientAdmin{User: u.ID, Client: c.ID})

	ch := database.Unscoped().Where("client_id in (?)", []string{"test1-id", "test2-id", "test3-id"})

	if err := ch.Delete(&models.Client{}).Error; err != nil {
		fmt.Print(err.Error())
	}

	uh := database.Unscoped().Where("email in (?)", []string{"test-1@client-admin-test.com", "test-2@client-admin-test.com"})

	if err := uh.Delete(&models.User{}).Error; err != nil {
		fmt.Print(err.Error())
	}
}

func before(database *db.Connection) {
	database.Create(&models.Client{Name: "client-admin-test1", ClientID: "test1-id", ClientSecret: "test1-secret"})
	database.Create(&models.Client{Name: "client-admin-test2", ClientID: "test2-id", ClientSecret: "test2-secret"})
	database.Create(&models.Client{Name: "client-admin-test3", ClientID: "test3-id", ClientSecret: "test3-secret"})

	for _, email := range []string{"test-1@client-admin-test.com", "test-2@client-admin-test.com"} {
		database.Create(&models.User{Email: &email})
	}

	var u models.User
	var c models.Client

	database.Where("client_id = ?", "test1-id").Find(&c)
	database.Where("email = ?", "test-1@client-admin-test.com").Find(&u)

	database.Create(&models.ClientAdmin{User: u.ID, Client: c.ID})
}

func TestFindClientAdminBadUser(t *testing.T) {
	_ = godotenv.Load("../.env")

	dbconf := db.Config{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DEBUG") == "true",
	}

	database, err := db.Open(dbconf)
	defer database.Close()

	if err != nil {
		panic(err)
	}

	before(database)
	defer after(database)

	logger := log.New("miritos")
	queue := make(chan activity.Message)

	buffer := make([]byte, 0)
	reader := bytes.NewReader(buffer)

	stub, err := http.NewRequest("GET", "/client-admins", reader)

	if err != nil {
		panic(err)
	}

	runtime := net.ServerRuntime{logger, dbconf, queue, nil}
	request, _ := runtime.Request(stub, &net.UrlParams{})

	database.Where("client_id = ?", "test1-id").Find(&request.Client)

	if err := FindClientAdmins(&request); err != nil {
		t.Log("successfully errored out w/o valid user")
		return
	}

	t.Fatalf("should not have passed w/o error")
}
func TestFindClientAdminsValidUser(t *testing.T) {
	_ = godotenv.Load("../.env")

	dbconf := db.Config{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DEBUG") == "true",
	}

	database, err := db.Open(dbconf)
	defer database.Close()

	if err != nil {
		panic(err)
	}

	before(database)
	defer after(database)

	logger := log.New("miritos")
	queue := make(chan activity.Message)

	buffer := make([]byte, 0)
	reader := bytes.NewReader(buffer)

	stub, err := http.NewRequest("GET", "/client-admins", reader)

	if err != nil {
		panic(err)
	}

	runtime := net.ServerRuntime{logger, dbconf, queue, nil}
	request, _ := runtime.Request(stub, &net.UrlParams{})

	database.Where("client_id = ?", "test1-id").Find(&request.Client)
	database.Where("email = ?", "test-1@client-admin-test.com").Find(&request.User)

	if err := FindClientAdmins(&request); err != nil {
		t.Fatal(err)
		return
	}

	t.Log("successfully passed w/ valid user")
}
