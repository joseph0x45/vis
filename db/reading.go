package db

import (
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

func (c *Conn) GetReadings() ([]models.Reading, error) {
	return nil, nil
}
