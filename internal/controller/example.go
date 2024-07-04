package controller

import (
	"net/http"

	"micros/internal/common"
	"micros/internal/log"
	"micros/internal/model"
	"micros/internal/service"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// Helloworld godoc
// @Summary ping example
// @Produce json
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Success 200 {string} Hello
// @Router /example/hello [get]
func Helloworld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, common.Success("hello world"))
}

// CreateStudent godoc
// @Summary student example
// @Produce json
// @Schemes
// @Param student body model.Student true "学生信息"
// @Description 创建学生
// @Tags example
// @Success 200 {string} Student "学生"
// @Router /student [put]
func CreateStudent(ctx *gin.Context) {
	var student model.Student
	if err := ctx.ShouldBindJSON(&student); err != nil {
		log.Logger().Errorf("参数错误 %s", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Error(common.BindJSONError, err))
		return 
	}
	if err := service.CreateStudent(&student); err != nil {
		ctx.JSON(http.StatusOK, common.Error(common.CommonServerError, err))
		return 
	}
	ctx.JSON(http.StatusOK, student)
}
