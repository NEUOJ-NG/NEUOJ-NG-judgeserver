package main

import (
	backendUtil "github.com/NEUOJ-NG/NEUOJ-NG-backend/util"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/config"
	c "github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/controller"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/mq"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/router"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// setup storage dir
	backendUtil.CreateDirOrPanic(config.GetSubmissionStoragePath())
	backendUtil.CreateDirOrPanic(config.GetTestCaseStoragePath())
	backendUtil.CreateDirOrPanic(config.GetExecutableStoragePath())

	// setup log
	backendUtil.SetupLog(true, true)

	// create Gin Engine with Logger and Recovery middleware
	app := gin.Default()

	// init router
	// test
	app.GET("/ping", c.Ping)
	// DOMJudge RESTful API
	// protected by HTTP Simple Auth
	v4 := app.Group("/api/v4", gin.BasicAuth(gin.Accounts{
		config.GetConfig().Judgehost.Username: config.GetConfig().Judgehost.Password,
	}))
	router.InitRouter(v4)

	// init message queue
	err := mq.InitConsumerMQ()
	if err != nil {
		log.Fatalf("failed to init consumer message queue: %s", err.Error())
		return
	} else {
		defer mq.ConsumerConnection.Close()
		defer mq.ConsumerChannel.Close()
	}

	// start hot update handler
	// config will be reloaded with SYSUSR1 signal
	backendUtil.SetupConfigHotUpdate()

	// start server with endless
	// server will reload with HUP signal
	// server will stop with INT signal
	server := endless.NewServer(
		config.GetConfig().App.Addr,
		app,
	)
	server.BeforeBegin = func(add string) {
		log.Info("NEUOJ-NG-judgeserver started")
		log.Infof("listen %v", config.GetConfig().App.Addr)
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
	log.Info("NEUOJ-NG-judgeserver terminated")
}
