package db

import (
	"fmt"

	"github.com/joseph0x45/vis/models"
)

func (c *Conn) InsertPurchase(purchase *models.Purchase) error {
	const query = `
    insert into purchases (
      id, timestamp, kwh, cost, date_str
    )
    values (
      :id, :timestamp, :kwh, :cost, :date_str
    );
  `
	if _, err := c.db.NamedExec(query, purchase); err != nil {
		return fmt.Errorf("Error while inserting purchase: %w", err)
	}
	return nil
}
