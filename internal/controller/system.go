package controller

import (
	"net/http"

	"micros/internal/common"

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
		ctx.JSON(http.StatusOK, common.Success(ctx.Request.Host))
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
