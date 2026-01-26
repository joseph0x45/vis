package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/joseph0x45/vis/models"
)

func (c *Conn) InsertReading(reading *models.Reading) error {
	const query = `
    insert into readings (
      id, timestamp, kwh, date_str
    )
    values (
      :id, :timestamp, :kwh, :date_str
    )
    on conflict(date_str) do update set
      timestamp = excluded.timestamp,
      kwh = excluded.kwh;
  `
	if _, err := c.db.NamedExec(query, reading); err != nil {
		return fmt.Errorf("Error while inserting reading: %w", err)
	}
	return nil
}

func (c *Conn) GetReadings(filterForDashboard bool) ([]models.Reading, error) {
	const query = "select * from readings order by timestamp desc"
	const dashboardQuery = "select * from readings order by timestamp asc"
	readings := []models.Reading{}
	var err error
	if filterForDashboard {
		err = c.db.Select(&readings, dashboardQuery)
	} else {
		err = c.db.Select(&readings, query)
	}
	if err != nil {
		return nil, fmt.Errorf("Error while getting readings: %w", err)
	}
	return readings, nil
}

func (c *Conn) GetLatestReading() (*models.Reading, error) {
	reading := &models.Reading{}
	const query = "select * from readings order by timestamp desc limit 1"
	if err := c.db.Get(reading, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting latest reading: %w", err)
	}
	return reading, nil
}
