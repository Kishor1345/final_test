// Package controllersqms provides the HTTP controller logic for the Quarters Management System (QMS).
// It manages the request lifecycle, including decryption of incoming data, security validation,
// quarters data retrieval, and the secure encryption of outgoing responses.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/qms
// --- Creator's Info ---
// Creator:  Elakiya
//
// Created On: 
//
// Last Modified By:
//
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

// QMSEUDetailsRequest represents the encrypted payload structure for Estate Unit detail requests.
// It expects the Data field to contain a session ID and encrypted JSON separated by "||".
type QMSEUDetailsRequest struct {
	Data string `json:"Data"`
}


// QMSEUDetails handles the HTTP request for retrieving Estate Unit (EU) specific information.
//
// The workflow includes:
// 1. Validating the POST method and reading the encrypted body.
// 2. Decrypting the payload using the session-based key (PID).
// 3. Extracting the user token and authorizing the request via the auth package.
// 4. Fetching EU-specific data from the database using the decrypted task ID.
// 5. Encrypting and returning the final JSON response to the client.
func QMSEUDetails(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	defer r.Body.Close()

	var req QMSEUDetailsRequest
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

		data, count, err := database.GetQMSEUDetailsFromDB(decryptedData)
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
