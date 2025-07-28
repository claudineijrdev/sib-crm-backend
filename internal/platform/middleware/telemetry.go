package middleware

import (
	"strconv"
	"time"

	"github.com/claudineijrdev/sib-crm-backend/internal/platform/telemetry"
	"github.com/gin-gonic/gin"
)

func TelemetryMiddleware(telemetryService telemetry.TelemetryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Criar span para a requisição
		span, ctx := telemetryService.StartSpan(c.Request.Context(), "http.request")
		defer span.End()
		
		// Adicionar contexto ao gin
		c.Request = c.Request.WithContext(ctx)
		
		// Adicionar tags básicas
		span.SetTag("http.method", c.Request.Method)
		span.SetTag("http.url", c.Request.URL.Path)
		span.SetTag("http.user_agent", c.Request.UserAgent())
		
		// Processar requisição
		c.Next()
		
		// Adicionar tags de resposta
		span.SetTag("http.status_code", strconv.Itoa(c.Writer.Status()))
		span.SetTag("http.response_size", strconv.Itoa(c.Writer.Size()))
		
		// Track metric de duração
		duration := time.Since(start).Milliseconds()
		telemetryService.TrackMetric(ctx, telemetry.Metric{
			Name:  "http.request.duration",
			Value: float64(duration),
			Tags: map[string]string{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
				"status": strconv.Itoa(c.Writer.Status()),
			},
		})
		
		// Track event de requisição
		telemetryService.TrackEvent(ctx, telemetry.Event{
			Name: "http.request.completed",
			Properties: map[string]interface{}{
				"method":      c.Request.Method,
				"path":        c.Request.URL.Path,
				"status_code": c.Writer.Status(),
				"duration_ms": duration,
			},
			Timestamp: time.Now(),
		})
		
		// Marcar erro se houver
		if len(c.Errors) > 0 {
			span.SetError(c.Errors.Last().Err)
		}
	}
} 