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

func (c *Conn) GetPurchases() ([]models.Purchase, error) {
	const query = "select * from purchases order by timestamp desc"
	purchases := []models.Purchase{}
	if err := c.db.Select(&purchases, query); err != nil {
		return nil, fmt.Errorf("Errorr while getting purchases: %w", err)
	}
	return purchases, nil
}
