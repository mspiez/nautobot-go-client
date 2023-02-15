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
	SitesEndpoint = "api/dcim/sites"
)

type Region struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Depth   int    `json:"_depth"`
	Display string `json:"display"`
}

type Site struct {
	ID                  string        `json:"id"`
	URL                 string        `json:"url"`
	Name                string        `json:"name"`
	Slug                string        `json:"slug"`
	Status              Status        `json:"status"`
	Region              Region        `json:"region"`
	Tenant              interface{}   `json:"tenant"`
	Facility            string        `json:"facility"`
	Asn                 int           `json:"asn"`
	TimeZone            string        `json:"time_zone"`
	Description         string        `json:"description"`
	PhysicalAddress     string        `json:"physical_address"`
	ShippingAddress     string        `json:"shipping_address"`
	Latitude            interface{}   `json:"latitude"`
	Longitude           interface{}   `json:"longitude"`
	ContactName         string        `json:"contact_name"`
	ContactPhone        string        `json:"contact_phone"`
	ContactEmail        string        `json:"contact_email"`
	Comments            string        `json:"comments"`
	Tags                []interface{} `json:"tags"`
	CustomFields        CustomFields  `json:"custom_fields"`
	Created             string        `json:"created"`
	LastUpdated         time.Time     `json:"last_updated"`
	CircuitCount        int           `json:"circuit_count"`
	DeviceCount         int           `json:"device_count"`
	PrefixCount         int           `json:"prefix_count"`
	RackCount           int           `json:"rack_count"`
	VirtualmachineCount int           `json:"virtualmachine_count"`
	VlanCount           int           `json:"vlan_count"`
	Display             string        `json:"display"`
}

func (c *Client) GetSites(ctx context.Context, filterOtions []OptionalFilter) ([]Site, error) {

	var Sites []*Site
	var SitesToValues []Site
	var lastPage bool
	options := url.Values{}

	link := fmt.Sprintf("%s/%s/", c.BaseURL, SitesEndpoint)
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
		var pageResults []*Site
		response, err := GetResources(ctx, link, c.nautobotToken)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(response.Results, &pageResults); err != nil {
			return nil, err
		}

		lastPage = response.isLastPage()
		Sites = append(Sites, pageResults...)
		link = fmt.Sprint(response.Next)
	}

	for _, v := range Sites {
		SitesToValues = append(SitesToValues, *v)
	}

	return SitesToValues, nil
}

func (c *Client) GetSite(ctx context.Context, siteID string) (Site, error) {

	link := fmt.Sprintf("%s/%s/%s/", c.BaseURL, SitesEndpoint, siteID)
	m, err := GetResource[Site](ctx, link, c.nautobotToken)

	if err != nil {
		log.Print(err)
		return Site{}, err
	}

	return m, nil
}

func (c *Client) PostSites(ctx context.Context, data any) ([]Site, error) {

	link := fmt.Sprintf("%s/%s/", c.BaseURL, SitesEndpoint)
	m, err := Post[Site](ctx, link, c.nautobotToken, data)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return m, nil
}

func (c *Client) PatchSites(ctx context.Context, data any) ([]Site, error) {

	link := fmt.Sprintf("%s/%s/", c.BaseURL, SitesEndpoint)
	m, err := PatchResources[Site](ctx, link, c.nautobotToken, data)

	if err != nil {
		log.Print(err)
		return m, err
	}

	return m, nil
}

func (c *Client) PatchSite(ctx context.Context, data any, id string) (Site, error) {

	link := fmt.Sprintf("%s/%s/%s/", c.BaseURL, SitesEndpoint, id)
	m, err := PatchResource[Site](ctx, link, c.nautobotToken, data)

	if err != nil {
		log.Print(err)
		return m, err
	}

	return m, nil
}

func (c *Client) DeleteSites(ctx context.Context, data any) error {

	link := fmt.Sprintf("%s/%s/", c.BaseURL, SitesEndpoint)
	err := Delete(ctx, link, c.nautobotToken, data)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (c *Client) DeleteSite(ctx context.Context, id string) error {

	link := fmt.Sprintf("%s/%s/%s/", c.BaseURL, SitesEndpoint, id)
	err := Delete(ctx, link, c.nautobotToken, []string{})

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
