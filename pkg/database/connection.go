package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/realbucksavage/robin/pkg/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	connectionFormat = "%s:%s@(%s:%d)/%s?parseTime=true"

	errNotConnected = errors.New("not connected")
)

func NewConnection(config Config) (conn *Connection, err error) {
	if config.MaxRetries <= 0 {
		config.MaxRetries = 5
	}

	for i := 1; i <= config.MaxRetries; i++ {
		c := fmt.Sprintf(connectionFormat, config.Username, config.Password, config.Host, config.Port, config.Database)

		db, err := gorm.Open("mysql", c)
		if err != nil {
			log.L.Errorf("database: connection failed: %s", err)

			// No need to sleep on the last attempt.
			if i != config.MaxRetries {
				time.Sleep(time.Duration(10) * time.Second)
			}
			continue
		}

		db.LogMode(config.LogMode)
		return &Connection{
			db: db,
		}, nil
	}

	return nil, fmt.Errorf("database: gave up after retying for %d times", config.MaxRetries)
}
