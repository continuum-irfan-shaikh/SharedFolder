package legacy

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
)

const partner = "1"
const uuidString = "123e4567-e89b-12d3-a456-426655440000"
const defaultMsg = `failed on unexpected value of result "%v"`

func init()  {
	logger.Load(config.Config.Log)
}

func TestUsecase_GetByPartner(t *testing.T) {
	RegisterTestingT(t)
	validUUID, _ := gocql.ParseUUID(uuidString)

	testCases := []struct {
		name        string
		mock        func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence)
		expected    []models.LegacyScriptInfo
		expectedErr string
	}{
		{
			name: "testCase 1 - GetByPartner server error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllScriptInfoByPartner(partner).Return([]models.LegacyScriptInfo{}, fmt.Errorf("err"))
				return migration, nil, nil
			},
			expectedErr: "server error. GetByPartner: partner 1 :err",
		},
		{
			name: "testCase 2 - GetByPartner not found error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllScriptInfoByPartner(partner).Return([]models.LegacyScriptInfo{}, gocql.ErrNotFound)
				return migration, nil, nil
			},
			expectedErr: "not found. GetByPartner: nothing found by partner 1",
		},
		{
			name: "testCase 3 - GetByPartner zero results",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllScriptInfoByPartner(partner).Return([]models.LegacyScriptInfo{}, nil)
				return migration, nil, nil
			},
			expected: []models.LegacyScriptInfo{},
		},
		{
			name: "testCase 3.5 - GetByPartner returned data with invalid definitionID ",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllScriptInfoByPartner(partner).Return([]models.LegacyScriptInfo{
					{DefinitionID: "invalid"},
				}, nil)

				return migration, nil, nil
			},
			expected: []models.LegacyScriptInfo{
				{DefinitionID: "invalid"},
			},
		},
		{
			name: "testCase 4 - definition Repo returned error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence) {

				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllScriptInfoByPartner(partner).Return([]models.LegacyScriptInfo{{
					DefinitionID: uuidString}}, nil)

				td := mocks.NewMockTaskDefinitionPersistence(ctl)
				td.EXPECT().
					GetByID(gomock.Any(), partner, validUUID).
					Return(models.TaskDefinitionDetails{}, fmt.Errorf("err"))
				return migration, nil, td
			},
			expected: []models.LegacyScriptInfo{{DefinitionID: uuidString}},
		},
	}

	for _, tc := range testCases {
		func() {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			m, _, td := tc.mock(ctl)
			uc := NewMigrationUsecase(m, td, logger.Log)
			actual, err := uc.GetByPartner(partner)
			if tc.expectedErr != "" {
				Ω(err).NotTo(BeNil(), fmt.Sprintf(defaultMsg, tc.name))
				Ω(err.Error()).To(Equal(tc.expectedErr), fmt.Sprintf(defaultMsg, tc.name))
				return
			}

			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, tc.name))
			Ω(actual).To(Equal(tc.expected), fmt.Sprintf(defaultMsg, tc.name))
		}()
	}
}

func TestUsecase_GetJobsByPartner(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		name        string
		mock        func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger)
		expected    []models.LegacyJobInfo
		expectedErr string
	}{
		{
			name: "testCase 1 - GetJobsByPartner server error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllJobsInfoByPartner(partner).Return([]models.LegacyJobInfo{}, fmt.Errorf("err"))
				return migration, nil
			},
			expectedErr: "server error. GetJobsByPartner: partner 1 :err",
		},
		{
			name: "testCase 2 - GetByPartner not found error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllJobsInfoByPartner(partner).Return([]models.LegacyJobInfo{}, gocql.ErrNotFound)
				return migration, nil
			},
			expectedErr: "not found. GetJobsByPartner: nothing found by partner 1",
		},
		{
			name: "testCase 3 - GetByPartner ok",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllJobsInfoByPartner(partner)
				return migration, nil
			},
		},
	}

	for _, tc := range testCases {
		func() {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			m, _ := tc.mock(ctl)
			uc := NewMigrationUsecase(m, nil, logger.Log)
			actual, err := uc.GetJobsByPartner(partner)
			if tc.expectedErr != "" {
				Ω(err).NotTo(BeNil(), fmt.Sprintf(defaultMsg, tc.name))
				Ω(err.Error()).To(Equal(tc.expectedErr), fmt.Sprintf(defaultMsg, tc.name))
				return
			}

			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, tc.name))
			Ω(actual).To(Equal(tc.expected), fmt.Sprintf(defaultMsg, tc.name))
		}()
	}
}

func TestUsecase_GetByJobID(t *testing.T) {
	RegisterTestingT(t)
	jobID := "id"
	testCases := []struct {
		name        string
		mock        func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger)
		expected    models.LegacyJobInfo
		expectedErr string
	}{
		{
			name: "testCase 1 - GetByPartner server error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetByJobID(partner, jobID).Return(models.LegacyJobInfo{}, fmt.Errorf("err"))
				return migration, nil
			},
			expectedErr: "server error. GetByJobID: partner 1 :err",
		},
		{
			name: "testCase 2 - GetByPartner not found error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetByJobID(partner, jobID).Return(models.LegacyJobInfo{}, gocql.ErrNotFound)
				return migration, nil
			},
			expectedErr: "not found. GetByJobID: nothing found by partner 1",
		},
		{
			name: "testCase 3 - GetByPartner ok",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetByJobID(partner, jobID)
				return migration, nil
			},
		},
	}

	for _, tc := range testCases {
		func() {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			m, _ := tc.mock(ctl)
			uc := NewMigrationUsecase(m, nil, logger.Log)
			actual, err := uc.GetByJobID(partner, jobID)
			if tc.expectedErr != "" {
				Ω(err).NotTo(BeNil(), fmt.Sprintf(defaultMsg, tc.name))
				Ω(err.Error()).To(Equal(tc.expectedErr), fmt.Sprintf(defaultMsg, tc.name))
				return
			}

			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, tc.name))
			Ω(actual).To(Equal(tc.expected), fmt.Sprintf(defaultMsg, tc.name))
		}()
	}
}

func Test_GetByScriptID(t *testing.T) {
	RegisterTestingT(t)
	validUUID, _ := gocql.ParseUUID(uuidString)

	testCases := []struct {
		name        string
		mock        func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence)
		expected    []models.LegacyScriptInfo
		expectedErr string
	}{
		{
			name: "testCase 1 - GetByPartner server error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllScriptInfoByPartner(partner).Return([]models.LegacyScriptInfo{}, fmt.Errorf("err"))
				return migration, nil, nil
			},
			expectedErr: "server error. GetByPartner: partner 1 :err",
		},
		{
			name: "testCase 2 - GetByPartner not found error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllScriptInfoByPartner(partner).Return([]models.LegacyScriptInfo{}, gocql.ErrNotFound)
				return migration, nil, nil
			},
			expectedErr: "not found. GetByPartner: nothing found by partner 1",
		},
		{
			name: "testCase 3 - GetByPartner returned data with invalid definitionID ",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllScriptInfoByPartner(partner).Return([]models.LegacyScriptInfo{
					{DefinitionID: "invalid"},
				}, nil)

				return migration, nil, nil
			},
			expected: []models.LegacyScriptInfo{
				{DefinitionID: "invalid"},
			},
		},
		{
			name: "testCase 4 - definition Repo returned error",
			mock: func(ctl *gomock.Controller) (repository.LegacyMigration, logger.Logger, models.TaskDefinitionPersistence) {
				migration := mocks.NewMockLegacyMigration(ctl)
				migration.EXPECT().GetAllScriptInfoByPartner(partner).Return([]models.LegacyScriptInfo{{
					DefinitionID: uuidString}}, nil)

				td := mocks.NewMockTaskDefinitionPersistence(ctl)
				td.EXPECT().
					GetByID(gomock.Any(), partner, validUUID).
					Return(models.TaskDefinitionDetails{}, fmt.Errorf("err"))
				return migration, nil, td
			},
			expected: []models.LegacyScriptInfo{{DefinitionID: uuidString}},
		},
	}

	for _, tc := range testCases {
		func() {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			m, _, td := tc.mock(ctl)
			uc := NewMigrationUsecase(m, td, logger.Log)
			actual, err := uc.GetByPartner(partner)
			if tc.expectedErr != "" {
				Ω(err).NotTo(BeNil(), fmt.Sprintf(defaultMsg, tc.name))
				Ω(err.Error()).To(Equal(tc.expectedErr), fmt.Sprintf(defaultMsg, tc.name))
				return
			}

			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, tc.name))
			Ω(actual).To(Equal(tc.expected), fmt.Sprintf(defaultMsg, tc.name))
		}()
	}
}
