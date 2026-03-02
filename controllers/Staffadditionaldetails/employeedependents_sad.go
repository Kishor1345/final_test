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
/* ============================
   REQUEST / RESPONSE
============================ */
// EmployeeDependentDetailsforsad represents the encrypted request wrapper
type EmployeeDependentDetailsforsad struct {
	Data string `json:"Data"`
}

// EmployeeEmployeedependentAPIResponse defines the standard API response
type EmployeeEmployeedependentAPIResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

/* ============================
   HANDLER
============================ */
// EmployeeDependentDetailssadhandler handles employee dependent details API.
func EmployeeDependentDetailssadhandler(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read full request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//Unmarshal encrypted request wrapper
	var req EmployeeDependentDetailsforsad
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
	if err := json.Unmarshal([]byte(decrypted), &payload); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	//Extract and validate token
	token, ok := payload["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token missing", http.StatusBadRequest)
		return
	}
	// Attach token to request header for authentication
	r.Header.Set("token", token)

	//Authentication and authorization checks
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}
	if err := auth.IsValidIDFromRequest(r); err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	//Extract and validate employee ID
	employeeID, ok := payload["employeeid"].(string)
	if !ok || employeeID == "" {
		http.Error(w, "employeeid missing", http.StatusBadRequest)
		return
	}

	//Fetch employee dependent details from DB
	data, err := databasesad.GetEmployeeDependentDetails(employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Build success response
	apiResp := EmployeeEmployeedependentAPIResponse{
		Status:  200,
		Message: "Success",
		Data:    data,
	}

	//Marshal and encrypt response
	respJSON, _ := json.Marshal(apiResp)
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
