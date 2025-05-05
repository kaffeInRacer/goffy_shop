// pkg/route/route.go
package route

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
	db *sql.DB
}

func NewRouter(db *sql.DB) *gin.Engine {
	r := &Router{
		Engine: gin.Default(),
		db:     db,
	}
	r.setupRoutes()
	return r.Engine
}

func (r *Router) setupRoutes() {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

}
