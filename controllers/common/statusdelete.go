// Package controllersnoc handles HTTP APIs for noc status delete
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On:13-12-2025
//
// Last Modified By:
//
// Last Modified Date:
package controllerscommon

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"log"
	"strings"

	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/lib/pq"
)

type OfficeOrderDeleteTaskRequest struct {
	Data string `json:"Data"`
}

type APIResponseForDeleteTask struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func OfficeOrderstatusdelete(w http.ResponseWriter, r *http.Request) {

	// 1️ Validate Method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// 2 Read Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req OfficeOrderDeleteTaskRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 3️ Split Data → pid||encrypted
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
		log.Printf("Key fetch failed: %v", err)
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// 5️ Decrypt payload
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		log.Printf("Decrypt error: %v", err)
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// 6️ Extract fields (SAFE TYPE HANDLING)
	token, _ := decryptedData["token"].(string)
	taskID, _ := decryptedData["task_id"].(string)

	processIDFloat, ok := decryptedData["process_id"].(float64)
	if !ok {
		http.Error(w, "Invalid process_id", http.StatusBadRequest)
		return
	}
	processID := int(processIDFloat)

	if token == "" || processID <= 0 || taskID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 7 Set token header
	r.Header.Set("token", token)

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 8 Validate token
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}

		// 9 DB connection

		db := credentials.GetDB()

		// Call stored procedure
		_, err = db.Exec(
			`CALL meivan.update_task_delete_dynamic($1, $2)`,
			processID,
			taskID,
		)

		// ERROR HANDLING
		if err != nil {

			statusCode := http.StatusInternalServerError
			message := "Task delete failed"

			if strings.Contains(err.Error(), "Task not found") {
				statusCode = http.StatusNotFound
				message = "Task not found for given task_id and process_id"
			}

			response := APIResponseForDeleteTask{
				Status:  statusCode,
				Message: message,
				Data:    nil,
			}

			responseJSON, _ := json.Marshal(response)
			encryptedResponse, _ := utils.EncryptAES(string(responseJSON), key)

			finalResp := map[string]string{
				"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
			}

			auth.SaveResponseLog(
				r,
				finalResp,
				statusCode,
				"application/json",
				len(responseJSON),
				string(body),
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(finalResp)
			return
		}

		//  SUCCESS RESPONSE
		response := APIResponseForDeleteTask{
			Status:  200,
			Message: "Task deleted successfully",
			Data:    nil,
		}

		responseJSON, _ := json.Marshal(response)
		encryptedResponse, _ := utils.EncryptAES(string(responseJSON), key)

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
