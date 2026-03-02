// Package controllersqms provides the HTTP controller logic for the Quarters Management System (QMS).
// It manages the request lifecycle, including decryption of incoming data, security validation,
// quarters data retrieval, and the secure encryption of outgoing responses.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/qms
// --- Creator's Info ---
// Creator:  Elakiya
// Created On: 
// Last Modified By:
// Last Modified Date:
package controllersqms

import (
	"Hrmodule/auth"
	database "Hrmodule/database/qms"
	"Hrmodule/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// QuartersStatusRequest represents the encrypted input structure for quarters status requests.
// The Data field must contain the PID and encrypted payload separated by "||".
type QuartersStatusRequest struct {
	Data string `json:"Data"`
}


// QuartersStatus handles the HTTP POST request to retrieve the status of quarters.
//
// This handler performs the following operations:
// 1. Validates the request method and decrypts the incoming AES payload.
// 2. Extracts and sets the authentication token from the decrypted data.
// 3. Verifies API and IP authorization.
// 4. Queries the database for quarters status information using section_id and status filters.
// 5. Encrypts the final result (status, message, and records) before sending the response.
func QuartersStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	defer r.Body.Close()

	var req QuartersStatusRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid, encrypted := parts[0], parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Key error", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encrypted, key)
	if err != nil {
		http.Error(w, "Decrypt failed", http.StatusUnauthorized)
		return
	}

	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	token := decryptedData["token"].(string)
	r.Header.Set("token", token)

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		data, count, err := database.GetQuartersStatusFromDB(decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := APIResponse{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"no_of_records": count,
				"records":       data,
			},
		}

		respJSON, _ := json.Marshal(resp)
		encryptedResp, _ := utils.EncryptAES(string(respJSON), key)

		finalResp := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResp),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(finalResp)

	})).ServeHTTP(w, r)
}
