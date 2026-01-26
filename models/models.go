package models

import "html/template"

type PageData struct {
	ChartData   map[string]any
	Content     template.HTML
	CurrentKWh  float64
	LastUpdated string
	Timestamp   int64
}

type Reading struct {
	ID         string  `db:"id"`
	Timestamp  int64   `db:"timestamp"`
	Kwh        float64 `db:"kwh"`
	DateString string  `db:"date_str"`
}

type Purchase struct {
	ID         string  `db:"id"`
	Timestamp  int64   `db:"timestamp"`
	Kwh        float64 `db:"kwh"`
	Cost       int     `db:"cost"`
	DateString string  `db:"date_str"`
	TimeString string
}
