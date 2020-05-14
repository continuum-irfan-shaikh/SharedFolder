package taskCounters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// Service is representation of TaskCounterService
type Service struct {
	repo repository.TaskCounter
}

// New is a function to return new TaskCounter service
func New(repo repository.TaskCounter) Service {
	return Service{
		repo: repo,
	}
}

// GetCountersByPartner ...
func (tc Service) GetCountersByPartner(w http.ResponseWriter, r *http.Request) {
	partnerID := mux.Vars(r)["partnerID"]

	counts, err := tc.repo.GetCounters(r.Context(), partnerID, gocql.UUID{})
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantGetPartnersForCounters, "TaskCounterService.GetCountersByPartner got error: %v", err)
	}

	common.RenderJSON(w, counts)
}

// GetCountersByPartnerAndEndpoint ...
func (tc Service) GetCountersByPartnerAndEndpoint(w http.ResponseWriter, r *http.Request) {
	endpointID, err := common.ExtractUUID("TaskCounterService.GetCountersByPartnerAndEndpoint", w, r, "managedEndpointID")
	if err != nil {
		return
	}
	partnerID := mux.Vars(r)["partnerID"]

	counts, err := tc.repo.GetCounters(r.Context(), partnerID, endpointID)
	if err != nil {
		logger.Log.WarnfCtx(r.Context(), "TaskCounterService.GetCountersByPartnerAndEndpoint got error: %v", err)
	}

	if len(counts) == 0 {
		counts = []models.TaskCount{{
			ManagedEndpointID: endpointID,
			Count:             0,
		}}
	}

	common.RenderJSON(w, counts)
}

// RecalculateAllCounters recalculates all tasks for endpoints and partners on demand
func (tc Service) RecalculateAllCounters(w http.ResponseWriter, r *http.Request) {
	type Errors struct {
		HasErrors bool `json:"HasErrors"`
	}

	var (
		globalErr error
		ctx       = r.Context()
	)

	partners, err := tc.repo.GetAllPartners(ctx)
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantGetPartnersForCounters, "RecalculateAllCounters: GetAllPartners: %v", err)
		common.RenderJSON(w, Errors{
			HasErrors: true,
		})
		return
	}
	logger.Log.DebugfCtx(r.Context(), "RecalculateAllCounters: GetAllPartners: got partners [%v]", partners)

	for currentPartner := range partners {
		err = tc.processCounters(ctx, currentPartner)
		if err != nil {
			globalErr = fmt.Errorf("%v %v", globalErr, err)
		}
	}

	common.RenderJSON(w, Errors{
		HasErrors: globalErr != nil,
	})
}

func (tc Service) processCounters(ctx context.Context, currentPartner string) (err error) {
	var (
		endpointsForUpdate   = make(map[gocql.UUID]bool)         // for current partner
		countersForResetting = map[gocql.UUID]models.TaskCount{} // this counters should be set to 0
	)

	taskCountsFromTasksTable, err := models.TaskPersistenceInstance.GetCountsByPartner(ctx, currentPartner)
	if err != nil {
		logger.Log.InfofCtx(ctx, "TaskPersistenceInstance.GetCountsByPartner: %v", err)
		return
	}

	// gocql.UUID{} means that all counter for each ManagedEndpoint of current partner will be fetched form DB
	taskCountsFromTaskCounterTable, err := tc.repo.GetCounters(ctx, currentPartner, gocql.UUID{})
	if err != nil {
		logger.Log.InfofCtx(ctx, "Service.GetCounters: %v", err)
		return
	}

	// prepare set of TaskCounters by ManagedEndpointIDs
	for _, c := range taskCountsFromTaskCounterTable {
		endpointsForUpdate[c.ManagedEndpointID] = false // is not updated yet
		countersForResetting[c.ManagedEndpointID] = c
	}

	for _, counter := range taskCountsFromTasksTable {
		currentCount := counter.Count
		currentEndpoint := counter.ManagedEndpointID

		if countersForResetting[currentEndpoint].Count == currentCount {
			endpointsForUpdate[currentEndpoint] = true // marks it like updated 'cause they equal already
			return
		}
		logger.Log.InfofCtx(ctx, "RecalculateAllCounters: processing partnerID [%s] and endpointID [%v]", currentPartner, currentEndpoint)

		err = tc.updateCounters(ctx, currentPartner, currentEndpoint, currentCount)
		if err != nil {
			logger.Log.InfofCtx(ctx, "RecalculateAllCounters: updateCounters: %v", err)
			return
		}
		endpointsForUpdate[currentEndpoint] = true
	}

	resetCounters := make([]models.TaskCount, 0)
	// sets counter for each particular ManagedEndpoint to 0
	for endpoint, isUpdated := range endpointsForUpdate {
		if !isUpdated {
			logger.Log.InfofCtx(ctx, "Removing TaskCount for partnerID [%s] and endpointID [%s]\n", currentPartner, endpoint.String())

			resetCounters = append(resetCounters, models.TaskCount{
				ManagedEndpointID: endpoint,
				Count:             countersForResetting[endpoint].Count,
			})
		}
	}

	if len(resetCounters) > 0 {
		err = tc.repo.DecreaseCounter(currentPartner, resetCounters, false)
		if err != nil {
			logger.Log.InfofCtx(ctx, "got error while removing count: %v\n", err)
		}
	}
	return err
}

// updateCounters sets `correctCount` as actual one for current partner and endpoint
func (tc Service) updateCounters(ctx context.Context, partner string, endpoint gocql.UUID, correctCount int) (err error) {
	counts, err := tc.repo.GetCounters(ctx, partner, endpoint)
	if err != nil {
		counts = []models.TaskCount{{
			ManagedEndpointID: endpoint,
			Count:             0,
		}}
	}

	err = tc.repo.DecreaseCounter(partner, counts, false)
	if err != nil {
		return err
	}

	err = tc.repo.IncreaseCounter(partner, []models.TaskCount{{ManagedEndpointID: endpoint, Count: correctCount}}, false)
	if err != nil {
		return err
	}

	return
}
