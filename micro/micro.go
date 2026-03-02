// Package micro handles external microservice orchestrations for workflows,
// including PDF generation, email queuing, and activity transitions.
//
// It is designed to be fully independent from controllers and direct database logic,
// serving as a gateway to centralized infrastructure services.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/micro
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 19-12-2025
package micro

import (
	"Hrmodule/utils"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// =====================================================
// Common Wrapper
// =====================================================

type EncryptedRequest struct {
	Data string `json:"Data"`
}

// =====================================================
// NEXT ACTIVITY (Request / Response)
// =====================================================

type NextActivityInput struct {
	Token           string
	ProcessID       int
	CurrentActivity int
	TaskID          *string
	IsApproved      int
	IsReturn        int
	Role            *string
	Conditions      map[string]string
	RequestedUser   *string
	ReturnToRole    *string
	ReturnToUser    *string
	SendBackToMe    bool
	SendBackToUser  *string
}

type nextActivityRequest struct {
	Token           string            `json:"token"`
	ProcessID       int               `json:"ProcessId"`
	CurrentActivity int               `json:"Current_Activity"`
	TaskID          *string           `json:"TaskId"`
	IsApproved      int               `json:"IsApproved"`
	IsReturn        int               `json:"IsReturn"`
	Role            *string           `json:"Role"`
	EmployeeID      *string           `json:"EmployeeId"`
	Conditions      map[string]string `json:"Conditions"`
	RequestedUser   *string           `json:"RequestedUser"`
	ReturnToRole    *string           `json:"ReturnToRole"`
	ReturnToUser    *string           `json:"ReturnToUser"`
	SendBackToMe    bool              `json:"SendBackToMe"`
	SendBackToUser  *string           `json:"SendBackToUser"`
}

type NextActivityOutput struct {
	NextActivitySeq int
	AssignedRole    string
	AssignTo        string
	EmailFlag       int
	RejectFlag      int
	TemplateID      *int
	FinalRoleName   string
}

type nextActivityResponse struct {
	Status  interface{} `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

type workflowData struct {
	Records []struct {
		NextActivities string  `json:"next_activities"`
		RoleNames      string  `json:"role_names"`
		EmailFlag      int     `json:"emailflag"`
		TemplateID     *int    `json:"template_id"`
		AssignTo       *string `json:"assign_to"`
		RejectFlag     int     `json:"reject_flag"`
	} `json:"Records"`
}

// =====================================================
// EMAIL QUEUE
// =====================================================

type emailQueueRequest struct {
	ProcessID  int    `json:"ProcessId"`
	TaskID     string `json:"TaskId"`
	TemplateID int    `json:"TemplateId"`
	Token      string `json:"token"`
}

type emailQueueResponse struct {
	Status  interface{} `json:"Status"`
	Message string      `json:"Message"`
}

// =====================================================
// PDF API
// =====================================================
// =====================================================
// PDF API RESPONSE STRUCTS
// =====================================================

type PDFPaths struct {
	EmployeeUserCopyPath string `json:"employee_user_copy_path"`
	OfficeCopyPath       string `json:"office_copy_path"`
	OrderNo              string `json:"order_no"`
	TaskID               string `json:"task_id"`
	UserCopyPathOffice   string `json:"user_copy_path_office"`
}
type PDFGenerateResponse struct {
	OfficeCopyPath       string `json:"office_copy_path"`
	UserCopyPathOffice   string `json:"user_copy_path_office"`
	EmployeeUserCopyPath string `json:"employee_user_copy_path"`
	OrderNo              string `json:"order_no"`
	TaskID               string `json:"task_id"`
}

type pdfAPIResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    PDFGenerateResponse `json:"data"`
}

// =====================================================
// Helpers
// =====================================================

func statusToString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.Itoa(int(t))
	case int:
		return strconv.Itoa(t)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// =====================================================
// NEXT ACTIVITY CALL
// NextActivity determines the succeeding step in a workflow process by communicating with the central workflow microservice.
//
// This function evaluates the current task state—including approval status, return flags, and business conditions—to
// calculate the next activity sequence, the target role, and any associated email or template triggers.
// All communication is secured using AES-GCM encryption.
//
// Parameters:
//   - input: A NextActivityInput struct containing the current process state and decision parameters.
//   - pid: The session identifier used to prefix encrypted data payloads.
//   - key: The 32-byte AES key used for payload encryption and response decryption.
//
// Returns:
//   - A pointer to NextActivityOutput containing calculated transition details.
//   - An error if encryption/decryption fails, the microservice call fails, or the workflow returns no records.
// =====================================================

func NextActivity(
	input NextActivityInput,
	pid string,
	key string,
) (*NextActivityOutput, error) {

	reqPayload := nextActivityRequest{
		Token:           input.Token,
		ProcessID:       input.ProcessID,
		CurrentActivity: input.CurrentActivity,
		TaskID:          input.TaskID,
		IsApproved:      input.IsApproved,
		IsReturn:        input.IsReturn,
		Role:            input.Role,
		Conditions:      input.Conditions,
		RequestedUser:   input.RequestedUser,
		ReturnToRole:    input.ReturnToRole,
		ReturnToUser:    input.ReturnToUser,
		SendBackToMe:    input.SendBackToMe,
		SendBackToUser:  input.SendBackToUser,
	}

	jsonData, _ := json.Marshal(reqPayload)
	encrypted, err := utils.EncryptAES(string(jsonData), key)
	if err != nil {
		return nil, err
	}

	finalReq, _ := json.Marshal(EncryptedRequest{
		Data: fmt.Sprintf("%s||%s", pid, encrypted),
	})

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Post(
		"https://wftest1.iitm.ac.in:5555/NextActivity",
		"application/json",
		strings.NewReader(string(finalReq)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var encResp EncryptedRequest
	if err := json.Unmarshal(body, &encResp); err != nil {
		return nil, err
	}

	parts := strings.Split(encResp.Data, "||")
	decrypted, err := utils.DecryptAES(parts[len(parts)-1], key)
	if err != nil {
		return nil, err
	}

	var apiResp nextActivityResponse
	if err := json.Unmarshal([]byte(decrypted), &apiResp); err != nil {
		return nil, err
	}

	if statusToString(apiResp.Status) != "200" &&
		statusToString(apiResp.Status) != "HRM-BUS-SUC-0001" {
		return nil, fmt.Errorf("NextActivity failed: %s", apiResp.Message)
	}

	raw, _ := json.Marshal(apiResp.Data)
	var wf workflowData
	if err := json.Unmarshal(raw, &wf); err != nil {
		return nil, err
	}

	if len(wf.Records) == 0 {
		return nil, fmt.Errorf("workflow returned no records")
	}

	rec := wf.Records[0]
	nextSeq, _ := strconv.Atoi(rec.NextActivities)

	assignTo := ""
	if rec.AssignTo != nil {
		assignTo = *rec.AssignTo
	}

	return &NextActivityOutput{
		NextActivitySeq: nextSeq,
		AssignedRole:    rec.RoleNames,
		AssignTo:        assignTo,
		EmailFlag:       rec.EmailFlag,
		RejectFlag:      rec.RejectFlag,
		TemplateID:      rec.TemplateID,
		FinalRoleName:   strings.TrimSpace(rec.RoleNames),
	}, nil
}

// =====================================================
// EMAIL QUEUE CALL
// EmailQueue triggers the asynchronous email notification service for a specific workflow task.
//
// This function packages the process and task identifiers into an encrypted payload
// and sends it to the central EmailQueue microservice. It handles secure communication
// using AES-GCM encryption and validates the microservice response against specific
// business success codes.
//
// Parameters:
//   - token: The authorization token required by the microservice.
//   - processID: The unique identifier for the workflow process.
//   - taskID: The specific task ID associated with the email trigger.
//   - templateID: The identifier for the email template to be used.
//   - pid: The session identifier used to prefix the encrypted data.
//   - key: The 32-byte AES key used for payload encryption and response decryption.
//
// Returns:
//   - nil if the email was successfully queued.
//   - An error if encryption fails, the network call fails, or the microservice
//     returns a non-success status code.
// =====================================================

func EmailQueue(
	token string,
	processID int,
	taskID string,
	templateID int,
	pid string,
	key string,
) error {

	reqPayload := emailQueueRequest{
		ProcessID:  processID,
		TaskID:     taskID,
		TemplateID: templateID,
		Token:      token,
	}

	jsonData, _ := json.Marshal(reqPayload)
	encrypted, err := utils.EncryptAES(string(jsonData), key)
	if err != nil {
		return err
	}

	finalReq, _ := json.Marshal(EncryptedRequest{
		Data: fmt.Sprintf("%s||%s", pid, encrypted),
	})

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Post(
		"https://wftest1.iitm.ac.in:5555/Emailqueue",
		"application/json",
		strings.NewReader(string(finalReq)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var encResp EncryptedRequest
	if err := json.Unmarshal(body, &encResp); err != nil {
		return err
	}

	parts := strings.Split(encResp.Data, "||")
	decrypted, err := utils.DecryptAES(parts[len(parts)-1], key)
	if err != nil {
		return err
	}

	var apiResp emailQueueResponse
	if err := json.Unmarshal([]byte(decrypted), &apiResp); err != nil {
		return err
	}

	if statusToString(apiResp.Status) != "200" &&
		statusToString(apiResp.Status) != "HRM-BUS-SUC-0001" {
		return fmt.Errorf("EmailQueue failed: %s", apiResp.Message)
	}

	return nil
}

// =====================================================
// PDF API CALL
// GeneratePDFCommon communicates with the centralized PDF generation microservice
// to create documents based on specific process and task IDs.
//
// The function performs the following steps:
// 1. Constructs a payload with the provided token, process ID, task ID, and status.
// 2. Encrypts the payload using AES-GCM and wraps it in the PID||Cipher format.
// 3. Executes a secure POST request to the internal PDF generation endpoint.
// 4. Decrypts and unmarshals the response to return the PDF metadata.
//
// Parameters:
//   - token: Authorization token for the request.
//   - processID: The unique identifier for the workflow process.
//   - taskID: The specific task ID for which the PDF is being generated.
//   - status: The current status of the task (e.g., "completed", "ongoing").
//   - pid: The session identifier used for key retrieval.
//   - key: The 32-byte AES key used for encryption and decryption.
//
// Returns:
//   - A pointer to PDFGenerateResponse containing the file path or generation details.
//   - An error if encryption, the network call, or the PDF generation fails.
//
// =====================================================
func GeneratePDFCommon(
	token string,
	processID int,
	taskID string,
	status string,
	pid string,
	key string,
) (*PDFGenerateResponse, error) {

	// Build payload
	reqPayload := map[string]interface{}{
		"token":      token,
		"process_id": processID,
		"task_id":    taskID,
		"status":     status,
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	// Encrypt
	encrypted, err := utils.EncryptAES(string(jsonData), key)
	if err != nil {
		return nil, err
	}

	// Wrap encrypted payload
	finalReq, _ := json.Marshal(EncryptedRequest{
		Data: fmt.Sprintf("%s||%s", pid, encrypted),
	})

	// HTTP Client
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// API Call
	resp, err := client.Post(
		"https://wftest1.iitm.ac.in:5555/pdfapinew",
		"application/json",
		strings.NewReader(string(finalReq)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Decrypt response
	var encResp EncryptedRequest
	if err := json.Unmarshal(body, &encResp); err != nil {
		return nil, err
	}

	parts := strings.Split(encResp.Data, "||")
	decrypted, err := utils.DecryptAES(parts[len(parts)-1], key)
	if err != nil {
		return nil, err
	}

	var apiResp pdfAPIResponse
	if err := json.Unmarshal([]byte(decrypted), &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Status != "success" {
		return nil, fmt.Errorf("PDF generation failed: %s", apiResp.Message)
	}

	return &apiResp.Data, nil
}
