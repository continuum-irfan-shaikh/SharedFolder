package triggerDefinition

import (
	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
)

// DefinitionService represents triggers definition Usecase
type DefinitionService struct {
	triggerRepo repository.TriggersRepo
}

// NewDefinitionService returns new DefinitionService Usecase
func NewDefinitionService(tr repository.TriggersRepo) *DefinitionService {
	return &DefinitionService{
		triggerRepo: tr,
	}
}

// ImportExternalTriggers make operations with new definitions and removing old ones
func (s *DefinitionService) ImportExternalTriggers(defs []e.TriggerDefinition) (err error) {
	if err = s.triggerRepo.TruncateDefinitions(); err != nil {
		return err
	}
	return s.triggerRepo.InsertDefinitions(defs)
}

// GetTriggerTypes returns Trigger definitions with types
func (s *DefinitionService) GetTriggerTypes() (triggerTypes []e.TriggerDefinition, err error) {
	return s.triggerRepo.GetAllDefinitionsNamesAndIDs()
}
