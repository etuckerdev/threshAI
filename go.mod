module github.com/etuckerdev/threshAI

go 1.21

replace (
	github.com/etuckerdev/Eidos => ./legacy/Eidos
	github.com/etuckerdev/NOUSx => ./legacy/NOUSx
)

require (
	github.com/cornelk/hashmap v1.0.8 // for memory vault
	github.com/go-cmd/cmd v1.4.2 // for task execution
	github.com/spf13/cobra v1.8.0
)

require github.com/spf13/pflag v1.0.5

require github.com/inconshreveable/mousetrap v1.1.0 // indirect
