package config

import "os"

// Конфигурация приложения
type Config struct {
	HTTPAddr     string // адрес, на котором запускается HTTP-сервер
	DBDSN        string // строка подключения к б/д Postgres
	KafkaBroker  string // адрес брокера Kafka -- Redpanda
	KafkaTopic   string // название топика для сообщений
	KafkaGroupID string // идентификатор группы потребителей
}

// Функция для чтения переменной окружения
func getenv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def // значение по умолчанию
}

// Загружаем конфигурацию (берем из окружения или используем стандартные значения)
func Load() Config {
	return Config{
		HTTPAddr:     getenv("HTTP_ADDR", ":8082"),
		DBDSN:        getenv("DB_DSN", "postgres://app:app@localhost:5432/orders"),
		KafkaBroker:  getenv("KAFKA_BROKER", "localhost:9093"),
		KafkaTopic:   getenv("KAFKA_TOPIC", "orders"),
		KafkaGroupID: getenv("KAFKA_GROUP", "order-service"),
	}
}
