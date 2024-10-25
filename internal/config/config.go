package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func (c *Config) Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка чтения из .env")
	}

	c.Port = os.Getenv("TODO_PORT")
}
