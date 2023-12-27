package generator

import (
	"encoding/json"
	"fmt"
	"github.com/WadhahJemai/psgen/internal/utils"
	"os"
)

func LoadConfig() *Config {
	var cfg *Config

	basePath := utils.GetConfigBasePath()
	cfgFilePath := fmt.Sprintf("%s/%s", basePath, utils.ConfigPath)

	if _, err := os.Stat(cfgFilePath); err != nil {

		key, errGen := utils.GenerateNewRandomKey()
		if errGen != nil {
			panic(fmt.Errorf("failed loading config: %w", errGen))
		}

		cfg = &Config{
			DbPath:      fmt.Sprintf("%s/%s", basePath, utils.DbPath),
			EncKey:      key,
			ExecTimeout: utils.DefaultExecutionTimeout,
		}

		file, errF := os.OpenFile(cfgFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if errF != nil {
			panic(fmt.Errorf("failed loading config: %w", errF))
		}

		bytes, errMarsh := json.Marshal(cfg)
		if errMarsh != nil {
			panic(fmt.Errorf("failed loading config: %w", errMarsh))
		}

		_, errW := file.Write(bytes)
		if errW != nil {
			panic(fmt.Errorf("failed loading config: %w", errW))
		}
	}

	bytes, err := os.ReadFile(cfgFilePath)
	if err != nil {
		panic(fmt.Errorf("failed loading config: %w", err))
	}

	if err := json.Unmarshal(bytes, &cfg); err != nil {
		panic(fmt.Errorf("failed loading config: %w", err))
	}

	return cfg
}
