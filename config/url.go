package config

type urlConfig struct {
	Submissions string `toml:"submissions"`
	Executables string `toml:"executables"`
	TestCases   string `toml:"test_cases"`
}
