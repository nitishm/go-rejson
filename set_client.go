package rejson

import (
	"github.com/Shivam010/go-rejson/clients"
	"github.com/Shivam010/go-rejson/rjs"
	goredis "github.com/go-redis/redis"
	redigo "github.com/gomodule/redigo/redis"
)

type RedisClient interface {
	SetClientInactive()
	SetRedigoClient(redigo.Conn)
	SetGoRedisClient(conn *goredis.Client)
}

func (r *Handler) SetClientInactive() {
	_t := &Handler{clientName: rjs.ClientInactive}
	r.clientName = _t.clientName
	r.implementation = _t.implementation
}

func (r *Handler) SetRedigoClient(conn redigo.Conn) {
	r.clientName = "redigo"
	r.implementation = &clients.Redigo{Conn: conn}
}

func (r *Handler) SetGoRedisClient(conn *goredis.Client) {
	r.clientName = "goredis"
	r.implementation = &clients.GoRedis{Conn: conn}
}
