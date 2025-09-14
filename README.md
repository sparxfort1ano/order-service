# Демонстрационный сервис с Kafka, PostgreSQL и кэшем

Что может?
- получать данные заказов из очереди сообщений Kafka Apache (Redpanda);
- сохранять их в базу данных (PostgreSQL);
- кэшировать в памяти для быстрого доступа;
- при перезапуске восстанавливать кэш из б/д;
- предоставлять REST API для получения заказа по его id.

---

## Эмулятор заказов

Здесь используется следующий скрипт-эмулятор:

1) .JSON -> топик orders:

$body = (Get-Content -Raw model.json | ConvertFrom-Json | ConvertTo-Json -Compress)
$body | docker exec -i ordersvc-redpanda rpk topic produce orders

2) Контейнеризация и запуск:

docker compose up -d; go run ./cmd/server

---

## Полезные ссылки

HTTP-сервис:  
  [http://localhost:8082] -- здесь открывается HTML-страница для поиска заказа (`/order/{id}` и `/ping` доступны как REST API)

Проверить, запущен ли сервер:
  [http://localhost:8082/ping] -- должно вернуть pong

Redpanda-консоль:
  [http://localhost:8080] -- для просмотра топиков и консьюмеров

---

## Используемые технологии

- Golang
- PostgreSQL (jsonb)
- Redpanda (API Kafka)
- Docker Compose
- HTML/CSS/JavaScript
