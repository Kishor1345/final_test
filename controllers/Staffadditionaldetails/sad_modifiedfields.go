// Package modelssad contains structs and queries for Modified feilds.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/Staffadditionaldetails
// --- Creator's Info ---
// Creator: Rovita
//
// Created On: 29-01-2026
//
// Last Modified By:
//  
// Last Modified Date:
package controllerssad

import (
	"Hrmodule/auth"
	databasesad "Hrmodule/database/Staffadditionaldetails"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"bytes"
)

// =====================
// REQUEST STRUCT
// =====================
type SadPersonalDetailsRequest struct {
	Data string `json:"Data"`
}

// =====================
// RESPONSE STRUCT
// =====================
type SadAPIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// =====================
// CONTROLLER
// =====================
func SadPersonalDetails(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	defer r.Body.Close()

	var req SadPersonalDetailsRequest
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

		data, count, err := databasesad.GetSadPersonalDetailsFromDB(decryptedData)
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

		auth.SaveResponseLog(
			r,
			finalResp,
			http.StatusOK,
			"application/json",
			len(respJSON),
			string(body),
		)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(finalResp)

	})).ServeHTTP(w, r)
}
