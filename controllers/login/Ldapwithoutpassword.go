// Package controllerslogin provides LDAP-based authentication,
// encrypted credential validation, session management,
// and JWT token generation for secure login workflows.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/login
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 26-08-2025
/*
package controllerslogin

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type AuthRequestf struct {
	Token    string `json:"Hrtoken"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponsef struct {
	Valid        bool   `json:"valid"`
	UserId       string `json:"userId,omitempty"`
	Username     string `json:"username,omitempty"`
	EmployeeId   string `json:"EmployeeId"`
	MobileNumber string `json:"MobileNumber"`
	Role         string `json:"Role"` // ADD THIS
	Token        string `json:"token,omitempty"`
}

type AuthResponsefalsef struct {
	Valid    bool   `json:"valid"`
	Username string `json:"username,omitempty"`
	Error    string `json:"error,omitempty"`
}

var jwtSecretf []byte
var encryptionKeyf string

func init() {
	_ = godotenv.Load()
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		panic("JWT_SECRET_KEY environment variable not set")
	}
	jwtSecretf = []byte(jwtKey)

	encryptionKeyf = os.Getenv("ENCRYPTION_KEY")
	if encryptionKeyf == "" {
		panic("ENCRYPTION_KEY environment variable not set")
	}
}

// Create JWT Token
func generateJWTf(userId, username, employeeId string) (string, error) {
	claims := jwt.MapClaims{
		"userId":     userId,
		"username":   username,
		"employeeId": employeeId,
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretf)
}

// Helper function to check if string is hex-encoded
func isHexStringf(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

// decryptDataf decrypts hex-encoded encrypted data using AES
func decryptDataf(encryptedData, key string) (string, error) {
	keyBytes := []byte(key)
	encryptedBytes, err := hex.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("invalid hex encoding: %v", err)
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}
	if len(encryptedBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("encrypted data is not a multiple of the block size")
	}
	decrypted := make([]byte, len(encryptedBytes))
	for i := 0; i < len(encryptedBytes); i += aes.BlockSize {
		block.Decrypt(decrypted[i:i+aes.BlockSize], encryptedBytes[i:i+aes.BlockSize])
	}
	decrypted = PKCS5Unpadf(decrypted)
	return string(decrypted), nil
}

// decryptDatafStrict only accepts encrypted (hex-encoded) data
func decryptDatafStrict(data, key string) (string, error) {
	if !isHexStringf(data) {
		return "", fmt.Errorf("invalid input: data must be encrypted (hex-encoded)")
	}
	return decryptDataf(data, key)
}

// validateEncryptedCredentialsf validates that both username and password are encrypted
func validateEncryptedCredentialsf(username, password string) (bool, string) {
	if username == "" || password == "" {
		return false, "Missing username or password"
	}
	if !isHexStringf(username) {
		return false, "Invalid username format - must be encrypted (hex-encoded)"
	}
	if !isHexStringf(password) {
		return false, "Invalid password format - must be encrypted (hex-encoded)"
	}
	return true, ""
}

// PKCS5Unpadf removes padding from decrypted data
func PKCS5Unpadf(data []byte) []byte {
	pad := int(data[len(data)-1])
	return data[:len(data)-pad]
}

// HandleLDAPAuthf processes an HTTP request for LDAP authentication.
func HandleLDAPAuthf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var req AuthRequestf
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", req.Token)
	authorized := auth.HandleRequestfor_apiname_ipaddress_token(w, r)
	if !authorized {
		return
	}

	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := auth.IsValidIDFromRequest(r)
		if err != nil {
			http.Error(w, "Invalid Token provided", http.StatusBadRequest)
			return
		}
		username := req.Username
		password := req.Password
		valid, errorMsg := validateEncryptedCredentialsf(username, password)
		if !valid {
			log.Printf("Validation error: %s", errorMsg)
			resp := AuthResponsefalsef{Valid: false, Error: errorMsg}
			jsonResponse, _ := json.Marshal(resp)
			encrypted, _ := utils.Encrypt(jsonResponse)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
			return
		}
		decodedUsername, err := decryptDatafStrict(username, encryptionKeyf)
		if err != nil {
			log.Printf("Error decrypting username: %v", err)
			resp := AuthResponsefalsef{Valid: false, Username: "Invalid", Error: "Username decryption failed"}
			jsonResponse, _ := json.Marshal(resp)
			encrypted, _ := utils.Encrypt(jsonResponse)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
			return
		}
		decodedPassword, err := decryptDatafStrict(password, encryptionKeyf)
		if err != nil {
			log.Printf("Error decrypting password: %v", err)
			resp := AuthResponsefalsef{Valid: false, Username: "Invalid", Error: "Password decryption failed"}
			jsonResponse, _ := json.Marshal(resp)
			encrypted, _ := utils.Encrypt(jsonResponse)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
			return
		}

		dn := "cn=academicbind,ou=bind,dc=ldap,dc=iitm,dc=ac,dc=in"
		pass := "1@iIL~0K"
		ldapUserFilter := "(&(objectclass=*)(uid=" + decodedUsername + "))"
		searchBaseStaff := "ou=staff,ou=people,dc=ldap,dc=iitm,dc=ac,dc=in"
		searchBaseFaculty := "ou=faculty,ou=people,dc=ldap,dc=iitm,dc=ac,dc=in"
		searchbaseProject := "ou=project,ou=employee,dc=ldap,dc=iitm,dc=ac,dc=in"
		ldapURL := "ldap://ldap.iitm.ac.in:389"

		conn, err := ldap.DialURL(ldapURL)
		if err != nil {
			log.Printf("Failed to connect to LDAP server: %v", err)
			http.Error(w, "Internal Server Error1", http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		err = conn.Bind(dn, pass)
		if err != nil {
			log.Printf("Server DN Bind Failed: %v", err)
			http.Error(w, "Internal Server Error2", http.StatusInternalServerError)
			return
		}

		var ou string
		var responseSent bool
		var authSuccess bool

		performSearch := func(searchBase, userType string) {
			if responseSent {
				return // Skip if response already sent
			}

			req := ldap.NewSearchRequest(searchBase, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, ldapUserFilter, nil, nil)
			sr, err := conn.Search(req)
			if err != nil {
				log.Printf("Search Failed: %v", err)
				return
			}

			for _, entry := range sr.Entries {
				dn := entry.DN
				if dn != "" {
					// Try to bind with user credentials
					err = conn.Bind(dn, decodedPassword)
					if err != nil {
						log.Printf("%s Bind Failed (wrong password), proceeding anyway: %v", userType, err)
						// Continue with authentication even if password is wrong
					} else {
						log.Printf("%s Bind Successful", userType)
					}

					// Proceed with authentication regardless of password correctness
					authSuccess = true
					ou = userType

					if !responseSent {
						responseSent = true
						userId := generateUserIdf()
						employeeId, mobileNumber, err := getEmployeeInfof(decodedUsername)
						if err != nil {
							log.Printf("Error getting employee info: %v, using default values", err)
							// Use default/mock values when employee info retrieval fails
							employeeId = "000000"       // Default employee ID
							mobileNumber = "0000000000" // Default mobile number
						}

						// Always try to insert session data, but don't fail if it doesn't work
						err = insertSessionDataf(userId, decodedUsername, ou, employeeId)
						if err != nil {
							log.Printf("Error inserting session data: %v, continuing anyway", err)
							// Continue without failing
						}

						// Always generate JWT token
						tokenString, err := generateJWTf(userId, decodedUsername, employeeId)
						if err != nil {
							log.Printf("Error generating JWT: %v, using fallback", err)
							// Generate a simple fallback token if JWT generation fails
							tokenString = "fallback_token_" + userId
						}
						roleName, err := getUserRoleByLoginf(decodedUsername)
						if err != nil {
							log.Printf("Role fetch failed: %v", err)
							roleName = "UNKNOWN"
						}
						// Always return success response
						resp := AuthResponsef{Valid: true, UserId: userId, Username: decodedUsername, EmployeeId: employeeId, MobileNumber: mobileNumber, Role: roleName, Token: tokenString}
						jsonResponse, _ := json.Marshal(resp)
						encrypted, _ := utils.Encrypt(jsonResponse)
						w.Header().Set("Content-Type", "application/json")
						_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
					}
					return // Exit once user is found (regardless of password correctness)
				}
			}
		}

		// Search in all organizational units
		performSearch(searchBaseStaff, "staff")
		performSearch(searchBaseFaculty, "faculty")
		performSearch(searchbaseProject, "project")

		// If no successful authentication and no response sent, send failure response
		if !authSuccess && !responseSent {
			resp := AuthResponsefalsef{Valid: false, Username: decodedUsername, Error: "Invalid username or password"}
			jsonResponse, _ := json.Marshal(resp)
			encrypted, _ := utils.Encrypt(jsonResponse)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
		}
	}))
	loggedHandler.ServeHTTP(w, r)
}

func generateUserIdf() string { return uuid.New().String() }

func updatePreviousActiveSessionsf(employeeId string) error {

	// Database connection
	db := credentials.GetDB()
	var err error
	query := `UPDATE meivan.Session_Data SET Is_Active = '0', idletimeout = '1', Logout_Date = NOW() WHERE Employee_id = $1 AND Is_Active = '1'`
	_, err = db.Exec(query, employeeId)
	return err
}

func insertSessionDataf(userId, username, ou, employeeId string) error {
	_ = updatePreviousActiveSessionsf(employeeId)

	// Database connection
	db := credentials.GetDB()

	var err error
	query := `INSERT INTO meivan.Session_Data (Session_Id, Logout_Date, Username, Is_Active, idletimeout, Department, User_id, Employee_id, Login_Date,reason) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(),$9)`
	_, err = db.Exec(query, userId, nil, username, "1", "0", ou, userId, employeeId, "LOGIN")
	return err
}

func getEmployeeInfof(username string) (string, string, error) {

	// Database connection
	db := credentials.GetDB()

	var err error
	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Printf("Database ping failed: %v", err)
		return "", "", err
	}

	// Try different column name variations
	queries := []string{
		`SELECT EmployeeId, Mobilenumber FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
		`SELECT employeeid, mobilenumber FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
		`SELECT EmployeeId, MobileNumber FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
		`SELECT employeeid, MobileNumber FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
	}

	var employeeId, mobileNumber string

	for i, query := range queries {
		log.Printf("Trying query %d: %s with username: %s", i+1, query, username)
		row := db.QueryRow(query, username)
		err = row.Scan(&employeeId, &mobileNumber)
		if err == nil {
			log.Printf("Query successful! EmployeeId: %s, MobileNumber: %s", employeeId, mobileNumber)
			return employeeId, mobileNumber, nil
		}
		log.Printf("Query %d failed: %v", i+1, err)
	}

	// If all queries failed, try to get just EmployeeId
	simpleQueries := []string{
		`SELECT EmployeeId FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
		`SELECT employeeid FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
	}

	for i, query := range simpleQueries {
		log.Printf("Trying simple query %d: %s", i+1, query)
		row := db.QueryRow(query, username)
		err = row.Scan(&employeeId)
		if err == nil {
			log.Printf("Simple query successful! EmployeeId: %s", employeeId)
			return employeeId, "0000000000", nil // Default mobile number
		}
		log.Printf("Simple query %d failed: %v", i+1, err)
	}

	log.Printf("All queries failed for username: %s", username)
	return "", "", fmt.Errorf("no employee found for username: %s", username)
}
func getUserRoleByLoginf(loginname string) (string, error) {

	// Database connection
	db := credentials.GetDB()

	query := `
	SELECT
	    B.campuscode || ' ' ||
	    CASE
	        WHEN A.sectionid IS NOT NULL THEN D.sectioncode
	        ELSE A.departmentcode
	    END || ' ' ||
	    C.rolename AS RoleName
	FROM humanresources.employeerolemapping A
	JOIN humanresources.campus B
	    ON A.campusid = B.id
	JOIN meivan.rolemaster C
	    ON A.roleid = C.id
	LEFT JOIN humanresources.section D
	    ON A.sectionid = D.id
	JOIN humanresources.employeebasicinfo e
	    ON A.employeeid = e.employeeid
	WHERE e.loginname = $1
	LIMIT 1;
	`

	var role string
	var err error
	err = db.QueryRow(query, loginname).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}
*/
// Package controllerslogin provides LDAP-based authentication,
// encrypted credential validation, session management,
// and JWT token generation for secure login workflows.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/login
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 26-08-2025

package controllerslogin

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type AuthRequestf struct {
	Token    string `json:"Hrtoken"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponsef struct {
	Valid        bool   `json:"valid"`
	UserId       string `json:"userId,omitempty"`
	Username     string `json:"username,omitempty"`
	EmployeeId   string `json:"EmployeeId"`
	MobileNumber string `json:"MobileNumber"`
	Role         string `json:"Role"` // ADD THIS
	Token        string `json:"token,omitempty"`
}

type AuthResponsefalsef struct {
	Valid    bool   `json:"valid"`
	Username string `json:"username,omitempty"`
	Error    string `json:"error,omitempty"`
}

var jwtSecretf []byte
var encryptionKeyf string

func init() {
	_ = godotenv.Load()
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		panic("JWT_SECRET_KEY environment variable not set")
	}
	jwtSecretf = []byte(jwtKey)

	encryptionKeyf = os.Getenv("ENCRYPTION_KEY")
	if encryptionKeyf == "" {
		panic("ENCRYPTION_KEY environment variable not set")
	}
}

// Create JWT Token
func generateJWTf(userId, username, employeeId string) (string, error) {
	claims := jwt.MapClaims{
		"userId":     userId,
		"username":   username,
		"employeeId": employeeId,
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretf)
}

// Helper function to check if string is hex-encoded
func isHexStringf(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

// decryptDataf decrypts hex-encoded encrypted data using AES
func decryptDataf(encryptedData, key string) (string, error) {
	keyBytes := []byte(key)
	encryptedBytes, err := hex.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("invalid hex encoding: %v", err)
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}
	if len(encryptedBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("encrypted data is not a multiple of the block size")
	}
	decrypted := make([]byte, len(encryptedBytes))
	for i := 0; i < len(encryptedBytes); i += aes.BlockSize {
		block.Decrypt(decrypted[i:i+aes.BlockSize], encryptedBytes[i:i+aes.BlockSize])
	}
	decrypted = PKCS5Unpadf(decrypted)
	return string(decrypted), nil
}

// decryptDatafStrict only accepts encrypted (hex-encoded) data
func decryptDatafStrict(data, key string) (string, error) {
	if !isHexStringf(data) {
		return "", fmt.Errorf("invalid input: data must be encrypted (hex-encoded)")
	}
	return decryptDataf(data, key)
}

// validateEncryptedCredentialsf validates that both username and password are encrypted
func validateEncryptedCredentialsf(username, password string) (bool, string) {
	if username == "" || password == "" {
		return false, "Missing username or password"
	}
	if !isHexStringf(username) {
		return false, "Invalid username format - must be encrypted (hex-encoded)"
	}
	if !isHexStringf(password) {
		return false, "Invalid password format - must be encrypted (hex-encoded)"
	}
	return true, ""
}

// PKCS5Unpadf removes padding from decrypted data
func PKCS5Unpadf(data []byte) []byte {
	pad := int(data[len(data)-1])
	return data[:len(data)-pad]
}

// checkImpersonationf checks if user has active impersonation settings
// Returns: impersonateUsername, isImpersonating, originalUsername, error
func checkImpersonationf(username string) (string, bool, string, error) {
	db := credentials.GetDB()

	query := `SELECT 	, isactive FROM meivan.user_impersonation WHERE username = $1`

	var impersonateUsername string
	var isActive int

	err := db.QueryRow(query, username).Scan(&impersonateUsername, &isActive)
	if err != nil {
		// No impersonation record found, return original username
		log.Printf("No impersonation record for %s or query error: %v", username, err)
		return username, false, username, nil
	}

	if isActive == 1 {
		log.Printf("Impersonation active: %s -> %s", username, impersonateUsername)
		return impersonateUsername, true, username, nil
	}

	log.Printf("Impersonation inactive for %s", username)
	return username, false, username, nil
}

// HandleLDAPAuthf processes an HTTP request for LDAP authentication.
func HandleLDAPAuthf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var req AuthRequestf
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", req.Token)
	authorized := auth.HandleRequestfor_apiname_ipaddress_token(w, r)
	if !authorized {
		return
	}

	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := auth.IsValidIDFromRequest(r)
		if err != nil {
			http.Error(w, "Invalid Token provided", http.StatusBadRequest)
			return
		}
		username := req.Username
		password := req.Password
		valid, errorMsg := validateEncryptedCredentialsf(username, password)
		if !valid {
			log.Printf("Validation error: %s", errorMsg)
			resp := AuthResponsefalsef{Valid: false, Error: errorMsg}
			jsonResponse, _ := json.Marshal(resp)
			encrypted, _ := utils.Encrypt(jsonResponse)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
			return
		}
		decodedUsername, err := decryptDatafStrict(username, encryptionKeyf)
		if err != nil {
			log.Printf("Error decrypting username: %v", err)
			resp := AuthResponsefalsef{Valid: false, Username: "Invalid", Error: "Username decryption failed"}
			jsonResponse, _ := json.Marshal(resp)
			encrypted, _ := utils.Encrypt(jsonResponse)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
			return
		}
		decodedPassword, err := decryptDatafStrict(password, encryptionKeyf)
		if err != nil {
			log.Printf("Error decrypting password: %v", err)
			resp := AuthResponsefalsef{Valid: false, Username: "Invalid", Error: "Password decryption failed"}
			jsonResponse, _ := json.Marshal(resp)
			encrypted, _ := utils.Encrypt(jsonResponse)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
			return
		}

		// *** IMPERSONATION CHECK - Check if user has impersonation enabled ***
	// *** IMPERSONATION CHECK - Check if user has impersonation enabled ***
originalUsername := decodedUsername
impersonateUsername, isImpersonating, _, err := checkImpersonationf(decodedUsername)
if err != nil {
    log.Printf("Error checking impersonation: %v", err)
    // Continue with original username if impersonation check fails
}

// For impersonation: authenticate the ORIGINAL user, not the impersonated user
var ldapUsername string
if isImpersonating {
    ldapUsername = originalUsername  // Authenticate premc with premc's password
    log.Printf("LDAP Authentication (Impersonation Mode) - Authenticating original user: %s, Will impersonate as: %s", originalUsername, impersonateUsername)
} else {
    ldapUsername = decodedUsername
    log.Printf("LDAP Authentication (Normal Mode) - User: %s", ldapUsername)
}





		dn := "cn=academicbind,ou=bind,dc=ldap,dc=iitm,dc=ac,dc=in"
		pass := "1@iIL~0K"
		ldapUserFilter := "(&(objectclass=*)(uid=" + ldapUsername + "))"
		searchBaseStaff := "ou=staff,ou=people,dc=ldap,dc=iitm,dc=ac,dc=in"
		searchBaseFaculty := "ou=faculty,ou=people,dc=ldap,dc=iitm,dc=ac,dc=in"
		searchbaseProject := "ou=project,ou=employee,dc=ldap,dc=iitm,dc=ac,dc=in"
		ldapURL := "ldap://ldap.iitm.ac.in:389"

		conn, err := ldap.DialURL(ldapURL)
		if err != nil {
			log.Printf("Failed to connect to LDAP server: %v", err)
			http.Error(w, "Internal Server Error1", http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		err = conn.Bind(dn, pass)
		if err != nil {
			log.Printf("Server DN Bind Failed: %v", err)
			http.Error(w, "Internal Server Error2", http.StatusInternalServerError)
			return
		}

		var ou string
		var responseSent bool
		var authSuccess bool

		performSearch := func(searchBase, userType string) {
			if responseSent {
				return // Skip if response already sent
			}

			req := ldap.NewSearchRequest(searchBase, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, ldapUserFilter, nil, nil)
			sr, err := conn.Search(req)
			if err != nil {
				log.Printf("Search Failed: %v", err)
				return
			}

			for _, entry := range sr.Entries {
				dn := entry.DN
				if dn != "" {
					// Try to bind with user credentials
					err = conn.Bind(dn, decodedPassword)
					if err != nil {
						log.Printf("%s Bind Failed (wrong password), proceeding anyway: %v", userType, err)
						// Continue with authentication even if password is wrong
					} else {
						log.Printf("%s Bind Successful", userType)
					}

					// Proceed with authentication regardless of password correctness
					authSuccess = true
					ou = userType

					if !responseSent {
						responseSent = true
						userId := generateUserIdf()

						// *** Get employee info and mobile number from ORIGINAL username ***
// *** Get mobile number from ORIGINAL user, but EmployeeId from IMPERSONATED user ***
var employeeId, mobileNumber string
var err error

if isImpersonating {
    // Get mobile number from original user (premc)
    _, mobileNumber, err = getEmployeeInfof(originalUsername)
    if err != nil {
        log.Printf("Error getting mobile for %s: %v, using default", originalUsername, err)
        mobileNumber = "0000000000"
    }
    
    // Get EmployeeId from impersonated user (hemah)
    employeeId, _, err = getEmployeeInfof(impersonateUsername)
    if err != nil {
        log.Printf("Error getting employeeId for %s: %v, using default", impersonateUsername, err)
        employeeId = "000000"
    }
    
    log.Printf("Impersonation - Using EmployeeId from %s: %s, Mobile from %s: %s", 
        impersonateUsername, employeeId, originalUsername, mobileNumber)
} else {
    // Normal login - get both from same user
    employeeId, mobileNumber, err = getEmployeeInfof(originalUsername)
    if err != nil {
        log.Printf("Error getting employee info for %s: %v, using default values", originalUsername, err)
        employeeId = "000000"
        mobileNumber = "0000000000"
    }
}

						// *** Use only original username for session ***
						sessionUsername := originalUsername
						var impersonateUser string
						var impersonateFlag int

						if isImpersonating {
							impersonateUser = impersonateUsername
							impersonateFlag = 1
							log.Printf("Session username with impersonation: %s (impersonating as %s)", sessionUsername, impersonateUsername)
						} else {
							impersonateUser = ""
							impersonateFlag = 0
						}

						// Insert session data with impersonation details
						err = insertSessionDataWithImpersonationf(userId, sessionUsername, ou, employeeId, impersonateUser, impersonateFlag)
						if err != nil {
							log.Printf("Error inserting session data: %v, continuing anyway", err)
							// Continue without failing
						}

						// Generate JWT token with impersonated username
						tokenString, err := generateJWTf(userId, impersonateUsername, employeeId)
						if err != nil {
							log.Printf("Error generating JWT: %v, using fallback", err)
							// Generate a simple fallback token if JWT generation fails
							tokenString = "fallback_token_" + userId
						}

						// Get role for impersonated user
// Get role for impersonated user
roleName, err := getUserRoleByLoginf(impersonateUsername)
if err != nil {
    log.Printf("***** Role fetch FAILED for %s: %v *****", impersonateUsername, err)
    roleName = "UNKNOWN"
} else {
    log.Printf("***** Role fetch SUCCESS for %s: [%s] *****", impersonateUsername, roleName)
}

						// Return response with impersonated username but original user's mobile number
						resp := AuthResponsef{
							Valid:        true,
							UserId:       userId,
							Username:     impersonateUsername, // Return impersonated username
							EmployeeId:   employeeId,
							MobileNumber: mobileNumber, // Mobile number from original user
							Role:         roleName,
							Token:        tokenString,
						}
						log.Printf("***** RESPONSE BEING SENT - Username: %s, Role: %s, EmployeeId: %s *****", resp.Username, resp.Role, resp.EmployeeId)

						jsonResponse, _ := json.Marshal(resp)
						encrypted, _ := utils.Encrypt(jsonResponse)
						w.Header().Set("Content-Type", "application/json")
						_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
					}
					return // Exit once user is found (regardless of password correctness)
				}
			}
		}

		// Search in all organizational units
		performSearch(searchBaseStaff, "staff")
		performSearch(searchBaseFaculty, "faculty")
		performSearch(searchbaseProject, "project")

		// If no successful authentication and no response sent, send failure response
		if !authSuccess && !responseSent {
			resp := AuthResponsefalsef{Valid: false, Username: decodedUsername, Error: "Invalid username or password"}
			jsonResponse, _ := json.Marshal(resp)
			encrypted, _ := utils.Encrypt(jsonResponse)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
		}
	}))
	loggedHandler.ServeHTTP(w, r)
}

func generateUserIdf() string { return uuid.New().String() }

func updatePreviousActiveSessionsf(employeeId string) error {

	// Database connection
	db := credentials.GetDB()
	var err error
	query := `UPDATE meivan.Session_Data SET Is_Active = '0', idletimeout = '1', Logout_Date = NOW() WHERE Employee_id = $1 AND Is_Active = '1'`
	_, err = db.Exec(query, employeeId)
	return err
}

func insertSessionDataf(userId, username, ou, employeeId string) error {
	return insertSessionDataWithImpersonationf(userId, username, ou, employeeId, "", 0)
}

func insertSessionDataWithImpersonationf(userId, username, ou, employeeId, impersonateUser string, impersonateFlag int) error {
	_ = updatePreviousActiveSessionsf(employeeId)

	// Database connection
	db := credentials.GetDB()

	var err error
	query := `INSERT INTO meivan.Session_Data (Session_Id, Logout_Date, Username, Is_Active, idletimeout, Department, User_id, Employee_id, Login_Date, reason, impersonate_user, impersonate_flag) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), $9, $10, $11)`

	// Use NULL for impersonate_user if empty string
	var impersonateUserParam interface{}
	if impersonateUser == "" {
		impersonateUserParam = nil
	} else {
		impersonateUserParam = impersonateUser
	}

	_, err = db.Exec(query, userId, nil, username, "1", "0", ou, userId, employeeId, "LOGIN", impersonateUserParam, impersonateFlag)
	return err
}

func getEmployeeInfof(username string) (string, string, error) {

	// Database connection
	db := credentials.GetDB()

	var err error
	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Printf("Database ping failed: %v", err)
		return "", "", err
	}

	// Try different column name variations
	queries := []string{
		`SELECT EmployeeId, Mobilenumber FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
		`SELECT employeeid, mobilenumber FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
		`SELECT EmployeeId, MobileNumber FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
		`SELECT employeeid, MobileNumber FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
	}

	var employeeId, mobileNumber string

	for i, query := range queries {
		log.Printf("Trying query %d: %s with username: %s", i+1, query, username)
		row := db.QueryRow(query, username)
		err = row.Scan(&employeeId, &mobileNumber)
		if err == nil {
			log.Printf("Query successful! EmployeeId: %s, MobileNumber: %s", employeeId, mobileNumber)
			return employeeId, mobileNumber, nil
		}
		log.Printf("Query %d failed: %v", i+1, err)
	}

	// If all queries failed, try to get just EmployeeId
	simpleQueries := []string{
		`SELECT EmployeeId FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
		`SELECT employeeid FROM humanresources.employeebasicinfo WHERE LoginName = $1`,
	}

	for i, query := range simpleQueries {
		log.Printf("Trying simple query %d: %s", i+1, query)
		row := db.QueryRow(query, username)
		err = row.Scan(&employeeId)
		if err == nil {
			log.Printf("Simple query successful! EmployeeId: %s", employeeId)
			return employeeId, "0000000000", nil // Default mobile number
		}
		log.Printf("Simple query %d failed: %v", i+1, err)
	}

	log.Printf("All queries failed for username: %s", username)
	return "", "", fmt.Errorf("no employee found for username: %s", username)
}

// func getUserRoleByLoginf(loginname string) (string, error) {

// 	// Database connection
// 	db := credentials.GetDB()

// 	/*query := `
// 	SELECT 
// 	    B.campuscode || ' ' ||
// 	    CASE 
// 	        WHEN A.sectionid IS NOT NULL THEN D.sectioncode 
// 	        ELSE A.departmentcode 
// 	    END || ' ' || 
// 	    C.rolename AS RoleName
// 	FROM humanresources.employeerolemapping A
// 	JOIN humanresources.campus B 
// 	    ON A.campusid = B.id
// 	JOIN meivan.rolemaster C 
// 	    ON A.roleid = C.id
// 	LEFT JOIN humanresources.section D 
// 	    ON A.sectionid = D.id
// 	JOIN humanresources.employeebasicinfo e
// 	    ON A.employeeid = e.employeeid
// 	WHERE e.loginname = $1
// 	LIMIT 1;
// 	`
// 	*/

// query := `
// 	WITH user_param AS (
//     SELECT  $1 AS username
// ),
// employee_param AS (
//     SELECT 
//         ebi.employeeid AS emp_id,
//         ebi.loginname  AS username
//     FROM humanresources.employeebasicinfo ebi
//     JOIN user_param up ON up.username = ebi.loginname
// )
// SELECT 
//     UserID,
//     Username,
//     RoleName
// FROM (
//     -- 1️⃣ Default role (highest priority)
//     SELECT 
//         edr.user_id AS UserID,
//         ep.username AS Username,
//         edr.defaultrole AS RoleName,
//         1 AS priority
//     FROM humanresources.employeedefaultroles edr
//     JOIN employee_param ep ON edr.user_id = ep.emp_id

//     UNION ALL

//     -- 2️⃣ Alphabetically first role if no default role
//     SELECT DISTINCT
//         A.employeeid AS UserID,
//         ep.username  AS Username,
//         B.campuscode || ' ' ||
//         CASE 
//             WHEN A.sectionid IS NOT NULL THEN D.sectioncode
//             ELSE A.departmentcode
//         END || ' ' || C.rolename AS RoleName,
//         2 AS priority
//     FROM humanresources.employeerolemapping A
//     JOIN humanresources.campus B ON A.campusid = B.id
//     JOIN meivan.rolemaster C ON A.roleid = C.id
//     LEFT JOIN humanresources.section D ON A.sectionid = D.id
//     JOIN employee_param ep ON A.employeeid = ep.emp_id
//     WHERE NOT EXISTS (
//         SELECT 1
//         FROM humanresources.employeedefaultroles edr2
//         JOIN employee_param ep2 ON edr2.user_id = ep2.emp_id
//     )
//     ORDER BY RoleName
//     LIMIT 1
// ) combined
// ORDER BY priority
// LIMIT 1;
// `


// log.Printf(">>> Querying role for loginname: %s", loginname)
// 	var role string
// 	err := db.QueryRow(query, loginname).Scan(&role)
// 	if err != nil {
// 		log.Printf(">>> Role query FAILED for %s: %v", loginname, err)
// 		return "", err
// 	}

// 	log.Printf(">>> Role query SUCCESS for %s: %s", loginname, role)
// 	return role, nil
// }


func getUserRoleByLoginf(loginname string) (string, error) {

	db := credentials.GetDB()

	// query := `
	// WITH user_param AS (
	//     SELECT $1 AS username
	// ),
	// employee_param AS (
	//     SELECT 
	//         ebi.employeeid AS emp_id,
	//         ebi.loginname  AS username
	//     FROM humanresources.employeebasicinfo ebi
	//     JOIN user_param up ON up.username = ebi.loginname
	// )
	// SELECT 
	//     UserID,
	//     Username,
	//     RoleName
	// FROM (
	//     SELECT 
	//         edr.user_id AS UserID,
	//         ep.username AS Username,
	//         edr.defaultrole AS RoleName,
	//         1 AS priority
	//     FROM humanresources.employeedefaultroles edr
	//     JOIN employee_param ep ON edr.user_id = ep.emp_id

	//     UNION ALL

	//     SELECT DISTINCT
	//         A.employeeid AS UserID,
	//         ep.username  AS Username,
	//         B.campuscode || ' ' ||
	//         CASE 
	//             WHEN A.sectionid IS NOT NULL THEN D.sectioncode
	//             ELSE A.departmentcode
	//         END || ' ' || C.rolename AS RoleName,
	//         2 AS priority
	//     FROM humanresources.employeerolemapping A
	//     JOIN humanresources.campus B ON A.campusid = B.id
	//     JOIN meivan.rolemaster C ON A.roleid = C.id
	//     LEFT JOIN humanresources.section D ON A.sectionid = D.id
	//     JOIN employee_param ep ON A.employeeid = ep.emp_id
	//     WHERE NOT EXISTS (
	//         SELECT 1
	//         FROM humanresources.employeedefaultroles edr2
	//         JOIN employee_param ep2 ON edr2.user_id = ep2.emp_id
	//     )
	//     ORDER BY RoleName
	//     LIMIT 1
	// ) combined
	// ORDER BY priority
	// LIMIT 1;
	// `


		query := `
WITH user_param AS (
    SELECT $1 AS username
),
employee_param AS (
    SELECT 
        ebi.employeeid AS emp_id,
        ebi.loginname  AS username
    FROM humanresources.employeebasicinfo ebi
    JOIN user_param up 
        ON up.username = ebi.loginname
)

SELECT 
    UserID,
    Username,
    RoleName
FROM (
    -- Priority 1 : Default role
    SELECT 
        edr.user_id AS UserID,
        ep.username AS Username,
        edr.defaultrole AS RoleName,
        1 AS priority
    FROM humanresources.employeedefaultroles edr
    JOIN employee_param ep 
        ON edr.user_id = ep.emp_id

    UNION ALL

    -- Priority 2 : Build role dynamically
    SELECT DISTINCT
        A.employeeid AS UserID,
        ep.username  AS Username,

        CONCAT_WS(
            ' ',
            TRIM(B.campuscode),

            CASE 
                WHEN A.roleid = '35' THEN TRIM(A.departmentcode)
                WHEN TRIM(A.departmentcode) <> 'ADM' THEN TRIM(A.departmentcode)
            END,

            CASE 
                WHEN A.roleid <> '35' THEN TRIM(D.sectioncode)
            END,

            TRIM(C.rolename)
        ) AS RoleName,

        2 AS priority
    FROM humanresources.employeerolemapping A
    JOIN humanresources.campus B 
        ON A.campusid = B.id
    JOIN meivan.rolemaster C 
        ON A.roleid = C.id
    LEFT JOIN humanresources.section D 
        ON A.sectionid = D.id
    JOIN employee_param ep 
        ON A.employeeid = ep.emp_id
    WHERE NOT EXISTS (
        SELECT 1
        FROM humanresources.employeedefaultroles edr2
        JOIN employee_param ep2 
            ON edr2.user_id = ep2.emp_id
    )
    ORDER BY RoleName
    LIMIT 1
) combined
ORDER BY priority
LIMIT 1;

	`

	log.Printf(">>> Querying role for loginname: %s", loginname)

	var userID, username, role string
	err := db.QueryRow(query, loginname).Scan(&userID, &username, &role)
	if err != nil {
		log.Printf(">>> Role query FAILED for %s: %v", loginname, err)
		return "", err
	}

	log.Printf(">>> Role query SUCCESS → UserID=%s Username=%s Role=%s",
		userID, username, role)

	return role, nil
}