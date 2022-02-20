package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
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

const uri = "https://api.octopus.energy"
const api_electricity_fmt = "/v1/electricity-meter-points/%s/meters/%s/consumption/"
const api_gas_fmt = "/v1/gas-meter-points/%s/meters/%s/consumption/"
const query_periodfrom = "period_from=%s"

func getReadings(u *url.URL, cost float32) ([]costing, error) {
	page := u.String()

	var c []costing

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

			c = append(c, costing{
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
			next_page.RawQuery = u.RawQuery

			page = next_page.String()
		} else {
			page = ""
		}

	}

	return c, nil
}

func GetConsumption(os octopus_settings) (costings, error) {
	var c costings

	u, err := url.Parse(uri)

	if err != nil {
		log.Fatal(err)
		return c, err
	}

	u.User = url.User(os.Apikey)
	u.Path = fmt.Sprintf(api_electricity_fmt, os.Electricity.Mpan, os.Electricity.Serial)

	t := time.Now()
	past := t.Add(time.Duration(-48) * time.Hour)
	q := u.Query()
	q.Set("period_from", past.Format(time.RFC3339))
	u.RawQuery = q.Encode()

	c.Electricity, err = getReadings(u, os.Electricity.Cost)

	u.Path = fmt.Sprintf(api_gas_fmt, os.Gas.Mprn, os.Gas.Serial)
	c.Gas, err = getReadings(u, os.Gas.Cost)

	return c, nil
}
