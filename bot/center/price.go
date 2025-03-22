package center

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Структура для данных о сделках
type TradeMessage struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Price     string `json:"p"` // Цена в виде строки
	Quantity  string `json:"q"` // Объем в виде строки
}

// Глобальные переменные
var (
	lastPrice string
	mu        sync.Mutex // Защита доступа к переменной lastPrice
)

// Обработчик для получения средней цены биткойна за один год
func BitcoinPricesOneYearHandler(c *gin.Context) {
	prices, err := GetBitcoinPricesOneYear()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	result := ""
	for date, price := range prices {
		result += fmt.Sprintf("%s - средняя цена биткоина в этот день: %.2f\n", date, price)
	}

	c.String(200, result)
}

// BitcoinPricesHandler обрабатывает запросы для получения средних цен биткоина.
func BitcoinPricesHandler(c *gin.Context) {
	prices, err := GetBitcoinPrices()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	result := ""
	for date, price := range prices {
		result += fmt.Sprintf("%s - средняя цена биткоина в этот день: %.2f\n", date, price)
	}

	c.String(200, result)
}

func PriceBtc() {
	url := "wss://stream.binance.com:9443/ws/btcusdt@trade"

	for {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Ошибка при установлении соединения: %v. Повторная попытка через 5 секунд...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Успешное подключение к WebSocket")
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Ошибка при чтении сообщения: %v. Повторное подключение...", err)
				break // Выход из цикла чтения для переподключения
			}
			handleTradeMessage(msg)
		}
	}
}

// Функция для обработки сообщений о сделках
func handleTradeMessage(msg []byte) {
	var trade TradeMessage
	err := json.Unmarshal(msg, &trade)
	if err != nil {
		log.Printf("Ошибка при разборе сообщения: %v\n", err)
		return
	}

	mu.Lock()
	lastPrice = trade.Price
	mu.Unlock()
}

// Handler для получения текущей цены BTC
func getPriceHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	if lastPrice == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Price not available yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": lastPrice}) // Возвращаем последнюю цену в формате JSON
}

// HTML-страница для отображения цены
func pricePageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// Функция для запуска сервера с использованием Gin
func StartServer() {
	// Создаем новый экземпляр рутера Gin с предустановленными настройками.
	router := gin.Default()
	// Загружаем HTML-шаблоны
	router.LoadHTMLFiles("temp/index.html") // Предполагается, что файл index.html находится в папке templates

	// крс, это открываются странички в браузере
	router.GET("/", pricePageHandler)                                   // Главная страница с HTML
	router.GET("/bitcoin_prices", BitcoinPricesHandler)                 // Для получения средней цены биткоина в день
	router.GET("/bitcoin_prices_one_year", BitcoinPricesOneYearHandler) // Для получения средней цены биткоина за 1 год
	router.GET("/price", getPriceHandler)                               // Для получения последней цены BTC
	router.Static("./temp", "./temp")                                   // где ./temp - это директория с вашими статическими файлами

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v\n", err)
	}
}
