// package controllerscriteria contains data structures and database access logic.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/criteria
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:09/01/2026
//Last Modify by : Ramya
// Last Modify On :12/02/2026
// package controllerscriteria

// import (
// 	"Hrmodule/auth"
// 	credentials "Hrmodule/dbconfig"
// 	"Hrmodule/micro"
// 	"Hrmodule/utils"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"strings"

// 	"github.com/google/uuid"
// 	_ "github.com/lib/pq"
// )

// // APICriteriainsertResponse defines the encrypted API response format
// // sent back to the client after insert/update
// type APICriteriainsertResponse struct {
// 	Status                int       `json:"status"`
// 	Message               string    `json:"message"`
// 	Response_task_id      uuid.UUID `json:"response_task_id"`
// 	Response_reference_no string    `json:"response_reference_no"`
// }

// // EncryptedRequest wraps encrypted payload sent from client
// // Format: PID||EncryptedData
// type EncryptedRequest struct {
// 	Data string `json:"Data"`
// }

// // CPCCriteria represents each criteria entry with CPC mapping
// type CPCCriteria struct {
// 	CriteriaID  string `json:"criteria_id"`
// 	Description string `json:"description"`
// 	CPCIDs      []int  `json:"cpc_ids"`
// }

// // CriteriaInsertRequest represents the decrypted request payload
// // received from client
// type CriteriaInsertRequest struct {

// 	// TASK
// 	TaskID       *uuid.UUID `json:"task_id"`
// 	ReferenceNO  *string    `json:"reference_no"`
// 	TypeOfSubmit string     `json:"typeofsubmit"`

// 	// MASTER
// 	ProcessID    int           `json:"process_id"`
// 	CriteriaData []CPCCriteria `json:"criteria_data"`
// 	Status       int           `json:"status"`

// 	AssignTo     string `json:"assign_to"`
// 	AssignedRole string `json:"assigned_role"`

// 	// ACTIVITY
// 	TaskStatusID         int `json:"task_status_id"`
// 	ActivitySeqNo        int `json:"activity_seq_no"`
// 	CurrentActivitySeqNo int `json:"p_current_activity_seq_no"`
// 	IsTaskReturn         int `json:"p_is_task_return"`
// 	IsTaskApproved       int `json:"p_is_task_approved"`

// 	Role           *string           `json:"p_role"`
// 	RequestedUser  *string           `json:"p_requested_user"`
// 	ReturnToRole   *string           `json:"p_return_to_role"`
// 	ReturnToUser   *string           `json:"p_return_to_user"`
// 	SendBackToMe   *bool             `json:"p_send_back_to_me"`
// 	SendBackToUser *string           `json:"p_send_back_to_user"`
// 	Conditions     map[string]string `json:"p_conditions"`

// 	// FLAGS
// 	EmailFlag  int     `json:"email_flag"`
// 	TemplateID NullInt `json:"template_id"`
// 	RejectFlag int     `json:"reject_flag"`
// 	RejectRole *string `json:"reject_role"`
// 	StatusFlag int     `json:"status_flag"`
// 	Badge      int     `json:"badge"`
// 	Priority   int     `json:"priority"`
// 	Starred    int     `json:"starred"`

// 	// COMMENTS
// 	UserID    string `json:"user_id"`
// 	UserRole  string `json:"user_role"`
// 	Remarks   string `json:"remarks"`
// 	UpdatedBy string `json:"updated_by"`

// 	// TOKEN
// 	Token string `json:"token"`
// 	PID   string `json:"P_id"`
// }

// // validate checks mandatory fields before processing request
// func (req *CriteriaInsertRequest) validate() error {
// 	if req.ProcessID == 0 {
// 		return fmt.Errorf("process_id is required")
// 	}
// 	if req.AssignTo == "" {
// 		return fmt.Errorf("assign_to is required")
// 	}
// 	if req.UserID == "" {
// 		return fmt.Errorf("user_id is required")
// 	}
// 	return nil
// }

// // sendErrorResponse sends encrypted error response to client
// func sendErrorResponse(w http.ResponseWriter, r *http.Request, pid, key string, statusCode int, message string) {
// 	response := APICriteriainsertResponse{
// 		Status:  statusCode,
// 		Message: message,
// 	}
// 	encryptAndRespond(w, r, pid, key, statusCode, response)
// }

// // encryptAndRespond encrypts payload and sends JSON response
// // Format: { "Data": "PID||EncryptedPayload" }
// func encryptAndRespond(w http.ResponseWriter, r *http.Request, pid, key string, statusCode int, payload interface{}) {
// 	responseJSON, err := json.MarshalIndent(payload, "", "  ")
// 	if err != nil {
// 		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
// 		return
// 	}

// 	encryptedResponse := string(responseJSON)
// 	if key != "" && pid != "" {
// 		if enc, err := utils.EncryptAES(string(responseJSON), key); err == nil {
// 			encryptedResponse = fmt.Sprintf("%s||%s", pid, enc)
// 		}
// 	}

// 	// Wrapper
// 	finalResp := map[string]string{"Data": encryptedResponse}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	json.NewEncoder(w).Encode(finalResp)
// }

// // NullInt wraps sql.NullInt64 for optional integer fields
// type NullInt struct {
// 	sql.NullInt64
// }

// // stringPtrOrNil returns nil if string is empty
// // Used for optional workflow parameters
// func stringPtrOrNil(s *string) *string {
// 	if s == nil || *s == "" {
// 		return nil
// 	}
// 	return s
// }

// // CriteriaInsert handles Criteria create / update API
// // Supports draft & submit flows with workflow integration
// func CriteriaInsert(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	//  Read body
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Unable to read body", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	//  Encrypted wrapper
// 	var encReq EncryptedRequest
// 	if err := json.Unmarshal(body, &encReq); err != nil {
// 		http.Error(w, "Invalid request format", http.StatusBadRequest)
// 		return
// 	}
// 	fmt.Println()
// 	// Split PID || encryptedData
// 	parts := strings.Split(encReq.Data, "||")
// 	if len(parts) != 2 {
// 		http.Error(w, "Invalid Data format", http.StatusBadRequest)
// 		return
// 	}

// 	pid := parts[0]
// 	encryptedPart := parts[1]

// 	// Get key
// 	key, err := utils.GetDecryptKey(pid)
// 	if err != nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Decrypt
// 	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
// 	if err != nil {
// 		http.Error(w, "Decryption failed", http.StatusUnauthorized)
// 		return
// 	}

// 	// Unmarshal request
// 	var req CriteriaInsertRequest
// 	if err := json.Unmarshal([]byte(decryptedJSON), &req); err != nil {
// 		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
// 		return
// 	}

// 	// Validate
// 	if err := req.validate(); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Token & authorization validation
// 	if req.Token != "" {
// 		r.Header.Set("token", req.Token)
// 	}
// 	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
// 		return
// 	}

// 	// =====================================================
// 	// 6 MAIN PROCESS
// 	// =====================================================
// 	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		if err := auth.IsValidIDFromRequest(r); err != nil {
// 			sendErrorResponse(w, r, pid, key, http.StatusBadRequest, "Invalid TOKEN provided")
// 			return
// 		}

// 		isSubmit := strings.EqualFold(req.TypeOfSubmit, "submit")
// 		var taskIDStr string

// 		if req.TaskID != nil {
// 			taskIDStr = req.TaskID.String()
// 		}

// 		taskIDStr = ""
// 		if req.TaskID != nil {
// 			taskIDStr = req.TaskID.String()
// 		}

// 		// If submit, call workflow microservice to fetch next activity
// 		if isSubmit {
// 			sendBackToMe := false
// 			if req.SendBackToMe != nil {
// 				sendBackToMe = *req.SendBackToMe
// 			}

// 			currentActivity := req.CurrentActivitySeqNo
// 			if currentActivity <= 0 {
// 				currentActivity = req.ActivitySeqNo
// 			}

// 			input := micro.NextActivityInput{
// 				Token:           req.Token,
// 				ProcessID:       req.ProcessID,
// 				CurrentActivity: currentActivity,
// 				TaskID:          stringPtrOrNil(&taskIDStr),
// 				IsApproved:      req.IsTaskApproved,
// 				IsReturn:        req.IsTaskReturn,
// 				Role:            stringPtrOrNil(req.Role),
// 				RequestedUser:   stringPtrOrNil(req.RequestedUser),
// 				ReturnToRole:    stringPtrOrNil(req.ReturnToRole),
// 				ReturnToUser:    stringPtrOrNil(req.ReturnToUser),
// 				SendBackToMe:    sendBackToMe,
// 				SendBackToUser:  stringPtrOrNil(req.SendBackToUser),
// 				Conditions: map[string]string{
// 					"choice":        "",
// 					"EmployeeGroup": "",
// 				},
// 			}

// 			out, err := micro.NextActivity(input, pid, key)
// 			if err != nil {
// 				sendErrorResponse(w, r, pid, key, http.StatusInternalServerError, err.Error())
// 				return
// 			}

// 			// Update request object with workflow response
// 			req.ActivitySeqNo = out.NextActivitySeq
// 			req.AssignedRole = out.AssignedRole
// 			req.AssignTo = out.AssignTo
// 			req.EmailFlag = out.EmailFlag
// 			req.RejectFlag = out.RejectFlag

// 			if out.TemplateID != nil {
// 				req.TemplateID = NullInt{sql.NullInt64{
// 					Int64: int64(*out.TemplateID),
// 					Valid: true,
// 				}}
// 			}
// 		}

// 		//Db function call for Insert/Update
// 		taskID, refNo, err := executeCriteriaInsert(&req)
// 		if err != nil {
// 			sendErrorResponse(w, r, pid, key, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		// Trigger email queue only on submit with valid template
// 		if isSubmit && req.EmailFlag > 0 && req.TemplateID.Valid {
// 			if err := micro.EmailQueue(
// 				req.Token,
// 				req.ProcessID,
// 				taskIDStr,
// 				int(req.TemplateID.Int64),
// 				pid,
// 				key,
// 			); err != nil {
// 				fmt.Printf("EmailQueue failed: %v\n", err)
// 			}
// 		}

// 		// Return success response based on draft / submit
// 		var resp APICriteriainsertResponse
// 		if req.TypeOfSubmit == "draft" {
// 			resp = APICriteriainsertResponse{
// 				Status:                200,
// 				Message:               "Draft Added Successfully",
// 				Response_task_id:      taskID,
// 				Response_reference_no: refNo,
// 			}
// 		} else {
// 			resp = APICriteriainsertResponse{
// 				Status:                200,
// 				Message:               "Criteria Insert Successfully",
// 				Response_task_id:      taskID,
// 				Response_reference_no: refNo,
// 			}
// 		}
// 		encryptAndRespond(w, r, pid, key, http.StatusOK, resp)
// 	})).ServeHTTP(w, r)
// }

// // executeCriteriaInsert executes PostgreSQL stored procedure
// // and retrieves OUT parameters (task_id, reference_no)
// func executeCriteriaInsert(req *CriteriaInsertRequest) (uuid.UUID, string, error) {

// 	// Database connection
// 	db := credentials.GetDB()

// 	// Convert CriteriaData slice into JSONB
// 	criteriaJSON, err := json.Marshal(req.CriteriaData)
// 	if err != nil {
// 		return uuid.Nil, "", err
// 	}

// 	var taskID uuid.UUID
// 	var referenceNo string

// 	// CALL stored procedure using QueryRow + Scan
// 	err = db.QueryRow(`
// 		CALL meivan.cmes_qms_full(
// 			$1,$2,$3,
// 			$4,$5,$6,$7,$8,
// 			$9,$10,$11,$12,$13,
// 			$14,$15,$16,$17,$18,
// 			$19,$20,$21,$22,$23,$24,$25
// 		)
// 	`,
// 		req.TaskID,      // INOUT uuid
// 		req.ReferenceNO, // INOUT text
// 		req.ProcessID,
// 		criteriaJSON,
// 		req.TypeOfSubmit,
// 		req.Status,
// 		req.AssignTo,
// 		req.AssignedRole,
// 		req.TaskStatusID,
// 		req.ActivitySeqNo,
// 		req.CurrentActivitySeqNo,
// 		req.IsTaskReturn,
// 		req.IsTaskApproved,
// 		req.EmailFlag,
// 		req.TemplateID,
// 		req.RejectFlag,
// 		req.RejectRole,
// 		req.StatusFlag,
// 		req.Badge,
// 		req.Priority,
// 		req.Starred,
// 		req.UserID,
// 		req.UserRole,
// 		req.Remarks,
// 		req.UpdatedBy,
// 	).Scan(&taskID, &referenceNo)

// 	if err != nil {
// 		return uuid.Nil, "", fmt.Errorf("SP execution failed: %w", err)
// 	}

// 	return taskID, referenceNo, nil
// }


package controllerscriteria

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/micro"
	"Hrmodule/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// APICriteriainsertResponse defines the encrypted API response format
// sent back to the client after insert/update
type APICriteriainsertResponse struct {
	Status                int       `json:"status"`
	Message               string    `json:"message"`
	Response_task_id      uuid.UUID `json:"response_task_id"`
	Response_reference_no string    `json:"response_reference_no"`
}

// EncryptedRequest wraps encrypted payload sent from client
// Format: PID||EncryptedData
type EncryptedRequest struct {
	Data string `json:"Data"`
}

// CPCCriteria represents each criteria entry with CPC mapping
type CPCCriteria struct {
	CriteriaID  string `json:"criteria_id"`
	Description string `json:"description"`
	CPCIDs      []int  `json:"cpc_ids"`
}

// CriteriaInsertRequest represents the decrypted request payload
// received from client
type CriteriaInsertRequest struct {

	// TASK
	TaskID       *uuid.UUID `json:"task_id"`
	ReferenceNO  *string    `json:"reference_no"`
	TypeOfSubmit string     `json:"typeofsubmit"`

	// MASTER
	ProcessID    int           `json:"process_id"`
	CriteriaData []CPCCriteria `json:"criteria_data"`
	Status       int           `json:"status"`

	AssignTo     string `json:"assign_to"`
	AssignedRole string `json:"assigned_role"`

	// ACTIVITY
	TaskStatusID         int `json:"task_status_id"`
	ActivitySeqNo        int `json:"activity_seq_no"`
	CurrentActivitySeqNo int `json:"p_current_activity_seq_no"`
	IsTaskReturn         int `json:"p_is_task_return"`
	IsTaskApproved       int `json:"p_is_task_approved"`

	Role           *string           `json:"p_role"`
	RequestedUser  *string           `json:"p_requested_user"`
	ReturnToRole   *string           `json:"p_return_to_role"`
	ReturnToUser   *string           `json:"p_return_to_user"`
	SendBackToMe   *bool             `json:"p_send_back_to_me"`
	SendBackToUser *string           `json:"p_send_back_to_user"`
	Conditions     map[string]string `json:"p_conditions"`

	// FLAGS
	EmailFlag  int     `json:"email_flag"`
	TemplateID NullInt `json:"template_id"`
	RejectFlag int     `json:"reject_flag"`
	RejectRole *string `json:"reject_role"`
	StatusFlag int     `json:"status_flag"`
	Badge      int     `json:"badge"`
	Priority   int     `json:"priority"`
	Starred    int     `json:"starred"`

	// COMMENTS
	UserID    string `json:"user_id"`
	UserRole  string `json:"user_role"`
	Remarks   string `json:"remarks"`
	UpdatedBy string `json:"updated_by"`

	// TOKEN
	Token string `json:"token"`
	PID   string `json:"P_id"`
}

// validate checks mandatory fields before processing request
func (req *CriteriaInsertRequest) validate() error {
	if req.ProcessID == 0 {
		return fmt.Errorf("process_id is required")
	}
	if req.AssignTo == "" {
		return fmt.Errorf("assign_to is required")
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	return nil
}

// sendErrorResponse sends encrypted error response to client
func sendErrorResponse(w http.ResponseWriter, r *http.Request, pid, key string, statusCode int, message string) {
	response := APICriteriainsertResponse{
		Status:  statusCode,
		Message: message,
	}
	encryptAndRespond(w, r, pid, key, statusCode, response)
}

// encryptAndRespond encrypts payload and sends JSON response
// Format: { "Data": "PID||EncryptedPayload" }
func encryptAndRespond(w http.ResponseWriter, r *http.Request, pid, key string, statusCode int, payload interface{}) {
	responseJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	encryptedResponse := string(responseJSON)
	if key != "" && pid != "" {
		if enc, err := utils.EncryptAES(string(responseJSON), key); err == nil {
			encryptedResponse = fmt.Sprintf("%s||%s", pid, enc)
		}
	}

	// Wrapper
	finalResp := map[string]string{"Data": encryptedResponse}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(finalResp)
}

// NullInt wraps sql.NullInt64 for optional integer fields
type NullInt struct {
	sql.NullInt64
}

// stringPtrOrNil returns nil if string is empty
// Used for optional workflow parameters
func stringPtrOrNil(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

// CriteriaInsert handles Criteria create / update API
// Supports draft & submit flows with workflow integration
func CriteriaInsert(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Encrypted wrapper
	var encReq EncryptedRequest
	if err := json.Unmarshal(body, &encReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	fmt.Println()
	// Split PID || encryptedData
	parts := strings.Split(encReq.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	// Get key
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decrypt
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Unmarshal request
	var req CriteriaInsertRequest
	if err := json.Unmarshal([]byte(decryptedJSON), &req); err != nil {
		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
		return
	}

	// Validate
	if err := req.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Token & authorization validation
	if req.Token != "" {
		r.Header.Set("token", req.Token)
	}
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// =====================================================
	// MAIN PROCESS
	// =====================================================
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := auth.IsValidIDFromRequest(r); err != nil {
			sendErrorResponse(w, r, pid, key, http.StatusBadRequest, "Invalid TOKEN provided")
			return
		}

		isSubmit := strings.EqualFold(req.TypeOfSubmit, "submit")

		// FIX 1: Removed the redundant `taskIDStr = ""` reset that was blanking
		// out the task ID for new records. Now set once and preserved correctly.
		var taskIDStr string
		if req.TaskID != nil {
			taskIDStr = req.TaskID.String()
		}

		// If submit, call workflow microservice to fetch next activity
		if isSubmit {
			sendBackToMe := false
			if req.SendBackToMe != nil {
				sendBackToMe = *req.SendBackToMe
			}

			currentActivity := req.CurrentActivitySeqNo
			if currentActivity <= 0 {
				currentActivity = req.ActivitySeqNo
			}

			input := micro.NextActivityInput{
				Token:           req.Token,
				ProcessID:       req.ProcessID,
				CurrentActivity: currentActivity,
				TaskID:          stringPtrOrNil(&taskIDStr),
				IsApproved:      req.IsTaskApproved,
				IsReturn:        req.IsTaskReturn,
				Role:            stringPtrOrNil(req.Role),
				RequestedUser:   stringPtrOrNil(req.RequestedUser),
				ReturnToRole:    stringPtrOrNil(req.ReturnToRole),
				ReturnToUser:    stringPtrOrNil(req.ReturnToUser),
				SendBackToMe:    sendBackToMe,
				SendBackToUser:  stringPtrOrNil(req.SendBackToUser),
				Conditions: map[string]string{
					"choice":        "",
					"EmployeeGroup": "",
				},
			}

			out, err := micro.NextActivity(input, pid, key)
			if err != nil {
				sendErrorResponse(w, r, pid, key, http.StatusInternalServerError, err.Error())
				return
			}

			// Update request object with workflow response
			req.ActivitySeqNo = out.NextActivitySeq
			req.AssignedRole = out.AssignedRole
			req.AssignTo = out.AssignTo
			req.EmailFlag = out.EmailFlag
			req.RejectFlag = out.RejectFlag

			if out.TemplateID != nil {
				req.TemplateID = NullInt{sql.NullInt64{
					Int64: int64(*out.TemplateID),
					Valid: true,
				}}
			}
		}

		// DB function call for Insert/Update
		taskID, refNo, err := executeCriteriaInsert(&req)
		if err != nil {
			sendErrorResponse(w, r, pid, key, http.StatusInternalServerError, err.Error())
			return
		}

		// FIX 2: Use the DB-returned taskID for the email queue (same pattern as
		// Ramya's code). This ensures a valid Task ID is always passed, including
		// for newly created records where req.TaskID was nil before the DB call.
		emailTaskID := taskIDStr
		if taskID != uuid.Nil {
			emailTaskID = taskID.String()
		}

		// Trigger email queue only on submit with valid template
		if isSubmit && req.EmailFlag > 0 && req.TemplateID.Valid {
			if err := micro.EmailQueue(
				req.Token,
				req.ProcessID,
				emailTaskID,
				int(req.TemplateID.Int64),
				pid,
				key,
			); err != nil {
				fmt.Printf("EmailQueue failed: %v\n", err)
			}
		}

		// Return success response based on draft / submit
		var resp APICriteriainsertResponse
		if req.TypeOfSubmit == "draft" {
			resp = APICriteriainsertResponse{
				Status:                200,
				Message:               "Draft Added Successfully",
				Response_task_id:      taskID,
				Response_reference_no: refNo,
			}
		} else {
			resp = APICriteriainsertResponse{
				Status:                200,
				Message:               "Criteria Insert Successfully",
				Response_task_id:      taskID,
				Response_reference_no: refNo,
			}
		}
		encryptAndRespond(w, r, pid, key, http.StatusOK, resp)
	})).ServeHTTP(w, r)
}

// executeCriteriaInsert executes PostgreSQL stored procedure
// and retrieves OUT parameters (task_id, reference_no)
func executeCriteriaInsert(req *CriteriaInsertRequest) (uuid.UUID, string, error) {

	// Database connection
	db := credentials.GetDB()

	// Convert CriteriaData slice into JSONB
	criteriaJSON, err := json.Marshal(req.CriteriaData)
	if err != nil {
		return uuid.Nil, "", err
	}

	var taskID uuid.UUID
	var referenceNo string

	// CALL stored procedure using QueryRow + Scan
	err = db.QueryRow(`
		CALL meivan.cmes_qms_full(
			$1,$2,$3,
			$4,$5,$6,$7,$8,
			$9,$10,$11,$12,$13,
			$14,$15,$16,$17,$18,
			$19,$20,$21,$22,$23,$24,$25
		)
	`,
		req.TaskID,      // INOUT uuid
		req.ReferenceNO, // INOUT text
		req.ProcessID,
		criteriaJSON,
		req.TypeOfSubmit,
		req.Status,
		req.AssignTo,
		req.AssignedRole,
		req.TaskStatusID,
		req.ActivitySeqNo,
		req.CurrentActivitySeqNo,
		req.IsTaskReturn,
		req.IsTaskApproved,
		req.EmailFlag,
		req.TemplateID,
		req.RejectFlag,
		req.RejectRole,
		req.StatusFlag,
		req.Badge,
		req.Priority,
		req.Starred,
		req.UserID,
		req.UserRole,
		req.Remarks,
		req.UpdatedBy,
	).Scan(&taskID, &referenceNo)

	if err != nil {
		return uuid.Nil, "", fmt.Errorf("SP execution failed: %w", err)
	}

	return taskID, referenceNo, nil
}