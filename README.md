# Go Redis Facade

![Lint](https://github.com/coopnorge/go-redis-facade/actions/workflows/lint.yml/badge.svg)
![Build](https://github.com/coopnorge/go-redis-facade/actions/workflows/master-test.yml/badge.svg)

Coop Redis Facade wraps simple interaction with 
Redis clients for CRUD operations by preventing
race conditions between multiple client instances
against singular instances of Redis.

If you are interested how Sync between clients works,
take a look at
[this post.](https://redis.io/docs/manual/patterns/distributed-locks/)


## Installation

```bash
$ go get -u github.com/coopnorge/go-datadog-lib
```

## Quick Start

Add this import line to the file you're working in:

```Go
import "github.com/coopnorge/go-datadog-lib"
```

We recommend create custom constructor.

```go
func NewRedisStorageFacade(cfg *config.MyAppConfig) *database.RedisFacade {
	// ...
}
```

Prepare configuration for Redis connection

```go
redisCfg := database.Config{
    Address:           "RedisAddress:RedisPort",
    Password:          "RedisPassword",
    Database:          "RedisDatabase",
    DialTimeout:       "RedisDialTimeout",
    ReadTimeout:       "RedisReadTimeout",
    WriteTimeout:      "RedisWriteTimeout",
    MaxRetries:        "RedisMaxRetries",
    PoolSize:          "RedisPoolSize",
    PoolTimeout:       "RedisPoolTimeout",
    EncryptionEnabled: "RedisEncryptorEnabled",
}

encrCfg := database.EncryptionConfig{
    RedisKeyURI: "RedisEncryptorKeyURI",
    Aad:         []byte("RedisEncryptorAad"),
}
```

Create Redis wrapper instances

```go
	encryptor, encryptorErr := database.NewEncryptionClient(encryptionConfig)
	if encryptorErr != nil {
		panic(fmt.Errorf("unable to create new encryption client for Redis Client, error: %w", encryptorErr))
	}

	redFacade, redFacadeErr := database.NewRedisFacade(redisConfig, encryptor)
	if redFacadeErr != nil {
		panic(fmt.Errorf("unable to create Redis Client, error: %w", redFacadeErr))
	}
```

## Use case

Then you can add client to your repository to work with Redis

```go
func NewUserRepository(s database.KeyValueStorage) *UserRepository {
	return &UserRepository{db: s}
}

func (r *UserRepository) Create(ctx context.Context, u model.User) error {
	u.ID = uuid.NewUUID()
    u.CreatedAt = time.Now()
    
    j, jErr := json.Marshal(cart)
    if jErr != nil {
        return jErr
    }
    
    return r.db.Save(ctx, u.ID, string(j), expirationTime)
}
```

## Mocks

To generate or update mocks use tools
[Eitri](https://github.com/Clink-n-Clank/Eitri)
or use directly
[Mockhandler](github.com/sanposhiho/gomockhandle)
