package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
)

type recoverer struct {
	l RequestLogger
}

type Recoverer interface {
	Handler() func(next http.Handler) http.Handler
}

func NewRecoverer(l RequestLogger) Recoverer {
	return &recoverer{
		l: l,
	}
}

func (rec *recoverer) Handler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logEntry := rec.l.FromContext(r.Context())

					err, ok := recovered.(error)
					if !ok {
						err = fmt.Errorf("%v", recovered)
					}

					wrapped := errors.Wrap(err, 2)

					logEntry.WithField("stackTrace", wrapped.StackFrames()).WithError(err).Error("recovered from panic")
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
