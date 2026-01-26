package db

import "github.com/joseph0x45/sad"

var migrations = []sad.Migration{
	{
		Version: 1,
		Name:    "create_readings",
		SQL: `
      create table readings (
        id text not null primary key,
        timestamp integer not null,
        kwh real not null,
        date_str text not null unique
      );
    `,
	},
	{
		Version: 2,
		Name:    "create_purchases",
		SQL: `
      create table purchases (
        id text not null primary key,
        timestamp integer not null,
        kwh real not null,
        cost integer not null,
        date_str text not null
      );
    `,
	},
}
