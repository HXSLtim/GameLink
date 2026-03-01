package ws

import "github.com/gin-gonic/gin"

// ServeCustomerServiceWS handles customer-service websocket connections.
// It reuses the common websocket auth/context and hub broadcast logic.
func (h *Handler) ServeCustomerServiceWS(c *gin.Context) {
	h.ServeWS(c)
}
