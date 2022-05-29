package main

import (
	"fmt"
	"log"
	"os"

	"energy.echo-moo.co.uk/internal/influx"
	"energy.echo-moo.co.uk/internal/octopus"
	"github.com/spf13/viper"
)

type config struct {
	Octopus octopus.OctopusSettings `json:"octopus"`
	Influx  influx.InfluxSettings   `json:"influx"`
}

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	c := config{}
	viper.Unmarshal(&c)

	octopus := octopus.New(c.Octopus)
	influx := influx.New(c.Influx)

	costings, err := octopus.GetConsumption()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	influx.WriteData(costings)

}
