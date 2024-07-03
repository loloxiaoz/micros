package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// Helloworld godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Hello
// @Router /example/hello [get]
func Helloworld(ctx *gin.Context)  {
	ctx.JSON(http.StatusOK,"hello world")
 }
 