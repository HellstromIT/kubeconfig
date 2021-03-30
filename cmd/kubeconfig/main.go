package main

import (
	"github.com/HellstromIT/kubeconfig/cmd/kubeconfig/internal/kc"
)

var version = "dev"

func main() {
	kc.Cli(version)
}
