package database

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

var stubConn *miniredis.Miniredis

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
	cfg := Config{Address: stubConn.Addr()}

	facadeClient0, facadeClient0Err := NewRedisFacade(cfg)
	facadeClient1, facadeClient1Err := NewRedisFacade(cfg)
	facadeClient2, facadeClient2Err := NewRedisFacade(cfg)
	if facadeClient0Err != nil || facadeClient1Err != nil || facadeClient2Err != nil {
		t.Fatal("unable to create one of redis facades")
	}

	// Act
	const testStoredKey = "race-update"
	const expectedStoredValue = "the-one-bar"

	// Create testStoredKey
	assert.Nil(t, facadeClient0.Save(context.Background(), testStoredKey, "bar", time.Minute))
	time.Sleep(time.Millisecond)

	// Try update testStoredKey - value
	go func() {
		assert.Nil(t, facadeClient1.Update(context.Background(), testStoredKey, "first-update", time.Minute))
		time.Sleep(time.Millisecond)
		assert.Nil(t, facadeClient1.Update(context.Background(), testStoredKey, expectedStoredValue, time.Minute))
	}()
	go func() {
		assert.Nil(t, facadeClient2.Update(context.Background(), testStoredKey, "second-update", time.Minute))
	}()

	time.Sleep(time.Millisecond * 500)

	assert.True(t, isRecordSame(facadeClient0, testStoredKey, expectedStoredValue, t), "expected to be found vale")
}

func TestRedisFacadeSaveWithLockInSameTime(t *testing.T) {
	cfg := Config{Address: stubConn.Addr()}

	validatorClient, validatorClientErr := NewRedisFacade(cfg)
	writeClient1, writeClient1Err := NewRedisFacade(cfg)
	writeClient2, writeClient2Err := NewRedisFacade(cfg)
	if validatorClientErr != nil || writeClient1Err != nil || writeClient2Err != nil {
		t.Fatal("unable to create one of redis facades")
	}

	// Act
	const testStoredKey = "race-write"

	// Create testStoredKey
	assert.Nil(t, validatorClient.Save(context.Background(), testStoredKey, "init", time.Minute))

	// Try update testStoredKey - value
	go func() {
		assert.Nil(t, writeClient1.Update(context.Background(), testStoredKey, "first-update", time.Minute))

		assert.True(t, isRecordSame(writeClient1, testStoredKey, "second-update", t), "expected to be found vale")
	}()
	go func() {
		assert.True(
			t,
			isRecordSame(writeClient1, testStoredKey, "init", t) || isRecordSame(writeClient1, testStoredKey, "first-update", t),
			"record can be already update, expected to find 'init' or 'first-update' value",
		)

		assert.Nil(t, writeClient2.Update(context.Background(), testStoredKey, "second-update", time.Minute))
	}()

	assert.True(t, isRecordSame(writeClient1, testStoredKey, "init", t), "expected to be found vale")

	time.Sleep(time.Second)

	assert.True(t, isRecordSame(writeClient1, testStoredKey, "second-update", t), "expected to be found vale")
}

func isRecordSame(cli *RedisFacade, testStoredKey, expectedRes string, t *testing.T) bool {
	res, resErr := cli.Find(context.Background(), testStoredKey)
	assert.Nil(t, resErr)

	return expectedRes == string(res)
}
