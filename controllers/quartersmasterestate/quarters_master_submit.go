// Package controllersquartersmasterestate handles HTTP APIs for Estate Master Submit.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/quartersmasterestate
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 20-01-2026
//
// Last Modified By:  Ramya M R
//
// Last Modified Date:  10-02-2026
//
// CHANGES:
// - Fixed user_id to preserve leading zeros in comments table (using VARCHAR instead of BIGINT)
// - Added comprehensive debug logging throughout the code
// - GRID TABLE ONLY: Changed category_idcategory_name, building_idbuilding_name, floor_idfloor_name (all VARCHAR)
// - Fixed is_servant_quarters to handle "yes"/"no" strings and convert to 1/0
// - Added campus_id to grid table
package controllersquartersmasterestate

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
	"strconv"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Custom type to handle is_servant_quarters conversion
type BoolAsInt int

func (b *BoolAsInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "yes", "true", "1":
			*b = 1
		case "no", "false", "0", "":
			*b = 0
		default:
			return fmt.Errorf("invalid boolean string: %s", s)
		}
		return nil
	}

	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*b = BoolAsInt(i)
		return nil
	}

	return fmt.Errorf("is_servant_quarters must be yes/no/true/false/1/0")
}

type GridItemInput struct {
	CategoryName      string    `json:"category_name"`
	BuildingName      string    `json:"building_name"`
	CampusID          int       `json:"campus_id"`
	QuartersType      string    `json:"quarters_type"`
	QuartersNumber    string    `json:"quarters_number"`
	QuartersStatus    string    `json:"quarters_status"`
	Street            string    `json:"street"`
	FloorName         string    `json:"floor_name"`
	Address           string    `json:"address"`
	IsServantQuarters BoolAsInt `json:"is_servant_quarters"`
	ServantQuartersNo string    `json:"servant_quartersno"`
	GarageCharges     float64   `json:"garagecharges"`
	EffectiveFrom     string    `json:"effectivefrom"`
	PlinthArea        float64   `json:"plinth_area"`
	LicenceFee        float64   `json:"licence_fee"`
	ServiceCharge     float64   `json:"service_charge"`
	CautionDeposit    float64   `json:"caution_deposit"`
	SwdCharge         float64   `json:"swd_charge"`
	DisplayName       string    `json:"display_name"`
}

type GridItemResponse struct {
	CategoryName   string  `json:"category_name"`
	BuildingName   string  `json:"building_name"`
	CampusID       int     `json:"campus_id"`
	QuartersType   string  `json:"quarters_type"`
	QuartersNumber string  `json:"quarters_number"`
	QuartersStatus string  `json:"quarters_status"`
	Street         string  `json:"street"`
	FloorName      string  `json:"floor_name"`
	PlinthArea     float64 `json:"plinth_area"`
	LicenceFee     float64 `json:"licence_fee"`
	ServiceCharge  float64 `json:"service_charge"`
	CautionDeposit float64 `json:"caution_deposit"`
	SwdCharge      float64 `json:"swd_charge"`
	DisplayName    string  `json:"display_name"`
}

type EstateMasterSubmitResponseData struct {
	TaskID      string `json:"task_id"`
	ReferenceNo string `json:"reference_no"`
	StatusCode  int    `json:"status_code"`
	StatusMsg   string `json:"status_msg"`
}

type EstateMasterSubmitResponse struct {
	Status  int                            `json:"status"`
	Message string                         `json:"message"`
	Data    EstateMasterSubmitResponseData `json:"Data"`
	PID     string                         `json:"P_id"`
}

type EncryptedRequest struct {
	Data string `json:"Data"`
}

type NullInt struct {
	sql.NullInt64
}

type EstateMasterSubmitRequest struct {
	TaskID       *uuid.UUID     `json:"task_id"`
	ReferenceNO  *string        `json:"reference_no"`
	TypeOfSubmit string         `json:"typeofsubmit"`
	ProcessID    int            `json:"process_id"`
	GridData     json.RawMessage `json:"grid_data"`
	AssignTo     string         `json:"assign_to"`
	AssignedRole string         `json:"assigned_role"`

	TaskStatusID         int               `json:"task_status_id"`
	ActivitySeqNo        int               `json:"activity_seq_no"`
	CurrentActivitySeqNo int               `json:"current_activity_seq_no"`
	IsTaskReturn         int               `json:"p_is_task_return"`
	IsTaskApproved       int               `json:"p_is_task_approved"`

	Role           *string           `json:"p_role"`
	RequestedUser  *string           `json:"p_requested_user"`
	ReturnToRole   *string           `json:"p_return_to_role"`
	ReturnToUser   *string           `json:"p_return_to_user"`
	SendBackToMe   *bool             `json:"p_send_back_to_me"`
	SendBackToUser *string           `json:"p_send_back_to_user"`
	Conditions     map[string]string `json:"p_conditions"`

	EmailFlag  int     `json:"email_flag"`
	TemplateID NullInt `json:"template_id"`
	RejectFlag int     `json:"reject_flag"`
	RejectRole *string `json:"reject_role"`
	Status     int     `json:"status"`
	Badge      int     `json:"badge"`
	Priority   int     `json:"priority"`
	Starred    int     `json:"starred"`

	UserID    string `json:"user_id"`
	UserRole  string `json:"user_role"`
	Remarks   string `json:"remarks"`
	Comments  string `json:"comments"`
	UpdatedBy string `json:"updated_by"`
	
	InitiatedBy *string `json:"initiated_by,omitempty"`
	InitiatedOn *string `json:"initiated_on,omitempty"`
	UpdatedOn   *string `json:"updated_on,omitempty"`

	Token string `json:"token"`
	PID   string `json:"P_id"`
}

func (req *EstateMasterSubmitRequest) validate() error {
	if req.ProcessID == 0 {
		return fmt.Errorf("process_id is required")
	}
	if req.AssignTo == "" {
		return fmt.Errorf("assign_to is required")
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if _, err := strconv.ParseInt(req.UserID, 10, 64); err != nil {
		return fmt.Errorf("user_id must be a valid number")
	}
	if len(req.GridData) > 0 && !json.Valid(req.GridData) {
		return fmt.Errorf("grid_data contains invalid JSON")
	}
	return nil
}

func (req *EstateMasterSubmitRequest) GetComments() string {
	if req.Comments != "" {
		return req.Comments
	}
	return req.Remarks
}

func QuartersMasterSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var dataWrapper map[string]interface{}
	if err := json.Unmarshal(body, &dataWrapper); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dataValue, ok := dataWrapper["data"]
	if !ok {
		dataValue, ok = dataWrapper["Data"]
	}
	if !ok {
		http.Error(w, "Missing 'Data' field", http.StatusBadRequest)
		return
	}

	dataStr, ok := dataValue.(string)
	if !ok {
		http.Error(w, "Data must be a string", http.StatusBadRequest)
		return
	}

	parts := strings.Split(dataStr, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	var req EstateMasterSubmitRequest
	if err := json.Unmarshal([]byte(decryptedJSON), &req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid decrypted JSON: %v", err), http.StatusBadRequest)
		return
	}

	req.PID = pid

	if err := req.validate(); err != nil {
		sendPlainError(w, err.Error())
		return
	}

	if req.Token != "" {
		r.Header.Set("token", req.Token)
	}

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleMainLogic(w, r, &req, pid, key)
	})).ServeHTTP(w, r)
}

func handleMainLogic(w http.ResponseWriter, r *http.Request, req *EstateMasterSubmitRequest, pid, key string) {
	if err := auth.IsValidIDFromRequest(r); err != nil {
		sendPlainError(w, "Invalid TOKEN provided")
		return
	}

	isSubmit := strings.EqualFold(req.TypeOfSubmit, "submit")
	var taskIDStr string
	if req.TaskID != nil {
		taskIDStr = req.TaskID.String()
	}

	if isSubmit {
		sendBackToMe := false
		if req.SendBackToMe != nil {
			sendBackToMe = *req.SendBackToMe
		}

		currentActivity := req.CurrentActivitySeqNo
		if currentActivity <= 0 {
			currentActivity = req.ActivitySeqNo
		}

		conditions := map[string]string{"-default": "true"}
		if req.Conditions != nil && len(req.Conditions) > 0 {
			conditions = req.Conditions
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
			Conditions:      conditions,
		}

		out, err := micro.NextActivity(input, pid, key)
		if err != nil {
			sendPlainError(w, err.Error())
			return
		}

		req.ActivitySeqNo = out.NextActivitySeq
		req.AssignedRole = out.AssignedRole
		req.AssignTo = out.AssignTo
		req.EmailFlag = out.EmailFlag
		req.RejectFlag = out.RejectFlag

		if out.TemplateID != nil {
			req.TemplateID = NullInt{
				sql.NullInt64{Int64: int64(*out.TemplateID), Valid: true},
			}
		}
	}

	responseData, err := executeEstateMasterInsert(req)
	if err != nil {
		sendPlainError(w, err.Error())
		return
	}

	emailTaskID := taskIDStr
	if responseData.TaskID != "" {
		emailTaskID = responseData.TaskID
	}

	if isSubmit && req.EmailFlag > 0 && req.TemplateID.Valid {
		_ = micro.EmailQueue(
			req.Token,
			req.ProcessID,
			emailTaskID,
			int(req.TemplateID.Int64),
			pid,
			key,
		)
	}

	msg := "Estate Record Processed Successfully"
	if strings.EqualFold(req.TypeOfSubmit, "draft") {
		msg = "Draft Added Successfully"
	}

	resp := EstateMasterSubmitResponse{
		Status:  200,
		Message: msg,
		PID:     pid,
		Data:    responseData,
	}

	responseJSON, _ := json.Marshal(resp)
	encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
	if err != nil {
		sendPlainError(w, "Failed to encrypt response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	finalResponse := map[string]string{
		"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
	}
	json.NewEncoder(w).Encode(finalResponse)
}

func executeEstateMasterInsert(req *EstateMasterSubmitRequest) (EstateMasterSubmitResponseData, error) {
	db := credentials.GetDB()

	var transformedGridData string
	if len(req.GridData) == 0 || string(req.GridData) == "null" || string(req.GridData) == "" {
		transformedGridData = "[]"
	} else {
		var inputGrid []GridItemInput
		if err := json.Unmarshal(req.GridData, &inputGrid); err != nil {
			return EstateMasterSubmitResponseData{}, fmt.Errorf("Invalid grid_data format: %w", err)
		}
		gridBytes, _ := json.Marshal(inputGrid)
		transformedGridData = string(gridBytes)
	}

	var oTaskID sql.NullString
	var oRefNo sql.NullString
	var oStatusCode sql.NullInt64
	var oStatusMsg sql.NullString

	err := db.QueryRow(`
		CALL meivan.estatemaster_insert(
			$1::varchar,  -- p_action
			$2::integer,  -- p_priority
			$3::varchar,  -- p_assign_to
			$4::varchar,  -- p_assigned_role
			$5::integer,  -- p_task_status_id
			$6::integer,  -- p_activity_seq_no
			$7::varchar,  -- p_user_role
			$8::varchar,  -- p_user_str
			$9::jsonb,    -- p_grid_data
			$10::uuid,    -- p_task_id
			$11::varchar, -- p_reference_no
			$12::text,    -- p_comments
			$13::integer, -- p_is_task_return
			$14::integer, -- p_is_task_approved
			$15::integer, -- p_email_flag
			$16::integer, -- p_template_id
			$17::integer, -- p_reject_flag
			$18::varchar, -- p_reject_role
			$19::integer, -- p_status_flag
			$20::integer, -- p_badge
			$21::integer, -- p_starred
			$22::integer, -- p_current_activity_seq_no
			NULL,         -- o_task_id (INOUT)
			NULL,         -- o_reference_no (INOUT)
			NULL,         -- o_status_code (INOUT)
			NULL          -- o_status_msg (INOUT)
		)
	`,
		req.TypeOfSubmit,
		req.Priority,
		req.AssignTo,
		req.AssignedRole,
		req.TaskStatusID,
		req.ActivitySeqNo,
		req.UserRole,
		req.UserID,
		transformedGridData,
		req.TaskID,
		req.ReferenceNO,
		req.GetComments(),
		req.IsTaskReturn,
		req.IsTaskApproved,
		req.EmailFlag,
		nullIntToPtr(req.TemplateID),
		req.RejectFlag,
		req.RejectRole,
		req.Status,
		req.Badge,
		req.Starred,
		req.CurrentActivitySeqNo,
	).Scan(&oTaskID, &oRefNo, &oStatusCode, &oStatusMsg)

	if err != nil {
		return EstateMasterSubmitResponseData{}, fmt.Errorf("Procedure call failed: %w", err)
	}

	taskID := ""
	if oTaskID.Valid {
		taskID = oTaskID.String
	}

	refNo := ""
	if oRefNo.Valid {
		refNo = oRefNo.String
	}

	statusCode := 200
	if oStatusCode.Valid {
		statusCode = int(oStatusCode.Int64)
	}

	statusMsg := "Success"
	if oStatusMsg.Valid {
		statusMsg = oStatusMsg.String
	}

	return EstateMasterSubmitResponseData{
		TaskID:      taskID,
		ReferenceNo: refNo,
		StatusCode:  statusCode,
		StatusMsg:   statusMsg,
	}, nil
}

func stringPtrOrNil(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

func nullIntToPtr(ni NullInt) *int {
	if ni.Valid {
		val := int(ni.Int64)
		return &val
	}
	return nil
}

func sendPlainError(w http.ResponseWriter, message string) {
	errorResp := map[string]string{
		"Data": message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorResp)
}