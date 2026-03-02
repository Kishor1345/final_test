// Package controllersquartersmasterestate handles HTTP APIs for Estate Master Details API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/quartersmasterestate
// --- Creator's Info ---
// Creator:  Ramya M R
//
// Created On: 19-01-2026
//
// Last Modified By:
//
// Last Modified Date:
package controllersquartersmasterestate

import (
	"Hrmodule/auth"
	databasequartersmasterestate "Hrmodule/database/quartersmasterestate"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// API response structure
type APIResponseEstateMasterDetails struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Request structure
type EstateMasterDetailsRequest struct {
	Data string `json:"Data"`
}

// EstateMasterDetails — Controller handler
func EstateMasterDetails(w http.ResponseWriter, r *http.Request) {

	// 1 Validate request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// 2 Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req EstateMasterDetailsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 3 Split PID and encrypted data
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	// 4 Get decrypt key
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Decryption key fetch failed", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// -------------------------------------------------
	// 5 Validation (Updated to be flexible)
	// -------------------------------------------------

	// Token validation
	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	// category_id validation
	// REMOVED: strict .(float64) check.
	// The DB layer will now handle parsing if it's a string OR a number.
	if cat, exists := decryptedData["Campus_Id"]; !exists || cat == nil {
		http.Error(w, "CampusId is required", http.StatusBadRequest)
		return
	}

	// -------------------------------------------------
	// 6 Authentication
	// -------------------------------------------------
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}

		// -------------------------------------------------
		// 7 DB call (Using the robust parsing logic)
		// -------------------------------------------------
		data, totalCount, err := databasequartersmasterestate.GetEstateMasterDetailsFromDB(decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// -------------------------------------------------
		// 8 Response formatting
		// -------------------------------------------------
		response := APIResponseEstateMasterDetails{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"no_of_records": totalCount,
				"records":       data,
			},
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Response marshal failed", http.StatusInternalServerError)
			return
		}

		encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
		if err != nil {
			http.Error(w, "Response encryption failed", http.StatusInternalServerError)
			return
		}

		finalResp := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
		}

		auth.SaveResponseLog(
			r,
			finalResp,
			http.StatusOK,
			"application/json",
			len(responseJSON),
			string(body),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(finalResp)

	})).ServeHTTP(w, r)
}
