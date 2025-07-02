package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Tool represents an MCP tool
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
	Dangerous   bool        `json:"dangerous,omitempty"`
}

// MCPServer represents the MCP server
type MCPServer struct {
	clients *Clients
}

// handleListTools returns the list of available tools
func (s *MCPServer) handleListTools(w http.ResponseWriter, r *http.Request) {
	tools := []Tool{
		{
			Name:        "albisyncservice_triggerprojectsync",
			Description: "TriggerProjectSync operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_triggercollectionsync",
			Description: "TriggerCollectionSync operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_triggerfullsync",
			Description: "TriggerFullSync operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_getsyncstatus",
			Description: "GetSyncStatus operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_cancelsyncjob",
			Description: "CancelSyncJob operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_getsynchistory",
			Description: "GetSyncHistory operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_getsyncconfig",
			Description: "GetSyncConfig operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_updatesyncconfig",
			Description: "UpdateSyncConfig operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_testconnection",
			Description: "TestConnection operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_getfieldmappings",
			Description: "GetFieldMappings operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_updatefieldmappings",
			Description: "UpdateFieldMappings operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_processwebhook",
			Description: "ProcessWebhook operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_getwebhooksecret",
			Description: "GetWebhookSecret operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "albisyncservice_rotatewebhooksecret",
			Description: "RotateWebhookSecret operation from ALBISyncService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_login",
			Description: "Login operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_refreshtoken",
			Description: "RefreshToken operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_validatetoken",
			Description: "ValidateToken operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_changepassword",
			Description: "ChangePassword operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_requestpasswordreset",
			Description: "RequestPasswordReset operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_resetpassword",
			Description: "ResetPassword operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_getactivesessions",
			Description: "GetActiveSessions operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_revokeallsessions",
			Description: "RevokeAllSessions operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_createapikey",
			Description: "CreateAPIKey operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "authservice_listapikeys",
			Description: "ListAPIKeys operation from AuthService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_createuser",
			Description: "CreateUser operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_getuser",
			Description: "GetUser operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_updateuser",
			Description: "UpdateUser operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_listusers",
			Description: "ListUsers operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_createproject",
			Description: "CreateProject operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_getproject",
			Description: "GetProject operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_updateproject",
			Description: "UpdateProject operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_listprojects",
			Description: "ListProjects operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_getprojectbyreferencenumber",
			Description: "GetProjectByReferenceNumber operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_createcollection",
			Description: "CreateCollection operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_getcollection",
			Description: "GetCollection operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_updatecollection",
			Description: "UpdateCollection operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_listcollections",
			Description: "ListCollections operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_listcollectionsbyproject",
			Description: "ListCollectionsByProject operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_createcalltask",
			Description: "CreateCallTask operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_getcalltask",
			Description: "GetCallTask operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_updatecalltask",
			Description: "UpdateCallTask operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_listcalltasks",
			Description: "ListCallTasks operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_listcalltasksbyassignee",
			Description: "ListCallTasksByAssignee operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_getalbisyncjob",
			Description: "GetALBISyncJob operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_listalbisyncjobs",
			Description: "ListALBISyncJobs operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_createalbisyncjob",
			Description: "CreateALBISyncJob operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_updatealbisyncjob",
			Description: "UpdateALBISyncJob operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_createalbiwebhook",
			Description: "CreateALBIWebhook operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_updatealbiwebhook",
			Description: "UpdateALBIWebhook operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "databaseservice_listalbiwebhooks",
			Description: "ListALBIWebhooks operation from DatabaseService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_getfilemetadata",
			Description: "GetFileMetadata operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_updatefilemetadata",
			Description: "UpdateFileMetadata operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_listfiles",
			Description: "ListFiles operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_generateuploadurl",
			Description: "GenerateUploadURL operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_generatedownloadurl",
			Description: "GenerateDownloadURL operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_copyfile",
			Description: "CopyFile operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_movefile",
			Description: "MoveFile operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_createfolder",
			Description: "CreateFolder operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_getstorageusage",
			Description: "GetStorageUsage operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_scanfile",
			Description: "ScanFile operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "filestorageservice_generatethumbnail",
			Description: "GenerateThumbnail operation from FileStorageService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "notificationservice_sendemail",
			Description: "SendEmail operation from NotificationService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "notificationservice_sendsms",
			Description: "SendSMS operation from NotificationService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "notificationservice_sendpushnotification",
			Description: "SendPushNotification operation from NotificationService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "notificationservice_getnotificationpreferences",
			Description: "GetNotificationPreferences operation from NotificationService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "notificationservice_updatenotificationpreferences",
			Description: "UpdateNotificationPreferences operation from NotificationService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "notificationservice_getnotificationhistory",
			Description: "GetNotificationHistory operation from NotificationService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "notificationservice_sendbulknotification",
			Description: "SendBulkNotification operation from NotificationService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "notificationservice_getnotificationtemplates",
			Description: "GetNotificationTemplates operation from NotificationService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "notificationservice_upsertnotificationtemplate",
			Description: "UpsertNotificationTemplate operation from NotificationService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_search",
			Description: "Search operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_getsuggestions",
			Description: "GetSuggestions operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_indexdocument",
			Description: "IndexDocument operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_updatedocument",
			Description: "UpdateDocument operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_deletedocument",
			Description: "DeleteDocument operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   true,
		},
		{
			Name:        "searchservice_bulkindex",
			Description: "BulkIndex operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_getfacets",
			Description: "GetFacets operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_savesearch",
			Description: "SaveSearch operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_getsavedsearches",
			Description: "GetSavedSearches operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_deletesavedsearch",
			Description: "DeleteSavedSearch operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   true,
		},
		{
			Name:        "searchservice_getsearchanalytics",
			Description: "GetSearchAnalytics operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "searchservice_reindexdocuments",
			Description: "ReindexDocuments operation from SearchService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_createworkflowdefinition",
			Description: "CreateWorkflowDefinition operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_updateworkflowdefinition",
			Description: "UpdateWorkflowDefinition operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_getworkflowdefinition",
			Description: "GetWorkflowDefinition operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_listworkflowdefinitions",
			Description: "ListWorkflowDefinitions operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_startworkflow",
			Description: "StartWorkflow operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_getworkflowinstance",
			Description: "GetWorkflowInstance operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_listworkflowinstances",
			Description: "ListWorkflowInstances operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_retryworkflow",
			Description: "RetryWorkflow operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_gettask",
			Description: "GetTask operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_listtasks",
			Description: "ListTasks operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_completetask",
			Description: "CompleteTask operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_reassigntask",
			Description: "ReassignTask operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_addtaskcomment",
			Description: "AddTaskComment operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_transitionworkflowstate",
			Description: "TransitionWorkflowState operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_getavailabletransitions",
			Description: "GetAvailableTransitions operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_getworkflowhistory",
			Description: "GetWorkflowHistory operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "workflowservice_getworkflowmetrics",
			Description: "GetWorkflowMetrics operation from WorkflowService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tools": tools,
	})
}

// handleRPC handles MCP RPC calls
func (s *MCPServer) handleRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Method string                 `json:"method"`
		Params map[string]interface{} `json:"params"`
		Auth   map[string]interface{} `json:"auth"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Authenticate request
	if err := s.authenticate(req.Auth); err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
		return
	}

	// Route to appropriate handler
	var result interface{}
	var err error

	switch req.Method {
	case "albisyncservice_triggerprojectsync":
		result, err = s.handleALBISyncServiceTriggerProjectSync(r.Context(), req.Params)
	case "albisyncservice_triggercollectionsync":
		result, err = s.handleALBISyncServiceTriggerCollectionSync(r.Context(), req.Params)
	case "albisyncservice_triggerfullsync":
		result, err = s.handleALBISyncServiceTriggerFullSync(r.Context(), req.Params)
	case "albisyncservice_getsyncstatus":
		result, err = s.handleALBISyncServiceGetSyncStatus(r.Context(), req.Params)
	case "albisyncservice_cancelsyncjob":
		result, err = s.handleALBISyncServiceCancelSyncJob(r.Context(), req.Params)
	case "albisyncservice_getsynchistory":
		result, err = s.handleALBISyncServiceGetSyncHistory(r.Context(), req.Params)
	case "albisyncservice_getsyncconfig":
		result, err = s.handleALBISyncServiceGetSyncConfig(r.Context(), req.Params)
	case "albisyncservice_updatesyncconfig":
		result, err = s.handleALBISyncServiceUpdateSyncConfig(r.Context(), req.Params)
	case "albisyncservice_testconnection":
		result, err = s.handleALBISyncServiceTestConnection(r.Context(), req.Params)
	case "albisyncservice_getfieldmappings":
		result, err = s.handleALBISyncServiceGetFieldMappings(r.Context(), req.Params)
	case "albisyncservice_updatefieldmappings":
		result, err = s.handleALBISyncServiceUpdateFieldMappings(r.Context(), req.Params)
	case "albisyncservice_processwebhook":
		result, err = s.handleALBISyncServiceProcessWebhook(r.Context(), req.Params)
	case "albisyncservice_getwebhooksecret":
		result, err = s.handleALBISyncServiceGetWebhookSecret(r.Context(), req.Params)
	case "albisyncservice_rotatewebhooksecret":
		result, err = s.handleALBISyncServiceRotateWebhookSecret(r.Context(), req.Params)
	case "authservice_login":
		result, err = s.handleAuthServiceLogin(r.Context(), req.Params)
	case "authservice_refreshtoken":
		result, err = s.handleAuthServiceRefreshToken(r.Context(), req.Params)
	case "authservice_validatetoken":
		result, err = s.handleAuthServiceValidateToken(r.Context(), req.Params)
	case "authservice_changepassword":
		result, err = s.handleAuthServiceChangePassword(r.Context(), req.Params)
	case "authservice_requestpasswordreset":
		result, err = s.handleAuthServiceRequestPasswordReset(r.Context(), req.Params)
	case "authservice_resetpassword":
		result, err = s.handleAuthServiceResetPassword(r.Context(), req.Params)
	case "authservice_getactivesessions":
		result, err = s.handleAuthServiceGetActiveSessions(r.Context(), req.Params)
	case "authservice_revokeallsessions":
		result, err = s.handleAuthServiceRevokeAllSessions(r.Context(), req.Params)
	case "authservice_createapikey":
		result, err = s.handleAuthServiceCreateAPIKey(r.Context(), req.Params)
	case "authservice_listapikeys":
		result, err = s.handleAuthServiceListAPIKeys(r.Context(), req.Params)
	case "databaseservice_createuser":
		result, err = s.handleDatabaseServiceCreateUser(r.Context(), req.Params)
	case "databaseservice_getuser":
		result, err = s.handleDatabaseServiceGetUser(r.Context(), req.Params)
	case "databaseservice_updateuser":
		result, err = s.handleDatabaseServiceUpdateUser(r.Context(), req.Params)
	case "databaseservice_listusers":
		result, err = s.handleDatabaseServiceListUsersREAL(r.Context(), req.Params)
	case "databaseservice_createproject":
		result, err = s.handleDatabaseServiceCreateProject(r.Context(), req.Params)
	case "databaseservice_getproject":
		result, err = s.handleDatabaseServiceGetProject(r.Context(), req.Params)
	case "databaseservice_updateproject":
		result, err = s.handleDatabaseServiceUpdateProject(r.Context(), req.Params)
	case "databaseservice_listprojects":
		result, err = s.handleDatabaseServiceListProjectsREAL(r.Context(), req.Params)
	case "databaseservice_getprojectbyreferencenumber":
		result, err = s.handleDatabaseServiceGetProjectByReferenceNumber(r.Context(), req.Params)
	case "databaseservice_createcollection":
		result, err = s.handleDatabaseServiceCreateCollection(r.Context(), req.Params)
	case "databaseservice_getcollection":
		result, err = s.handleDatabaseServiceGetCollection(r.Context(), req.Params)
	case "databaseservice_updatecollection":
		result, err = s.handleDatabaseServiceUpdateCollection(r.Context(), req.Params)
	case "databaseservice_listcollections":
		result, err = s.handleDatabaseServiceListCollections(r.Context(), req.Params)
	case "databaseservice_listcollectionsbyproject":
		result, err = s.handleDatabaseServiceListCollectionsByProject(r.Context(), req.Params)
	case "databaseservice_createcalltask":
		result, err = s.handleDatabaseServiceCreateCallTask(r.Context(), req.Params)
	case "databaseservice_getcalltask":
		result, err = s.handleDatabaseServiceGetCallTask(r.Context(), req.Params)
	case "databaseservice_updatecalltask":
		result, err = s.handleDatabaseServiceUpdateCallTask(r.Context(), req.Params)
	case "databaseservice_listcalltasks":
		result, err = s.handleDatabaseServiceListCallTasks(r.Context(), req.Params)
	case "databaseservice_listcalltasksbyassignee":
		result, err = s.handleDatabaseServiceListCallTasksByAssignee(r.Context(), req.Params)
	case "databaseservice_getalbisyncjob":
		result, err = s.handleDatabaseServiceGetALBISyncJob(r.Context(), req.Params)
	case "databaseservice_listalbisyncjobs":
		result, err = s.handleDatabaseServiceListALBISyncJobs(r.Context(), req.Params)
	case "databaseservice_createalbisyncjob":
		result, err = s.handleDatabaseServiceCreateALBISyncJob(r.Context(), req.Params)
	case "databaseservice_updatealbisyncjob":
		result, err = s.handleDatabaseServiceUpdateALBISyncJob(r.Context(), req.Params)
	case "databaseservice_createalbiwebhook":
		result, err = s.handleDatabaseServiceCreateALBIWebhook(r.Context(), req.Params)
	case "databaseservice_updatealbiwebhook":
		result, err = s.handleDatabaseServiceUpdateALBIWebhook(r.Context(), req.Params)
	case "databaseservice_listalbiwebhooks":
		result, err = s.handleDatabaseServiceListALBIWebhooks(r.Context(), req.Params)
	case "filestorageservice_getfilemetadata":
		result, err = s.handleFileStorageServiceGetFileMetadata(r.Context(), req.Params)
	case "filestorageservice_updatefilemetadata":
		result, err = s.handleFileStorageServiceUpdateFileMetadata(r.Context(), req.Params)
	case "filestorageservice_listfiles":
		result, err = s.handleFileStorageServiceListFiles(r.Context(), req.Params)
	case "filestorageservice_generateuploadurl":
		result, err = s.handleFileStorageServiceGenerateUploadURL(r.Context(), req.Params)
	case "filestorageservice_generatedownloadurl":
		result, err = s.handleFileStorageServiceGenerateDownloadURL(r.Context(), req.Params)
	case "filestorageservice_copyfile":
		result, err = s.handleFileStorageServiceCopyFile(r.Context(), req.Params)
	case "filestorageservice_movefile":
		result, err = s.handleFileStorageServiceMoveFile(r.Context(), req.Params)
	case "filestorageservice_createfolder":
		result, err = s.handleFileStorageServiceCreateFolder(r.Context(), req.Params)
	case "filestorageservice_getstorageusage":
		result, err = s.handleFileStorageServiceGetStorageUsage(r.Context(), req.Params)
	case "filestorageservice_scanfile":
		result, err = s.handleFileStorageServiceScanFile(r.Context(), req.Params)
	case "filestorageservice_generatethumbnail":
		result, err = s.handleFileStorageServiceGenerateThumbnail(r.Context(), req.Params)
	case "notificationservice_sendemail":
		result, err = s.handleNotificationServiceSendEmail(r.Context(), req.Params)
	case "notificationservice_sendsms":
		result, err = s.handleNotificationServiceSendSMS(r.Context(), req.Params)
	case "notificationservice_sendpushnotification":
		result, err = s.handleNotificationServiceSendPushNotification(r.Context(), req.Params)
	case "notificationservice_getnotificationpreferences":
		result, err = s.handleNotificationServiceGetNotificationPreferences(r.Context(), req.Params)
	case "notificationservice_updatenotificationpreferences":
		result, err = s.handleNotificationServiceUpdateNotificationPreferences(r.Context(), req.Params)
	case "notificationservice_getnotificationhistory":
		result, err = s.handleNotificationServiceGetNotificationHistory(r.Context(), req.Params)
	case "notificationservice_sendbulknotification":
		result, err = s.handleNotificationServiceSendBulkNotification(r.Context(), req.Params)
	case "notificationservice_getnotificationtemplates":
		result, err = s.handleNotificationServiceGetNotificationTemplates(r.Context(), req.Params)
	case "notificationservice_upsertnotificationtemplate":
		result, err = s.handleNotificationServiceUpsertNotificationTemplate(r.Context(), req.Params)
	case "searchservice_search":
		result, err = s.handleSearchServiceSearch(r.Context(), req.Params)
	case "searchservice_getsuggestions":
		result, err = s.handleSearchServiceGetSuggestions(r.Context(), req.Params)
	case "searchservice_indexdocument":
		result, err = s.handleSearchServiceIndexDocument(r.Context(), req.Params)
	case "searchservice_updatedocument":
		result, err = s.handleSearchServiceUpdateDocument(r.Context(), req.Params)
	case "searchservice_deletedocument":
		result, err = s.handleSearchServiceDeleteDocument(r.Context(), req.Params)
	case "searchservice_bulkindex":
		result, err = s.handleSearchServiceBulkIndex(r.Context(), req.Params)
	case "searchservice_getfacets":
		result, err = s.handleSearchServiceGetFacets(r.Context(), req.Params)
	case "searchservice_savesearch":
		result, err = s.handleSearchServiceSaveSearch(r.Context(), req.Params)
	case "searchservice_getsavedsearches":
		result, err = s.handleSearchServiceGetSavedSearches(r.Context(), req.Params)
	case "searchservice_deletesavedsearch":
		result, err = s.handleSearchServiceDeleteSavedSearch(r.Context(), req.Params)
	case "searchservice_getsearchanalytics":
		result, err = s.handleSearchServiceGetSearchAnalytics(r.Context(), req.Params)
	case "searchservice_reindexdocuments":
		result, err = s.handleSearchServiceReindexDocuments(r.Context(), req.Params)
	case "workflowservice_createworkflowdefinition":
		result, err = s.handleWorkflowServiceCreateWorkflowDefinition(r.Context(), req.Params)
	case "workflowservice_updateworkflowdefinition":
		result, err = s.handleWorkflowServiceUpdateWorkflowDefinition(r.Context(), req.Params)
	case "workflowservice_getworkflowdefinition":
		result, err = s.handleWorkflowServiceGetWorkflowDefinition(r.Context(), req.Params)
	case "workflowservice_listworkflowdefinitions":
		result, err = s.handleWorkflowServiceListWorkflowDefinitions(r.Context(), req.Params)
	case "workflowservice_startworkflow":
		result, err = s.handleWorkflowServiceStartWorkflow(r.Context(), req.Params)
	case "workflowservice_getworkflowinstance":
		result, err = s.handleWorkflowServiceGetWorkflowInstance(r.Context(), req.Params)
	case "workflowservice_listworkflowinstances":
		result, err = s.handleWorkflowServiceListWorkflowInstances(r.Context(), req.Params)
	case "workflowservice_retryworkflow":
		result, err = s.handleWorkflowServiceRetryWorkflow(r.Context(), req.Params)
	case "workflowservice_gettask":
		result, err = s.handleWorkflowServiceGetTask(r.Context(), req.Params)
	case "workflowservice_listtasks":
		result, err = s.handleWorkflowServiceListTasks(r.Context(), req.Params)
	case "workflowservice_completetask":
		result, err = s.handleWorkflowServiceCompleteTask(r.Context(), req.Params)
	case "workflowservice_reassigntask":
		result, err = s.handleWorkflowServiceReassignTask(r.Context(), req.Params)
	case "workflowservice_addtaskcomment":
		result, err = s.handleWorkflowServiceAddTaskComment(r.Context(), req.Params)
	case "workflowservice_transitionworkflowstate":
		result, err = s.handleWorkflowServiceTransitionWorkflowState(r.Context(), req.Params)
	case "workflowservice_getavailabletransitions":
		result, err = s.handleWorkflowServiceGetAvailableTransitions(r.Context(), req.Params)
	case "workflowservice_getworkflowhistory":
		result, err = s.handleWorkflowServiceGetWorkflowHistory(r.Context(), req.Params)
	case "workflowservice_getworkflowmetrics":
		result, err = s.handleWorkflowServiceGetWorkflowMetrics(r.Context(), req.Params)
	default:
		http.Error(w, "Unknown method", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}
