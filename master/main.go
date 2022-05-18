package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/master/config"
	"github.com/AlaricGilbert/argos-core/master/dal"
	"github.com/AlaricGilbert/argos-core/master/handlers"
	master "github.com/AlaricGilbert/argos-core/master/kitex_gen/master/argosmaster"
	"github.com/cloudwego/kitex/server"
	"github.com/gin-gonic/gin"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.AddHook(lfshook.NewHook(
		fmt.Sprintf("logs/%s_master.log", time.Now().Format(time.RFC3339)),
		&logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: true,
		},
	))

	argos.SetLogger(logger)

	dal.InitDatabase()
	go startGinServer()

	addr, _ := net.ResolveTCPAddr("tcp", "localhost:4222")

	svr := master.NewServer(new(ArgosMasterImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

func startGinServer() {
	r := gin.Default()

	task := r.Group("task")
	task.GET("/list", handlers.GetTasks)
	task.POST("/write", handlers.WriteTask)

	r.GET("status", handlers.GetStatus)

	query := r.Group("query")
	query.GET("/time", handlers.QueryByTime)
	query.GET("/ip", handlers.QueryByIP)
	query.GET("/tx", handlers.QueryByTx)
	r.Run(config.WebListenAddr) // listen and serve on 0.0.0.0:8080
}
