package types

import "time"

type Costings struct {
	Electricity []Costing
	Gas         []Costing
}

type Costing struct {
	Reading float32
	Cost    float32
	Date    time.Time
}
