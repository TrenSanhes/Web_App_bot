function fetchAveragePrice() {
    fetch('/bitcoin_prices')
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка при получении средней цены');
            }
            return response.text();
        })
        .then(data => {
            document.getElementById("average-price").innerText = data; // Обновляем блок с средней ценой
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}

function fetchCurrentPrice() {
    fetch('/price')
        .then(response => {
            if (!response.ok) {
                throw new Error('Сервер еще не обновил цену');
            }
            return response.json();
        })
        .then(data => {
            document.getElementById("current-price").innerText = `Текущая цена BTC: ${data.price}`; // Обновление текущей цены
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}

// Запускаем получение текущей цены при загрузке страницы
document.addEventListener('DOMContentLoaded', (event) => {
    fetchCurrentPrice();
});

function fetchPricesOneYear() {
    fetch('/bitcoin_prices_one_year')
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка при получении цен за 1 год');
            }
            return response.text();
        })
        .then(data => {
            document.getElementById("one-year-prices").innerText = data; // Обновляем блок с ценами за 1 год
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}