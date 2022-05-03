package sendgrid_test

import (
	_ "embed"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/nakatanakatana/slack-workflowbot/client/sendgrid"
	"github.com/sendgrid/rest"
)

//go:embed testdata/bounces_delete_success.json
var bounceDeleteSuccessBody string

//go:embed testdata/bounces_delete_notfound.json
var bounceDeleteNotFoundBody string

//go:embed testdata/bounces_get_success.json
var bounceGetSuccessBody string

//go:embed testdata/bounces_get_notfound.json
var bounceGetNotFoundBody string

func bounceDeleteSuccess(_ rest.Request) (*rest.Response, error) {
	resp := &rest.Response{
		StatusCode: 204,
		Body:       bounceDeleteSuccessBody,
	}

	return resp, nil
}

func bounceDeleteNotFound(_ rest.Request) (*rest.Response, error) {
	resp := &rest.Response{
		StatusCode: 404,
		Body:       bounceDeleteNotFoundBody,
	}

	return resp, nil
}

func bounceGetSuccess(_ rest.Request) (*rest.Response, error) {
	resp := &rest.Response{
		StatusCode: 200,
		Body:       bounceGetSuccessBody,
	}

	return resp, nil
}

func bounceGetNotFound(_ rest.Request) (*rest.Response, error) {
	resp := &rest.Response{
		StatusCode: 200,
		Body:       bounceGetNotFoundBody,
	}

	return resp, nil
}

func TestDeleteBounce(t *testing.T) {
	t.Parallel()

	t.Run("delete success", func(t *testing.T) {
		t.Parallel()

		cli := sendgrid.New("", sendgrid.API(bounceDeleteSuccess))
		_, _, err := cli.DeleteBounce("nakatanakatana@gmail.com")
		if err != nil {
			t.Fail()
		}
	})

	t.Run("delete not found", func(t *testing.T) {
		t.Parallel()

		cli := sendgrid.New("", sendgrid.API(bounceDeleteNotFound))

		_, _, err := cli.DeleteBounce("nakatanakatana@gmail.com")
		if !errors.Is(err, sendgrid.ErrNotFound) {
			t.Fail()
		}
	})

	t.Run("get success", func(t *testing.T) {
		t.Parallel()

		cli := sendgrid.New("", sendgrid.API(bounceGetSuccess))
		result, _, err := cli.GetBounce("nakatanakatana@gmail.com")
		if err != nil {
			t.Fail()
		}
		if result.Created != 1443651125 ||
			!result.CreatedAt.Equal(
				time.Date(2015, time.October, 1, 7, 12, 5, 0, time.Local),
			) ||
			result.Email != "bounce1@test.com" ||
			!strings.HasPrefix(result.Reason, "550 5.1.1 The email account that you") ||
			result.Status != "5.1.1" {
			t.Fail()
		}
	})

	t.Run("get not found", func(t *testing.T) {
		t.Parallel()

		cli := sendgrid.New("", sendgrid.API(bounceGetNotFound))
		_, _, err := cli.GetBounce("nakatanakatana@gmail.com")
		if !errors.Is(err, sendgrid.ErrNotFound) {
			t.Fail()
		}
	})
}
