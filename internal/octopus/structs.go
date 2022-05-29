package octopus

type ElectricitySettings struct {
	Mpan   string  `json:"mpan"`
	Serial string  `json:"serial"`
	Cost   float32 `json:"cost"`
}

type GasSettings struct {
	Mprn   string  `json:"mprn"`
	Serial string  `json:"serial"`
	Cost   float32 `json:"cost"`
}

type OctopusSettings struct {
	Apikey      string              `json:"apikey"`
	Electricity ElectricitySettings `json:"electricity"`
	Gas         GasSettings         `json:"gas"`
}
