package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Wa4h1h/psgen/internal/utils"
)

func LoadConfig() (*Config, error) {
	var cfg *Config

	basePath := utils.GetConfigBasePath()
	cfgFilePath := fmt.Sprintf("%s/%s", basePath, utils.ConfigPath)
	logPath := fmt.Sprintf("%s/%s", basePath, utils.LogsPath)

	if _, err := os.Stat(cfgFilePath); err != nil {
		key, errGen := utils.GenerateNewRandomKey()
		if errGen != nil {
			return nil, errGen
		}

		cfg = &Config{
			DbPath:      fmt.Sprintf("%s/%s", basePath, utils.DbPath),
			EncKey:      key,
			ExecTimeout: utils.DefaultExecutionTimeout,
			LogsPath:    logPath,
		}

		file, errF := os.OpenFile(cfgFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if errF != nil {
			return nil, fmt.Errorf("creating config file error: %w", errF)
		}
		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()

		bytes, errMarsh := json.Marshal(cfg)
		if errMarsh != nil {
			return nil, fmt.Errorf("marshal config error: %w", errMarsh)
		}

		_, errW := file.Write(bytes)
		if errW != nil {
			return nil, fmt.Errorf("writing config to file error: %w", errW)
		}

		errP := utils.CreateFolder(logPath)
		if errP != nil {
			return nil, errP
		}
	} else {
		bytes, err := os.ReadFile(cfgFilePath)
		if err != nil {
			return nil, fmt.Errorf("reading config file error: %w", err)
		}

		if err := json.Unmarshal(bytes, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal config error: %w", err)
		}
	}

	return cfg, nil
}

func (c *Config) WriteLogs(bytes []byte) error {
	logFile := fmt.Sprintf("%s/logs_%s.log", c.LogsPath, time.Now().Format("2006-01-02 15:04:05"))

	logs, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening log file error: %w", err)
	}
	defer func() {
		if err := logs.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := logs.Write(bytes); err != nil {
		return fmt.Errorf("writing to log file error: %w", err)
	}

	return nil
}
