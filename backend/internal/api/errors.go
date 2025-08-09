package api

import (
    "errors"
    "net/http"
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
    }
    rid := requestIDFromContext(c)
    c.JSON(code, gin.H{"error": errorPayload{Code: errCode, Message: message, Details: details, TraceID: rid}})
}

// respondOK standardizes success envelopes (optional wrapper for future metadata).
func respondOK(c *gin.Context, payload interface{}) {
    c.JSON(http.StatusOK, gin.H{"data": payload, "trace_id": requestIDFromContext(c)})
}

func requestIDFromContext(c *gin.Context) string {
    if v, ok := c.Get("request_id"); ok {
        if s, ok2 := v.(string); ok2 { return s }
    }
    return ""
}
