package database

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketPrice struct {
	GUID      uuid.UUID `gorm:"primaryKey" json:"guid"`
	AssetName string    `gorm:"column:asset_name" json:"asset_name"`
	PriceUsdt string    `gorm:"column:price_usdt" json:"price_usdt"`
	Volume    string    `gorm:"column:volume" json:"volume"`
	Rate      string    `gorm:"column:volume" json:"rate"`
	Timestamp uint64
}

type MarketPriceViewer interface {
	QueryMarketPriceByAsset() ([]*MarketPrice, error)
}

type MarketPriceDB interface {
	MarketPriceViewer

	StoreMarketPrice([]*MarketPrice) error
	//Transaction(fn func(db *marketPriceDB) error) error
	GetMarketPriceByAssetName(assetName string) ([]*MarketPrice, error)
}

type marketPriceDB struct {
	gorm *gorm.DB
}

func NewMarketPriceDB(db *gorm.DB) MarketPriceDB {
	return &marketPriceDB{
		gorm: db,
	}
}

func (m *marketPriceDB) QueryMarketPriceByAsset() ([]*MarketPrice, error) {
	var marketPriceList []*MarketPrice
	if err := m.gorm.Table("market_price").Find(&marketPriceList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return marketPriceList, nil
}

func (m *marketPriceDB) StoreMarketPrice(marketPriceList []*MarketPrice) error {
	return m.gorm.Table("market_price").CreateInBatches(marketPriceList, len(marketPriceList)).Error
}

//func (m *marketPriceDB) Transaction(fn func(db *marketPriceDB) error) error {
//	return m.gorm.Transaction(func(tx *gorm.DB) error {
//		return fn(&marketPriceDB{
//			gorm: tx,
//		})
//	})
//}

// GetMarketPriceByAssetName 根据资产名称获取市场价格信息。
// 参数:
//
//	assetName: 资产名称，用于查询市场价。
//
// 返回值:
//
//	[]*MarketPrice: 市场价格信息的切片，如果找不到相关记录，则返回空切片。
//	error: 如果数据库操作出错，则返回错误信息。
func (m *marketPriceDB) GetMarketPriceByAssetName(assetName string) ([]*MarketPrice, error) {
	// 初始化市场价格信息切片。
	var marketPriceList []*MarketPrice

	// 使用gorm.DB的Table方法指定要查询的表，Where方法添加查询条件，Find方法执行查询。
	// 如果查询过程中发生错误，则检查错误类型。
	if err := m.gorm.Table("market_price").Where("asset_name = ?", assetName).Find(&marketPriceList).Error; err != nil {
		// 如果错误为记录未找到，则返回nil, nil表示没有找到相关记录，但不是错误。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// 如果是其他错误，则返回该错误。
		return nil, err
	}

	// 如果查询成功，则返回查询结果。
	return marketPriceList, nil
}
