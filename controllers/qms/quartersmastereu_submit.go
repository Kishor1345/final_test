// Package controllersqms provides the HTTP controller logic for the Quarters Management System (QMS).
// It specifically handles the submission and drafting of Estate Unit (EU) quarters information.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/qms
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 19-12-2025
package controllersqms

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
	"time"

	_ "github.com/lib/pq"
)

// ======================
// Robust Null Types
// ======================

// NullTime is a custom wrapper for sql.NullTime to handle multiple JSON date formats during unmarshaling.
type NullTime struct{ sql.NullTime }

// UnmarshalJSON implements the json.Unmarshaler interface for NullTime.
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		nt.Valid = false
		return nil
	}
	layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02", "02-01-2006"}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			nt.Time, nt.Valid = t, true
			return nil
		}
	}
	return nil
}

// NullString is a custom wrapper for sql.NullString to handle null or empty JSON strings.
type NullString struct{ sql.NullString }

// UnmarshalJSON implements the json.Unmarshaler interface for NullString.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		ns.Valid = false
		return nil
	}
	ns.String = strings.Trim(string(data), `"`)
	ns.Valid = true
	return nil
}

// NullInt is a custom wrapper for sql.NullInt64 to handle string-based numeric JSON input.
type NullInt struct{ sql.NullInt64 }

// UnmarshalJSON implements the json.Unmarshaler interface for NullInt.
func (ni *NullInt) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		ni.Valid = false
		return nil
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	ni.Int64, ni.Valid = val, true
	return nil
}

// ======================
// Structs
// ======================

// QmseuResponseData contains the success metadata returned after a QMS EU operation.
type QmseuResponseData struct {
	TaskID        string `json:"task_id"`
	ReferenceNo   string `json:"reference_no"`
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// QmseuInsertResponse defines the standard structure for API responses in the QMS EU module.
type QmseuInsertResponse struct {
	Status  int               `json:"Status"`
	Message string            `json:"Message"`
	Data    QmseuResponseData `json:"Data"`
}

// QmseuInsertRequest encapsulates all possible fields for inserting or updating
// an Estate Unit quarters record, including workflow and decision parameters.
type QmseuInsertRequest struct {
	Token              string     `json:"token"`
	PID                string     `json:"P_id"`
	TypeOfSubmit       string     `json:"typeofsubmit"`
	TaskID             NullString `json:"p_task_id"`
	ProcessID          int        `json:"p_process_id"`
	QuartersNumber     string     `json:"p_quartersnumber"`
	Address            string     `json:"p_address"`
	Floor              NullInt    `json:"p_floor"`
	Street             NullString `json:"p_street"`
	QuartersCategoryID int        `json:"p_quarterscategoryid"`
	BuildingTypeID     int        `json:"p_buildingtypeid"`
	PlinthArea         float64    `json:"p_plintharea"`
	QuartersStatus     string     `json:"p_quartersstatus"`
	ResidentID         string     `json:"p_resident_id"`
	ResidentName       string     `json:"p_resident_name"`
	Department         string     `json:"p_department"`
	Designation        string     `json:"p_designation"`
	OccupiedDate       NullTime   `json:"p_ocupieddate"`
	ContactNo          NullString `json:"p_contactno"`
	LicenceFee         float64    `json:"p_licencefee"`
	SWDCharges         float64    `json:"p_swdcharges"`
	ServiceCharges     float64    `json:"p_servicecharges"`
	GarageCharges      float64    `json:"p_garagecharges"`
	CautionDeposit     float64    `json:"p_cautiondeposit"`
	EffectiveFrom      NullTime   `json:"p_effectivefrom"`
	ValidTill          NullTime   `json:"p_validtill"`
	InitiatedBy        string     `json:"p_initiatedby"`
	Badge              int        `json:"p_badge"`
	Priority           int        `json:"p_priority"`
	Starred            int        `json:"p_starred"`
	Remarks            NullString `json:"p_remarks"`

	TaskStatusID         int        `json:"p_task_status_id"`
	CurrentActivitySeqNo int        `json:"p_current_activity"`
	IsTaskApproved       int        `json:"p_is_task_approved"`
	IsTaskReturn         int        `json:"p_is_task_return"`
	Role                 NullString `json:"p_role"`
	RequestedUser        NullString `json:"p_requested_user"`
	ReturnToRole         NullString `json:"p_return_to_role"`
	ReturnToUser         NullString `json:"p_return_to_user"`
	SendBackToMe         *bool      `json:"p_send_back_to_me"`
	SendBackToUser       NullString `json:"p_send_back_to_user"`

	UserID   string `json:"p_user_id"`
	UserRole string `json:"p_user_role"`

	ActivitySeqNo int    `json:"-"`
	AssignedRole  string `json:"p_assigned_role"`
	AssignTo      string `json:"p_assign_to"`

	EmailFlag  int `json:"-"`
	RejectFlag int `json:"-"`
	TemplateID int `json:"-"`
}

// ======================
// Main Handler
// ======================

// QmseuSubmit manages the lifecycle of an Estate Unit quarters submission.
//
// The function performs the following steps:
// 1. Decrypts and unmarshals the incoming request (supporting PID||Ciphertext format).
// 2. Authorizes the request via the token and security middleware.
// 3. If 'submit' is requested, it calculates the next workflow step using the NextActivity microservice.
// 4. If 'draft' is requested, it persists current data without advancing the workflow.
// 5. Invokes the meivan.qmseu_submit stored procedure to save the record.
// 6. Triggers an email queue entry if the operation is a successful submission with email flags.
func QmseuSubmit(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req QmseuInsertRequest
	var pid, key string
	var isEncrypted bool

	var encReq micro.EncryptedRequest
	if err := json.Unmarshal(body, &encReq); err == nil && strings.Contains(encReq.Data, "||") {
		isEncrypted = true
		parts := strings.Split(encReq.Data, "||")
		pid = parts[0]
		key, _ = utils.GetDecryptKey(pid)
		decrypted, _ := utils.DecryptAES(parts[len(parts)-1], key)
		decrypted = strings.Trim(decrypted, "\x00")
		json.Unmarshal([]byte(decrypted), &req)
	} else {
		isEncrypted = false
		json.Unmarshal(body, &req)
		pid = req.PID
		key, _ = utils.GetDecryptKey(pid)
	}

	if req.Token != "" {
		r.Header.Set("token", req.Token)
	}
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			sendErrorResponse(w, r, pid, key, isEncrypted, 401, "Invalid Token", "", "", 401, "Unauthorized")
			return
		}

		isSubmit := strings.EqualFold(req.TypeOfSubmit, "submit")

		if isSubmit {
			sendBackToMe := false
			if req.SendBackToMe != nil {
				sendBackToMe = *req.SendBackToMe
			}

			currentActivity := req.CurrentActivitySeqNo
			if currentActivity <= 0 {
				currentActivity = 1
			}

			input := micro.NextActivityInput{
				Token:           req.Token,
				ProcessID:       req.ProcessID,
				CurrentActivity: currentActivity,
				TaskID:          toNullString(req.TaskID),
				IsApproved:      req.IsTaskApproved,
				IsReturn:        req.IsTaskReturn,
				Role:            toNullString(req.Role),
				Conditions: map[string]string{
					"QuartersStatus": req.QuartersStatus,
				},
				RequestedUser:  toNullString(req.RequestedUser),
				ReturnToRole:   toNullString(req.ReturnToRole),
				ReturnToUser:   toNullString(req.ReturnToUser),
				SendBackToMe:   sendBackToMe,
				SendBackToUser: toNullString(req.SendBackToUser),
			}

			out, err := micro.NextActivity(input, pid, key)
			if err != nil {
				sendErrorResponse(w, r, pid, key, isEncrypted, 500, "Workflow API Error", "", "", 500, err.Error())
				return
			}

			req.ActivitySeqNo = out.NextActivitySeq
			req.AssignedRole = out.AssignedRole
			req.AssignTo = out.AssignTo
			req.EmailFlag = out.EmailFlag
			req.RejectFlag = out.RejectFlag
			if out.TemplateID != nil {
				req.TemplateID = *out.TemplateID
			}
		} else {
			req.ActivitySeqNo = req.CurrentActivitySeqNo
			if req.AssignedRole == "" {
				req.AssignedRole = req.UserRole
			}
			req.EmailFlag = 0
			req.TemplateID = 0
		}

		taskID, refNo, statusCode, statusMsg, err := executeSP_Qmseu(&req)
		if err != nil {
			sendErrorResponse(w, r, pid, key, isEncrypted, 500, "Database Error", "", "", 500, err.Error())
			return
		}

		if isSubmit && statusCode == 200 && req.EmailFlag > 0 {
			_ = micro.EmailQueue(req.Token, req.ProcessID, taskID, req.TemplateID, pid, key)
		}

		msg := "Draft saved successfully"
		if isSubmit {
			msg = "Record submitted successfully"
		}
		if statusCode == 409 {
			msg = statusMsg
		}

		sendSuccessResponse(w, r, pid, key, isEncrypted, QmseuInsertResponse{
			Status:  200,
			Message: msg,
			Data: QmseuResponseData{
				TaskID:        taskID,
				ReferenceNo:   refNo,
				StatusCode:    statusCode,
				StatusMessage: statusMsg,
			},
		})
	})).ServeHTTP(w, r)
}

// ======================
// SP Logic
// ======================

// executeSP_Qmseu calls the meivan.qmseu_submit PostgreSQL procedure.
// It maps the input request fields to the procedure parameters and captures
// the OUT parameters including the new task ID, reference number, and status codes.
func executeSP_Qmseu(req *QmseuInsertRequest) (string, string, int, string, error) {

	db := credentials.GetDB()

	var outTaskID sql.NullString = req.TaskID.NullString
	var outRefNo sql.NullString
	var outCode sql.NullInt32
	var outMsg sql.NullString
	var err error
	query := `CALL meivan.qmseu_submit($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44)`

	err = db.QueryRow(query,
		&outTaskID, req.ProcessID, req.QuartersNumber, req.Address,
		toSQLInt(req.Floor), toSQLString(req.Street), req.QuartersCategoryID, req.BuildingTypeID,
		req.PlinthArea, req.QuartersStatus, req.ResidentID, req.ResidentName,
		req.Department, req.Designation, toSQLTime(req.OccupiedDate), toSQLString(req.ContactNo),
		req.LicenceFee, req.SWDCharges, req.ServiceCharges, req.GarageCharges,
		req.CautionDeposit, toSQLTime(req.EffectiveFrom), toSQLTime(req.ValidTill),
		req.AssignTo, req.AssignedRole, req.TaskStatusID, req.ActivitySeqNo,
		req.InitiatedBy, req.Badge, req.Priority, req.Starred, toSQLString(req.Remarks),
		req.CurrentActivitySeqNo, req.IsTaskApproved, req.IsTaskReturn,
		req.EmailFlag, req.TemplateID, req.RejectFlag, "",
		req.UserID, req.UserRole, &outRefNo, &outCode, &outMsg,
	).Scan(&outTaskID, &outRefNo, &outCode, &outMsg)

	if err != nil {
		return "", "", 500, "", err
	}
	return outTaskID.String, outRefNo.String, int(outCode.Int32), outMsg.String, nil
}

// Helpers

func toNullString(ns NullString) *string {
	if ns.Valid && ns.String != "" {
		return &ns.String
	}
	return nil
}
func toSQLTime(nt NullTime) interface{} {
	if nt.Valid {
		return nt.Time
	}
	return nil
}
func toSQLString(ns NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}
func toSQLInt(ni NullInt) interface{} {
	if ni.Valid {
		return ni.Int64
	}
	return nil
}

func sendErrorResponse(w http.ResponseWriter, r *http.Request, pid, key string, isEncrypted bool, httpCode int, msg, tid, ref string, spCode int, spMsg string) {
	resp := QmseuInsertResponse{Status: httpCode, Message: msg, Data: QmseuResponseData{TaskID: tid, ReferenceNo: ref, StatusCode: spCode, StatusMessage: spMsg}}
	respond(w, r, pid, key, isEncrypted, httpCode, resp)
}

func sendSuccessResponse(w http.ResponseWriter, r *http.Request, pid, key string, isEncrypted bool, resp QmseuInsertResponse) {
	respond(w, r, pid, key, isEncrypted, 200, resp)
}

func respond(w http.ResponseWriter, r *http.Request, pid, key string, isEncrypted bool, status int, payload interface{}) {
	if isEncrypted {
		jsonBytes, _ := json.Marshal(payload)
		enc, _ := utils.EncryptAES(string(jsonBytes), key)
		final := map[string]string{"Data": fmt.Sprintf("%s||%s", pid, enc)}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(final)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(payload)
	}
}
