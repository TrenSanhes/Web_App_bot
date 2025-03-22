package center

import (
	"fmt"
	"net/http" // Импортируем пакет net/http для обработки HTTP-запросов.
	"time"     // Импортируем пакет time для работы с временными метками.

	"github.com/go-resty/resty/v2" // Импортируем библиотеку resty для выполнения HTTP-запросов.
)

const (
	coinGeckoAPI_5 = "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart"
	number         = 365
	age            = 5
)

// Глобальная переменная для кэша
var cachedPrices map[string]float64
var lastFetchTime time.Time
var cacheDuration = time.Hour // продолжительность кэширования

// PriceResponse представляет структуру для ответа API.
type PriceResponse struct {
	Prices [][]float64 `json:"prices"`
}

// GetBitcoinPrices возвращает среднюю дневную цену биткоина за последние 5 лет.
func GetBitcoinPrices() (map[string]float64, error) {
	// Проверяем, нужно ли обновить кэш
	if cachedPrices != nil && time.Since(lastFetchTime) < cacheDuration {
		return cachedPrices, nil
	}

	// Если кэш устарел, делаем запрос к API
	days := number * age
	url := fmt.Sprintf("%s?vs_currency=usd&days=%d", coinGeckoAPI_5, days)

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

	// Сохраняем данные в кэш
	cachedPrices = averagePrices
	lastFetchTime = time.Now() // обновляем время последнего запроса

	return averagePrices, nil
}
