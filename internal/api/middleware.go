package api

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/logging"
)

// RequestIDMiddleware garantiza que todo request que pasa por la API tenga
// un id de correlación (request_id):
//   - si quien llama (el frontend, u otro servicio) ya mandó uno en el
//     header X-Request-ID, lo reusamos, para poder seguir un mismo request
//     a través de varios servicios;
//   - si no vino ninguno, generamos uno nuevo acá.
//
// Ese id se: (1) guarda en el context.Context del request, para que
// cualquier log emitido mientras se procesa el request lo incluya;
// (2) devuelve en el header de la respuesta, para que el frontend pueda
// loguear el mismo id y así correlacionar sus logs con los del backend;
// (3) loguea en una única línea estructurada de "access log" por request,
// con método, path, status y latencia.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(logging.HeaderRequestID)
		if requestID == "" {
			requestID = generateRequestID()
		}

		ctx := logging.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(logging.HeaderRequestID, requestID)

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		logging.FromContext(ctx).Info("http_request",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", latency.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

func generateRequestID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "unknown-request-id"
	}
	return hex.EncodeToString(buf)
}
