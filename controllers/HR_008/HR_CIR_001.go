//package controllerscircular contains data structures and database access logic for the Circular Submit.
//
//Path :/var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/controllers/hr_008
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:17/01/2026
package controllerscircular

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
	// "time"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// APICriteriainsertResponse defines the encrypted API response format
// sent back to the client after insert/update
type APICircularinsertResponse struct {
	Status                int       `json:"status"`
	Message               string    `json:"message"`
	Response_task_id      uuid.UUID `json:"response_task_id"`
	Response_reference_no string    `json:"response_reference_no"`
}

// EncryptedRequest wraps encrypted payload sent from client
// Format: PID||EncryptedData
type EncryptedRequestForCircularSubmit struct {
	Data string `json:"Data"`
}

type GridTableCircular struct{
	QuartersCategory string `json:"quarters_category"`
	QuartersNo       string `json:"quarters_no"`
	Floor            string `json:"quarters_floor"`
	Location         string `json:"location"`
	Password         string `json:"password"`
	KeyBox           string `json:"key_box"`
	FirstChoice      string `json:"first_choice"`
	SecondChoice     string `json:"second_choice"`
	ThirdChoice      string `json:"third_choice"`
	Campus           string `json:"campus"`
}

// CriteriaInsertRequest represents the decrypted request payload
// received from client
type CircularInsertRequest struct {

	TaskID *uuid.UUID `json:"task_id"`

	CriteriaType string `json:"criteria_type"`
	CircularFor  string `json:"circular_for"`

	OpenDateRegistration *string `json:"open_date_registration"`
	LastDateCancellation *string `json:"last_date_cancellation"`

	TemplateID NullInt `json:"template_id"`
	ProcessID  int  `json:"process_id"`
	ActivitySeqNo int `json:"activity_seq_no"`

	AssignTo     string `json:"assign_to"`
	AssignedRole string `json:"assigned_role"`

	IsTaskReturn   int `json:"p_is_task_return"`
	IsTaskApproved int `json:"p_is_task_approved"`
	RejectFlag     int `json:"reject_flag"`
	RejectRole     *string `json:"reject_role"`

	Priority  int `json:"priority"`
	Starred   int `json:"starred"`
	Badge     int `json:"badge"`
	EmailFlag int `json:"email_flag"`

	TaskStatusID int `json:"task_status_id"`
	Status       int `json:"status"`
	UserName     string `json:"user_name"`

	NoOfOpenDays         int `json:"no_of_open_days"`
	NoOfCancellationDays int `json:"no_of_cancellation_days"`

	ReferenceNO *string `json:"reference_no"`

	QuarterData  []GridTableCircular `json:"quarter_data"`
 

	// ACTIVITY
	CurrentActivitySeqNo int `json:"p_current_activity_seq_no"`
	Role           *string           `json:"p_role"`
	RequestedUser  *string           `json:"p_requested_user"`
	ReturnToRole   *string           `json:"p_return_to_role"`
	ReturnToUser   *string           `json:"p_return_to_user"`
	SendBackToMe   *bool             `json:"p_send_back_to_me"`
	SendBackToUser *string           `json:"p_send_back_to_user"`
	Conditions     map[string]string `json:"p_conditions"`
	

	// TEMPLATE
	OrderDate     *string    `json:"order_date"`
	HeaderHTML    string     `json:"header_html"`
	ToColumn      string     `json:"to_column"`
	Subject       string     `json:"subject"`
	ReferenceNoTemp string    `json:"reference_no_temp"`
	BodyHTML      string     `json:"body_html"`
	SignatureHTML string     `json:"signature_html"`
	CCTo          string     `json:"cc_to"`
	FooterHTML    string     `json:"footer_html"`

	// COMMENTS
	UserID   string    `json:"user_id"`
	UserRole string `json:"user_role"`
	Remarks  string `json:"remarks"`

	// TOKEN
	Token        string `json:"token"`
	PID          string `json:"P_id"`
	TypeOfSubmit string `json:"submit_type"`
}


// validate checks mandatory fields before processing request
func (req *CircularInsertRequest) validate() error {
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
	response := APICircularinsertResponse{
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

func (ni *NullInt) UnmarshalJSON(data []byte) error {
    // Handle null
    if string(data) == "null" {
        ni.Valid = false
        ni.Int64 = 0
        return nil
    }

    var temp int64
    if err := json.Unmarshal(data, &temp); err != nil {
        return err
    }

    ni.Int64 = temp
    ni.Valid = true
    return nil
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
func CircularInsert(w http.ResponseWriter, r *http.Request) {

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
	var encReq EncryptedRequestForCircularSubmit
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
	var req CircularInsertRequest
	if err := json.Unmarshal([]byte(decryptedJSON), &req); err != nil {
		 fmt.Printf("Unmarshal error: %v\n", err)
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
				Conditions: req.Conditions,
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
		taskID, refNo, err := executeCircularInsert(&req)
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
		var resp APICircularinsertResponse
		if req.TypeOfSubmit == "draft" {
			resp = APICircularinsertResponse{
				Status:                200,
				Message:               "Draft Added Successfully",
				Response_task_id:      taskID,
				Response_reference_no: refNo,
			}
		} else {
			resp = APICircularinsertResponse{
				Status:                200,
				Message:               "Circular Insert Successfully",
				Response_task_id:      taskID,
				Response_reference_no: refNo,
			}
		}
		encryptAndRespond(w, r, pid, key, http.StatusOK, resp)
	})).ServeHTTP(w, r)
}

// executeCriteriaInsert executes PostgreSQL stored procedure
// and retrieves OUT parameters (task_id, reference_no)
func executeCircularInsert(req *CircularInsertRequest) (uuid.UUID, string, error) {

	// Database connection
	db := credentials.GetDB()

	// Convert GridTableCircular slice into JSONB
	QuarterJSON, err := json.Marshal(req.QuarterData)
	if err != nil {
		return uuid.Nil, "", err
	}

	var taskID uuid.UUID
	var referenceNo string



	// CALL stored procedure using QueryRow + Scan
     // 39
	fmt.Println(string(QuarterJSON))
	err = db.QueryRow(`
    CALL meivan.circular_submit(
        $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
        $11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
        $21,$22,$23,$24,$25,$26,$27,$28,$29,$30,
        $31,$32,$33,$34,$35,$36,$37,$38,$39
    )
`,
    req.TaskID,                   // 1
    req.CriteriaType,             // 2
    req.CircularFor,              // 3
    req.OpenDateRegistration,     // 4
    req.LastDateCancellation,     // 5
    req.TemplateID,               // 6
    req.ProcessID,                // 7
    req.ActivitySeqNo,            // 8
    req.CurrentActivitySeqNo,     // 9
    req.AssignTo,                 // 10
    req.AssignedRole,             // 11
    req.IsTaskReturn,             // 12
    req.IsTaskApproved,           // 13
    req.RejectFlag,               // 14
    req.RejectRole,               // 15
    req.Priority,                 // 16
    req.Starred,                  // 17
    req.Badge,                    // 18
    req.EmailFlag,                // 19
    req.TaskStatusID,             // 20
    req.Status,                   // 21
    req.UserName,                 // 22
    req.NoOfOpenDays,             // 23
    req.NoOfCancellationDays,     // 24

    QuarterJSON,                  // 25  p_quarter_data (jsonb)

    req.UserRole,                 // 26
    req.Remarks,                  // 27
    req.UserID,                   // 28

    req.ReferenceNoTemp,          // 29  p_reference_temp
    req.OrderDate,                // 30
    req.HeaderHTML,               // 31
    req.ToColumn,                 // 32
    req.Subject,                  // 33
    req.ReferenceNO,              // 34  p_reference (INOUT)
    req.BodyHTML,                 // 35
    req.SignatureHTML,            // 36
    req.CCTo,                     // 37
    req.FooterHTML,               // 38
    req.TypeOfSubmit,             // 39
).Scan(&taskID, &referenceNo)

	if err != nil {
		return uuid.Nil, "", fmt.Errorf("SP execution failed: %w", err)
	}

	return taskID, referenceNo, nil
}