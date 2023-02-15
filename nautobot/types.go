package nautobot

import (
	"encoding/json"
	"fmt"
)

type Status struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type CustomFields struct {
}

type successResponse struct {
	Count    int             `json:"count"`
	Next     interface{}     `json:"next"`
	Previous interface{}     `json:"previous"`
	Results  json.RawMessage `json:"results"`
}

func (res *successResponse) isLastPage() bool {
	switch v := res.Next.(type) {
	case string:
		fmt.Println("Next page: ", v)
		return false
	default:
		return true
	}
}
