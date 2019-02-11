package model

type TestCase struct {
	MD5SumInput  string `json:"md5sum_input"`
	MD5SumOutput string `json:"md5sum_output"`
	TestCaseID   int    `json:"testcaseid"`
	Rank         int    `json:"rank"`
}
