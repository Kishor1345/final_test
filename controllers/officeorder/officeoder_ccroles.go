// Package controllersofficeorder provides APIs for office order related operations.
// It includes endpoints used to retrieve CC role details for office orders.
//
//Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 20-11-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 20-11-2025
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

// APICcRolesResponse defines the standard API response structure
// for CC roles, including status, message, and response data.
type APICcRolesResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// CcRolesRequest represents the incoming request body
// containing encrypted request data.
type CcRolesRequest struct {
	Data string `json:"Data"`
}

// CcRoles handles POST requests to fetch CC role details.
// It validates the HTTP method, decrypts the incoming request,
// extracts and verifies the authentication token, performs
// authorization checks, fetches CC role data from the database,
// encrypts the response, and sends it back to the client.
func CcRoles(w http.ResponseWriter, r *http.Request) {
    
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}
    // Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
    // Unmarshal request JSON into request structure
	var req CcRolesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
    // Split request data to extract PID and encrypted payload
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]
    // Fetch decryption key using PID
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Decryption key fetch failed", http.StatusUnauthorized)
		return
	}
    // Decrypt the request payload
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}
    // Unmarshal decrypted JSON into a map
	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}
    // Extract token from decrypted data
	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)
    // Perform API name, IP address, and token validation
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}
    // Log request information and process the request
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate token and request ID
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}
		// Fetch CC roles data from the database

		data, totalCount, err := databaseofficeorder.GetCcRolesFromDB(decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Prepare API response

		response := APICcRolesResponse{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"No Of Records": totalCount,
				"Records":       data,
			},
		}
        // Marshal response into JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Response marshal failed", http.StatusInternalServerError)
			return
		}
        // Encrypt the response JSON
		encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
		if err != nil {
			http.Error(w, "Response encryption failed", http.StatusInternalServerError)
			return
		}
        // Prepare final encrypted response
		finalResp := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
		}

		// Save response log

		auth.SaveResponseLog(
			r,
			finalResp,
			http.StatusOK,
			"application/json",
			len(responseJSON),
			string(body),
		)
        // Write response to client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(finalResp)

	})).ServeHTTP(w, r)
}
