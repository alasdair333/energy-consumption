package influx

import (
	"fmt"

	"energy.echo-moo.co.uk/internal/types"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Influx struct {
	settings InfluxSettings
}

func New(s InfluxSettings) *Influx {
	return &Influx{
		settings: s,
	}
}

func (i *Influx) WriteData(c types.Costings) {

	client := influxdb2.NewClient(i.settings.Token, i.settings.Token)
	defer client.Close()

	// get non-blocking write client
	writeAPI := client.WriteAPI(i.settings.Org, i.settings.Bucket)

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
