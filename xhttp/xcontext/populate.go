package xcontext

import (
	"context"
	"net/http"
	"time"
)

// Populate accepts any number of go-kit request functions and returns an Alice-style constructor that
// uses the request functions to build a context.  The resulting context is then assocated with the request
// prior to the next http.Handler being invoked.
//
// This function mimics the behavior of go-kit's transport/http package without requiring and endpoint with
// encoding and decoding.
func Populate(timeout time.Duration, rf ...func(context.Context, *http.Request) context.Context) func(http.Handler) http.Handler {
	if timeout > 0 || len(rf) > 0 {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				ctx := request.Context()
				for _, f := range rf {
					ctx = f(ctx, request)
				}

				if timeout > 0 {
					var cancel func()
					ctx, cancel = context.WithTimeout(ctx, timeout)
					defer cancel()
				}

				next.ServeHTTP(response, request.WithContext(ctx))
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return next
	}
}
