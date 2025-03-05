package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"github.com/qiaopengjun5162/web3-market-services/common/retry"
	"github.com/qiaopengjun5162/web3-market-services/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	gorm             *gorm.DB
	MarketPrice      MarketPriceDB
	OfficialCoinRate OfficialCoinRateDB
}

func NewDB(ctx context.Context, dbConfig config.DBConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Name)
	if dbConfig.Port != 0 {
		dsn += fmt.Sprintf(" port=%d", dbConfig.Port)
	}
	if dbConfig.User != "" {
		dsn += fmt.Sprintf(" user=%s", dbConfig.User)
	}
	if dbConfig.Password != "" {
		dsn += fmt.Sprintf(" password=%s", dbConfig.Password)
	}

	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
		CreateBatchSize:        3_000,
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	gorm, err := retry.Do[*gorm.DB](context.Background(), 10, retryStrategy, func() (*gorm.DB, error) {
		gorm, err := gorm.Open(postgres.Open(dsn), &gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return gorm, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &DB{
		gorm:             gorm,
		MarketPrice:      NewMarketPriceDB(gorm),
		OfficialCoinRate: NewOfficialCoinRateDB(gorm),
	}, nil
}

func (db *DB) Transaction(fn func(db *DB) error) error {
	return db.gorm.Transaction(func(tx *gorm.DB) error {
		return fn(&DB{
			gorm:             tx,
			MarketPrice:      NewMarketPriceDB(tx),
			OfficialCoinRate: NewOfficialCoinRateDB(tx),
		})
	})
}

func (db *DB) Close() error {
	dbConn, err := db.gorm.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	return dbConn.Close()
}

func (db *DB) ExecuteSQLMigration(migrationsFolder string) error {
	return filepath.Walk(migrationsFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to walk through migrations folder: %s", path))
		}
		if info.IsDir() {
			return nil
		}

		fileContent, readErr := os.ReadFile(filepath.Clean(path))
		log.Info("path: ", path)
		log.Info("fileContent: ", fileContent)
		if fileContent == nil {
			return errors.Wrap(readErr, fmt.Sprintf("failed to read file: %s", path))
		}
		//if fileContent, readErr := os.ReadFile(path); readErr != nil {
		if readErr != nil {
			return errors.Wrap(readErr, fmt.Sprintf("failed to read file: %s", path))
		}

		execErr := db.gorm.Exec(string(fileContent)).Error
		if execErr != nil {
			return errors.Wrap(execErr, fmt.Sprintf("failed to execute migration: %s", path))
		}
		return nil
	})

}
