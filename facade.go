package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// ErrorType is the type of error that is returned
type ErrorType int

const (
	// ErrorTypeInvalid invalid error type
	ErrorTypeInvalid ErrorType = iota
	// KeyMissing error type
	KeyMissing
	// KeyExist error type
	KeyExist
	// SyncError will occur with resource locking
	SyncError
	// CommandError error type
	CommandError
	// OperationError error type
	OperationError
	// BytesConvertionError error type
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
	// Delete value in storage by key
	Delete(ctx context.Context, key string) (count int64, err error)
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

	rSync *redsync.Redsync

	sync.Mutex
	redMutexKey string
	redMutex    *redsync.Mutex
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

	return &RedisFacade{
		c:     *client,
		rSync: redsync.New(goredis.NewPool(client)),
	}, nil
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

	set, setErr := rf.doInSync(ctx, key, value, expiration, rf.c.Set)
	if setErr != nil {
		return setErr
	}

	return handleCommandError("Set", set)
}

// Update value in storage by key and update expiration
func (rf *RedisFacade) Update(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	set, setErr := rf.doInSync(ctx, key, value, expiration, rf.c.Set)
	if setErr != nil {
		return setErr
	}

	return handleCommandError("Set", set)
}

// Delete value in storage by key
func (rf *RedisFacade) Delete(ctx context.Context, key string) (count int64, err error) {
	if lockErr := rf.lockAcquire(ctx, key); lockErr != nil {
		return count, lockErr
	}

	rawResult := rf.c.Del(ctx, key)
	if rawResult == nil {
		return 0, &StorageFacadeError{
			Type:    CommandError,
			Details: fmt.Sprintf("unexpected redis error, operation (%s) command response is nil", "Del"),
		}
	}
	redisVal, errRedisVal := rawResult.Result()
	err = rf.handleValueError(errRedisVal, "Del", key)
	if lockErr := rf.lockRelease(ctx); lockErr != nil {
		return count, lockErr
	}
	return redisVal, err
}

// Find in storage by key
func (rf *RedisFacade) Find(ctx context.Context, key string) (b []byte, err error) {
	if lockErr := rf.lockAcquire(ctx, key); lockErr != nil {
		return nil, lockErr
	}

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

	if lockErr := rf.lockRelease(ctx); lockErr != nil {
		return nil, lockErr
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

// doInSync function execution with resource locking in redis
func (rf *RedisFacade) doInSync(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
	cmd func(context.Context, string, interface{}, time.Duration) *redis.StatusCmd,
) (*redis.StatusCmd, error) {
	if lockErr := rf.lockAcquire(ctx, key); lockErr != nil {
		return nil, lockErr
	}

	result := cmd(ctx, key, value, expiration)

	if lockErr := rf.lockRelease(ctx); lockErr != nil {
		return result, lockErr
	}

	return result, nil
}

// lockAcquire will lock resource
func (rf *RedisFacade) lockAcquire(ctx context.Context, prefix string) error {
	rf.Lock()
	defer rf.Unlock()

	newUUID, errNewUUID := uuid.NewUUID()
	if errNewUUID != nil {
		return &StorageFacadeError{
			Type:    SyncError,
			Details: fmt.Sprintf("unable to generate unique lock uuid, err: %v", errNewUUID),
		}
	}

	rf.redMutexKey = fmt.Sprintf("%s_", newUUID.String())
	rf.redMutex = rf.rSync.NewMutex(rf.redMutexKey)
	if lockErr := rf.redMutex.LockContext(ctx); lockErr != nil {
		return &StorageFacadeError{
			Type:    SyncError,
			Details: fmt.Sprintf("unable to acquire lock for redis resource (key:%s, lock:%s), err: %v", prefix, rf.redMutexKey, lockErr),
		}
	}

	return nil
}

// lockRelease will remove lock from resource
func (rf *RedisFacade) lockRelease(ctx context.Context) error {
	rf.Lock()
	defer rf.Unlock()

	if rf.redMutex == nil {
		return nil
	}

	if _, unlockErr := rf.redMutex.UnlockContext(ctx); unlockErr != nil {
		return &StorageFacadeError{
			Type:    SyncError,
			Details: fmt.Sprintf("unable to unlock redis resource (lock:%s), err: %s", rf.redMutexKey, unlockErr.Error()),
		}
	}

	rf.redMutexKey = ""
	rf.redMutex = nil

	return nil
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
