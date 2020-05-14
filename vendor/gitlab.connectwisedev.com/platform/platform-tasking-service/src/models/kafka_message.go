package models

import "github.com/gocql/gocql"

// KafkaMessage defines kafka message model for WS
type KafkaMessage struct {
	PartnerID           string      `json:"partnerID"`
	UID                 string      `json:"uid"`
	TaskID              gocql.UUID  `json:"taskID"`
	IsRequiredNOCAccess bool        `json:"isRequiredNOCAccess"`
	Entity              interface{} `json:"entity"`
}

// TaskIsUpdatedMessage shows if common fields of task are updated or not
type TaskIsUpdatedMessage struct {
	PartnerID           string      `json:"partnerID"`
	TaskID              gocql.UUID  `json:"taskID"`
}