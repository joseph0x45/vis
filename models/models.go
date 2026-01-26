package models

import "html/template"

type PageData struct {
	Content     template.HTML
	CurrentKWh  int
	LastUpdated string
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
}
