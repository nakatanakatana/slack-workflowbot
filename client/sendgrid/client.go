package sendgrid

import (
	"errors"

	"github.com/sendgrid/rest"
	lib "github.com/sendgrid/sendgrid-go"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrUnexpected = errors.New("unexpected error")
)

type Client struct {
	apiKey string
	host   string
	api    API
}

type Option interface {
	Apply(*Client)
}

type Host string

func (h Host) Apply(c *Client) {
	c.host = string(h)
}

type API func(rest.Request) (*rest.Response, error)

func (a API) Apply(c *Client) {
	c.api = a
}

func New(token string, options ...Option) *Client {
	cli := &Client{
		token,
		"https://api.sendgrid.com",
		lib.API,
	}
	for _, o := range options {
		o.Apply(cli)
	}

	return cli
}
