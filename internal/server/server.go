package server

import (
	"fmt"
	"time"

	"micros/api"
	"micros/internal/config"
	"micros/internal/controller"
	"micros/pkg/registry"

	fileSwagger "github.com/swaggo/files"
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
	conf *config.Conf
}

//NewServer 新建http服务
func NewServer(conf *config.Conf) *Server {
	server := new(Server)
	//conf
	server.conf = conf
	server.engine = gin.New()

	//assemble
	server.assemble()

	server.engine.Use(statBefore())
	server.engine.Use(statAfter())

	//apis
	server.engine.GET("/system/health", healthHandler())
	server.engine.GET("/hello", controller.Helloworld)

	//service discovery
	node := &registry.Node{Id: "1", Address: "127.0.0.1", Port: 8080}
	service := &registry.Service{Name: conf.Profile, Nodes: []*registry.Node{node}}
	registry.DefaultRegistry.Register(service, registry.RegisterTTL(time.Minute*5))
	return server
}

func (s *Server) assemble() {
	//apiDoc switch
	if s.conf.IsAPIDoc() {
		api.SwaggerInfo.BasePath = "/api/v1"
		s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(fileSwagger.Handler))
	}
	//profile switch
	if s.conf.IsProfile() {
		pprof.Register(s.engine) 
	}
	//monitor switch
	if s.conf.IsMonitor() {
		s.engine.GET("/system/metrics", prometheusHandler())
	}

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