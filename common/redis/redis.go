package redis

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// Redis return
type Redis struct {
	p        *pool.Pool
	maxConns int
	dbNum    int
	conn     string
	password string
}

// NewRedis return Redis
func NewRedis() *Redis {
	return &Redis{}
}

// Connect return Redis
func (r *Redis) Connect(conn, password string, maxConns, dbNum int) (err error) {
	r.conn = conn
	r.dbNum = dbNum
	r.password = password
	r.maxConns = maxConns
	dialFunc := func(network, addr string) (c *redis.Client, err error) {
		c, err = redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}

		if r.password != "" {
			if err := c.Cmd("AUTH", r.password).Err; err != nil {
				c.Close()
				return nil, err
			}
		}

		err = c.Cmd("SELECT", r.dbNum).Err
		if err != nil {
			c.Close()
			return nil, err
		}

		return
	}

	p, err := pool.NewCustom("tcp", r.conn, r.maxConns, dialFunc)
	if err != nil {
		return err
	}
	r.p = p
	return nil
}

func (r *Redis) cmd(cmd string, args ...interface{}) *redis.Resp {
	c, err := r.p.Get()
	if err != nil {
		return nil
	}
	defer r.p.Put(c)
	return c.Cmd(cmd, args...)

}

// Get return
func (r *Redis) Get(key string) *redis.Resp {
	return r.cmd("Get", key)
}

// Set return
func (r *Redis) Set(key string, data interface{}) *redis.Resp {
	return r.cmd("Set", key, data)
}
