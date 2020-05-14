package models

//go:generate mockgen -destination=../mocks/mocks-gomock/templateCache_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/models TemplateCache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/coocood/freecache"
	"github.com/gocql/gocql"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

const (
	scriptType = "script"
	windows    = "windows"
	linux      = "linux"

	globalPartnerID                 = "00000000-0000-0000-0000-000000000000"
	scriptNameKey                   = "TKS_TEMPLATE_NAME_%v"
	scriptNamePartnerKey            = "TKS_PARTNER_%v_TEMPLATE_NAME_%v"
	defaultExpectedExecutionTimeSec = 300
	expectedExecutionTimeField      = "expectedExecutionTimeSec"
	expirationTime                  = 24 * 60 * 60 // hours * minutes * seconds
)

// Template structure contains general data about task definition template
type Template struct {
	PartnerID          string     `json:"partnerId"`
	OriginID           gocql.UUID `json:"originId"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Type               string     `json:"type"`
	Categories         []string   `json:"categories"`
	Tags               []string   `json:"tags"`
	IsParameterized    bool       `json:"isParameterized"`
	IsRequireNOCAccess bool       `json:"isRequireNOCAccess"`
	Engine             string     `json:"engine"`
}

// TemplateDetails structure contains data about particular task definition template
type TemplateDetails struct {
	PartnerID                string     `json:"partnerId"`
	OriginID                 gocql.UUID `json:"originId"`
	Name                     string     `json:"name"`
	Description              string     `json:"description"`
	Type                     string     `json:"type"`
	Categories               []string   `json:"categories"`
	Tags                     []string   `json:"tags"`
	CreatedAt                time.Time  `json:"createdAt"`
	CreatedBy                gocql.UUID `json:"createdBy"`
	JSONSchema               string     `json:"JSONSchema"`
	UISchema                 string     `json:"UISchema"`
	IsParameterized          bool       `json:"isParameterized"`
	ExpectedExecutionTimeSec int        `json:"expectedExecutionTimeSec"`
	SuccessMessage           string     `json:"successMessage"`
	FailureMessage           string     `json:"failureMessage"`
	IsRequireNOCAccess       bool       `json:"isRequireNOCAccess"`
	Engine                   string     `json:"engine"`
}

// TemplateCache interface to perform actions with task_definition_template
type TemplateCache interface {
	GetAllTemplatesDetails(ctx context.Context, partnerID string) ([]TemplateDetails, error)
	GetAllTemplates(ctx context.Context, partnerID string, hasNOCAccess bool) ([]Template, error)
	GetByOriginID(ctx context.Context, partnerID string, templateID gocql.UUID, hasNOCAccess bool) (TemplateDetails, error)
	GetByType(ctx context.Context, partnerID string, taskType string, hasNOCAccess bool) ([]Template, error)
	ExistsWithName(ctx context.Context, partnerID, scriptName string) bool
	CalculateExpectedExecutionTimeSec(ctx context.Context, task Task) int
}

// TemplateNotFoundError returns in case of Template with particular ID was not found in the DB for particular partner
type TemplateNotFoundError struct {
	OriginID  gocql.UUID
	PartnerID string
}

// Error - error interface implementation for TemplateNotFoundError type
func (err TemplateNotFoundError) Error() string {
	return fmt.Sprintf(`Template with OriginID=%s not found for partner %s`, err.OriginID, err.PartnerID)
}

// TemplatesByTypeNotFoundError returns in case of Template with particular type was not found in the DB for particular partner
type TemplatesByTypeNotFoundError struct {
	Type      string
	PartnerID string
}

func (err TemplatesByTypeNotFoundError) Error() string {
	return fmt.Sprintf(`Templates with Type=%s are not found for partner %s`, err.Type, err.PartnerID)
}

// TemplateCacheLocal is an implementation of TemplateCache interface for Local Cache
type TemplateCacheLocal struct{}

var (
	// TemplateCacheInstance is an instance presented TemplateCacheLocal
	TemplateCacheInstance TemplateCache = TemplateCacheLocal{}
	// TemplatesCache Task Definition Templates cache
	TemplatesCache persistency.Cache
	// ModelsDecoder is a package wide json decoder
	ModelsDecoder = json.NewDecoder
	httpClient    = &http.Client{Timeout: time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second}
)

// LoadTemplatesCache init cache for Task Definition Templates.
func LoadTemplatesCache() {
	TemplatesCache = freecache.NewCache(config.Config.TDTCacheSettings.Size)
	debug.SetGCPercent(config.Config.TDTCacheSettings.GCPercent)
	ctx := transactionID.NewContext()

	if err := loadCacheData(ctx); err != nil {
		logger.Log.WarnfCtx(ctx, "Couldn't load templates to cache, err: %v", err)
	}

	go func(ctx context.Context) {
		for {
			time.Sleep(time.Duration(config.Config.TDTCacheSettings.ReloadIntervalSec) * time.Second)
			if err := loadCacheData(ctx); err != nil {
				logger.Log.WarnfCtx(ctx, "LoadTemplatesCache: error while loading cache data: ", err)
			}
		}
	}(ctx)
}

// ExistsWithName returns true if the template with this name exists
func (TemplateCacheLocal) ExistsWithName(ctx context.Context, partnerID, scriptName string) bool {
	// first checking global list
	key := fmt.Sprintf(scriptNameKey, scriptName)
	if _, err := TemplatesCache.Get([]byte(key)); err != nil {
		// then partner only
		key = fmt.Sprintf(scriptNamePartnerKey, partnerID, scriptName)
		if _, err := TemplatesCache.Get([]byte(key)); err != nil {
			return false
		}
	}
	return true
}

// GetAllTemplatesDetails returns all templates from the cache or from the Scripting MS
func (TemplateCacheLocal) GetAllTemplatesDetails(ctx context.Context, partnerID string) ([]TemplateDetails, error) {
	var emptyUUID gocql.UUID
	generalTemplates, genErr := getTemplatesDetailsFromCache([]byte(emptyUUID.String()))
	if genErr != nil {
		logger.Log.WarnfCtx(ctx, "Couldn't retrieve general templates, err: %v", genErr)
	}

	partnerTemplates, partnerErr := getTemplatesDetailsFromCache([]byte(partnerID))
	if partnerErr != nil {
		logger.Log.WarnfCtx(ctx, "Couldn't retrieve partner specific templates, err: %v", partnerErr)
	}

	templates := append(generalTemplates, partnerTemplates...)

	if len(templates) == 0 {
		return templates, fmt.Errorf("couldn't retrieve templates")
	}

	return templates, nil
}

// GetAllTemplates returns all templates from the cache or from the Scripting MS
func (TemplateCacheLocal) GetAllTemplates(ctx context.Context, partnerID string, hasNOCAccess bool) ([]Template, error) {
	var emptyUUID gocql.UUID
	generalTemplates, genErr := getTemplatesWithoutDetails([]byte(emptyUUID.String()), hasNOCAccess)
	if genErr != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantGetTemplatesForExecutionMS, "Couldn't retrieve general templates, err: %v", genErr)
	}

	partnerTemplates, partnerErr := getTemplatesWithoutDetails([]byte(partnerID), hasNOCAccess)
	if partnerErr != nil {
		logger.Log.WarnfCtx(ctx, "Couldn't retrieve partner specific templates, err: %v", partnerErr)
	}

	templates := append(generalTemplates, partnerTemplates...)

	if len(templates) == 0 {
		return templates, fmt.Errorf("couldn't retrieve templates")
	}

	return templates, nil
}

// GetByType returns all templates that been found by TaskType
func (TemplateCacheLocal) GetByType(ctx context.Context, partnerID, taskType string, hasNOCAccess bool) ([]Template, error) {
	var templatesByType []Template
	templates, err := TemplateCacheInstance.GetAllTemplates(ctx, partnerID, hasNOCAccess)
	if err != nil {
		return []Template{}, fmt.Errorf("No Task Definition Templates with Type: %v for Partner: %s, err: %v",
			taskType, partnerID, err)
	}

	for _, template := range templates {
		if template.IsRequireNOCAccess && !hasNOCAccess {
			// if current script only for NOC user but user has not access to it - skip
			continue
		}

		if template.Type == taskType {
			templatesByType = append(templatesByType, template)
		}
	}

	return templatesByType, nil
}

// CalculateExpectedExecutionTimeSec calculates expected execution time for a task execution
func (tcl TemplateCacheLocal) CalculateExpectedExecutionTimeSec(ctx context.Context, task Task) (eeTime int) {
	eeTime = defaultExpectedExecutionTimeSec

	schedule := task.Schedule
	if schedule.Regularity == apiModels.RunNow && !schedule.EndRunTime.IsZero() {
		eeTime = int(schedule.EndRunTime.UTC().Sub(time.Now().UTC()).Seconds())
		return
	}

	if strings.Contains(task.Parameters, expectedExecutionTimeField) {
		var customScript apiModels.CustomScript

		if err := json.Unmarshal([]byte(task.Parameters), &customScript); err != nil {
			logger.Log.WarnfCtx(ctx, "TemplateCacheLocal.CalculateExpectedExecutionTimeSec: invalid task.Parameters: %v", err)
			return
		}

		eeTime = customScript.ExpectedExecutionTimeSec + getExtraTime(task)
		return
	}

	templateDetails, err := tcl.GetByOriginID(ctx, task.PartnerID, task.OriginID, true)
	if err != nil {
		if task.Type == TaskTypeScript {
			logger.Log.WarnfCtx(ctx, "TemplateCacheLocal.CalculateExpectedExecutionTimeSec: can't get script by OriginID(%s), err=%v", task.OriginID, err)
		}
		templateDetails = TemplateDetails{ExpectedExecutionTimeSec: expirationTime}
	}

	eeTime = templateDetails.ExpectedExecutionTimeSec + getExtraTime(task)
	return
}

func getExtraTime(task Task) (extraTime int) {
	if !task.Schedule.BetweenEndTime.IsZero() && task.Schedule.Regularity != apiModels.RunNow {
		if !task.Schedule.Repeat.RunTime.IsZero() {
			return int(task.Schedule.BetweenEndTime.Sub(task.Schedule.Repeat.RunTime).Seconds())
		}
		return int(task.Schedule.BetweenEndTime.Sub(task.Schedule.StartRunTime).Seconds())
	}
	return 0
}

// GetByOriginID returns TemplateDetails that been found by Origin ID and partner ID
func (TemplateCacheLocal) GetByOriginID(ctx context.Context, partnerID string, originID gocql.UUID, hasNOCAccess bool) (TemplateDetails, error) {
	templates, err := TemplateCacheInstance.GetAllTemplatesDetails(ctx, partnerID)
	if err != nil {
		return TemplateDetails{}, fmt.Errorf("No Task Definition TemplateDetails with Origin ID: %v for Partner: %s, err: %v",
			originID.String(), partnerID, err)
	}

	for _, template := range templates {
		if template.IsRequireNOCAccess && !hasNOCAccess {
			continue
		}

		if template.OriginID == originID {
			return template, nil
		}
	}

	return TemplateDetails{}, TemplateNotFoundError{originID, partnerID}
}

// loadCacheData load templates from Scripting MS to cache by particular Partner ID and template type
func loadCacheData(ctx context.Context) error {
	executionMsURL := fmt.Sprintf("%s/scripts?internal=false", config.Config.ScriptingMsURL)
	req, err := http.NewRequest(http.MethodGet, executionMsURL, nil)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "Couldn't get Templates from URL %v, err: %v", executionMsURL, err)
		return err
	}
	req.Header.Set(transactionID.Key, transactionID.FromContext(ctx))

	httpResp, err := httpClient.Do(req)
	if err != nil {
		logger.Log.ErrfCtx(ctx, "Couldn't get Templates from URL %v, err: %v", executionMsURL, err)
		return err
	}
	defer func() {
		if err := httpResp.Body.Close(); err != nil {
			logger.Log.WarnfCtx(ctx, "loadCacheData: error while closing body: %v", err)
		}
	}()

	var scripts []apiModels.Script
	err = ModelsDecoder(httpResp.Body).Decode(&scripts)
	if err != nil {
		logger.Log.WarnfCtx(ctx, "Couldn't parse Templates from URL %v, err: %v", executionMsURL, err)
		return err
	}

	result := make(map[string][]TemplateDetails)
	for _, script := range scripts {
		processScript(ctx, script, result)
	}

	var counter int
	for partnerID, templates := range result {
		counter = setScriptToCache(ctx, templates, partnerID, counter)
	}
	logger.Log.DebugfCtx(ctx, "Successfully loaded %d templates", counter)

	return nil
}

func setScriptToCache(ctx context.Context, templates []TemplateDetails, partnerID string, counter int) int {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(templates)
	if err != nil {
		logger.Log.ErrfCtx(ctx, errorcode.ErrorCantDecodeInputData, "Couldn't encode Templates %v, err: %v", templates, err)
		return counter
	}

	if err = TemplatesCache.Set([]byte(partnerID), buffer.Bytes(), 0); err != nil {
		logger.Log.WarnfCtx(ctx, "Couldn't set templates %v for partner %s, err: %v", templates, partnerID, err)
		return counter
	}

	for _, template := range templates {
		if err = TemplatesCache.Set([]byte(fmt.Sprintf(scriptNamePartnerKey, template.PartnerID, template.Name)), buffer.Bytes(), 0); err != nil {
			logger.Log.WarnfCtx(ctx, "Couldn't set templates %v for partner %s, err: %v", templates, partnerID, err)
			return counter
		}
	}

	return counter + len(templates)
}

func processScript(ctx context.Context, script apiModels.Script, result map[string][]TemplateDetails) {
	uuid, parseErr := gocql.ParseUUID(script.ID)
	if parseErr != nil {
		logger.Log.WarnfCtx(ctx, "Couldn't parse Template UUID %v, err: %v", script.ID, parseErr)
		return
	}

	if script.Engine == Bash {
		script.Engine = linux
	}
	if script.Engine == CMD || script.Engine == Powershell {
		script.Engine = windows
	}

	tdt := TemplateDetails{
		OriginID:                 uuid,
		Name:                     script.Name,
		Description:              script.Description,
		PartnerID:                script.PartnerID,
		Type:                     scriptType,
		Categories:               script.Categories,
		Tags:                     script.Tags,
		CreatedAt:                script.CreatedAt,
		JSONSchema:               script.JSONSchema,
		UISchema:                 script.UISchema,
		IsParameterized:          script.JSONSchema != "",
		ExpectedExecutionTimeSec: script.ExpectedExecutionTimeSec,
		SuccessMessage:           script.SuccessMessage,
		FailureMessage:           script.FailureMessage,
		IsRequireNOCAccess:       script.NOCVisibleOnly,
		Engine:                   script.Engine,
	}

	if templates, ok := result[script.PartnerID]; ok {
		result[script.PartnerID] = append(templates, tdt)
	} else {
		result[script.PartnerID] = []TemplateDetails{tdt}
	}

	if script.PartnerID == globalPartnerID {
		if err := TemplatesCache.Set([]byte(fmt.Sprintf(scriptNameKey, script.Name)), []byte("true"), 0); err != nil {
			logger.Log.WarnfCtx(ctx, "Couldn't set template names %v, err: %v", script.Name, err)
		}
	}
}

func getTemplatesDetailsFromCache(key []byte) ([]TemplateDetails, error) {
	var foundTemplates []TemplateDetails

	foundTemplatesBin, err := TemplatesCache.Get(key)
	if err != nil {
		return foundTemplates, fmt.Errorf("can't find templates in cache by key %s", string(key))
	}

	err = ModelsDecoder(bytes.NewReader(foundTemplatesBin)).Decode(&foundTemplates)
	if err != nil {
		return foundTemplates, fmt.Errorf("can't decode templates from cache: %v", err)
	}

	return foundTemplates, nil
}

func getTemplatesWithoutDetails(key []byte, hasNOCAccess bool) ([]Template, error) {
	templatesFull, err := getTemplatesDetailsFromCache(key)
	if err != nil {
		return nil, err
	}

	result := make([]Template, 0, len(templatesFull))
	for _, template := range templatesFull {
		if template.IsRequireNOCAccess && !hasNOCAccess {
			// if current script only for NOC user but user has not access to it - skip
			continue
		}
		result = append(result, Template{
			PartnerID:       template.PartnerID,
			OriginID:        template.OriginID,
			Name:            template.Name,
			Description:     template.Description,
			Type:            template.Type,
			Categories:      template.Categories,
			Tags:            template.Tags,
			IsParameterized: template.IsParameterized,
			Engine:          template.Engine,
		})
	}
	return result, nil
}
