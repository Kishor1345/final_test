// Package controllersofficeorder provides APIs for office order related operations.
// This file contains the OfficeOrderHistory API used to retrieve
// office order history details.
//
// The API accepts an encrypted request payload, decrypts the data using
// a PID-based key, validates authentication and authorization details,
// fetches order history records from the database, and returns the
// response in encrypted format without altering the existing logic.
//
//Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 29-10-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 29-10-2025
package controllersofficeorder

import (
	"Hrmodule/auth"
	databaseofficeorder "Hrmodule/database/officeorder"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIResponseOrderHistory defines the standard response structure
// for the office order history API.
type APIResponseOrderHistory struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}
// OrderHistoryRequest represents the request body structure
// containing encrypted input data.
type OrderHistoryRequest struct {
	Data string `json:"Data"`
}

// OfficeOrderHistory handles POST requests to fetch office order history.
// It validates the request method, decrypts the incoming payload,
// verifies the authentication token and request details, retrieves
// order history data from the database, encrypts the response,
// and returns the final result to the client.
func OfficeOrderHistory(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req OrderHistoryRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid encrypted data format", http.StatusBadRequest)
		return
	}
	pid := parts[0]
	enc := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Key fetch failed", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(enc, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
		return
	}

	token := decryptedData["token"].(string)
	r.Header.Set("token", token)

	// Validate token, api, ip etc.
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}

		data, count, err := databaseofficeorder.GetOrderHistoryFromDB(decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := APIResponseOrderHistory{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"No Of Records": count,
				"Records":       data,
			},
		}

		respJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Marshal error", http.StatusInternalServerError)
			return
		}

		encResp, err := utils.EncryptAES(string(respJSON), key)
		if err != nil {
			http.Error(w, "Encrypt error", http.StatusInternalServerError)
			return
		}

		final := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encResp),
		}

		auth.SaveResponseLog(
			r,
			final,
			http.StatusOK,
			"application/json",
			len(respJSON),
			string(body),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(final)

	})).ServeHTTP(w, r)
}
