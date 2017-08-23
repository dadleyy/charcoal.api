package models

import "fmt"

// Client records represent application consumers of the charcoal api.
type Client struct {
	Common
	Name         string `json:"name"`
	ClientID     string `json:"-"`
	ClientSecret string `json:"-"`
	RedirectURI  string `json:"redirect_uri"`
	Description  string `json:"description"`
	System       bool   `json:"system"`
}

// URL returns the api url that will provide details for the client.
func (client *Client) URL() string {
	return fmt.Sprintf("/clients?filter[id]=eq(%d)", client.ID)
}

// Type returns unique type string to identify this kind of record.
func (client *Client) Type() string {
	return "application/vnd.miritos.client+json"
}
