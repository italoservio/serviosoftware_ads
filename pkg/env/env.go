package env

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Env struct {
	MONGODB_URI              string
	AUTH_SECRET              string
	SERVER_BASE_URL          string
	NETIFY_BASE_URL          string
	NETIFY_DATA_FEED_API_KEY string
}

func Load() *Env {
	env := &Env{}

	env.MONGODB_URI = os.Getenv("MONGODB_URI")
	env.AUTH_SECRET = os.Getenv("AUTH_SECRET")
	env.SERVER_BASE_URL = os.Getenv("SERVER_BASE_URL")
	env.NETIFY_BASE_URL = os.Getenv("NETIFY_BASE_URL")
	env.NETIFY_DATA_FEED_API_KEY = os.Getenv("NETIFY_DATA_FEED_API_KEY")

	if env.MONGODB_URI == "" {
		panic("variavel de ambiente MONGODB_URI nao definida")
	}

	if env.AUTH_SECRET == "" {
		panic("variavel de ambiente AUTH_SECRET nao definida")
	}

	if env.SERVER_BASE_URL == "" {
		panic("variavel de ambiente SERVER_BASE_URL nao definida")
	}

	if env.NETIFY_BASE_URL == "" {
		panic("variavel de ambiente NETIFY_BASE_URL nao definida")
	}

	if env.NETIFY_DATA_FEED_API_KEY == "" {
		panic("variavel de ambiente NETIFY_DATA_FEED_API_KEY nao definida")
	}

	return env
}
