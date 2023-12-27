package generator

import (
	"github.com/WadhahJemai/psgen/internal/store"
	"time"
)

type GenFlags struct {
	lowerCase bool
	upperCase bool
	digits    bool
	special   bool
	length    int
}

type Cli struct {
	genFlagSet  *GenFlags
	store       store.Store
	key         string
	execTimeout time.Duration
}
