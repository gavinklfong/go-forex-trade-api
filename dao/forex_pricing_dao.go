package dao

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/gavinklfong/go-forex-trade-api/model"
)

type ForexPricingDaoImpl struct {
	entries map[string]model.ForexPricing
}

func NewForexPricingDao(csvFilePath string) (ForexPricingDao, error) {
	entries, err := initializeForexPricing(csvFilePath)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to initialize forex pricing with %s, error: %v", csvFilePath, err))
		return nil, err
	}
	return &ForexPricingDaoImpl{entries}, nil
}

func initializeForexPricing(filePath string) (map[string]model.ForexPricing, error) {
	records, err := readCsvFile(filePath)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to read CSV %s", filePath))
		return nil, err
	}

	entries := make(map[string]model.ForexPricing)
	for _, v := range records {

		buyPip, err := strconv.ParseFloat(v[2], 32)
		if err != nil {
			return nil, fmt.Errorf("BuyPip %s is not a decimal", v[2])
		}

		sellPip, err := strconv.ParseFloat(v[3], 32)
		if err != nil {
			return nil, fmt.Errorf("SellPip %s is not a decimal", v[3])
		}

		newEntry := model.ForexPricing{
			BaseCurrency:    v[0],
			CounterCurrency: v[1],
			BuyPip:          float32(buyPip),
			SellPip:         float32(sellPip),
		}

		entries[fmt.Sprintf("%s-%s", v[0], v[1])] = newEntry
	}

	return entries, nil
}

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to read input file %s, %s", filePath, err))
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to parse file as CSV for %s, %s", filePath, err))
		return nil, err
	}

	return records, nil
}

func (s *ForexPricingDaoImpl) GetPricingByCurrencyPair(base, counter string) *model.ForexPricing {
	pricing, exist := s.entries[fmt.Sprintf("%s-%s", base, counter)]
	if !exist {
		return nil
	}
	return &pricing
}
