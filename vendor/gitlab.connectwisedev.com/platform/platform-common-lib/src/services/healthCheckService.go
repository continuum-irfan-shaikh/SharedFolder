package services

import (
	aModel "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/healthCheck"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/services/model"
)

//HealthCheckServiceFactoryImpl returns the concrete implementation of Factory
type HealthCheckServiceFactoryImpl struct {
}

//GetHealthCheckService : A factory function to create an instance of HealthCheck Service
func (HealthCheckServiceFactoryImpl) GetHealthCheckService(f model.HealthCheckDependencies) model.HealthCheckService {
	return healthCheckServiceImpl{
		f: f,
	}
}

//healthCheckServiceImpl returns the concrete implementation of HealthCheckService
type healthCheckServiceImpl struct {
	f model.HealthCheckDependencies
}

func (h healthCheckServiceImpl) GetHealthCheck(healthCheck model.HealthCheck) (aModel.HealthCheck, error) {
	logger.Get().Debug(healthCheck.Version.ServiceName, "Retrieving Health Information for %v", healthCheck)
	dal := h.f.GetHealthCheckDal(h.f)
	return dal.GetHealthCheck(healthCheck)
}
