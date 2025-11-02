package main

import (
	"github.com/navetacandra/gowle/internal/config"
)

func main() {
	appConfig := config.GowleConfig{}
	appConfig.Load()
}
