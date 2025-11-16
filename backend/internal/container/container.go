package container

import (
	"os"

	"github.com/gin-gonic/gin"

	"gamelink/internal/config"
	"gamelink/internal/lifecycle"
	"gamelink/internal/logging"
)

// Application aggregates the top-level dependencies used by the entrypoint.
type Application struct {
	Engine    *gin.Engine
	Config    config.AppConfig
	Lifecycle *lifecycle.Manager
}

// NewApplication builds the dependency graph using the generated wire factory.
func NewApplication() (*Application, error) {
	logging.Init(os.Getenv("LOG_LEVEL"))
	return initializeApplication()
}
