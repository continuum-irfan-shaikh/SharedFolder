package repository

import "github.com/gocql/gocql"

//go:generate mockgen -destination=../mocks/mocks-gomock/profiles_mock.go  -package=mocks -source=./profiles.go

// Profiles represents interface to retrieve info about ProfileID from agent-config service that stored in cassandra
type Profiles interface {
	GetByTaskID(taskID gocql.UUID) (gocql.UUID, error)
	Insert(taskID, profileID gocql.UUID) error
	Delete(taskID gocql.UUID) error
}
