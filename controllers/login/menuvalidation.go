package controllerslogin

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)
type MenuValidationTerminalLog struct {
	PID      string `json:"P_id"`
	Path     string `json:"path"`
	RoleName string `json:"role_name"`
	Token    string `json:"token"`
	Username string `json:"username"`
}

// Incoming encrypted request for menu validation
type MenuValidationEncryptedRequest struct {
	Data string `json:"Data"` // pid||encryptedPayload
}

// Decrypted payload structure
type MenuValidationPayload struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	RoleName string `json:"role_name"`
	Path     string `json:"path"`
}

// API response
type MenuValidationResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	HasRole int    `json:"has_role"` // 1 or 0
}

// Validate user role and menu path access from DB
func ValidateUserRole(username string, sessionID string, roleName string, path string) (int, error) {

	// Database connection
	db := credentials.GetDB()

	// Query to check if user has the role with active session and menu path access
	// query := `
	// 	SELECT 
	// 		CASE 
	// 			WHEN COUNT(*) > 0 THEN 1
	// 			ELSE 0
	// 		END AS has_requested_role
	// 	FROM (
	// 		SELECT DISTINCT
	// 			A.Employeeid AS UserID,
	// 			e.loginname AS Username,
	// 			B.campuscode || ' ' ||
	// 			CASE 
	// 				WHEN A.sectionid IS NOT NULL THEN D.sectioncode 
	// 				ELSE A.departmentcode 
	// 			END || ' ' || 
	// 			TRIM(C.rolename) AS RoleName
	// 		FROM humanresources.employeerolemapping A
	// 		JOIN humanresources.campus B ON A.campusid = B.id
	// 		JOIN meivan.rolemaster C ON A.roleid = C.id
	// 		LEFT JOIN humanresources.section D ON A.sectionid = D.id
	// 		JOIN humanresources.employeebasicinfo e ON A.employeeid = e.employeeid
	// 		JOIN meivan.session_data s ON e.loginname = s.username
	// 		WHERE e.loginname = $1
	// 		  AND s.session_id = $2
	// 		  AND s.is_active = 1
	// 	) subquery
	// 	JOIN sirion.menumaster m ON subquery.RoleName = m.role_names
	// 	WHERE subquery.RoleName = $3
	// 	  AND m.path = $4`


// 		query := `
// 		WITH user_roles AS (
//     SELECT DISTINCT
//         B.campuscode || ' ' ||
//         CASE
//             WHEN A.sectionid IS NOT NULL THEN D.sectioncode
//             ELSE A.departmentcode
//         END || ' ' || TRIM(C.rolename) AS default_role_name,
//         eb.loginname
//     FROM humanresources.employeerolemapping A
//     JOIN humanresources.campus B ON A.campusid = B.id
//     JOIN meivan.rolemaster C ON A.roleid = C.id
//     LEFT JOIN humanresources.section D ON A.sectionid = D.id
//     JOIN humanresources.employeebasicinfo eb ON A.employeeid = eb.employeeid
// ),
// valid_session AS (
//     SELECT 1
//     FROM meivan.session_data s
//     WHERE s.username = $1
//       AND s.session_id = $2
//       AND s.is_active = 1
// )
// SELECT
//     CASE
//         WHEN EXISTS (
//             SELECT 1
//             FROM sirion.component_master ct
//             WHERE ct.path = $4
//               AND ct.status = 1
//               AND EXISTS (SELECT 1 FROM valid_session)
//               AND (
//                     (
//                         ct.component_type IN ('Initiator', 'Additionaldetails','Dashboard','TaskSummary','Inbox')
//                         AND EXISTS (
//                             SELECT 1
//                             FROM sirion.role_rights rr
//                             WHERE rr.component_id = ct.id
//                               AND rr.role_name = $3
//                               AND rr.status = 1
//                         )
//                     )
//                     OR
//                     (
//                         ct.component_type = 'Approver'
//                         AND EXISTS (
//                             SELECT 1
//                             FROM user_roles ur
//                             JOIN sirion.role_rights rr
//                                 ON rr.role_name = ur.default_role_name
//                                AND rr.component_id = ct.id
//                                AND rr.status = 1
//                             WHERE ur.loginname = $1
//                         )
//                     )
//                 )
//         )
//         THEN 1
//         ELSE 0
//     END AS has_requested_role;
// `

	query := `
WITH user_roles AS (
    SELECT DISTINCT
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
        ) AS default_role_name,
        eb.loginname
    FROM humanresources.employeerolemapping A
    JOIN humanresources.campus B 
        ON A.campusid = B.id
    JOIN meivan.rolemaster C 
        ON A.roleid = C.id
    LEFT JOIN humanresources.section D 
        ON A.sectionid = D.id
    JOIN humanresources.employeebasicinfo eb 
        ON A.employeeid = eb.employeeid
),

valid_session AS (
    SELECT 1
    FROM meivan.session_data s
    WHERE s.username = $1
      AND s.session_id = $2
      AND s.is_active = 1
)

SELECT
    CASE
        WHEN EXISTS (
            SELECT 1
            FROM sirion.component_master ct
            WHERE ct.path = $4
              AND ct.status = 1
              AND EXISTS (SELECT 1 FROM valid_session)
              AND (
                    (
                        ct.component_type IN
                        ('Initiator','Additionaldetails','Dashboard','TaskSummary','Inbox')
                        AND EXISTS (
                            SELECT 1
                            FROM sirion.role_rights rr
                            WHERE rr.component_id = ct.id
                              AND rr.role_name = $3
                              AND rr.status = 1
                        )
                    )
                    OR
                    (
                        ct.component_type = 'Approver'
                        AND EXISTS (
                            SELECT 1
                            FROM user_roles ur
                            JOIN sirion.role_rights rr
                                ON rr.role_name = ur.default_role_name
                               AND rr.component_id = ct.id
                               AND rr.status = 1
                            WHERE ur.loginname = $1
                        )
                    )
                )
        )
        THEN 1
        ELSE 0
    END AS has_requested_role;


`

	var hasRole int
	err := db.QueryRow(query, username, sessionID, roleName, path).Scan(&hasRole)
	if err != nil {
		return 0, fmt.Errorf("query error: %v", err)
	}

	return hasRole, nil
}

func MenuValidationHandler(w http.ResponseWriter, r *http.Request) {

	// Only POST allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read full body
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// Parse JSON for Data
	var encReq MenuValidationEncryptedRequest
	if err := json.Unmarshal(rawBody, &encReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Split: pid || encryptedPayload
	parts := strings.Split(encReq.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedInput := parts[1]

	// Fetch AES Key from DB using P_ID
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Failed to fetch decryption key", http.StatusUnauthorized)
		return
	}

	// Decrypt AES → JSON
	decryptedJSON, err := utils.DecryptAES(encryptedInput, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Parse decrypted JSON
	var payload MenuValidationPayload
	if err := json.Unmarshal([]byte(decryptedJSON), &payload); err != nil {
		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
		return
	}
// ---- PRINT REQUEST IN TERMINAL ----
reqLog := MenuValidationTerminalLog{
	PID:      pid,
	Path:     payload.Path,
	RoleName: payload.RoleName,
	Token:    payload.Token,
	Username: payload.Username,
}

reqJSON, _ := json.MarshalIndent(reqLog, "", "  ")
fmt.Println("\n================ MENU VALIDATION REQUEST ================")
fmt.Println(string(reqJSON))
fmt.Println("=========================================================\n")

	// Validate required fields
	if payload.Token == "" || payload.Username == "" || payload.RoleName == "" || payload.Path == "" {
		http.Error(w, "Missing required fields in decrypted payload", http.StatusBadRequest)
		return
	}

	// Inject token into header
	r.Header.Set("token", payload.Token)

	// Step: validate token + IP
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Wrap next stage
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validate token
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN provided", http.StatusBadRequest)
			return
		}

		// Validate user role and menu path - using pid as session_id
		hasRole, err := ValidateUserRole(payload.Username, pid, payload.RoleName, payload.Path)



		var apiResp MenuValidationResponse
		if err != nil {
			apiResp = MenuValidationResponse{
				Status:  500,
				Message: "Failed to validate role: " + err.Error(),
				HasRole: 0,
			}
		} else {
			if hasRole == 1 {
				apiResp = MenuValidationResponse{
					Status:  200,
					Message: fmt.Sprintf("User %s has the role: %s and access to path: %s", payload.Username, payload.RoleName, payload.Path),
					HasRole: 1,
				}
			} else {
				apiResp = MenuValidationResponse{
					Status:  200,
					Message: fmt.Sprintf("User %s does not have the role: %s or access to path: %s", payload.Username, payload.RoleName, payload.Path),
					HasRole: 0,
				}
			}
		}
// ---- PRINT RESPONSE IN TERMINAL ----
respJSON, _ := json.MarshalIndent(apiResp, "", "  ")
fmt.Println("**************** MENU VALIDATION RESPONSE ****************")
fmt.Println(string(respJSON))
fmt.Println("*********************************************************\n")
		// Encrypt response JSON
		jsonResp, _ := json.Marshal(apiResp)
		encryptedResp, _ := utils.EncryptAES(string(jsonResp), key)

		// Final output: pid||encryptedResponse
		finalResponse := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResp),
		}

		// Write JSON output
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(finalResponse)

	})).ServeHTTP(w, r)
}
