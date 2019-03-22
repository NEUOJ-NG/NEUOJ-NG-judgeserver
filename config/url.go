package config

type urlConfig struct {
	Username    string `toml:"username"`
	Password    string `toml:"password"`
	Submissions string `toml:"submissions"`
	Executables string `toml:"executables"`
	TestCases   string `toml:"test_cases"`
	Judgings    string `toml:"judgings"`
	JudgingRuns string `toml:"judging_runs"`
}
