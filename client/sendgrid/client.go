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

type APIKey string

func (a APIKey) Apply(c *Client) {
	c.apiKey = string(a)
}

type Host string

func (h Host) Apply(c *Client) {
	c.host = string(h)
}

type API func(rest.Request) (*rest.Response, error)

func (a API) Apply(c *Client) {
	c.api = a
}

func New(options ...Option) *Client {
	cli := &Client{
		"",
		"https://api.sendgrid.com",
		lib.API,
	}
	for _, o := range options {
		o.Apply(cli)
	}

	return cli
}

func (cli *Client) GetRequest(endpoint string, options ...Option) rest.Request {
	config := &Client{}
	for _, o := range options {
		o.Apply(config)
	}

	apiKey := config.apiKey
	if apiKey == "" {
		apiKey = cli.apiKey
	}

	host := config.host
	if host == "" {
		host = cli.host
	}

	return lib.GetRequest(apiKey, endpoint, host)
}
