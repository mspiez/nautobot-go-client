package nautobot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type OptionalFilter struct {
	Field string
	Value interface{}
}

type Client struct {
	nautobotToken string
	BaseURL       string
	Offset        int
	Limit         int
	HTTPClient    *http.Client
}

func NewClient(options ...func(*Client)) *Client {
	client := &Client{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	for _, o := range options {
		o(client)
	}
	return client
}

type PatchRequestError struct {
	StatusCode     string
	Message        string
	Err            error
	DetailNotFound bool
}

func (r *PatchRequestError) Error() string {
	return fmt.Sprintf("Status Code: %s; Message: %v", r.StatusCode, r.Message)
}

type DetailNotFound struct {
	Detail string
}

func WithToken(nautobotToken string) func(*Client) {
	return func(s *Client) {
		s.nautobotToken = nautobotToken
	}
}

func WithBaseURL(BaseURL string) func(*Client) {
	return func(s *Client) {
		s.BaseURL = BaseURL
	}
}

func WithOffset(offset int) func(*Client) {
	return func(s *Client) {
		if offset < 0 {
			log.Println("Offset set to 50. Wrong Offset value given: ", offset)
			s.Offset = 50
		} else {
			s.Offset = offset
		}
	}
}

func WithLimit(limit int) func(*Client) {
	return func(s *Client) {
		if limit < 0 {
			log.Println("Limit set to 50. Wrong Limit value given: ", limit)
			s.Limit = 50
		} else {
			s.Limit = limit
		}
	}
}

func GetResource[T any](ctx context.Context, url string, nautobotToken string) (T, error) {

	var m T
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return m, err
	}

	setHeaders(req, nautobotToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return m, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return m, err
	}

	if res.StatusCode != http.StatusOK {
		return m, fmt.Errorf("Error while retrieving data from %s. Status: %s. Message: %s", url, res.Status, string(body))
	}

	res.Body.Close()

	return parseJSON[T](body)
}

func GetResources(ctx context.Context, url string, nautobotToken string) (successResponse, error) {

	var m successResponse
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return m, err
	}

	setHeaders(req, nautobotToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return m, err
	}

	if res.StatusCode != http.StatusOK {
		return m, fmt.Errorf("Error while retrieving data from %s. Status: %s", url, res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return m, err
	}
	res.Body.Close()

	return m, nil

}

func Post[T any](ctx context.Context, url string, nautobotToken string, data any) ([]T, error) {

	var m []T
	b, err := toJSON(data)

	if err != nil {
		return nil, err
	}

	byteReader := bytes.NewReader(b)
	req, err := http.NewRequestWithContext(ctx, "POST", url, byteReader)

	if err != nil {
		return nil, err
	}

	setHeaders(req, nautobotToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return m, err
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf("Error while posting data to %s. Status: %s. Message: %s", url, res.Status, string(body))
	}

	res.Body.Close()

	if err := json.Unmarshal(body, &m); err != nil {
		return m, err
	}

	return m, nil
}

func PatchResources[T any](ctx context.Context, url string, nautobotToken string, data any) ([]T, error) {
	fmt.Println(url)
	var m []T
	b, err := toJSON(data)

	if err != nil {
		return m, err
	}

	byteReader := bytes.NewReader(b)
	req, err := http.NewRequestWithContext(ctx, "Patch", url, byteReader)

	if err != nil {
		return m, err
	}

	setHeaders(req, nautobotToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return m, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return m, err
	}

	if res.StatusCode != 200 {

		return m, fmt.Errorf("Error while updating resource %s. Status: %s", url, res.Status)
	}

	res.Body.Close()

	if err := json.Unmarshal(body, &m); err != nil {
		return m, err
	}

	return m, nil
}

func PatchResource[T any](ctx context.Context, url string, nautobotToken string, data any) (T, error) {

	var m T
	b, err := toJSON(data)

	if err != nil {
		return m, err
	}

	byteReader := bytes.NewReader(b)
	req, err := http.NewRequestWithContext(ctx, "Patch", url, byteReader)

	if err != nil {
		return m, err
	}

	setHeaders(req, nautobotToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return m, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return m, err
	}

	if res.StatusCode != 200 {
		errResponse := PatchRequestError{
			StatusCode: res.Status,
			Err:        errors.New("Error in response"),
			Message:    string(body),
		}
		var detaiNotFound DetailNotFound
		if err := json.Unmarshal(body, &detaiNotFound); err != nil {
			return m, &errResponse
		}
		if detaiNotFound.Detail == "Not found." {
			errResponse.DetailNotFound = true
		}

		return m, &errResponse
	}

	res.Body.Close()

	if err != nil {
		return m, err
	}

	return parseJSON[T](body)
}

func Delete(ctx context.Context, url string, nautobotToken string, data any) error {

	b, err := toJSON(data)

	if err != nil {
		return err
	}

	byteReader := bytes.NewReader(b)

	req, err := http.NewRequestWithContext(ctx, "Delete", url, byteReader)

	if err != nil {
		return err
	}

	setHeaders(req, nautobotToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("Error while updating resource %s. Status: %s", url, res.Status)
	} else {
		log.Println("Objects removed.")
	}

	return nil
}

func parseJSON[T any](s []byte) (T, error) {
	var r T
	if err := json.Unmarshal(s, &r); err != nil {
		return r, err
	}
	return r, nil
}

func toJSON(T any) ([]byte, error) {
	return json.Marshal(T)
}

func setHeaders(req *http.Request, nautobotToken string) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", nautobotToken))
	req.Header.Set("Content-Type", "application/json")
}
