package rejson

import (
	"context"
	"github.com/nitishm/go-rejson/v4/clients"
	"github.com/nitishm/go-rejson/v4/rjs"
)

// RedisClient provides interface for Client handling in the ReJSON Handler
type RedisClient interface {
	SetClientInactive()
	SetRedigoClient(conn clients.RedigoClientConn)
	SetGoRedisClient(conn clients.GoRedisClientConn)
}

// SetClientInactive resets the handler and unset any client, set to the handler
func (r *Handler) SetClientInactive() {
	_t := &Handler{clientName: rjs.ClientInactive}
	r.clientName = _t.clientName
	r.implementation = _t.implementation
}

// SetRedigoClient sets Redigo (https://github.com/gomodule/redigo/redis) client
// to the handler
func (r *Handler) SetRedigoClient(conn clients.RedigoClientConn) {
	r.clientName = "redigo"
	r.implementation = &clients.Redigo{Conn: conn}
}

// SetGoRedisClient sets Go-Redis (https://github.com/go-redis/redis) client to
// the handler. It is left for backward compatibility.
func (r *Handler) SetGoRedisClient(conn clients.GoRedisClientConn) {
	r.SetGoRedisClientWithContext(context.Background(), conn)
}

// SetGoRedisClientWithContext sets Go-Redis (https://github.com/go-redis/redis) client to
// the handler with a global context for the connection
func (r *Handler) SetGoRedisClientWithContext(ctx context.Context, conn clients.GoRedisClientConn) {
	r.clientName = "goredis"
	r.implementation = clients.NewGoRedisClient(ctx, conn)
}
