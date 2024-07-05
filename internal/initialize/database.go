package initialize

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/wopoczynski/todoapp/internal/database"
)

type DBConfig struct {
	DSN string `env:"URL"`
}

func DB(cfg DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return db, nil
}

func Automigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&database.TodoModel{})
	if err != nil {
		return err
	}
	return nil
}
