package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/repository"
	"gamelink/internal/service"
	adminservice "gamelink/internal/service/admin"
)

func newTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	installTestErrorHandler(r)
	return r
}

func installTestErrorHandler(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, adminservice.ErrNotFound), errors.Is(err, repository.ErrNotFound):
			status = http.StatusNotFound
		case errors.Is(err, adminservice.ErrValidation), errors.Is(err, service.ErrOrderInvalidTransition):
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}

		c.Abort()
		writeJSONError(c, status, err.Error())
	})
}
