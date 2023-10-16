package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis .redis cache
type Redis struct {
	ctx  context.Context
	conn redis.UniversalClient
}

// RedisOpts redis 连接属性
type RedisOpts struct {
	Host        string `yml:"host" json:"host"`
	Password    string `yml:"password" json:"password"`
	Database    int    `yml:"database" json:"database"`
	MaxIdle     int    `yml:"max_idle" json:"max_idle"`
	MaxActive   int    `yml:"max_active" json:"max_active"`
	IdleTimeout int    `yml:"idle_timeout" json:"idle_timeout"` // second
}

// NewRedis 实例化
func NewRedis(ctx context.Context, opts *RedisOpts) *Redis {
	conn := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:           []string{opts.Host},
		DB:              opts.Database,
		Password:        opts.Password,
		ConnMaxIdleTime: time.Second * time.Duration(opts.IdleTimeout),
		MinIdleConns:    opts.MaxIdle,
	})
	return &Redis{ctx: ctx, conn: conn}
}

// SetConn 设置 conn
func (r *Redis) SetConn(conn redis.UniversalClient) {
	r.conn = conn
}

// SetRedisCtx 设置 redis ctx 参数
func (r *Redis) SetRedisCtx(ctx context.Context) {
	r.ctx = ctx
}

// Get 获取一个值
func (r *Redis) Get(key string) interface{} {
	return r.GetContext(r.ctx, key)
}

// GetContext 获取一个值
func (r *Redis) GetContext(ctx context.Context, key string) interface{} {
	result := r.conn.Get(ctx, key)
	if result.Err() != nil {
		return ""
	}
	return result.Val()
}

// Set 设置一个值
func (r *Redis) Set(key string, val interface{}, timeout time.Duration) error {
	return r.SetContext(r.ctx, key, val, timeout)
}

// SetContext 设置一个值
func (r *Redis) SetContext(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return r.conn.SetEx(ctx, key, val, timeout).Err()
}

// IsExist 判断 key 是否存在
func (r *Redis) IsExist(key string) bool {
	return r.IsExistContext(r.ctx, key)
}

// IsExistContext 判断 key 是否存在
func (r *Redis) IsExistContext(ctx context.Context, key string) bool {
	result, _ := r.conn.Exists(ctx, key).Result()

	return result > 0
}

// Delete 删除
func (r *Redis) Delete(key string) error {
	return r.DeleteContext(r.ctx, key)
}

// DeleteContext 删除
func (r *Redis) DeleteContext(ctx context.Context, key string) error {
	return r.conn.Del(ctx, key).Err()
}
