package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sendgrid/rest"
)

type GetBounceResult struct {
	CreatedAt time.Time
	Created   int64  `json:"created"`
	Email     string `json:"email"`
	Reason    string `json:"reason"`
	Status    string `json:"status"`
}

type DeleteBounceResult struct{}

type BounceManager interface {
	GetBounce(email string, options ...Option) (*GetBounceResult, *rest.Response, error)
	DeleteBounce(email string, options ...Option) (*DeleteBounceResult, *rest.Response, error)
}

var _ BounceManager = &Client{}

func (cli *Client) GetBounce(email string, options ...Option) (*GetBounceResult, *rest.Response, error) {
	// https://docs.sendgrid.com/api-reference/bounces-api/retrieve-a-bounce
	request := cli.GetRequest(
		fmt.Sprintf("/v3/suppression/bounces/%s", email),
		options...,
	)
	request.Method = "GET"

	response, err := cli.api(request)
	if err != nil {
		return nil, nil, err
	}

	var results []GetBounceResult

	err = json.Unmarshal([]byte(response.Body), &results)
	if err != nil {
		return nil, response, fmt.Errorf("unmarshal failed:%w", err)
	}

	if len(results) == 0 {
		return nil, response, ErrNotFound
	}

	if len(results) != 1 {
		return nil, response, ErrUnexpected
	}

	result := results[0]
	result.CreatedAt = time.Unix(result.Created, 0)

	return &result, response, nil
}

func (cli *Client) DeleteBounce(email string, options ...Option) (*DeleteBounceResult, *rest.Response, error) {
	// https://docs.sendgrid.com/api-reference/bounces-api/delete-a-bounce
	request := cli.GetRequest(
		fmt.Sprintf("/v3/suppression/bounces/%s", email),
		options...,
	)
	request.Method = "DELETE"

	response, err := cli.api(request)
	if err != nil {
		return nil, nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, response, ErrNotFound
	}

	if response.StatusCode != http.StatusNoContent {
		return nil, response, ErrUnexpected
	}

	var result DeleteBounceResult

	return &result, response, nil
}
