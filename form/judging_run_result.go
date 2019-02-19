package form

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
