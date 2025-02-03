package database_test

import (
	"context"

	database "github.com/coopnorge/go-redis-facade"
)

func Example() {
	encryptionConfig := database.EncryptionConfig{
		RedisKeyURI: "",
		Aad:         []byte{},
	}
	encryptor, err := database.NewEncryptionClient(encryptionConfig)
	if err != nil {
		panic(err)
	}

	redisConfig := database.Config{
		Address:           "",
		Password:          "",
		Database:          0,
		DialTimeout:       0,
		WriteTimeout:      0,
		ReadTimeout:       0,
		MaxRetries:        0,
		PoolSize:          0,
		PoolTimeout:       0,
		EncryptionEnabled: false,
	}
	redFacade, err := database.NewRedisFacade(redisConfig, encryptor)
	if err != nil {
		panic(err)
	}

	value, err := redFacade.Find(context.Background(), "key")
	if err != nil {
		panic(err)
	}

	println(value)
}
