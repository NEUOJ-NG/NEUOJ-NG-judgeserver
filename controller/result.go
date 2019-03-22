package controller

import (
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/config"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/form"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// update compile results
func PutJudgings(ctx *gin.Context) {
	id := ctx.Param("id")
	judgehost := ctx.PostForm("judgehost")
	compileSuccess := ctx.PostForm("compile_success")
	outputCompile := ctx.PostForm("output_compile")

	log.Debug(id)
	log.Debug(judgehost)
	log.Debug(compileSuccess)
	log.Debug(outputCompile)

	// PUT compile result to backend
	// create request with basic auth
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.SetBasicAuth(config.GetConfig().URL.Username, config.GetConfig().URL.Password)
			return nil
		},
	}
	judgingForm := url.Values{
		"judgehost":       {judgehost},
		"compile_success": {compileSuccess},
		"output_compile":  {outputCompile},
	}
	req, err := http.NewRequest(
		"PUT",
		config.GetConfig().URL.Judgings+id,
		strings.NewReader(judgingForm.Encode()),
	)
	if err != nil {
		log.Errorf("failed to create PUT request for compile result: %v", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	// IMPORTANT: remember to set content type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.GetConfig().URL.Username, config.GetConfig().URL.Password)

	// do the request
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("failed to do PUT request for compile result: %v", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// check status code
	if resp.StatusCode != http.StatusOK {
		log.Errorf("failed to do PUT request for compile result (status %v)", resp.StatusCode)
		respBodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Error(string(respBodyBytes))
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func PostJudgingRuns(ctx *gin.Context) {
	var judgingRunResult form.JudgingRunResult
	if err := ctx.ShouldBind(&judgingRunResult); err != nil {
		log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		log.Debugf(
			"judging run result for judging %s test case %s is %s",
			judgingRunResult.JudgingID,
			judgingRunResult.TestCaseID,
			judgingRunResult.RunResult,
		)
		log.Debugf("runtime: %s", judgingRunResult.Runtime)
		log.Debugf("output_system: %s", judgingRunResult.OutputSystem)

		// send judging run result to backend
		// create request with basic auth
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				req.SetBasicAuth(config.GetConfig().URL.Username, config.GetConfig().URL.Password)
				return nil
			},
		}
		resultForm := judgingRunResult.ConvertToForm()
		req, err := http.NewRequest(
			"POST",
			config.GetConfig().URL.JudgingRuns,
			strings.NewReader(resultForm.Encode()),
		)
		if err != nil {
			log.Errorf("failed to create POST request for run result: %v", err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// IMPORTANT: remember to set content type
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(config.GetConfig().URL.Username, config.GetConfig().URL.Password)

		// do the request
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			log.Errorf("failed to do POST request for run result: %v", err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// check status code
		if resp.StatusCode != http.StatusOK {
			log.Errorf("failed to do POST request for run result (status %v)", resp.StatusCode)
			respBodyBytes, _ := ioutil.ReadAll(resp.Body)
			log.Error(string(respBodyBytes))
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}
}
