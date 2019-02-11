package model

type Task struct {
	JudgingID           int        `json:"judgingid"`
	SubmitID            int        `json:"submitid"`
	CID                 int        `json:"cid"`
	TeamID              int        `json:"teamid"`
	ProbID              int        `json:"probid"`
	LangID              string     `json:"langid"`
	RejudgingID         int        `json:"rejudgingid"`
	EntryPoint          string     `json:"entry_point"`
	OrigSubmitID        int        `json:"origsubmitid"`
	MaxRuntime          float32    `json:"maxruntime"`
	MemLimit            int        `json:"memlimit"`
	OutputLimit         int        `json:"outputlimit"`
	Run                 string     `json:"run"`
	RunMD5Sum           string     `json:"run_md5sum"`
	Compare             string     `json:"compare"`
	CompareMD5Sum       string     `json:"compare_md5sum"`
	CompareArgs         string     `json:"compare_args"`
	CompileScript       string     `json:"compile_script"`
	CompileScriptMD5Sum string     `json:"compile_script_md5sum"`
	TestCases           []TestCase `json:"testcases"`
}
