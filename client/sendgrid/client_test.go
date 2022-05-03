package sendgrid_test

import (
	"log"
	"testing"

	"github.com/nakatanakatana/slack-workflowbot/client/sendgrid"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("without Options", func(t *testing.T) {
		t.Parallel()

		cli := sendgrid.New("sg-token")
		log.Println("cli", cli)
	})

	t.Run("with Host Options", func(t *testing.T) {
		t.Parallel()

		cli := sendgrid.New("sg-token", sendgrid.Host("http://localhost:8111"))
		log.Println("cli", cli)
	})
}
