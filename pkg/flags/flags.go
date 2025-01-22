package flags

import (
	"log"
	"sync"

	"github.com/spf13/pflag"
)

var (
	brutalFlag  int
	quantumFlag bool
	initOnce    sync.Once
)

func Init(flags *pflag.FlagSet) {
	initOnce.Do(func() {
		log.Println("Initializing flags package")
		flags.IntVarP(&brutalFlag, "brutal", "b", 0, "Brutalization intensity (0-3)")
		flags.BoolVarP(&quantumFlag, "quantum", "q", false, "Enable 11D vector smearing")
	})
}

func Brutal() int {
	return brutalFlag
}

func Quantum() bool {
	return quantumFlag
}
