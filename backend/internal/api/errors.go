package api

import (
    "context"
    "errors"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
)

// Domain / sentinel errors
var (
    ErrNotFound     = errors.New("not_found")
    ErrValidation   = errors.New("validation_error")
    ErrConflict     = errors.New("conflict")
    ErrInternal     = errors.New("internal_error")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrTimeout      = errors.New("timeout")
)

type errorPayload struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
    TraceID string      `json:"trace_id,omitempty"`
}

// respondError writes a standardized error envelope.
func respondError(c *gin.Context, err error, message string, details interface{}) {
    code := http.StatusInternalServerError
    errCode := "internal_error"
    switch {
    case errors.Is(err, ErrValidation):
        code = http.StatusBadRequest; errCode = "validation_error"
    case errors.Is(err, ErrNotFound):
        code = http.StatusNotFound; errCode = "not_found"
    case errors.Is(err, ErrConflict):
        code = http.StatusConflict; errCode = "conflict"
    case errors.Is(err, ErrUnauthorized):
        code = http.StatusUnauthorized; errCode = "unauthorized"
    case errors.Is(err, ErrForbidden):
        code = http.StatusForbidden; errCode = "forbidden"
    case errors.Is(err, ErrTimeout):
        code = http.StatusGatewayTimeout; errCode = "timeout"
    }
    rid := requestIDFromContext(c)
    c.JSON(code, gin.H{"error": errorPayload{Code: errCode, Message: message, Details: details, TraceID: rid}})
}

// respondOK standardizes success envelopes (optional wrapper for future metadata).
func respondOK(c *gin.Context, payload interface{}) {
    c.JSON(http.StatusOK, gin.H{"data": payload, "trace_id": requestIDFromContext(c)})
}

func respondCreated(c *gin.Context, payload interface{}) {
    c.JSON(http.StatusCreated, gin.H{"data": payload, "trace_id": requestIDFromContext(c)})
}

func requestIDFromContext(c *gin.Context) string {
    if v, ok := c.Get("request_id"); ok {
        if s, ok2 := v.(string); ok2 { return s }
    }
    return ""
}

// timeoutMiddleware enforces per-request timeouts; cancels context passed to handlers.
func (r *Router) timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
        defer cancel()
        c.Request = c.Request.WithContext(ctx)
        done := make(chan struct{})
        go func() { c.Next(); close(done) }()
        select {
        case <-done:
            return
        case <-ctx.Done():
            // Only write if not already written
            if !c.Writer.Written() {
                respondError(c, ErrTimeout, "request timed out", nil)
            }
            c.Abort()
        }
    }
}
