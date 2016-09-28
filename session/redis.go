package session

import "fmt"

type RedisStore struct {
}

func (store RedisStore) SetString(key string, value string) error {
	return nil
}

func (store RedisStore) SetInt(key string, value int) error {
	return nil
}

func (store RedisStore) GetInt(key string) (int, error) {
	return -1, nil
}

func (store RedisStore) GetString(key string) (string, error) {
	fmt.Printf("whoa\n")
	return "", nil
}
