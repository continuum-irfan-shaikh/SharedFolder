package patching

import (
	"time"

	"github.com/google/uuid"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/agent"
)

// WindowsAssessmentResult structure represents the Kafka message
type WindowsAssessmentResult struct {
	agent.BrokerEnvelope
	Message WindowsAssessmentMessage `json:"message"`
}

// WindowsAssessmentMessage contains information about OS patches
type WindowsAssessmentMessage struct {
	InstalledPatches []WindowsPatchDetails `json:"installedPatches"`
	MissingPatches   []WindowsPatchDetails `json:"missingPatches"`
	CatalogOnlyIDs   []string              `json:"catalogOnlyIDs"`
}

// WindowsPatchDetails is a definition of OS patch
type WindowsPatchDetails struct {
	LastDeploymentChangeTime time.Time   `json:"lastDeploymentChangeTime"`
	KBArticleIDs             []string    `json:"kbArticleIDs"`
	Categories               []string    `json:"categories"`
	CategoryIDs              []uuid.UUID `json:"categoryIDs"`
	UpdateID                 uuid.UUID   `json:"updateID"`
	RebootBehavior           string      `json:"rebootBehavior"`
	Title                    string      `json:"title"`
	Description              string      `json:"description"`
	Type                     string      `json:"type"`
	MsrcSeverity             string      `json:"msrcSeverity"`
	SupportURL               string      `json:"supportUrl"`
	RevisionNumber           int32       `json:"revisionNumber"`
	IsInstalled              bool        `json:"isInstalled"`
	IsDownloaded             bool        `json:"isDownloaded"`
	IsUninstallable          bool        `json:"isUninstallable"`
	PendingReboot            bool        `json:"pendingReboot"`
}
