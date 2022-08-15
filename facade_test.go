package database

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/coopnorge/scan-and-pay-redis-facade/internal/generated/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var stubConn *miniredis.Miniredis

type stubEncryption struct{}

func (m *stubEncryption) Encrypt(value []byte) ([]byte, error) {
	return []byte(fmt.Sprintf("encrypted-%s", string(value))), nil
}

func (m *stubEncryption) Decrypt(ciphertext []byte) ([]byte, error) {
	asString := string(ciphertext)
	withoutPrefix := strings.TrimPrefix(asString, "encrypted-")
	return []byte(withoutPrefix), nil
}

func getPreparedMocks(t *testing.T) *mock_database.MockEncryption {
	ctrl := gomock.NewController(t)
	mockEncryptor := mock_database.NewMockEncryption(ctrl)
	ctrl.Finish()

	return mockEncryptor
}

func TestMain(m *testing.M) {
	tmpRedisServer, tmpRedisServerErr := miniredis.Run()
	if tmpRedisServerErr != nil {
		panic(fmt.Sprintf("stub connection error: %s", tmpRedisServerErr))
	}

	stubConn = tmpRedisServer

	code := m.Run()
	os.Exit(code)
}

func TestRedisFacadeSaveWithLock(t *testing.T) {
	redisValueExpiryTime := 30 * time.Second
	redisCfg := Config{Address: stubConn.Addr(), EncryptionEnabled: true, DialTimeout: 2 * time.Minute}

	mockBaseEncryption := getPreparedMocks(t)
	facadeBaseClient, facadeBaseClientErr := NewRedisFacade(redisCfg, mockBaseEncryption)
	if facadeBaseClientErr != nil {
		t.Fatal("unable to create one of redis facades")
	}

	// Act
	const testStoredKey = "race-update"
	const maxRedisInst = 2

	// Create testStoredKey
	mockBaseEncryption.EXPECT().Encrypt([]byte("bar")).Return([]byte("encrypted-bar"), nil)
	assert.Nil(t, facadeBaseClient.Save(context.Background(), testStoredKey, "bar", redisValueExpiryTime))
	time.Sleep(time.Millisecond)

	var wg sync.WaitGroup
	for i := 1; i <= maxRedisInst; i++ {
		wg.Add(1)
		go func(i int, t *testing.T) {
			time.Sleep(time.Duration(i) * time.Millisecond)

			defer func() {
				wg.Done()
			}()

			mockEncryptor := getPreparedMocks(t)
			facadeClient, facadeClientErr := NewRedisFacade(redisCfg, mockEncryptor)
			if facadeClientErr != nil {
				t.Errorf("unable to create one of redis facades")
				return
			}

			updateVal := fmt.Sprintf("update-%d", i)
			encryptUpdateVal := fmt.Sprintf("encrypt-update-%d", i)
			mockEncryptor.EXPECT().Encrypt([]byte(updateVal)).Return([]byte(encryptUpdateVal), nil)
			assert.Nil(t, facadeClient.Update(context.Background(), testStoredKey, updateVal, redisValueExpiryTime))
		}(i, t)
	}

	wg.Wait()

	mockBaseEncryption.EXPECT().Decrypt([]byte(fmt.Sprintf("encrypt-update-%d", maxRedisInst))).Return([]byte(fmt.Sprintf("update-%d", maxRedisInst)), nil)
	assert.True(t, isRecordSame(facadeBaseClient, testStoredKey, fmt.Sprintf("update-%d", maxRedisInst), t), "expected to be found vale")
}

func isRecordSame(cli *RedisFacade, testStoredKey, expectedRes string, t *testing.T) bool {
	res, resErr := cli.Find(context.Background(), testStoredKey)
	assert.Nil(t, resErr)

	return assert.Equal(t, expectedRes, res)
}

func TestEncryptionDisabled(t *testing.T) {
	mockEncryptor := getPreparedMocks(t)

	cfg := Config{Address: stubConn.Addr(), EncryptionEnabled: false, DialTimeout: 2 * time.Minute}

	validatorClient, validatorClientErr := NewRedisFacade(cfg, mockEncryptor)
	assert.Nil(t, validatorClientErr)

	err := validatorClient.Save(context.TODO(), "something", "val", time.Minute)
	assert.Nil(t, err)

	val, err := validatorClient.Find(context.TODO(), "something")
	assert.Nil(t, err)
	assert.Equal(t, "val", val)
}
