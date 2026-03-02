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
	"log"
	"net/http"
	"strings"
)

// EmployeeDependentDetailsRequest represents the encrypted request wrapper.
type EmployeeDependentDetailsRequest struct {
	Data string `json:"Data"`
}

// DependentAPIResponse defines the standard API response structure
type DependentAPIResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// EmployeeDependentDetailsHandler handles employee dependent details API.
func EmployeeDependentDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {

		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON request
	var req EmployeeDependentDetailsRequest
	if err := json.Unmarshal(body, &req); err != nil {

		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	//Split and decrypt
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

	// Decrypt request payload
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {

		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Convert decrypted JSON to map
	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {

		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {

		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	//Authentication check
	log.Println("Performing authentication check...")
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {

		return
	}

	// Validate token from request
	if err := auth.IsValidIDFromRequest(r); err != nil {

		http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
		return
	}

	//Extract employee ID and business logic
	employeeID, ok := decryptedData["employeeid"].(string)
	if !ok || employeeID == "" {

		http.Error(w, "Missing 'employeeid' in request data", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching dependent details for employee: %s", employeeID)
	data, err := databasesad.FetchEmployeeDependentDetails(employeeID)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare success response
	response := DependentAPIResponse{
		Status:  200,
		Message: "Success",
		Data:    data,
	}

	//Marshal response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {

		http.Error(w, "Response marshal failed", http.StatusInternalServerError)
		return
	}

	//  Encrypt the response
	encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
	if err != nil {

		http.Error(w, "Response encryption failed", http.StatusInternalServerError)
		return
	}

	finalResp := map[string]string{
		"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
	}

	//Save exactly what is sent to client
	auth.SaveResponseLog(
		r,
		finalResp,          // only final response
		http.StatusOK,      // status code
		"application/json", // content type
		len(responseJSON),  // size
		string(body),       // original request
	)

	//Send encrypted response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(finalResp); err != nil {
		return
	}

}
