package form

import "net/url"

type JudgingRunResult struct {
	JudgingID    string `form:"judgingid" binding:"required"`
	TestCaseID   string `form:"testcaseid" binding:"required"`
	RunResult    string `form:"runresult" binding:"required"`
	Runtime      string `form:"runtime" binding:"required"`
	Judgehost    string `form:"judgehost" binding:"required"`
	OutputRun    string `form:"output_run"`
	OutputError  string `form:"output_error"`
	OutputSystem string `form:"output_system"`
	OutputDiff   string `form:"output_diff"`
}

func (result JudgingRunResult) ConvertToForm() *url.Values {
	return &url.Values{
		"judgingid":     {result.JudgingID},
		"testcaseid":    {result.TestCaseID},
		"runresult":     {result.RunResult},
		"runtime":       {result.Runtime},
		"judgehost":     {result.Judgehost},
		"output_run":    {result.OutputRun},
		"output_error":  {result.OutputError},
		"output_system": {result.OutputSystem},
		"output_diff":   {result.OutputDiff},
	}
}
