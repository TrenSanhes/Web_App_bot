package center

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	coinGeckoAPI = "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart"
	oneYear      = 365 // Количество дней в одном году
)

// Получить средние цены биткойна за один год
func GetBitcoinPricesOneYear() (map[string]float64, error) {
	oneYearDays := oneYear
	url := fmt.Sprintf("%s?vs_currency=usd&days=%d", coinGeckoAPI, oneYearDays)

	client := resty.New()
	resp, err := client.R().SetResult(&PriceResponse{}).Get(url)

	if err != nil {
		return nil, err
	}

	result := resp.Result().(*PriceResponse)

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("ошибка при получении данных: %v", resp.Status())
	}

	averagePrices := make(map[string]float64)
	dailyPrices := make(map[string][]float64)

	for _, price := range result.Prices {
		t := time.Unix(int64(price[0]/1000), 0).Format("02.01.2006")
		dailyPrices[t] = append(dailyPrices[t], price[1])
	}

	for date, prices := range dailyPrices {
		total := 0.0
		for _, p := range prices {
			total += p
		}
		averagePrices[date] = total / float64(len(prices))
	}

	return averagePrices, nil
}
