package main

import (
	"fmt"

	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/env"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/procParser"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/services"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/services/model"
)

type healthCheckDependencyImpl struct {
	services.HealthCheckServiceFactoryImpl
	services.HealthCheckDalFactoryImpl
	services.VersionFactoryImpl
	env.FactoryEnvImpl
	procParser.ParserFactoryImpl
}

func main() {
	h := services.HealthCheckServiceFactoryImpl{}
	s := h.GetHealthCheckService(healthCheckDependencyImpl{})
	model.StrartTime = time.Now()
	health, _ := s.GetHealthCheck(model.HealthCheck{
		Version: model.Version{
			SolutionName:    "SolutionName",
			ServiceName:     "ServiceName",
			ServiceProvider: "ContinuumLLC",
			Major:           "1",
			Minor:           "1",
			Patch:           "11",
		},
		ListenPort: ":8081",
	})

	fmt.Println(health)
}
