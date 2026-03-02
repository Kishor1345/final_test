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
	"io"
	"net/http"
	"strings"
)

// Request structure for encrypted document details input
type EmployeePersonalDetailsRequest_sad struct {
	Data string `json:"Data"`
	Type string `json:"Type"`
}

// Standard API response structure
type APIResponses_sad struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// Handler to fetch PersonalDetails details
func EmployeePersonalDetailsHandler_sad(w http.ResponseWriter, r *http.Request) {

	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// Parse JSON request
	var req EmployeePersonalDetailsRequest_sad
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Split PID and encrypted payload
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid, encryptedPart := parts[0], parts[1]

	// Get AES decryption key using PID
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Key fetch failed", http.StatusUnauthorized)
		return
	}

	// Decrypt request data
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Convert decrypted JSON to map
	var decrypted map[string]interface{}
	json.Unmarshal([]byte(decryptedJSON), &decrypted)

	// Extract token and set it in request header
	token := decrypted["token"].(string)
	employeeID := decrypted["employeeid"].(string)
	reqType := decrypted["Type"].(string)

	r.Header.Set("token", token)

	//Authentication check
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Fetch  details from database
	data, err := databasesad.FetchEmployeePersonalDetails_sad(employeeID, reqType)
	if err != nil {
		resp := APIResponses_sad{Status: 400, Message: err.Error()}
		j, _ := json.Marshal(resp)
		enc, _ := utils.EncryptAES(string(j), key)
		json.NewEncoder(w).Encode(map[string]string{"Data": pid + "||" + enc})
		return
	}

	// Prepare API response
	resp := APIResponses_sad{Status: 200, Message: "Success", Data: data}
	j, _ := json.Marshal(resp)
	enc, _ := utils.EncryptAES(string(j), key)

	//Send to client
	json.NewEncoder(w).Encode(map[string]string{"Data": pid + "||" + enc})
}
