package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

const (
	dbUserEnv     = "DB_USER"
	dbPasswordEnv = "DB_PASSWORD"
	dbHostEnv     = "DB_HOST"
	dbPortEnv     = "DB_PORT"
	dbNameEnv     = "DB_NAME"

	kafkaBrokerEnv = "KAFKA_BROKER"
	kafkaTopicEnv  = "KAFKA_TOPIC"

	appPortEnv = "APP_PORT"
)

type Config struct {
	AppPort     string
	DatabaseURL string // postgres://...
	KafkaBroker string
	KafkaTopic  string
}

// .env reader and parser
func Load() (*Config, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}

	// Slice of required keys to check
	requiredKeys := []string{dbUserEnv, dbPasswordEnv, dbHostEnv, dbPortEnv, dbNameEnv, kafkaBrokerEnv, kafkaTopicEnv, appPortEnv}
	for _, key := range requiredKeys {
		if myEnv[key] == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		myEnv[dbUserEnv],
		myEnv[dbPasswordEnv],
		myEnv[dbHostEnv],
		myEnv[dbPortEnv],
		myEnv[dbNameEnv])

	return &Config{
		AppPort:     myEnv[appPortEnv],
		DatabaseURL: databaseURL,
		KafkaBroker: myEnv[kafkaBrokerEnv],
		KafkaTopic:  myEnv[kafkaTopicEnv],
	}, nil
}
