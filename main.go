package main

import (
	"github.com/NEUOJ-NG/NEUOJ-NG-backend/util"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/config"
	c "github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/controller"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/router"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// setup storage dir
	util.CreateDirOrPanic(config.GetSubmissionStoragePath())
	util.CreateDirOrPanic(config.GetTestCaseStoragePath())
	util.CreateDirOrPanic(config.GetExecutableStoragePath())

	// setup log
	util.SetupLog(true, true)

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

	// start hot update handler
	// config will be reloaded with SYSUSR1 signal
	util.SetupConfigHotUpdate()

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
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("failed to start server")
		log.Fatal(err)
	}
	log.Info("NEUOJ-NG-judgeserver terminated")
}
