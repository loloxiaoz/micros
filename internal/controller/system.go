package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @BasePath /api/v1

// Health godoc
// @Summary health check
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Health
// @Router /system/health [get]
func Health(ctx *gin.Context){
		resp := &struct {
			Timestamp   time.Time `json:"timestamp"`
			Environment string    `json:"environment"`
			Host        string    `json:"host"`
			Status      string    `json:"status"`
		}{
			Timestamp:   time.Now(),
			Host:        ctx.Request.Host,
			Status:      "ok",
		}
		ctx.JSON(http.StatusOK, resp)
}


// @BasePath /api/v1

// Monitor godoc
// @Summary monitor
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Monitor
// @Router /system/monitor [get]
func Monitor(ctx *gin.Context) {
	h := promhttp.Handler()
		h.ServeHTTP(ctx.Writer, ctx.Request)
	
}
