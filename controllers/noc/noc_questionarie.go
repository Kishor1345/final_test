// Package controllersnoc handles the HTTP controller logic for the NOC (No Objection Certificate) module.
// It manages the request processing, security validation, and response formatting for questionnaire-related services.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/noc
//
// --- Creator's Info ---
// Creator: Elakiya
// Created On: 21-01-2026
// Last Modified By: Vaishnavi
// Last Modified Date: 22-01-2026
package controllersnoc

import (
	"Hrmodule/auth"
	database "Hrmodule/database/noc"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// APIResponseForQuestionnaire represents the standard structure for a questionnaire API response.
// It contains the data payload, session identifier (P_id), numeric status, and a message string.
type APIResponseForQuestionnaire struct {
	Data    interface{} `json:"Data"`
	P_id    string      `json:"P_id"`
	Status  int         `json:"Status"`
	Message string      `json:"message"`
}

// QuestionnaireRequest represents the encrypted incoming request payload.
// The Data field is expected to be in the format "PID||EncryptedJSON".
type QuestionnaireRequest struct {
	Data string `json:"Data"`
}


// QuestionnaireCertificateHandler manages the HTTP lifecycle for retrieving dynamic questionnaire certificates.
//
// This handler executes the following workflow:
// 1. Enforces the POST method and reads the request body.
// 2. Extracts the session ID (PID) and decrypts the AES-encrypted payload.
// 3. Validates the user token and verifies IP/API authorization through the auth package.
// 4. Retrieves dynamic questionnaire data from the database.
// 5. Packages, marshals, and re-encrypts the result into the standardized APIResponse format.
func QuestionnaireCertificateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req QuestionnaireRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("JSON error: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Key fetch failed", http.StatusUnauthorized)
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

	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token missing", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data, totalCount, err := database.QuestionnaireDynamicFromDB(decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := APIResponseForQuestionnaire{
			Data: map[string]interface{}{
				"No Of Records": totalCount,
				"Records":       data,
			},
			P_id:    pid,
			Status:  200,
			Message: "Success",
		}

		respJSON, _ := json.MarshalIndent(response, "", "  ")
		encryptedResp, _ := utils.EncryptAES(string(respJSON), key)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResp),
		})

	})).ServeHTTP(w, r)
}
