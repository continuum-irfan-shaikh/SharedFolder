package mocks

import (
	"encoding/json"
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
	"github.com/golang/mock/gomock"
)

// ExecutionResultPersistenceConf is a function type for configuring custom mock of TaskDefinitionPersistance interface
type ExecutionResultPersistenceConf func(er *MockExecutionResultPersistence) *MockExecutionResultPersistence

// ExecutionResultViewPersistenceConf is a function type for configuring custom mock ExecutionResultViewPersistence interface
type ExecutionResultViewPersistenceConf func(ev *MockExecutionResultViewPersistence) *MockExecutionResultViewPersistence

// TaskInstancePersistenceConf is a function type for configuring custom mock of TaskInstancePersistence interface
type TaskInstancePersistenceConf func(ti *MockTaskInstancePersistence) *MockTaskInstancePersistence

// TaskSummaryPersistenceConf is a function type for configuring custom mock of TaskSummaryPersistence interface
type TaskSummaryPersistenceConf func(ts *MockTaskSummaryPersistence) *MockTaskSummaryPersistence

// TaskPersistenceConf is a function type for configuring custom mock of TaskPersistence interface
type TaskPersistenceConf func(tp *MockTaskPersistence) *MockTaskPersistence

// TemplateCacheConf is a function type for configuring custom mock of TemplateCache interface
type TemplateCacheConf func(tc *MockTemplateCache) *MockTemplateCache

// TaskCounterConf is a function type for configuring custom mock of TaskCounter interface
type TaskCounterConf func(tc *MockTaskCounterPersistence) *MockTaskCounterPersistence

// TaskDefinitionConf is a function type for configuring custom mock of TaskDefinitionPersistence interface
type TaskDefinitionConf func(td *MockTaskDefinitionPersistence) *MockTaskDefinitionPersistence

// ExecutionExpirationPersistenceConf is a function type for configuring custom mock of ExecutionExpirationPersistence interface
type ExecutionExpirationPersistenceConf func(ee *MockExecutionExpirationPersistence) *MockExecutionExpirationPersistence

// UserSitesConf is a function type for configuring custom mock of UserSites interface
type UserSitesConf func(us *MockUserSitesPersistence) *MockUserSitesPersistence

// TriggerUCConf is a function type for configuring custom mock of UserSites interface
type TriggerUCConf func(us *MockUsecase) *MockUsecase

// TargetRepoConf is a function type for configuring custom mock of UserSites interface
type TargetRepoConf func(us *MockTargetsRepo) *MockTargetsRepo

// TaskCounterPersistenceConf is a function type for configuring custom mock of TaskCounterPersistence interface
type TaskCounterPersistenceConf func(tc *MockTaskCounterPersistence) *MockTaskCounterPersistence

// DynamicGroupsConf is a func to configure DG mock
type DynamicGroupsConf func(d *MockDynamicGroups) *MockDynamicGroups

// UserUC ..
type UserUC func(u *MockUserUC) *MockUserUC

// AssetConf is a function type for configuring custom mock of Asset interface
type AssetConf func(a *MockAsset) *MockAsset

// DGRepoConf ..
type DGRepoConf func(u *MockDynamicGroupRepo) *MockDynamicGroupRepo

// SitesRepoConf ..
type SitesRepoConf func(a *MockSiteRepo) *MockSiteRepo

// IsErrorMessageValid is used to check if the received error message from handler is the same as expected
// formats from {"message":"gotErrMsg"}. Used fot testing functions
func IsErrorMessageValid(expectedErrorMsg string, gotErrorMsgBody []byte) (isValid bool) {
	translator, err := translation.New("en-us")
	if err != nil {
		return
	}
	var msg common.Message
	err = json.Unmarshal(gotErrorMsgBody, &msg)
	if err != nil {
		return
	}

	expectedErrorMsgTranslated := translator.Translate(expectedErrorMsg)
	if expectedErrorMsgTranslated != msg.Message {
		return
	}
	return true
}

// TranslateErrorMessage is a function to translate error message received from constant to String translated to en-US language
func TranslateErrorMessage(errMsgConstant string) (res string) {
	translator, err := translation.New("en-us")
	if err != nil {
		return
	}
	return fmt.Sprintf(translator.Translate(errMsgConstant))
}

// AnyN returns slice on gomock.Any's
// example: batchMock.EXPECT().Query(gomock.Any(), AnyN(2)...).Times(2)
func AnyN(n int) []interface{} {
	m := make([]interface{}, n)
	for i := 0; i < n; i++ {
		m[i] = gomock.Any()
	}

	return m
}
