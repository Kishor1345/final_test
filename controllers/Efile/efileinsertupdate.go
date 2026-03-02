// // Package controllersefile contains structs and queries for InsertUpdateModules.
//path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/controllers/Efile
// // --- Creator's Info ---
// // Creator: Rovita
// // Created On: 24-11-2025
// // Last Modified By: AI Assistant
// // Last Modified Date: 27-11-2025
// // This api is to InsertUpdate data into category_role_map with intelligent upsert logic.
package controllersefile
import (
	"Hrmodule/auth"
	databaseefile "Hrmodule/database/Efile"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIResponseInsertUpdateModulesMaster defines the structure of the JSON response
type APIResponseInsertUpdateModulesMaster struct {
	Status             int    `json:"Status"`
	Message            string `json:"message"`
	ModulesInserted    int64  `json:"modules_inserted,omitempty"`
	ModulesActivated   int64  `json:"modules_activated,omitempty"`
	ModulesDeactivated int64  `json:"modules_deactivated,omitempty"`
	PId                string `json:"P_id,omitempty"`
}

// ErrorResponseEfile defines the structure for error responses
type ErrorResponseEfile struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// InsertUpdateModulesMasterRequest defines the expected request structure
type InsertUpdateModulesMasterRequest struct {
	Data string `json:"Data"`
}

// Decrypted data structure NOTE: includes user_id
type InsertUpdateModulesDecryptedRequestData struct {
	Token    string `json:"token"`
	PId      string `json:"P_id"`
	ModuleID string `json:"module_id"` // comma-separated values or empty
	RoleName string `json:"role_name"`
	UserID   string `json:"user_id"` // required used as created_by / updated_by
}

// sendErrorResponseefile sends JSON error response
func sendErrorResponseefile(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponseefile := ErrorResponseEfile{
		Status:  statusCode,
		Message: message,
	}
	_ = json.NewEncoder(w).Encode(errorResponseefile)
}

// EfileInsertandUpdate handles insert/activate/deactivate requests
func EfileInsertandUpdate(w http.ResponseWriter, r *http.Request) {

	// Validate request method
	if r.Method != http.MethodPost {
		sendErrorResponseefile(w, http.StatusMethodNotAllowed, "Method not allowed, use POST")
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendErrorResponseefile(w, http.StatusBadRequest, "Unable to read body")
		return
	}
	defer r.Body.Close()

	var req InsertUpdateModulesMasterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		sendErrorResponseefile(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	//Split & decrypt "P_id||encrypted"
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		sendErrorResponseefile(w, http.StatusBadRequest, "Invalid Data format")
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		sendErrorResponseefile(w, http.StatusUnauthorized, "Decryption key fetch failed")
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		sendErrorResponseefile(w, http.StatusUnauthorized, "Decryption failed")
		return
	}

	var decryptedData InsertUpdateModulesDecryptedRequestData
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		sendErrorResponseefile(w, http.StatusBadRequest, "Invalid decrypted data")
		return
	}

	//Validate decrypted data fields
	if decryptedData.Token == "" {
		sendErrorResponseefile(w, http.StatusBadRequest, "token is required in decrypted data")
		return
	}
	if decryptedData.RoleName == "" {
		sendErrorResponseefile(w, http.StatusBadRequest, "role_name is required in decrypted data")
		return
	}
	if decryptedData.UserID == "" {
		sendErrorResponseefile(w, http.StatusBadRequest, "user_id is required in decrypted data")
		return
	}
	if decryptedData.PId != "" && decryptedData.PId != pid {
		sendErrorResponseefile(w, http.StatusBadRequest, "P_id mismatch")
		return
	}

	// Put token into header for auth helpers
	r.Header.Set("token", decryptedData.Token)

	//Validate token via auth helpers
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		// helper already writes response on failure
		return
	}
	if err := auth.IsValidIDFromRequest(r); err != nil {
		sendErrorResponseefile(w, http.StatusBadRequest, "Invalid TOKEN")
		return
	}

	// Use user_id from decrypted payload (created_by / updated_by)
	userID := decryptedData.UserID

	//  Business logic: Insert/Activate/Deactivate
	var inserted, activated, deactivated int64

	if strings.TrimSpace(decryptedData.ModuleID) == "" {
		// Special case: deactivate ALL modules for the role
		deactivated, err = databaseefile.DeactivateAllModulesForRole(
			decryptedData.RoleName,
			userID,
		)
		if err != nil {
			sendErrorResponseefile(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// parse comma-separated module IDs
		rawIDs := strings.Split(decryptedData.ModuleID, ",")
		validModuleIDs := []string{}
		for _, id := range rawIDs {
			t := strings.TrimSpace(id)
			if t != "" {
				validModuleIDs = append(validModuleIDs, t)
			}
		}
		if len(validModuleIDs) == 0 {
			sendErrorResponseefile(w, http.StatusBadRequest, "No valid module IDs provided")
			return
		}

		// Replace modules (insert new, activate selected, deactivate others)
		inserted, activated, deactivated, err = databaseefile.ReplaceModulesForRole(
			validModuleIDs,
			decryptedData.RoleName,
			userID,
		)
		if err != nil {
			sendErrorResponseefile(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	//Build response message
	var message string
	if strings.TrimSpace(decryptedData.ModuleID) == "" {
		message = fmt.Sprintf("All modules deactivated for role '%s'. Total affected: %d", decryptedData.RoleName, deactivated)
	} else {
		message = fmt.Sprintf("Successfully processed modules: %d inserted, %d activated, %d deactivated", inserted, activated, deactivated)
	}

	response := APIResponseInsertUpdateModulesMaster{
		Status:             200,
		Message:            message,
		ModulesInserted:    inserted,
		ModulesActivated:   activated,
		ModulesDeactivated: deactivated,
		PId:                pid,
	}

	//Encrypt response
	responseJSON, err := json.Marshal(response)
	if err != nil {
		sendErrorResponseefile(w, http.StatusInternalServerError, "Response marshal failed")
		return
	}

	encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
	if err != nil {
		sendErrorResponseefile(w, http.StatusInternalServerError, "Response encryption failed")
		return
	}

	finalResp := map[string]string{
		"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
	}

	//Save response log (original request body included)
	auth.SaveResponseLog(
		r,
		finalResp,
		http.StatusOK,
		"application/json",
		len(responseJSON),
		string(body),
	)

	//Send final output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(finalResp)
}
