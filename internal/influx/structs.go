package influx

type InfluxSettings struct {
	Token  string `json:"token"`
	Bucket string `json:"bucket"`
	Org    string `json:"org"`
	Url    string `json:"url"`
}
