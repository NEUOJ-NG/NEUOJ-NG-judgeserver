package config

import (
	"encoding/json"
	"errors"
	"sync"
)

type judgehostConfig struct {
	Username      string `toml:"username"`
	Password      string `toml:"password"`
	Configuration string `toml:"configuration"`
}

var (
	judgehostConfiguration map[string]interface{}
	judgehostConfigOnce    sync.Once
)

func GetJudgehostConfiguration(name string, getAll bool) (interface{}, error) {
	judgehostConfigOnce.Do(func() {
		judgehostConfiguration = make(map[string]interface{})
		err := json.Unmarshal(
			[]byte(GetConfig().Judgehost.Configuration),
			&judgehostConfiguration)
		if err != nil {
			panic(err)
		}
	})

	if getAll {
		return judgehostConfiguration, nil
	}

	if v, ok := judgehostConfiguration[name]; ok {
		return v, nil
	} else {
		return nil, errors.New("no configuration named " + name)
	}
}
