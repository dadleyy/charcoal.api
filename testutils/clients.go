package testutils

import "github.com/dadleyy/charcoal.api/util"
import "github.com/dadleyy/charcoal.api/models"

func CreateClient(out *models.Client, name string, system bool) error {
	db := NewDB()
	defer db.Close()

	c := models.Client{
		ClientID:     util.RandStringBytesMaskImprSrc(20),
		ClientSecret: util.RandStringBytesMaskImprSrc(40),
		Name:         name,
		System:       system,
	}

	if e := db.Create(&c).Error; e != nil {
		return e
	}

	if e := db.First(out, c.ID).Error; e != nil {
		return e
	}

	return nil
}
