package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// ErrorType is the type of error that is returned
type ErrorType int

const (
	//ErrorTypeInvalid invalid error type
	ErrorTypeInvalid ErrorType = iota
	//KeyMissing error type
	KeyMissing
	//KeyExist error type
	KeyExist
	//CommandError error type
	CommandError
	//OperationError error type
	OperationError
	//BytesConvertionError error type
	BytesConvertionError
)

// StorageFacadeError is a custom error type used for error handling
type StorageFacadeError struct {
	Type    ErrorType
	Details string
}

func (e *StorageFacadeError) Error() string {
	return e.Details
}

// KeyValueStorage defines interface
type KeyValueStorage interface {
	// Save value in storage by key with expiration
	Save(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Update value in storage by key and update expiration
	Update(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Find in storage by key
	Find(ctx context.Context, key string) ([]byte, error)
	// FindKeys in storage by given pattern
	FindKeys(ctx context.Context, pattern string) ([]string, error)
	// Close connection
	Close() error
}

// RedisFacade storage
type RedisFacade struct {
	c redis.Client
}

// Config has all data to connect to redis
type Config struct {
	Address      string
	Password     string
	Database     int
	DialTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	MaxRetries   int
	PoolSize     int
	PoolTimeout  time.Duration
}

// NewRedisFacade makes new connection to key value storage
func NewRedisFacade(config Config) (*RedisFacade, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Address,
		Password:     config.Password,
		DB:           config.Database,
		DialTimeout:  config.DialTimeout,
		WriteTimeout: config.WriteTimeout,
		ReadTimeout:  config.ReadTimeout,
		MaxRetries:   config.MaxRetries,
		PoolSize:     config.PoolSize,
		PoolTimeout:  config.PoolTimeout,
	})

	res := client.Ping(context.Background())
	if err := handleCommandError("Ping", res); err != nil {
		return nil, err
	}

	return &RedisFacade{c: *client}, nil
}

// Save value in storage by key with expiration
func (rf *RedisFacade) Save(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	result := rf.c.Keys(ctx, key)
	keys, err := result.Result()
	err = rf.handleValueError(err, "Keys", key)
	if err != nil {
		return &StorageFacadeError{Type: OperationError, Details: err.Error()}
	}

	if len(keys) != 0 {
		return &StorageFacadeError{Type: KeyExist, Details: fmt.Sprintf("Key (%s) already exists", key)}
	}

	set := rf.c.Set(ctx, key, value, expiration)

	return handleCommandError("Set", set)
}

// Update value in storage by key and update expiration
func (rf *RedisFacade) Update(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	set := rf.c.Set(ctx, key, value, expiration)

	return handleCommandError("Set", set)
}

// Find in storage by key
func (rf *RedisFacade) Find(ctx context.Context, key string) (b []byte, err error) {
	rawResult := rf.c.Get(ctx, key)
	if rawResult == nil {
		return nil, &StorageFacadeError{
			Type:    CommandError,
			Details: fmt.Sprintf("unexpected redis error, operation (%s) command response is nil", "Get"),
		}
	}

	redisVal, errRedisVal := rawResult.Result()
	err = rf.handleValueError(errRedisVal, "Get", key)
	if err != nil {
		return []byte{}, err
	}
	if redisVal == "" {
		return []byte{}, nil
	}

	if b, err = rawResult.Bytes(); err != nil {
		return nil, &StorageFacadeError{
			Type:    BytesConvertionError,
			Details: err.Error(),
		}
	}

	return b, nil
}

// FindKeys in storage by given pattern
func (rf *RedisFacade) FindKeys(ctx context.Context, pattern string) (redisKeysVal []string, err error) {
	result := rf.c.Keys(ctx, pattern)
	if result == nil {
		return nil, &StorageFacadeError{
			Type:    CommandError,
			Details: fmt.Sprintf("unexpected redis error, operation (%s) command response is nil", "Keys"),
		}
	}

	redisKeysVal, err = result.Result()
	err = rf.handleValueError(err, "Keys", pattern)
	if err != nil {
		return []string{}, err
	}

	return redisKeysVal, nil
}

// Close connection
func (rf *RedisFacade) Close() error {
	return rf.c.Close()
}

func handleCommandError(operation string, set *redis.StatusCmd) error {
	if set == nil {
		return &StorageFacadeError{
			Type:    CommandError,
			Details: fmt.Sprintf("unexpected redis error, operation (%s) command response is nil", operation),
		}
	}

	if set.Err() != nil {
		return &StorageFacadeError{
			Type:    CommandError,
			Details: set.Err().Error(),
		}
	}

	return nil
}

func (rf *RedisFacade) handleValueError(err error, op, key string) error {
	if err == redis.Nil {
		return &StorageFacadeError{
			Type:    KeyMissing,
			Details: fmt.Sprintf("redis error, operation (%s) - key (%s) does not exist", op, key),
		}
	} else if err != nil {
		return errors.Wrap(err, fmt.Sprintf("redis error, operation (%s) failed", op))
	}

	return err
}
