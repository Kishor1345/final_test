// Package controllersefile contains structs and queries for InsertModules.
//path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/controllers/Efile
// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to insert data into category_role_map.
package controllersefile
import (
	"Hrmodule/auth"
	databaseefile "Hrmodule/database/Efile"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIResponseInsertModulesMaster defines the structure of the JSON response
type APIResponseInsertModulesMaster struct {
	Status       int    `json:"Status"`
	Message      string `json:"message"`
	RowsAffected int64  `json:"rows_affected,omitempty"`
	PId          string `json:"P_id,omitempty"`
}

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// InsertModulesMasterRequest defines the expected request structure
type InsertModulesMasterRequest struct {
	Data string `json:"Data"`
}

// InsertModulesDecryptedRequestData defines the structure of decrypted data
type InsertModulesDecryptedRequestData struct {
	Token    string `json:"token"`
	PId      string `json:"P_id"`
	ModuleID string `json:"module_id"` // Can be comma-separated values
	RoleName string `json:"role_name"`
	Status   string `json:"status"`
}

// sendErrorResponse sends a JSON error response with the specified status code and message
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponse := ErrorResponse{
		Status:  statusCode,
		Message: message,
	}
	json.NewEncoder(w).Encode(errorResponse)
}


//InsertModules handles the Bank Master API request.
func InsertModules(w http.ResponseWriter, r *http.Request) {
	//Validate Request Method
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed, use POST")
		return
	}

	//Read and Parse Request Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Unable to read body")
		return
	}
	defer r.Body.Close()

	var req InsertModulesMasterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	//Split and decrypt
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid Data format")
		return
	}
	pid := parts[0]
	encryptedPart := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Decryption key fetch failed")
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Decryption failed")
		return
	}

	var decryptedData InsertModulesDecryptedRequestData
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid decrypted data")
		return
	}


	// Validate required fields
	if decryptedData.Token == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Token not found in decrypted data")
		return
	}

	if decryptedData.ModuleID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "module_id is required in decrypted data")
		return
	}

	if decryptedData.RoleName == "" {
		sendErrorResponse(w, http.StatusBadRequest, "role_name is required in decrypted data")
		return
	}

	if decryptedData.Status == "" {
		sendErrorResponse(w, http.StatusBadRequest, "status is required in decrypted data")
		return
	}

	// Check for P_id consistency
	if decryptedData.PId != "" && decryptedData.PId != pid {
		fmt.Printf("P_id mismatch: received %s, expected %s\n", decryptedData.PId, pid)
		sendErrorResponse(w, http.StatusBadRequest, "P_id mismatch")
		return
	}

	r.Header.Set("token", decryptedData.Token)

	//Authentication check
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	if err := auth.IsValidIDFromRequest(r); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid TOKEN")
		return
	}

	//Parse module IDs from comma-separated string
	moduleIDs := strings.Split(decryptedData.ModuleID, ",")
	// Trim spaces from each module ID
	for i, moduleID := range moduleIDs {
		moduleIDs[i] = strings.TrimSpace(moduleID)
	}

	// Validate that we have at least one valid module ID
	validModuleIDs := []string{}
	for _, moduleID := range moduleIDs {
		if moduleID != "" {
			validModuleIDs = append(validModuleIDs, moduleID)
		}
	}

	if len(validModuleIDs) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "No valid module IDs provided")
		return
	}

	//Business logic - insert multiple records into category_role_map
	rowsAffected, err := databaseefile.InsertMultipleCategoryRoleMap(
		validModuleIDs,
		decryptedData.RoleName,
		decryptedData.Status,
	)
	if err != nil {
		// Check if it's a duplicate error
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "all modules already exist") {
			sendErrorResponse(w, http.StatusConflict, err.Error())
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create response
	response := APIResponseInsertModulesMaster{
		Status:       200,
		Message:      fmt.Sprintf("Successfully inserted %d records", rowsAffected),
		RowsAffected: rowsAffected,
		PId:          pid,
	}

	//Marshal & encrypt before sending
	responseJSON, err := json.Marshal(response)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Response marshal failed")
		return
	}

	encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Response encryption failed")
		return
	}

	finalResp := map[string]string{
		"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
	}

	// Save exactly what is sent to client
	auth.SaveResponseLog(
		r,
		finalResp,
		http.StatusOK,
		"application/json",
		len(responseJSON),
		string(body),
	)

	//Send to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalResp)
}

