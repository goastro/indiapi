package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"

	"github.com/rickbassham/logging"
)

type contextKey struct {
	name string
}

var (
	logEntryCtxKey = &contextKey{"LogEntry"}
)

// RequestLogger represents a way to log all requests.
type RequestLogger interface {
	FromContext(context.Context) logging.Logger
	Handler() func(next http.Handler) http.Handler
}

type requestLogger struct {
	l logging.Logger
}

type requestLog struct {
	RequestID string    `json:"requestId"`
	URL       string    `json:"url"`
	Method    string    `json:"method"`
	StartTime time.Time `json:"startTime"`
	TraceID   string    `json:"traceId"`
	Bytes     int       `json:"bytes,omitempty"`
	Status    int       `json:"status,omitempty"`
	Duration  float32   `json:"duration,omitempty"`
}

// NewRequestLogger creates a new RequestLogger to log all requests.
func NewRequestLogger(log logging.Logger) RequestLogger {
	return &requestLogger{
		l: log,
	}
}

func (l *requestLogger) FromContext(ctx context.Context) logging.Logger {
	log := ctx.Value(logEntryCtxKey)
	if log == nil {
		return l.l
	}

	return log.(logging.Logger)
}

func (l *requestLogger) Handler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			traceID := uuid.New()

			request := requestLog{
				RequestID: middleware.GetReqID(r.Context()),
				URL:       r.URL.String(),
				Method:    r.Method,
				StartTime: t1,
				TraceID:   traceID.String(),
			}

			entry := l.l.WithField("request", request)

			entry.Info("request started")

			ww := &logWriter{
				w: w,
			}

			defer func(entry logging.Logger, request requestLog, ww *logWriter) {
				request.Bytes = ww.bytesWritten
				request.Status = ww.statusCode
				request.Duration = float32(time.Now().Sub(t1)) / float32(time.Second)

				entry = entry.WithField("request", request)

				if ww.statusCode >= 400 {
					entry.Warn("request complete")
				} else {
					entry.Info("request complete")
				}
			}(entry, request, ww)

			ctx := context.WithValue(r.Context(), logEntryCtxKey, entry)

			next.ServeHTTP(ww, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

type logWriter struct {
	w http.ResponseWriter

	statusCode   int
	bytesWritten int
}

func (w *logWriter) Header() http.Header {
	return w.w.Header()
}

func (w *logWriter) Write(data []byte) (int, error) {
	written, err := w.w.Write(data)
	w.bytesWritten += written
	return written, err
}

func (w *logWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}
