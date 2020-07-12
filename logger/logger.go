package logger

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const key = "logger"
const traceparent = "traceparent"

func MiddleWare(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newLogger := logger.With(zap.String(traceparent, r.Header.Get(traceparent)))
			nr := r.WithContext(context.WithValue(r.Context(), key, newLogger))
			next.ServeHTTP(w, nr)
		})
	}
}

func L(ctx context.Context) *zap.Logger {
	v := ctx.Value(key)
	if v == nil {
		return zap.NewExample()
	}

	l, ok := v.(*zap.Logger)

	if ok {
		return l
	}
	return zap.NewExample()
}
