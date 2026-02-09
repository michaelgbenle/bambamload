package redisService

import (
	"bambamload/logger"
	"bambamload/types"
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Redis struct {
	Client *redis.Client //nolint:typecheck
}

type RedisService interface {
	Ping() error
	GetRedisClient() *redis.Client //nolint:typecheck
	RedisLock() *redislock.Client  //nolint:typecheck
	RunWithLock(key string, ttl time.Duration, job func())
	SetSession(sessionInfo types.RedisSessionInfo) error
	GetSession(token string) (types.RedisSessionInfo, error)
	DeleteSession(token string) error
	SetValue(key string, value interface{}, expiration int) error
	GetValue(key string, target interface{}) error
	PushToQueue(queue string, msg any) error
}

func NewRedisService() Redis {

	logger.Logger.Info("Connecting to redis server ...")
	start := time.Now()

	//single instance
	client := redis.NewClient(&redis.Options{ //nolint:typecheck
		Addr: os.Getenv("REDIS_ADDRESS"),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Logger.Fatalf("failed to connect to redis: %v", err)
	}

	logger.Logger.Infof("Connected to redis server after %v seconds", time.Since(start).Seconds())

	return Redis{
		Client: client,
	}
}

// Ping ...
func (r Redis) Ping() error {
	if err := r.Client.Ping(ctx).Err(); err != nil {
		logger.Logger.Fatalf("[Ping]failed to connect to redis: %v", err)
		return err
	}
	return nil
}

func (r Redis) GetRedisClient() *redis.Client { //nolint:typecheck
	return r.Client
}
func (r Redis) RedisLock() *redislock.Client { //nolint:typecheck
	return redislock.New(r.Client) //nolint:typecheck
}

// SetSession saves a user's session - refresh token and information
func (r Redis) SetSession(sessionInfo types.RedisSessionInfo) error {
	var b bytes.Buffer

	if encodeErr := gob.NewEncoder(&b).Encode(&sessionInfo); encodeErr != nil {
		logger.Logger.Errorf("[SetSession]failed to encode session: %v", encodeErr)
		return encodeErr
	}

	setErr := r.Client.Set(context.Background(), sessionInfo.Token, b.Bytes(), sessionInfo.Expiry.Sub(time.Now())).Err()
	if setErr != nil {
		logger.Logger.Errorf("[SetSession]failed to set session: %v", setErr)
		return setErr
	}

	return nil
}

// GetSession returns a user's session information specified by the  token
func (r Redis) GetSession(token string) (types.RedisSessionInfo, error) {

	sessBytes, getErr := r.Client.Get(context.Background(), token).Bytes()
	if getErr != nil {
		logger.Logger.Errorf("[GetSession]failed to get session: %v", getErr)
		return types.RedisSessionInfo{}, getErr
	}

	sessByteReader := bytes.NewReader(sessBytes)
	var sessionInfo types.RedisSessionInfo

	if decodeErr := gob.NewDecoder(sessByteReader).Decode(&sessionInfo); decodeErr != nil {
		return types.RedisSessionInfo{}, decodeErr
	}

	return sessionInfo, nil
}

// DeleteSession deletes a user's session from redis
func (r Redis) DeleteSession(token string) error {
	delErr := r.Client.Del(context.Background(), token).Err()
	if delErr != nil {
		logger.Logger.Errorf("[DeleteSession]failed to delete session: %v", delErr)
		return delErr
	}

	return nil
}

func (r Redis) SetValue(key string, value interface{}, expiration int) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		logger.Logger.Errorf("[SetValue]failed to encode value: %v", err)
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	err = r.Client.Set(ctx, key, jsonValue, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		logger.Logger.Errorf("[SetValue]failed to set value: %v", err)
		return fmt.Errorf("failed to set value in Redis: %v", err)
	}

	return nil
}

func (r Redis) GetValue(key string, target interface{}) error {
	jsonValue, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		logger.Logger.Errorf("[GetValue]failed to find key: %v", key)
		return err
	}

	if err = json.Unmarshal([]byte(jsonValue), target); err != nil {
		logger.Logger.Errorf("[GetValue]failed to decode value: %v", err)
		return fmt.Errorf("failed to unmarshal value: %v", err)
	}

	return nil
}

func (r Redis) RunWithLock(key string, ttl time.Duration, job func()) {
	lock, err := r.RedisLock().Obtain(ctx, key, ttl, nil)
	if errors.Is(err, redislock.ErrNotObtained) { //nolint:typecheck
		logger.Logger.Infof("Another instance of %s is running", key)
		return
	} else if err != nil {
		logger.Logger.Errorf("Error acquiring %s: %v\n", key, err)
		return
	}

	defer func() {
		if err = lock.Release(ctx); err != nil {
			logger.Logger.Errorf("Error releasing %s: %v\n", key, err)
		}
	}()

	job()
}

func (r Redis) PushToQueue(queue string, msg any) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Logger.Errorf("Unable to marshal message: %s", err)
		return err
	}
	return r.Client.RPush(ctx, queue, string(msgBytes)).Err()
}
