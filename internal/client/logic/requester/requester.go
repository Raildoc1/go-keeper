package requester

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"syscall"
)

var (
	ErrServerUnavailable = errors.New("server unavailable")
	ErrBadRequest        = errors.New("bad request")
)

type Requester struct {
	host string
}

func New(host string) *Requester {
	return &Requester{
		host: host,
	}
}

func (r *Requester) Post(path string, body any) (*resty.Response, error) {
	url := r.createURL(path)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := resty.New().
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyBytes).
		Post(url)

	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return nil, fmt.Errorf("%w: %w", err, ErrServerUnavailable)
		}
		return nil, fmt.Errorf("request failed: %w, err")
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		return nil, ErrBadRequest
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return resp, nil
}

func Post[T any](r *Requester, path string, body any) (T, error) {
	var zero T

	resp, err := r.Post(path, body)
	if err != nil {
		return zero, err
	}

	var res T
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return zero, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

func (r *Requester) createURL(path string) string {
	return fmt.Sprintf("http://%s%s", r.host, path)
}
