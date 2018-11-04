package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/data/model"
	"github.com/dstpierre/gosaas/engine"
)

// API is the starting point of our API.
// Responsible for routing the request to the correct handler
type API struct {
	DB            *data.DB
	Logger        func(http.Handler) http.Handler
	Authenticator func(http.Handler) http.Handler
	Throttler     func(http.Handler) http.Handler
	RateLimiter   func(http.Handler) http.Handler
	User          *engine.Route
}

// NewAPI returns a production API with all middlewares
func NewAPI() *API {
	return &API{
		Logger:        engine.Logger,
		Authenticator: engine.Authenticator,
		Throttler:     engine.Throttler,
		RateLimiter:   engine.RateLimiter,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/public/") {
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, engine.ContextOriginalPath, r.URL.Path)

	ctx = context.WithValue(ctx, engine.ContextDatabase, a.DB)

	var next *engine.Route
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	if head == "" || head == "buy" {
		next = newBuy()
	} else {
		next = newError(fmt.Errorf("path not found"), http.StatusNotFound)
	}

	ctx = context.WithValue(ctx, engine.ContextMinimumRole, next.MinimumRole)

	// make sure we are authenticating all calls
	next.Handler = a.Authenticator(next.Handler)

	if next.Logger {
		next.Handler = a.Logger(next.Handler)
	}

	next.Handler = a.RateLimiter(next.Handler)
	next.Handler = a.Throttler(next.Handler)

	next.Handler.ServeHTTP(w, r.WithContext(ctx))
}

func newError(err error, statusCode int) *engine.Route {
	return &engine.Route{
		Logger:      true,
		MinimumRole: model.RoleUser,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			engine.Respond(w, r, statusCode, err)
		}),
	}
}
