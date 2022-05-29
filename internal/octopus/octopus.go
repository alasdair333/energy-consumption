package octopus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"energy.echo-moo.co.uk/internal/types"
)

const (
	uri                 = "https://api.octopus.energy"
	api_electricity_fmt = "/v1/electricity-meter-points/%s/meters/%s/consumption/"
	api_gas_fmt         = "/v1/gas-meter-points/%s/meters/%s/consumption/"
	query_periodfrom    = "period_from=%s"
)

type meter_points struct {
	Consumpton    float32 `json:"consumption"`
	IntervalStart string  `json:"interval_start"`
	IntervalEnd   string  `json:"interval_end"`
}

type consumptionResp struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Readings []meter_points `json:"results"`
}

type Octopus struct {
	settings OctopusSettings
}

func New(s OctopusSettings) *Octopus {
	return &Octopus{
		settings: s,
	}
}

func (o *Octopus) getReadings(u *url.URL, cost float32) ([]types.Costing, error) {
	page := u.String()

	var c []types.Costing

	for page != "" {
		response, err := http.Get(page)
		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		report := consumptionResp{}

		json.Unmarshal([]byte(responseData), &report)

		fmt.Println("We have:", report.Count)
		fmt.Println("Next page:", report.Next)
		fmt.Println("Previous page:", report.Previous)
		fmt.Println("Results:", len(report.Readings))

		for _, r := range report.Readings {

			t, err := time.Parse(time.RFC3339, r.IntervalEnd)

			if err != nil {
				return c, err
			}

			c = append(c, types.Costing{
				Reading: r.Consumpton,
				Cost:    r.Consumpton * cost,
				Date:    t,
			})
		}

		if report.Next != "" {
			next_page, err := url.Parse(report.Next)
			if err != nil {
				return nil, err
			}
			next_page.User = u.User
			//next_page.RawQuery = u.RawQuery

			page = next_page.String()
		} else {
			page = ""
		}

	}

	return c, nil
}

func (o *Octopus) GetConsumption() (types.Costings, error) {
	var c types.Costings

	u, err := url.Parse(uri)

	if err != nil {
		log.Fatal(err)
		return c, err
	}

	u.User = url.User(o.settings.Apikey)
	u.Path = fmt.Sprintf(api_electricity_fmt, o.settings.Electricity.Mpan, o.settings.Electricity.Serial)

	t := time.Now()
	past := t.Add(time.Duration(-168) * time.Hour)
	q := u.Query()
	q.Set("period_from", past.Format(time.RFC3339))
	u.RawQuery = q.Encode()

	c.Electricity, err = o.getReadings(u, o.settings.Electricity.Cost)

	u.Path = fmt.Sprintf(api_gas_fmt, o.settings.Gas.Mprn, o.settings.Gas.Serial)
	c.Gas, err = o.getReadings(u, o.settings.Gas.Cost)

	return c, nil
}
