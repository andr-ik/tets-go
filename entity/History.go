package entity

import "time"

type History struct {
	Id     int
	FlatId int
	Date   time.Time
	Price  float64
}
