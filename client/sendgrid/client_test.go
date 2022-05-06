package sendgrid_test

import (
	"log"
	"strings"
	"testing"

	"github.com/nakatanakatana/slack-workflowbot/client/sendgrid"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("without Options", func(t *testing.T) {
		t.Parallel()

		cli := sendgrid.New()
		log.Println("cli", cli)
	})

	t.Run("check default parameters", func(t *testing.T) {
		t.Parallel()

		const endpoint = "/endpoint/"
		cli := sendgrid.New()

		req := cli.GetRequest(endpoint)
		if !strings.HasPrefix(req.BaseURL, "https://api.sendgrid.com") {
			t.Fail()
		}
	})

	t.Run("with Host Options", func(t *testing.T) {
		t.Parallel()

		cli := sendgrid.New(sendgrid.Host("http://localhost:8111"))
		log.Println("cli", cli)
	})
}

func TestGetRequest(t *testing.T) {
	const (
		defaultAPIKey = "defaultAPIKey"
		defaultHost   = "http://default.host"
		endpoint      = "/endpoint/"
	)

	t.Parallel()

	cli := sendgrid.New(
		sendgrid.APIKey(defaultAPIKey),
		sendgrid.Host(defaultHost),
	)

	t.Run("without Options", func(t *testing.T) {
		t.Parallel()

		req := cli.GetRequest(endpoint)
		if req.BaseURL != defaultHost+endpoint {
			t.Fail()
		}

		authHeader, ok := req.Headers["Authorization"]
		if !ok || !strings.HasSuffix(authHeader, defaultAPIKey) {
			t.Fail()
		}
	})

	t.Run("with Options", func(t *testing.T) {
		t.Parallel()

		const (
			apiKey = "apiKey"
			host   = "http://host"
		)

		req := cli.GetRequest(endpoint, sendgrid.APIKey(apiKey), sendgrid.Host(host))
		if req.BaseURL != host+endpoint {
			t.Fail()
		}

		authHeader, ok := req.Headers["Authorization"]
		if !ok || !strings.HasSuffix(authHeader, apiKey) {
			t.Fail()
		}
	})
}
