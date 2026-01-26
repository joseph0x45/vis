package models

import "html/template"

type PageData struct {
	Content     template.HTML
	CurrentKWh  int
	LastUpdated string
}
