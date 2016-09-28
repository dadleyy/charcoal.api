package models

type Client struct {
	Common
	Name string `json:"name"`
	ClientID string `json:"client_id"`
	ClientSecret string `json:"-"`
	RedirectUri string `json:"redirect_uri"`
}

func (target *Client) Marshal() interface{} {
	return target
}
