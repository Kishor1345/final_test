// Package controllerscommon contains structs and queries for BloodGroup.
//path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/controllers/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all BloodGroup.
package controllerscommon
import (
	"Hrmodule/auth"
	databasesad "Hrmodule/database/common"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIResponseBloodGroupMaster defines the structure of the JSON response
type APIResponseBloodGroupMaster struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
	PId     string      `json:"P_id,omitempty"`
}

// BloodGroupMasterRequest defines the expected request structure
type BloodGroupMasterRequest struct {
	Data string `json:"Data"`
}

// BloodGroupDecryptedRequestData defines the structure of decrypted data
type BloodGroupDecryptedRequestData struct {
	Token string `json:"token"`
	PId   string `json:"P_id"`
}

//BloodGroupMaster handles the Bank Master API request.
func BloodGroupMaster(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// Read full request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal encrypted request wrapper
	var req BloodGroupMasterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Split PID and encrypted payload using "||" separator
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}
	pid := parts[0]
	encryptedPart := parts[1]

	// Fetch AES decryption key using PID
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Decryption key fetch failed", http.StatusUnauthorized)
		return
	}

	// Decrypt request payload using AES
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Unmarshal encrypted request wrapper
	var decryptedData BloodGroupDecryptedRequestData
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if decryptedData.Token == "" {
		http.Error(w, "Token not found in decrypted data", http.StatusBadRequest)
		return
	}

	// Check for P_id consistency
	if decryptedData.PId != "" && decryptedData.PId != pid {
		fmt.Printf("P_id mismatch: received %s, expected %s\n", decryptedData.PId, pid)
		http.Error(w, "P_id mismatch", http.StatusBadRequest)
		return
	}

	r.Header.Set("token", decryptedData.Token)

	//Authentication check
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	if err := auth.IsValidIDFromRequest(r); err != nil {
		http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
		return
	}

	//Business logic - fetch Bank master data
	data, err := databasesad.GetBloodGroupMaster()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build success response with data and PID
	response := APIResponseBloodGroupMaster{
		Status:  200,
		Message: "Success",
		Data:    data,
		PId:     pid, // Include P_id in the response
	}

	// Marshal & encrypt before sending
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Response marshal failed", http.StatusInternalServerError)
		return
	}

	// Encrypt response JSON using AES encryption
	encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
	if err != nil {
		http.Error(w, "Response encryption failed", http.StatusInternalServerError)
		return
	}

	// Wrap encrypted response with PID
	finalResp := map[string]string{
		"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
	}

	// Save encrypted response into audit log
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