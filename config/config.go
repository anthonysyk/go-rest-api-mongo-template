package config

import (
	"fmt"
	"os"
)

type Config struct {
	FileStorageRootPath string
	MongoURI            string
	AuthSecret          string
}

func New() Config {
	fileStorageRootPath := os.Getenv("FILE_STORAGE_ROOT_PATH")
	mongoURI := os.Getenv("MONGO_URI")
	authSecret := os.Getenv("AUTH_SECRET")
	fmt.Println("test", mongoURI)
	return Config{
		FileStorageRootPath: fileStorageRootPath,
		MongoURI:            mongoURI,
		AuthSecret:          authSecret,
	}
}
