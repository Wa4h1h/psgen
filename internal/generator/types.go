package generator

import (
	"github.com/WadhahJemai/psgen/internal/store"
)

type GenFlags struct {
	lowerCase bool
	upperCase bool
	digits    bool
	special   bool
	length    int
}

type Cli struct {
	genFlagSet *GenFlags
	config     *Config
	store      store.Store
}

type Config struct {
	EncKey      string `json:"enc_key"`
	ExecTimeout uint   `json:"execution_timeout"`
	DbPath      string `json:"db_path"`
	LogsPath    string `json:"logs_path"`
}
