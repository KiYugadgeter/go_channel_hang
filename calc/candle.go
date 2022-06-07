package calc

import (
	"time"
)

type Candle struct {
	Count    int       `json:"Count,omitempty"`
	Open     float64   `json:"Open"`
	Close    float64   `json:"Close"`
	Low      float64   `json:"Low"`
	High     float64   `json:"High"`
	Datetime time.Time `json:"Datetime"`
}
