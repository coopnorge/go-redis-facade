package database

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	mockEncryptor0 := getPreparedMocks(t)
	mockEncryptor1 := getPreparedMocks(t)
	mockEncryptor2 := getPreparedMocks(t)

	cfg := Config{Address: stubConn.Addr(), EncryptionEnabled: true}

	facadeClient0, facadeClient0Err := NewRedisFacade(cfg, mockEncryptor0)
	facadeClient1, facadeClient1Err := NewRedisFacade(cfg, mockEncryptor1)
	facadeClient2, facadeClient2Err := NewRedisFacade(cfg, mockEncryptor2)
	if facadeClient0Err != nil || facadeClient1Err != nil || facadeClient2Err != nil {
		t.Fatal("unable to create one of redis facades")
	}

	// Act
	const testStoredKey = "race-update"
	const expectedStoredValue = "the-one-bar"

	// Create testStoredKey
	mockEncryptor0.EXPECT().Encrypt([]byte("bar")).Return([]byte("encrypted-bar"), nil)
	assert.Nil(t, facadeClient0.Save(context.Background(), testStoredKey, "bar", time.Minute))
	time.Sleep(time.Millisecond)

	// Try update testStoredKey - value
	go func() {
		mockEncryptor1.EXPECT().Encrypt([]byte("first-update")).Return([]byte("encrypted-first-update"), nil)
		assert.Nil(t, facadeClient1.Update(context.Background(), testStoredKey, "first-update", time.Minute))
		time.Sleep(time.Millisecond)
		mockEncryptor1.EXPECT().Encrypt([]byte(expectedStoredValue)).Return([]byte("encrypted-the-one-bar"), nil)
		assert.Nil(t, facadeClient1.Update(context.Background(), testStoredKey, expectedStoredValue, time.Minute))
	}()
	go func() {
		mockEncryptor2.EXPECT().Encrypt([]byte("second-update")).Return([]byte("encrypted-second-update"), nil)
		assert.Nil(t, facadeClient2.Update(context.Background(), testStoredKey, "second-update", time.Minute))
	}()

	time.Sleep(time.Millisecond * 500)

	mockEncryptor0.EXPECT().Decrypt([]byte("encrypted-the-one-bar")).Return([]byte(expectedStoredValue), nil)
	assert.True(t, isRecordSame(facadeClient0, testStoredKey, expectedStoredValue, t), "expected to be found vale")
}

func TestRedisFacadeSaveWithLockInSameTime(t *testing.T) {
	mockEncryptor0 := &stubEncryption{}
	mockEncryptor1 := &stubEncryption{}
	mockEncryptor2 := &stubEncryption{}
	cfg := Config{Address: stubConn.Addr(), EncryptionEnabled: true}

	validatorClient, validatorClientErr := NewRedisFacade(cfg, mockEncryptor0)
	writeClient1, writeClient1Err := NewRedisFacade(cfg, mockEncryptor1)
	writeClient2, writeClient2Err := NewRedisFacade(cfg, mockEncryptor2)
	if validatorClientErr != nil || writeClient1Err != nil || writeClient2Err != nil {
		t.Fatal("unable to create one of redis facades")
	}

	// Act
	const testStoredKey = "race-write"

	// Create testStoredKey
	assert.Nil(t, validatorClient.Save(context.Background(), testStoredKey, "init", time.Minute))
	assert.True(t, isRecordSame(validatorClient, testStoredKey, "init", t), "unexpected stored value")

	// Try update testStoredKey - value
	go func() {
		assert.Nil(t, writeClient1.Update(context.Background(), testStoredKey, "first-update", time.Minute))
		assert.True(t, isRecordSame(validatorClient, testStoredKey, "first-update", t), "unexpected stored value")
	}()
	go func() {
		assert.True(
			t,
			isRecordSame(writeClient2, testStoredKey, "init", t) || isRecordSame(writeClient2, testStoredKey, "first-update", t),
			"unexpected stored value",
		)
		assert.Nil(t, validatorClient.Update(context.Background(), testStoredKey, "second-update", time.Minute))
	}()

	time.Sleep(time.Second)
	assert.True(t, isRecordSame(validatorClient, testStoredKey, "second-update", t), "unexpected stored value")
}

func isRecordSame(cli *RedisFacade, testStoredKey, expectedRes string, t *testing.T) bool {
	res, resErr := cli.Find(context.Background(), testStoredKey)
	assert.Nil(t, resErr)

	t.Log(fmt.Sprintf("Validating stored value in redis by key (%s) => Expected: %s - Stored: %s", testStoredKey, expectedRes, string(res)))

	return expectedRes == string(res)
}

func TestEncryptionDisabled(t *testing.T) {
	mockEncryptor := getPreparedMocks(t)

	cfg := Config{Address: stubConn.Addr(), EncryptionEnabled: false}

	validatorClient, validatorClientErr := NewRedisFacade(cfg, mockEncryptor)
	assert.Nil(t, validatorClientErr)

	err := validatorClient.Save(context.TODO(), "something", "val", time.Minute)
	assert.Nil(t, err)

	val, err := validatorClient.Find(context.TODO(), "something")
	assert.Nil(t, err)
	assert.Equal(t, "val", val)
}
