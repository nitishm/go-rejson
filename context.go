package rejson

import (
	"context"
	"github.com/nitishm/go-rejson/v4/clients"
	"github.com/nitishm/go-rejson/v4/rjs"
)

// SetContext helps redis-clients, provide use of command level context
// in the ReJSON commands.
// Currently, only go-redis@v8 supports command level context, therefore
// a separate method is added to support it, maintaining the support for
// other clients and for backward compatibility. (nitishm/go-rejson#46)
func (r *Handler) SetContext(ctx context.Context) *Handler {
	if r == nil {
		return r // nil
	}

	if r.clientName == rjs.ClientGoRedis {
		if old, ok := r.implementation.(*clients.GoRedis); ok {
			return &Handler{
				clientName:     r.clientName,
				implementation: clients.NewGoRedisClient(ctx, old.Conn),
			}
		}
	}

	// for other clients, context is of no use, hence return same
	return r
}
