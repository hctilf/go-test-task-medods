package postgres

import (
	"fmt"
	"sync"

	"github.com/hctilf/go-test-task-medods/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	databaseInstance *gorm.DB
	once             sync.Once
)

func GetDB(conf *config.Config) *gorm.DB {
	once.Do(func() {
		databaseInstance = newDB(conf)
	})

	return databaseInstance
}

func newDB(conf *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Database,
		conf.Postgres.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil
	}

	return db
}

func DoMigrate(db *gorm.DB, entities ...interface{}) error {
	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			return err
		}
	}

	return nil
}
