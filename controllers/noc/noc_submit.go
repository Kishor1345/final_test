// Package controllersnoc handles the HTTP controller logic for the No Objection Certificate (NOC) module.
// It orchestrates the submission process, including data decryption, workflow transitions,
// stored procedure execution, dynamic template rendering, and file uploads.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/noc
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 08-01-2026
// Last Modified By: Sridharan
// Last Modified Date: 12-01-2026
package controllersnoc

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/micro"
	"Hrmodule/utils"
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

/* ==========================================================================
   1. ROBUST NULL TYPES
   ========================================================================== */

type NullTime struct{ sql.NullTime }

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		nt.Valid = false
		return nil
	}
	layouts := []string{time.RFC3339, "2006-01-02", "2006-01-02 15:04:05", "02-01-2006"}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			nt.Time, nt.Valid = t, true
			return nil
		}
	}
	return nil
}

type NullString struct{ sql.NullString }

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		ns.Valid = false
		return nil
	}
	ns.String = strings.Trim(string(data), `"`)
	ns.Valid = true
	return nil
}

type NullInt struct{ sql.NullInt64 }

func (ni *NullInt) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
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

/* ==========================================================================
   2. STRUCTS
   ========================================================================== */

type NocResponseData struct {
	TaskID        string `json:"task_id"`
	OrderNo       string `json:"order_no"`
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

type NocSubmitResponse struct {
	Status  int             `json:"Status"`
	Message string          `json:"Message"`
	Data    NocResponseData `json:"Data"`
}

type NocSubmitRequest struct {
	Token        string `json:"token"`
	PID          string `json:"P_id"`
	TypeOfSubmit string `json:"typeofsubmit"`

	TaskID          NullString `json:"p_task_id"`
	ProcessID       int        `json:"p_process_id"`
	EmployeeID      string     `json:"p_employee_id"`
	EmployeeName    string     `json:"p_employee_name"`
	CertificateType int        `json:"p_certificate_type"`
	Purpose         int        `json:"p_purpose"`
	ServiceAtIITM   *string    `json:"p_service_at_iitm"`
	Choice          string     `json:"p_choice"`

	// Workflow Assignment Fields
	AssignTo             string            `json:"p_assign_to"`
	AssignedRole         string            `json:"p_assigned_role"`
	TaskStatusID         int               `json:"p_task_status_id"`
	IsTaskReturn         int               `json:"p_is_task_return"`
	IsTaskApproved       int               `json:"p_is_task_approved"`
	CurrentActivitySeqNo int               `json:"p_current_activity_seq_no"`
	Role                 NullString        `json:"p_role"`
	RequestedUser        NullString        `json:"p_requested_user"`
	ReturnToRole         NullString        `json:"p_return_to_role"`
	ReturnToUser         NullString        `json:"p_return_to_user"`
	SendBackToMe         *bool             `json:"p_send_back_to_me"`
	SendBackToUser       NullString        `json:"p_send_back_to_user"`
	Conditions           map[string]string `json:"p_conditions"`

	InitiatedBy string `json:"p_initiated_by"`
	Status      int    `json:"p_status"`
	Badge       int    `json:"p_badge"`
	Priority    int    `json:"p_priority"`
	Starred     int    `json:"p_starred"`

	// HTML/Display Content
	HeaderHTML    string   `json:"p_header_html"`
	OrderDate     NullTime `json:"p_order_date"`
	ToColumn      string   `json:"p_to_column"`
	Subject       string   `json:"p_subject"`
	Reference     string   `json:"p_reference"`
	BodyHTML      string   `json:"p_body_html"`
	SignatureHTML string   `json:"p_signature_html"`
	CCTo          string   `json:"p_cc_to"`
	FooterHTML    string   `json:"p_footer_html"`

	Comment     string `json:"p_comment"`
	CommentUser string `json:"p_comment_user"`
	UserRole    string `json:"p_user_role"`
	RejectRole  string `json:"p_reject_role"`

	// Child Tables Data
	FaData              json.RawMessage `json:"p_fa_data"`
	IsDeclaration       json.RawMessage `json:"isdeclaration"`
	HsProgrammeOrCourse *string         `json:"p_hs_programme_or_course"`
	HsTypeOfProgramme   *string         `json:"p_hs_type_of_programme"`
	HsDuration          *string         `json:"p_hs_duration"`
	HsUniversity        *string         `json:"p_hs_univ"`
	HsProspectus        *int            `json:"p_hs_prospectus"`
	HsDutyAffect        *int            `json:"p_hs_duty_affect"`
	HsHodRec            *int            `json:"p_hs_hod_rec"`
	HsModeOfStudy       *string         `json:"p_hs_mode_of_study"`
	HsAcademicYear      *string         `json:"p_hs_academic_year"`
	HsStartOfProgram    *string         `json:"p_hs_start_of_program"`
	HsTypeOfStudy       *string         `json:"p_hs_type_of_study"`

	OaInstName   *string `json:"p_oa_inst_name"`
	OaInstAddr   *string `json:"p_oa_inst_addr"`
	OaPostName   *string `json:"p_oa_post_name"`
	OaIsHon      *int    `json:"p_oa_is_hon"`
	OaHonDetails *string `json:"p_oa_hon_details"`
	OaDuration   *int    `json:"p_oa_duration"`
	// OaFrom       NullTime `json:"p_oa_from"`
	// OaTo         NullTime `json:"p_oa_to"`

	OaFrom *string `json:"p_oa_from"`
	OaTo   *string `json:"p_oa_to"`

	PassNocFor *string  `json:"p_pass_noc_for"`
	PassNo     *string  `json:"p_pass_no"`
	PassIssue  NullTime `json:"p_pass_issue"`
	PassValid  NullTime `json:"p_pass_valid"`
	//PassDecl   *int     `json:"p_pass_decl"`
	PassOther *string `json:"p_pass_other"`

	ResAddr    *string         `json:"p_res_addr"`
	ResPurpose *string         `json:"p_res_purpose"`
	ResOther   *string         `json:"p_res_other"`
	ResDepInfo json.RawMessage `json:"p_res_dep_info"`
	//ResDecl    *int            `json:"p_res_decl"`

	ServDepName    *string `json:"p_serv_dep_name"`
	ServSchoolName *string `json:"p_serv_school_name"`
	ServSchoolAddr *string `json:"p_serv_school_addr"`
	ServAcadYear   *string `json:"p_serv_acad_year"`
	ServOther      *string `json:"p_serv_other"`

	VisaPassNo  *string         `json:"p_visa_pass_no"`
	VisaIssue   NullTime        `json:"p_visa_issue"`
	VisaValid   NullTime        `json:"p_visa_valid"`
	VisaFrom    NullTime        `json:"p_visa_from"`
	VisaTo      NullTime        `json:"p_visa_to"`
	VisaCountry *string         `json:"p_visa_country"`
	VisaState   *string         `json:"p_visa_state"`
	VisaCity    *string         `json:"p_visa_city"`
	VisaPurpose *string         `json:"p_visa_purpose"`
	VisaFin     *string         `json:"p_visa_fin_assist"`
	VisaProj    *string         `json:"p_visa_proj_type"`
	VisaOther   *string         `json:"p_visa_other"`
	Templates   json.RawMessage `json:"p_templates_json"`

	// Logic Calculated Fields
	ActivitySeqNo int `json:"-"`
	EmailFlag     int `json:"-"`
	TemplateID    int `json:"-"`
	RejectFlag    int `json:"-"`
}

type FileUploadPayload struct {
	Token      string `json:"token"`
	PID        string `json:"P_id"`
	ProcessID  int    `json:"ProcessId"`
	TaskID     string `json:"TaskId"`
	EmployeeID string `json:"EmployeeId,omitempty"`
	Category   string `json:"Category"`
}

// =====================================================
// DYNAMIC TEMPLATE STRUCTS
// =====================================================

type DynamicTemplateRenderRequest struct {
	Token      string            `json:"token"`
	PID        string            `json:"P_id"`
	ProcessID  int               `json:"ProcessId"`
	Mode       string            `json:"Mode"`
	TaskID     string            `json:"TaskId"`
	Conditions map[string]string `json:"Conditions"`
}

type DynamicTemplateRenderRecord struct {
	ResID            interface{} `json:"res_id"`
	ResHeaderHTML    []byte      `json:"res_header_html"`
	ResToColumn      string      `json:"res_to_column"`
	ResSubject       string      `json:"res_subject"`
	ResReference     string      `json:"res_reference"`
	ResBodyHTML      string      `json:"res_body_html"`
	ResSignatureHTML string      `json:"res_signature_html"`
	ResCcTo          string      `json:"res_cc_to"`
	ResFooterHTML    string      `json:"res_footer_html"`
	ResOrderDate     string      `json:"res_order_date"`
	ResTemplateType  string      `json:"res_template_type"`
}

type DynamicTemplateRenderResponse struct {
	Status  int    `json:"Status"`
	Message string `json:"message"`
	Data    struct {
		NoOfRecords int                           `json:"No Of Records"`
		Records     []DynamicTemplateRenderRecord `json:"Records"`
	} `json:"Data"`
	PID string `json:"P_id"`
}

type TemplateObject struct {
	ID              interface{} `json:"id"`
	CertificateType string      `json:"certificate_type"`
	Purpose         string      `json:"purpose"`
	OrderNo         string      `json:"order_no"`
	OrderDate       string      `json:"order_date"`
	HeaderHTML      []byte      `json:"header_html"`
	ToColumn        string      `json:"to_column"`
	Subject         string      `json:"subject"`
	Reference       string      `json:"reference"`
	BodyHTML        string      `json:"body_html"`
	SignatureHTML   string      `json:"signature_html"`
	CcTo            string      `json:"cc_to"`
	FooterHTML      string      `json:"footer_html"`
}

type DynamicTemplateWriteRequest struct {
	Token         string           `json:"token"`
	PID           string           `json:"P_id"`
	ProcessID     int              `json:"ProcessId"`
	Mode          string           `json:"Mode"`
	TaskID        string           `json:"TaskId"`
	TemplatesJson []TemplateObject `json:"TemplatesJson"`
}

type DynamicTemplateWriteResponse struct {
	Status  int    `json:"Status"`
	Message string `json:"message"`
	PID     string `json:"P_id"`
}

/* ==========================================================================
   3. MAIN HANDLER
   ========================================================================== */

func NocSubmit(w http.ResponseWriter, r *http.Request) {

	var req NocSubmitRequest
	var pid, key string
	var nextRoleName string
	var isEncrypted = true

	// =====================================================
	// 1️⃣ READ REQUEST (JSON OR MULTIPART)
	// =====================================================
	var encryptedData string
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {

		if err := r.ParseMultipartForm(50 << 20); err != nil {
			http.Error(w, "Invalid multipart form", http.StatusBadRequest)
			return
		}

		encryptedData = r.FormValue("Data")
		if encryptedData == "" {
			http.Error(w, "Missing Data field", http.StatusBadRequest)
			return
		}

	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var encReq micro.EncryptedRequest
		if err := json.Unmarshal(body, &encReq); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		encryptedData = encReq.Data
	}

	// =====================================================
	// 2️⃣ SPLIT PID & CIPHERTEXT
	// =====================================================
	parts := strings.Split(encryptedData, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format (expect PID||cipher)", http.StatusBadRequest)
		return
	}

	pid = parts[0]
	cipher := parts[1]

	// =====================================================
	// 3️⃣ DECRYPT
	// =====================================================
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Invalid PID", http.StatusUnauthorized)
		return
	}

	plainJSON, err := utils.DecryptAES(cipher, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// =====================================================
	// 4️⃣ UNMARSHAL
	// =====================================================
	if err := json.Unmarshal([]byte(plainJSON), &req); err != nil {
		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
		return
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

		r.ParseMultipartForm(50 << 20)

		isSubmit := strings.EqualFold(req.TypeOfSubmit, "submit")
		fmt.Printf("🔵 START: %s request for NOC Purpose: %d (EMP: %s)\n", req.TypeOfSubmit, req.Purpose, req.EmployeeID)

		// Workflow Calculation
		if isSubmit {
			sendBackToMe := false
			if req.SendBackToMe != nil {
				sendBackToMe = *req.SendBackToMe
			}

			currentActivity := req.CurrentActivitySeqNo

			workflowConditions := make(map[string]string)
			if req.Conditions != nil {
				for k, v := range req.Conditions {
					workflowConditions[k] = v
				}
			}

			input := micro.NextActivityInput{
				Token:           req.Token,
				ProcessID:       req.ProcessID,
				CurrentActivity: currentActivity,
				TaskID:          toNullString(req.TaskID),
				IsApproved:      req.IsTaskApproved,
				IsReturn:        req.IsTaskReturn,
				Role:            toNullString(req.Role),
				Conditions:      workflowConditions,
				RequestedUser:   toNullString(req.RequestedUser),
				ReturnToRole:    toNullString(req.ReturnToRole),
				ReturnToUser:    toNullString(req.ReturnToUser),
				SendBackToMe:    sendBackToMe,
				SendBackToUser:  toNullString(req.SendBackToUser),
			}

			out, err := micro.NextActivity(input, pid, key)
			if err != nil {
				fmt.Printf("Workflow API Failure: %v\n", err)
				sendErrorResponse(w, r, pid, key, isEncrypted, 500, "Workflow API Error", "", "", 500, err.Error())
				return
			}

			b, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println("🔶 NextActivity API FULL Response:")
			fmt.Println(string(b))

			req.ActivitySeqNo = out.NextActivitySeq
			req.AssignedRole = out.AssignedRole
			req.AssignTo = out.AssignTo
			req.EmailFlag = out.EmailFlag
			req.RejectFlag = out.RejectFlag

			nextRoleName = strings.TrimSpace(out.AssignedRole)

			if out.TemplateID != nil {
				req.TemplateID = *out.TemplateID
			}
			fmt.Printf("Workflow API Success: Next Seq: %d, Target Role: %s\n", req.ActivitySeqNo, nextRoleName)
		} else {
			req.ActivitySeqNo = req.CurrentActivitySeqNo
			if req.AssignedRole == "" {
				req.AssignedRole = req.UserRole
			}
		}

		// Execute Procedure
		fmt.Printf("Calling SP for EMP: %s\n", req.EmployeeID)

		// =====================================================
		// 🗓️ AUTO SET ORDER DATE (YYYY-MM-DD) ON SUBMIT
		// =====================================================
		//if strings.EqualFold(req.TypeOfSubmit, "submit") {

		loc, _ := time.LoadLocation("Asia/Kolkata")

		req.OrderDate = NullTime{
			sql.NullTime{
				Time:  time.Now().In(loc),
				Valid: true,
			},
		}
		//}

		taskID, orderNo, code, msg, err := executeSP_NocSubmit(&req)
		if err != nil {
			fmt.Printf("DB Error: %v\n", err)
			sendErrorResponse(w, r, pid, key, isEncrypted, 500, "Database Error", "", "", 500, err.Error())
			return
		}

		fmt.Printf("DB Success | Code: %d | Msg: %s | TaskID: %s\n", code, msg, taskID)

		// =====================================================
		// 🆕 DYNAMIC TEMPLATE PROCESSING (Refined for Drafts)
		// =====================================================
		if code == 200 {
			var templatesJson []TemplateObject
			isDraft := strings.EqualFold(req.TypeOfSubmit, "draft")

			// Logic:
			// 1. Activity 1 + Submit: Call Render (API) -> then Write
			// 2. Activity 2+ + (Submit or Draft): Use req.Templates -> then Write

			if req.CurrentActivitySeqNo == 1 && isSubmit {
				// --- Activity 1: Only on SUBMIT ---
				fmt.Println("🔷 Activity 1 detected - Calling DynamicTemplate (Mode: render)")

				renderResp, err := callDynamicTemplateRender(req.Token, pid, req.ProcessID, taskID, req.Choice, req.Conditions, key)
				if err != nil {
					fmt.Printf("⚠️ DynamicTemplate Render failed: %v\n", err)
				} else if renderResp != nil && renderResp.Data.NoOfRecords > 0 {
					for _, record := range renderResp.Data.Records {
						templatesJson = append(templatesJson, TemplateObject{
							ID:              record.ResID,
							CertificateType: strconv.Itoa(req.CertificateType),
							Purpose:         record.ResTemplateType,
							OrderNo:         orderNo,
							OrderDate:       formatDateFromRFC3339(record.ResOrderDate),
							HeaderHTML:      record.ResHeaderHTML,
							ToColumn:        record.ResToColumn,
							Subject:         record.ResSubject,
							Reference:       record.ResReference,
							BodyHTML:        record.ResBodyHTML,
							SignatureHTML:   record.ResSignatureHTML,
							CcTo:            record.ResCcTo,
							FooterHTML:      record.ResFooterHTML,
						})
					}
					fmt.Printf("✅ %d Templates mapped from render response\n", len(templatesJson))
				}

			} else if req.CurrentActivitySeqNo > 1 && (isSubmit || isDraft) {
				// --- Activity 2+: Works for both SUBMIT and DRAFT ---
				fmt.Printf("🔷 Activity %d detected (%s) - Using existing p_templates_json\n", req.CurrentActivitySeqNo, req.TypeOfSubmit)

				if len(req.Templates) > 0 {
					if err := json.Unmarshal(req.Templates, &templatesJson); err != nil {
						fmt.Printf("⚠️ Failed to parse p_templates_json: %v\n", err)
					} else {
						fmt.Printf("✅ Using %d existing template(s)\n", len(templatesJson))
					}
				}
			}

			// ✅ Final Step: Call "write" if we have templates to save
			if len(templatesJson) > 0 {
				if err := callDynamicTemplateWrite(req.Token, pid, req.ProcessID, taskID, templatesJson, key); err != nil {
					fmt.Printf("⚠️ DynamicTemplate Write failed: %v\n", err)
				} else {
					fmt.Println("✅ Templates saved successfully")
				}
			}
		}
		// =====================================================
		// 📄 AUTO PDF GENERATION WHEN WORKFLOW COMPLETED
		// =====================================================
		if isSubmit &&
			code == 200 &&
			strings.EqualFold(nextRoleName, "COMPLETED") {

			go func() {
				_, err := micro.GeneratePDFCommon(
					req.Token,
					req.ProcessID,
					taskID,
					"completed",
					pid,
					key,
				)

				if err != nil {
					fmt.Printf("⚠️ PDF generation failed for TaskID %s: %v\n", taskID, err)
				} else {
					fmt.Printf("✅ PDF generated successfully for TaskID %s\n", taskID)
				}
			}()
		}

		// =====================================================
		// MULTIPLE FILE UPLOAD
		// =====================================================
		if r.MultipartForm != nil && len(r.MultipartForm.File["file"]) > 0 {

			for _, fh := range r.MultipartForm.File["file"] {

				file, err := fh.Open()
				if err != nil {
					fmt.Printf("⚠️ Cannot open %s\n", fh.Filename)
					continue
				}

				func() {
					defer file.Close()

					payload := FileUploadPayload{
						Token:     req.Token,
						PID:       pid,
						ProcessID: req.ProcessID,
						TaskID:    taskID,
						Category:  "ongoing",
					}

					if strings.EqualFold(nextRoleName, "COMPLETED") {
						payload.Category = "completed"
						payload.EmployeeID = req.EmployeeID
					}

					if err := callFileUploadAPI(file, fh, payload, pid, key); err != nil {
						fmt.Printf("⚠️ FileUpload failed (%s): %v\n", fh.Filename, err)
					} else {
						fmt.Printf("✅ Uploaded: %s\n", fh.Filename)
					}
				}()
			}
		}

		// Handle Email Queue
		if isSubmit && code == 200 && req.EmailFlag > 0 {
			_ = micro.EmailQueue(req.Token, req.ProcessID, taskID, req.TemplateID, pid, key)
		}

		successMsg := "Draft saved successfully"
		if isSubmit {
			successMsg = "NOC record submitted successfully"
		}
		if code != 200 {
			successMsg = msg
		}

		sendSuccessResponse(w, r, pid, key, isEncrypted, NocSubmitResponse{
			Status:  200,
			Message: successMsg,
			Data: NocResponseData{
				TaskID:        taskID,
				OrderNo:       orderNo,
				StatusCode:    code,
				StatusMessage: msg,
			},
		})
	})).ServeHTTP(w, r)
}

/* ==========================================================================
   4. DATABASE EXECUTION
   ========================================================================== */

func nullJSON(j json.RawMessage) interface{} {
	if len(j) == 0 || string(j) == "null" {
		return nil
	}
	return j
}

func executeSP_NocSubmit(req *NocSubmitRequest) (string, string, int, string, error) {

	// Database connection
	db := credentials.GetDB()

	var outTaskID sql.NullString = req.TaskID.NullString
	var outOrderNo sql.NullString
	var outCode sql.NullInt32
	var outMsg sql.NullString
	var err error
	// Updated to $80 placeholders to match the SQL procedure signature
	query := `CALL meivan.noc_submit(
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
		$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
		$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,
		$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,
		$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,
		$51,$52,$53,$54,$55,$56,$57,$58,$59,$60,
		$61,$62,$63,$64,$65,$66,$67,$68,$69,$70,
		$71,$72,$73,$74,$75,$76,$77,$78,$79)`

	// err = db.QueryRow(query,
	// 	&outTaskID,               // 1: p_out_task_id (INOUT)
	// 	&outOrderNo,              // 2: p_out_order_no (INOUT)
	// 	req.ProcessID,            // 3
	// 	req.EmployeeID,           // 4
	// 	req.EmployeeName,         // 5
	// 	req.CertificateType,      // 6
	// 	req.Purpose,              // 7
	// 	req.ServiceAtIITM,        // 8
	// 	req.AssignTo,             // 9
	// 	req.AssignedRole,         // 10
	// 	req.TaskStatusID,         // 11
	// 	req.ActivitySeqNo,        // 12
	// 	req.IsTaskReturn,         // 13
	// 	req.IsTaskApproved,       // 14
	// 	nullInt(req.EmailFlag),   // 15
	// 	nullInt(req.TemplateID),  // 16
	// 	nullInt(req.RejectFlag),  // 17
	// 	req.RejectRole,           // 18
	// 	req.InitiatedBy,          // 19
	// 	req.Status,               // 20
	// 	req.Badge,                // 21
	// 	req.Priority,             // 22
	// 	req.Starred,              // 23
	// 	toSQLDate(req.OrderDate), // 24
	// 	req.Choice,               // 25
	// 	req.Comment,              // 26
	// 	req.CommentUser,          // 27
	// 	req.UserRole,             // 28
	// 	req.CommentUser,          // 29: p_current_user
	// 	req.CurrentActivitySeqNo, // 30
	// 	nullJSON(req.FaData),     // 31
	// 	nullJSON(req.IsDeclaration),
	// 	req.HsProgrammeOrCourse, // 32
	// 	req.HsTypeOfProgramme,   // 33: ADDED THIS (WAS MISSING)
	// 	req.HsDuration,          // 34
	// 	req.HsUniversity,        // 35
	// 	req.HsProspectus,        // 36 (Integer)
	// 	req.HsDutyAffect,        // 37 (Integer)
	// 	req.HsHodRec,            // 38 (Integer)
	// 	req.HsModeOfStudy,       // 39
	// 	req.HsAcademicYear,      // 40
	// 	req.HsStartOfProgram,    // 41
	// 	req.HsTypeOfStudy,       // 42
	// 	req.OaInstName,          // 43
	// 	req.OaInstAddr,          // 44
	// 	req.OaPostName,          // 45
	// 	req.OaIsHon,             // 46
	// 	req.OaHonDetails,        // 47
	// 	req.OaDuration,          // 48
	// 	// toSQLDate(req.OaFrom),    // 49
	// 	// toSQLDate(req.OaTo),      // 50
	// 	req.OaFrom, // 49
	// 	req.OaTo,   // 50

	// 	req.PassNocFor,           // 51
	// 	req.PassNo,               // 52
	// 	toSQLDate(req.PassIssue), // 53
	// 	toSQLDate(req.PassValid), // 54
	// 	//req.PassDecl,             // 55
	// 	req.PassOther,            // 56
	// 	req.ResAddr,              // 57
	// 	req.ResPurpose,           // 58
	// 	req.ResOther,             // 59
	// 	nullJSON(req.ResDepInfo), // 60
	// 	//req.ResDecl,              // 61
	// 	req.ServDepName,          // 62
	// 	req.ServSchoolName,       // 63
	// 	req.ServSchoolAddr,       // 64
	// 	req.ServAcadYear,         // 65
	// 	req.ServOther,            // 66
	// 	req.VisaPassNo,           // 67
	// 	toSQLDate(req.VisaIssue), // 68
	// 	toSQLDate(req.VisaValid), // 69
	// 	toSQLDate(req.VisaFrom),  // 70
	// 	toSQLDate(req.VisaTo),    // 71
	// 	req.VisaCountry,          // 72
	// 	req.VisaState,            // 73
	// 	req.VisaCity,             // 74
	// 	req.VisaPurpose,          // 75
	// 	req.VisaFin,              // 76
	// 	req.VisaProj,             // 77
	// 	req.VisaOther,            // 78
	// 	&outCode,                 // 79: p_status_code (INOUT)
	// 	&outMsg,                  // 80: p_status_message (INOUT)
	// ).Scan(&outTaskID, &outOrderNo, &outCode, &outMsg)
	err = db.QueryRow(query,
		&outTaskID,                  // 1
		&outOrderNo,                 // 2
		req.ProcessID,               // 3
		req.EmployeeID,              // 4
		req.EmployeeName,            // 5
		req.CertificateType,         // 6
		req.Purpose,                 // 7
		req.ServiceAtIITM,           // 8
		req.AssignTo,                // 9
		req.AssignedRole,            // 10
		req.TaskStatusID,            // 11
		req.ActivitySeqNo,           // 12
		req.IsTaskReturn,            // 13
		req.IsTaskApproved,          // 14
		nullInt(req.EmailFlag),      // 15
		nullInt(req.TemplateID),     // 16
		nullInt(req.RejectFlag),     // 17
		req.RejectRole,              // 18
		req.InitiatedBy,             // 19
		req.Status,                  // 20
		req.Badge,                   // 21
		req.Priority,                // 22
		req.Starred,                 // 23
		toSQLDate(req.OrderDate),    // 24
		req.Choice,                  // 25
		req.Comment,                 // 26
		req.CommentUser,             // 27
		req.UserRole,                // 28
		req.CommentUser,             // 29
		req.CurrentActivitySeqNo,    // 30
		nullJSON(req.FaData),        // 31
		nullJSON(req.IsDeclaration), // 32 ✅
		req.HsProgrammeOrCourse,     // 33
		req.HsTypeOfProgramme,       // 34
		req.HsDuration,              // 35
		req.HsUniversity,            // 36
		req.HsProspectus,            // 37
		req.HsDutyAffect,            // 38
		req.HsHodRec,                // 39
		req.HsModeOfStudy,           // 40
		req.HsAcademicYear,          // 41
		req.HsStartOfProgram,        // 42
		req.HsTypeOfStudy,           // 43
		req.OaInstName,              // 44
		req.OaInstAddr,              // 45
		req.OaPostName,              // 46
		req.OaIsHon,                 // 47
		req.OaHonDetails,            // 48
		req.OaDuration,              // 49
		req.OaFrom,                  // 50
		req.OaTo,                    // 51
		req.PassNocFor,              // 52
		req.PassNo,                  // 53
		toSQLDate(req.PassIssue),    // 54
		toSQLDate(req.PassValid),    // 55
		req.PassOther,               // 56
		req.ResAddr,                 // 57
		req.ResPurpose,              // 58
		req.ResOther,                // 59
		nullJSON(req.ResDepInfo),    // 60
		req.ServDepName,             // 61
		req.ServSchoolName,          // 62
		req.ServSchoolAddr,          // 63
		req.ServAcadYear,            // 64
		req.ServOther,               // 65
		req.VisaPassNo,              // 66
		toSQLDate(req.VisaIssue),    // 67
		toSQLDate(req.VisaValid),    // 68
		toSQLDate(req.VisaFrom),     // 69
		toSQLDate(req.VisaTo),       // 70
		req.VisaCountry,             // 71
		req.VisaState,               // 72
		req.VisaCity,                // 73
		req.VisaPurpose,             // 74
		req.VisaFin,                 // 75
		req.VisaProj,                // 76
		req.VisaOther,               // 77
		&outCode,                    // 78
		&outMsg,                     // 79
	).Scan(&outTaskID, &outOrderNo, &outCode, &outMsg)

	if err != nil {
		return "", "", 500, "", err
	}
	return outTaskID.String, outOrderNo.String, int(outCode.Int32), outMsg.String, nil
}

/* ==========================================================================
   5. DYNAMIC TEMPLATE API CALLS
   ========================================================================== */

func callDynamicTemplateRender(
	token string,
	pid string,
	processID int,
	taskID string,
	choice string,
	reqConditions map[string]string,
	key string,
) (*DynamicTemplateRenderResponse, error) {

	// Build Conditions from request p_conditions and p_choice
	conditions := map[string]string{
		"purpose": choice, // Use p_choice directly
	}

	// Get certificate_type from p_conditions if available
	if certType, ok := reqConditions["certificate_type"]; ok {
		conditions["certificate_type"] = certType
	}

	// Get EmployeeGroup from p_conditions if available
	if empGroup, ok := reqConditions["EmployeeGroup"]; ok {
		conditions["EmployeeGroup"] = empGroup
	}

	payload := DynamicTemplateRenderRequest{
		Token:      token,
		PID:        pid,
		ProcessID:  processID,
		Mode:       "render",
		TaskID:     taskID,
		Conditions: conditions,
	}

	return callDynamicTemplateAPI(payload, pid, key)
}

func callDynamicTemplateWrite(
	token string,
	pid string,
	processID int,
	taskID string,
	templatesJson []TemplateObject,
	key string,
) error {

	payload := DynamicTemplateWriteRequest{
		Token:         token,
		PID:           pid,
		ProcessID:     processID,
		Mode:          "write",
		TaskID:        taskID,
		TemplatesJson: templatesJson,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	encryptedData, err := utils.EncryptAES(string(jsonData), key)
	if err != nil {
		return err
	}

	finalData := fmt.Sprintf("%s||%s", pid, encryptedData)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}

	reqBody := map[string]string{"Data": finalData}
	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://wftest1.iitm.ac.in:5555/DynamicTemplate", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DynamicTemplate write failed: %s", string(respBody))
	}

	fmt.Println("✅ DynamicTemplate Write API called successfully")
	return nil
}

func callDynamicTemplateAPI(
	payload DynamicTemplateRenderRequest,
	pid string,
	key string,
) (*DynamicTemplateRenderResponse, error) {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	encryptedData, err := utils.EncryptAES(string(jsonData), key)
	if err != nil {
		return nil, err
	}

	finalData := fmt.Sprintf("%s||%s", pid, encryptedData)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}

	reqBody := map[string]string{"Data": finalData}
	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://wftest1.iitm.ac.in:5555/DynamicTemplate", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DynamicTemplate API failed: %s", string(respBody))
	}

	// Decrypt response
	var encResp struct {
		Data string `json:"Data"`
	}
	if err := json.Unmarshal(respBody, &encResp); err != nil {
		return nil, err
	}

	parts := strings.Split(encResp.Data, "||")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid response format")
	}

	decrypted, err := utils.DecryptAES(parts[1], key)
	if err != nil {
		return nil, err
	}

	var renderResp DynamicTemplateRenderResponse
	if err := json.Unmarshal([]byte(decrypted), &renderResp); err != nil {
		return nil, err
	}

	return &renderResp, nil
}

func formatDateFromRFC3339(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try other formats
		t, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return dateStr // Return as-is if parsing fails
		}
	}
	return t.Format("2006-01-02")
}

/* ==========================================================================
   7. ADDITIONAL HELPER FUNCTIONS
   ========================================================================== */

func toNullString(ns NullString) *string {
	if ns.Valid && ns.String != "" {
		return &ns.String
	}
	return nil
}

func toSQLDate(nt NullTime) interface{} {
	if nt.Valid {
		return nt.Time
	}
	return nil
}

func sendErrorResponse(w http.ResponseWriter, r *http.Request, pid, key string, isEncrypted bool, httpCode int, msg, tid, ref string, spCode int, spMsg string) {
	resp := NocSubmitResponse{
		Status:  httpCode,
		Message: msg,
		Data: NocResponseData{
			TaskID:        tid,
			OrderNo:       ref,
			StatusCode:    spCode,
			StatusMessage: spMsg,
		},
	}
	respond(w, r, pid, key, isEncrypted, httpCode, resp)
}

func sendSuccessResponse(w http.ResponseWriter, r *http.Request, pid, key string, isEncrypted bool, resp NocSubmitResponse) {
	respond(w, r, pid, key, isEncrypted, 200, resp)
}

func respond(w http.ResponseWriter, r *http.Request, pid, key string, isEncrypted bool, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if isEncrypted {
		j, _ := json.Marshal(payload)
		enc, _ := utils.EncryptAES(string(j), key)
		json.NewEncoder(w).Encode(map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, enc),
		})
	} else {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(payload)
	}
}

func nullInt(v int) interface{} {
	if v == 0 {
		return nil
	}
	return v
}

func callFileUploadAPI(
	file multipart.File,
	fileHeader *multipart.FileHeader,
	payload FileUploadPayload,
	pid string,
	key string,
) error {

	// 1️⃣ Marshal payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// 2️⃣ Encrypt payload
	encryptedData, err := utils.EncryptAES(string(jsonData), key)
	if err != nil {
		return err
	}

	finalData := fmt.Sprintf("%s||%s", pid, encryptedData)

	// 3️⃣ Prepare multipart body
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// file part
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// Data part
	if err := writer.WriteField("Data", finalData); err != nil {
		return err
	}

	writer.Close()

	// 4️⃣ HTTP client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}

	req, err := http.NewRequest(
		"POST",
		"https://wftest1.iitm.ac.in:5555/FileUpload",
		&body,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("file upload failed: %s", string(respBody))
	}

	fmt.Println("✅ FileUpload API called successfully")
	return nil
}
