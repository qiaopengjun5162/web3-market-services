package database

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfficialCoinRate struct {
	GUID      uuid.UUID `gorm:"primaryKey" json:"guid"`
	AssetName string    `gorm:"column:asset_name" json:"asset_name"`
	BaseAsset string    `gorm:"column:base_asset" json:"base_asset"`
	Price     string    `gorm:"column:price" json:"price"`
	Timestamp uint64
}

type OfficialCoinRateView interface {
	QueryOfficialCoinRateByAsset() ([]*OfficialCoinRate, error)
}

type OfficialCoinRateDB interface {
	OfficialCoinRateView

	StoreOfficialCoinRate([]*OfficialCoinRate) error
	//GetOfficialCoinRateByAssetName(assetName string) ([]*OfficialCoinRate, error)
}

type officialCoinRateDB struct {
	gorm *gorm.DB
}

func NewOfficialCoinRateDB(db *gorm.DB) OfficialCoinRateDB {
	return &officialCoinRateDB{
		gorm: db,
	}
}

func (o *officialCoinRateDB) QueryOfficialCoinRateByAsset() ([]*OfficialCoinRate, error) {
	var officialCoinRateList []*OfficialCoinRate
	err := o.gorm.Table("official_coin_rate").Find(&officialCoinRateList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return officialCoinRateList, nil
}

func (o *officialCoinRateDB) StoreOfficialCoinRate(officialCoinRateList []*OfficialCoinRate) error {
	return o.gorm.Table("official_coin_rate").CreateInBatches(officialCoinRateList, len(officialCoinRateList)).Error
}

func (o *officialCoinRateDB) GetOfficialCoinRateByAssetName(assetName string) ([]*OfficialCoinRate, error) {
	var officialCoinRateList []*OfficialCoinRate
	err := o.gorm.Table("official_coin_rate").Where("asset_name = ?", assetName).Find(&officialCoinRateList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return officialCoinRateList, nil
}
