package nautobot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
)

const (
	InterfacesStatusEndpoint = "api/plugins/interfaces-telemetry/interfaces-status"
)

type InterfacesStatus struct {
	Display         string    `json:"display"`
	ID              string    `json:"id"`
	URL             string    `json:"url"`
	InterfaceName   string    `json:"interface_name"`
	InterfaceID     string    `json:"interface_id"`
	InterfaceStatus string    `json:"interface_status"`
	DeviceName      string    `json:"device_name"`
	DeviceID        string    `json:"device_id"`
	Notes           string    `json:"notes_url"`
	Created         string    `json:"created"`
	LastUpdated     time.Time `json:"last_updated"`
}

func (c *Client) GetInterfacesStatus(ctx context.Context, filterOtions []OptionalFilter) ([]InterfacesStatus, error) {

	var IfacesStatus []*InterfacesStatus
	var IfacesStatusToValues []InterfacesStatus
	var lastPage bool
	options := url.Values{}

	link := fmt.Sprintf("%s/%s/", c.BaseURL, InterfacesStatusEndpoint)
	for _, filterOp := range filterOtions {

		if filterValueString, ok := filterOp.Value.(string); ok {
			options.Add(filterOp.Field, filterValueString)
		}
		if filterValueInt, ok := filterOp.Value.(int); ok {
			options.Add(filterOp.Field, strconv.Itoa(filterValueInt))
		}
	}

	if len(options) != 0 {
		link = fmt.Sprintf("%s?%s", link, options.Encode())
	}

	for lastPage != true {
		var pageResults []*InterfacesStatus
		response, err := GetResources(ctx, link, c.nautobotToken)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(response.Results, &pageResults); err != nil {
			return nil, err
		}

		lastPage = response.isLastPage()
		IfacesStatus = append(IfacesStatus, pageResults...)
		link = fmt.Sprint(response.Next)
	}

	for _, v := range IfacesStatus {
		IfacesStatusToValues = append(IfacesStatusToValues, *v)
	}

	return IfacesStatusToValues, nil
}

// Interface Status Slug like: r2__ethernet1
func (c *Client) GetInterfaceStatus(ctx context.Context, interfaceStatusSlug string) (InterfacesStatus, error) {

	link := fmt.Sprintf("%s/%s/%s/", c.BaseURL, InterfacesStatusEndpoint, interfaceStatusSlug)
	m, err := GetResource[InterfacesStatus](ctx, link, c.nautobotToken)

	if err != nil {
		log.Print(err)
		return InterfacesStatus{}, err
	}

	return m, nil
}

func (c *Client) PostInterfacesStatus(ctx context.Context, data any) ([]InterfacesStatus, error) {

	link := fmt.Sprintf("%s/%s/", c.BaseURL, InterfacesStatusEndpoint)
	m, err := Post[InterfacesStatus](ctx, link, c.nautobotToken, data)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return m, nil
}

// Interface Status Slug like: r2__ethernet1
func (c *Client) PatchInterfaceStatus(ctx context.Context, interfaceStatusSlug string, data any) (InterfacesStatus, error) {

	link := fmt.Sprintf("%s/%s/%s/", c.BaseURL, InterfacesStatusEndpoint, interfaceStatusSlug)
	m, err := PatchResource[InterfacesStatus](ctx, link, c.nautobotToken, data)

	if err != nil {
		log.Print(err)
		return m, err
	}

	return m, nil
}

func (c *Client) DeleteInterfaceStatus(ctx context.Context, interfaceStatusSlug string) error {

	link := fmt.Sprintf("%s/%s/%s/", c.BaseURL, InterfacesStatusEndpoint, interfaceStatusSlug)
	err := Delete(ctx, link, c.nautobotToken, []string{})

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
