package redis

import (
	"context"
	"github.com/go-redis/redis"
)

type redisInst struct {
	inst *redis.Client
	ctx  context.Context
}

func newRedisInst() (*redisInst, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // host:port of the redis server
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &redisInst{client, context.TODO()}, nil
}

func (r *redisInst) setStrVal(key string, value string) error {
	_, err := r.inst.Set(key, value, 0).Result()
	return err
}

func (r *redisInst) getStrVal(key string) string {
	return r.inst.Get(key).Val()
}

func (r *redisInst) setDictVal(key, k, v string) error {
	interf := make(map[string]interface{}, 1)
	interf[k] = v
	_, err := r.inst.HMSet(key, interf).Result()
	return err
}

func (r *redisInst) getDictVal(key string) []interface{} {
	res := r.inst.HMGet(key).Val()
	return res

}
