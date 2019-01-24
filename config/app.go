package config

import "path/filepath"

const (
	SubmissionDir = "submissions"
	TestCaseDir = "test_cases"
	ExecutableDir = "executables"
)

type appConfig struct {
	Addr          string   `toml:"addr"`
	LogFile       string   `toml:"log_file"`
	LogLevel      string   `toml:"log_level"`
	StoragePath   string   `toml:"storage_path"`
}

func GetSubmissionStoragePath() string {
	return filepath.Join(GetConfig().App.StoragePath, SubmissionDir)
}

func GetTestCaseStoragePath() string {
	return filepath.Join(GetConfig().App.StoragePath, TestCaseDir)
}

func GetExecutableStoragePath() string {
	return filepath.Join(GetConfig().App.StoragePath, ExecutableDir)
}
