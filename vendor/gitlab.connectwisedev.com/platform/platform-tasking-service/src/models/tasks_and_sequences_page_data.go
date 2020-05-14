package models

import (
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
)

// TasksAndSequencesPageData stores data which is received from GraphQL MS
// for exporting Tasks and Sequences page data which is filtered by widget filters
type TasksAndSequencesPageData struct {
	Tasking struct {
		Summary struct {
			Data struct {
				List []struct {
					Cursor string                    `json:"cursor"`
					Node   apiModels.TaskSummaryData `json:"node"`
				} `json:"edges"`
			} `json:"taskingSummaryPageList"`
		} `json:"TaskingSummaryPage"`
	} `json:"data"`
}
