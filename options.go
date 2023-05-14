package socks

import "github.com/kayabe/socks/s5"

func WithUserPW(username string, password string) func(*Client) {
	return func(client *Client) {
		client.Authentication = &s5.AuthUserPW{Username: []byte(username), Password: []byte(password)}
	}
}

func WithVersion(version any) func(*Client) {
	return func(client *Client) {
		client.ProtocolVersion = version
	}
}
