//package controllerscircular contains data structures and database access logic for the Circular details.
//
//Path :/var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/controllers/hr_008
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:16/01/2026
package controllerscircular
import (
	"Hrmodule/auth"
	databasecircular "Hrmodule/database/HR_008"
	"Hrmodule/utils"
	"encoding/json"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIResponseCriteriaMasterDataFetchForApproval defines
// encrypted API response structure for Criteria Master fetch
type APIResponseCircularDetailDataFetchForApproval struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// CriteriaMasterDataFetchForApprovalRequest wraps encrypted client request
// Expected format: "PID||EncryptedData"
type CircularDetailFetchForApprovalRequest struct {
	Data string `json:"Data"`
}

// CriteriaMasterDataFetchForApproval handles fetching Criteria Master data for approval
func CircularDetailFetchForApproval(w http.ResponseWriter, r *http.Request) {

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
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	defer r.Body.Close()

	// Unmarshal encrypted request wrapper
	var req CircularDetailFetchForApprovalRequest
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

	// Convert decrypted JSON string into map for dynamic access
	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// Extract token from decrypted payload
	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}

		// Fetch Criteria Master data from database layer
		// Returns records list and total record count
		data, totalCount, err := databasecircular.DatabaseCircularDetailFetchForApproval(decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Build success response with record count and data
		// Response format:
		response := APIResponseCircularDetailDataFetchForApproval{
			Status:  200,
			Message : "Circular details fetched for Approval",
			Data: map[string]interface{}{
				"no_of_records": totalCount,
				"records":       data,
			},
		}

		// Marshal response struct to JSON
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

		//  Send to client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(finalResp)

	})).ServeHTTP(w, r)
}
