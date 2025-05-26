package requester

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go-keeper/internal/client/logic/requester/options"
	"net/http"
	"syscall"
)

const (
	NoStatusCode = -1
)

var (
	ErrServerUnavailable    = errors.New("server unavailable")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

type Requester struct {
	host string
	ops  []options.Option
}

func New(host string, ops []options.Option) *Requester {
	return &Requester{
		host: host,
		ops:  ops,
	}
}

func (r *Requester) Post(path string, body any) (*resty.Response, error) {
	url := r.createURL(path)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req := resty.New().
		R().
		SetHeader("Content-Type", "application/json")

	for _, op := range r.ops {
		hs, err := op.GetHeaders()
		if err != nil {
			return nil, fmt.Errorf("failed to get headers: %w", err)
		}
		for h, v := range hs {
			req = req.SetHeader(h, v)
		}
	}

	resp, err := req.
		SetBody(bodyBytes).
		Post(url)

	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return nil, fmt.Errorf("%w: %w", err, ErrServerUnavailable)
		}
		return nil, fmt.Errorf("request failed: %w, err")
	}

	return resp, nil
}

func (r *Requester) Get(path string) (*resty.Response, error) {
	url := r.createURL(path)

	req := resty.New().
		R().
		SetHeader("Content-Type", "application/json")

	for _, op := range r.ops {
		hs, err := op.GetHeaders()
		if err != nil {
			return nil, fmt.Errorf("failed to get headers: %w", err)
		}
		for h, v := range hs {
			req = req.SetHeader(h, v)
		}
	}

	resp, err := req.Get(url)

	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return nil, fmt.Errorf("%w: %w", err, ErrServerUnavailable)
		}
		return nil, fmt.Errorf("request failed: %w, err")
	}

	return resp, nil
}

func Post[T any](r *Requester, path string, body any) (result T, statusCode int, err error) {
	var zero T

	resp, err := r.Post(path, body)
	if err != nil {
		return zero, NoStatusCode, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		break
	default:
		return zero, resp.StatusCode(), ErrUnexpectedStatusCode
	}

	var res T
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return zero, NoStatusCode, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, http.StatusOK, nil
}

func Get[T any](r *Requester, path string) (result T, statusCode int, err error) {
	var zero T

	resp, err := r.Get(path)
	if err != nil {
		return zero, NoStatusCode, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		break
	default:
		return zero, resp.StatusCode(), ErrUnexpectedStatusCode
	}

	var res T
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return zero, NoStatusCode, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, http.StatusOK, nil
}

func (r *Requester) createURL(path string) string {
	return fmt.Sprintf("http://%s%s", r.host, path)
}
