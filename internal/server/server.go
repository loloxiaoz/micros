package server

import (
	"fmt"
	"time"

	"micros/api"
	"micros/internal/controller"
	"micros/pkg/registry"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

const (
	beginTime = "beginTime"
	endTime   = "endTime"
)

//Server 服务
type Server struct {
	engine *gin.Engine
}

//NewServer 新建http服务
func NewServer(name string) *Server {
	server := new(Server)
	//handler
	server.engine = gin.New()
	server.engine.Use(statBefore())
	server.engine.Use(statAfter())
	//prometheus
	server.engine.GET("/system/metrics", prometheusHandler())
	//health
	server.engine.GET("/system/health", healthHandler())
	//swagger
	api.SwaggerInfo.BasePath = "/api/v1"

	server.engine .GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//apis
	server.engine.GET("/hello", controller.Helloworld)

	//service discovery
	node := &registry.Node{Id: "1", Address: "127.0.0.1", Port: 8080}
	service := &registry.Service{Name: name, Nodes: []*registry.Node{node}}
	registry.DefaultRegistry.Register(service, registry.RegisterTTL(time.Minute*5))
	return server
}

func (s *Server) assemble() {
	pprof.Register(s.engine) // register pprof to gin


}

//Run 运行http服务
func (s *Server) Run() {
	err := s.engine.Run(":8090")
	if (err!=nil) {
		fmt.Printf("err is %v", err)
	}
}

func statBefore() gin.HandlerFunc {
	return func(c *gin.Context) {
		curTime := time.Now().UnixNano() / 1000000
		c.Set(beginTime, curTime)
	}
}

func statAfter() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
		}()
		c.Next()
	}
}