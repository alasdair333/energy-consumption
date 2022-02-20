package main

import (
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func write_data(is influx_settings, c costings) {

	client := influxdb2.NewClient(is.Url, is.Token)
	defer client.Close()

	// get non-blocking write client
	writeAPI := client.WriteAPI(is.Org, is.Bucket)

	for _, point := range c.Electricity {

		fmt.Println("Consumption:", point.Reading)
		fmt.Println("Cost:", point.Cost)
		fmt.Println("Time:", point.Date)

		p := influxdb2.NewPointWithMeasurement("Octopus").
			AddTag("Type", "Electricity").
			AddField("consumption", point.Reading).
			AddField("cost", point.Cost).
			SetTime(point.Date)
		// write point asynchronously
		writeAPI.WritePoint(p)
		// Flush writes
		writeAPI.Flush()
	}

	for _, point := range c.Gas {

		fmt.Println("Consumption:", point.Reading)
		fmt.Println("Cost:", point.Cost)
		fmt.Println("Time:", point.Date)

		p := influxdb2.NewPointWithMeasurement("Octopus").
			AddTag("Type", "Gas").
			AddField("consumption", point.Reading).
			AddField("cost", point.Cost).
			SetTime(point.Date)
		// write point asynchronously
		writeAPI.WritePoint(p)
		// Flush writes
		writeAPI.Flush()
	}

}
