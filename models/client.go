package models

type Client struct {
	Common
	Name string `json:"name"`
	ClientID string `json:"client_id"`
	ClientSecret string `json:"-"`
}

func (target *Client) Marshal() interface{} {
	return target
}
