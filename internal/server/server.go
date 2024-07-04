package server

import (
	"context"
	"micros/api"
	"micros/internal/config"
	"micros/internal/controller"
	"micros/internal/log"
	"net/http"

	fileSwagger "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Server 服务
type Server struct {
	engine *gin.Engine
	conf   *config.Conf
	server *http.Server
}

// New 新建http服务
func New(conf *config.Conf) *Server {
	server := new(Server)
	//conf
	server.conf = conf
	server.engine = gin.New()

	//middleware
	server.engine.Use(logHandler())

	//assemble
	v1 := server.engine.Group("api/v1")
	server.assemble(v1)
	v1.GET("/system/health", controller.Health)
	v1.GET("/example/hello", controller.Helloworld)

	return server
}

func (s *Server) assemble(group *gin.RouterGroup) {
	//apiDoc switch
	if s.conf.IsAPIDoc() {
		api.SwaggerInfo.BasePath = "/api/v1"
		group.GET("/swagger/*any", ginSwagger.WrapHandler(fileSwagger.Handler))
	}
	//profile switch
	if s.conf.IsProfile() {
		pprof.Register(s.engine)
	}
	//monitor switch
	if s.conf.IsMonitor() {
		group.GET("/system/monitor", controller.Monitor)
	}

}

// Run 运行http服务
func (s *Server) Run() error {
	log.Logger().Info("server starting")
	s.server = &http.Server{
		Addr:    s.conf.Addr,
		Handler: s.engine,
	}
	return s.server.ListenAndServe()
}

// Shutdown 停止http服务
func (s *Server) Shutdown(ctx context.Context) error {
	log.Logger().Info("server shutting down")
	return s.server.Shutdown(ctx)
}
