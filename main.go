package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type costings struct {
	Electricity []costing
	Gas         []costing
}

type costing struct {
	Reading float32
	Cost    float32
	Date    time.Time
}

type electricity_settings struct {
	Mpan   string  `json:"mpan"`
	Serial string  `json:"serial"`
	Cost   float32 `json:"cost"`
}

type gas_settings struct {
	Mprn   string  `json:"mprn"`
	Serial string  `json:"serial"`
	Cost   float32 `json:"cost"`
}

type octopus_settings struct {
	Apikey      string               `json:"apikey"`
	Electricity electricity_settings `json:"electricity"`
	Gas         gas_settings         `json:"gas"`
}

type influx_settings struct {
	Token  string `json:"token"`
	Bucket string `json:"bucket"`
	Org    string `json:"org"`
	Url    string `json:"url"`
}

type config struct {
	Octopus octopus_settings `json:"octopus"`
	Influx  influx_settings  `json:"influx"`
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

	costings, err := GetConsumption(c.Octopus)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	write_data(c.Influx, costings)

}
