// // Package controllerssad handles API logic for Staff Additional Details.
// //path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/controllers/Staffadditionaldetails
// // --- Creator's Info ---
// // Creator: Kishorekumar
// // Created On: 29-01-2026
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
)

/* =============================
   REQUEST STRUCT
============================= */

// EmployeeContactRequest represents the encrypted request wrapper
type EmployeeContactRequest struct {
	Data string `json:"Data"`
}

/* =============================
   RESPONSE
============================= */

// APIResponseForContactDetails defines the standard API response
type APIResponseForContactDetails struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

/* =============================
   HANDLER
============================= */

// Employeequartersaddresshandler handles employee dependent details API.
func Employeequartersaddresshandler(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read full request body
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	
	//Unmarshal encrypted request wrapper
	var req EmployeeContactRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	
	// Split PID and encrypted payload using "||"	
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encrypted := parts[1]

	
	// Fetch AES decryption key using PID
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Key error", http.StatusUnauthorized)
		return
	}

	
	// Decrypt request payload
	decrypted, err := utils.DecryptAES(encrypted, key)
	if err != nil {
		http.Error(w, "Decrypt error", http.StatusUnauthorized)
		return
	}

	
	//Unmarshal decrypted payload into map
	var payload map[string]interface{}
	json.Unmarshal([]byte(decrypted), &payload)

	token := payload["token"].(string)
	employeeID := payload["employeeid"].(string)

	r.Header.Set("token", token)

	
	//Authentication and authorization checks
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}
	if err := auth.IsValidIDFromRequest(r); err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	
	//Fetch employee dependent details from DB
	data, err := databasesad.GetEmployeeContactDetails(employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
	//Build success response
	resp := APIResponseForContactDetails{
		Status:  200,
		Message: "Success",
		Data:    data,
	}

	
	//Marshal and encrypt response
	respJSON, _ := json.Marshal(resp)
	encryptedResp, _ := utils.EncryptAES(string(respJSON), key)

	
	//Wrap encrypted response with PID
	finalResp := map[string]string{
		"Data": fmt.Sprintf("%s||%s", pid, encryptedResp),
	}

	
	//Save encrypted response to audit log
	auth.SaveResponseLog(
		r,
		finalResp,
		http.StatusOK,
		"application/json",
		len(respJSON),
		string(body),
	)
	
	//Send encrypted response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalResp)
}
