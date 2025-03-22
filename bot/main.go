package main

import (
	"log"
	"my_bot_Go/bot/center"
	"os"
	"os/signal"
	"syscall"
)

//  Для получения средней цены биткоина в день: http://localhost:8080/bitcoin_prices
//  Для получения последней цены BTC: http://localhost:8080/price
//  Для загрузки HTML-страницы: http://localhost:8080/

// Главная функция
func main() {
	// Запускаем WebSocket-клиент в goroutine
	go center.PriceBtc()
	// Запускаем HTTP-серве
	center.StartServer()

	// Канал для обработки сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Ждем сигнала завершения
	sig := <-signalChan
	log.Printf("Получен сигнал завершения: %s", sig)
}
