package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Connection struct {
	db *gorm.DB
}

func (c *Connection) Close() error {
	if c.db == nil {
		return errNotConnected
	}

	return c.db.Close()
}

func (c *Connection) Migrate(models ...interface{}) error {
	if c.db == nil {
		return errNotConnected
	}

	if err := c.db.AutoMigrate(models...).Error; err != nil {
		return fmt.Errorf("database: automigrate: %s", err)
	}
	return nil
}

func (c *Connection) Db() (*gorm.DB, error) {
	if c.db == nil {
		return nil, errNotConnected
	}

	return c.db, nil
}
