// // Package controllerssad handles API logic for Staff Additional Details.
// //path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/controllers/Staffadditionaldetails
// // --- Creator's Info ---
// // Creator: Rovita
// // Created On: 11-11-2025
// // Last Modified By:
// // Last Modified Date:
// // Description: API to insert or update active employee personal basic details.
package controllerssad
import (
	"Hrmodule/auth"
	databasesad "Hrmodule/database/Staffadditionaldetails"
	"Hrmodule/micro"
	modelssad "Hrmodule/models/Staffadditionaldetails"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

/* ===============================
   Wrapper + Response
   =============================== */

type MasterSadRequestWrapper struct {
	Data string `json:"Data"`
}

type APIResponseMasterSad struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

//ONLY action from frontend
type SadActionWrapper struct {
	ActionType string `json:"action_type"` // saveasdraft | submit
}

/* ===============================
   MAIN API
   =============================== */


//SadBasicDetails handles the Bank Master API request.
func SadBasicDetails(w http.ResponseWriter, r *http.Request) {

	//Method check
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var wrapper MasterSadRequestWrapper
	if err := json.Unmarshal(body, &wrapper); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//Split PID || encrypted
	parts := strings.Split(wrapper.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	//Decrypt
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Invalid PID", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	//Parse model
	var req modelssad.MasterSadRequest
	if err := json.Unmarshal([]byte(decryptedJSON), &req); err != nil {
		http.Error(w, "Invalid payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	//Read ONLY action
	var act SadActionWrapper
	if err := json.Unmarshal([]byte(decryptedJSON), &act); err != nil {
		http.Error(w, "Invalid action_type", http.StatusBadRequest)
		return
	}

	switch strings.ToLower(strings.TrimSpace(act.ActionType)) {
	case "saveasdraft":
		req.ActionType = "draft"
	case "submit":
		req.ActionType = "submit"
	default:
		http.Error(w, "Invalid action_type value", http.StatusBadRequest)
		return
	}

	// Ensure ProcessID
	if req.ProcessID == 0 {
		if pidInt, err := strconv.Atoi(pid); err == nil {
			req.ProcessID = pidInt
		}
	}

	//Auth
	token := extractToken(decryptedJSON)
	r.Header.Set("token", token)

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}
	if err := auth.IsValidIDFromRequest(r); err != nil {
		http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
		return
	}

	/* =====================================================
	   WORKFLOW ONLY FOR SUBMIT
	   ===================================================== */

	if strings.EqualFold(req.ActionType, "submit") {

		//Ignore frontend workflow OUTPUT fields
		req.AssignTo = nil
		req.AssignedRole = nil
		req.EmailFlag = nil
		req.TemplateID = nil
		req.RejectFlag = nil
		req.RejectRole = nil

		//Conditions from frontend
		conditions := map[string]string{
			"EmployeeGroup": req.EmployeeGroup,
			"PBM":           req.PBM,
		}

		input := micro.NextActivityInput{
			Token:           token,
			ProcessID:       req.ProcessID,
			CurrentActivity: derefInt(req.ActivitySeqNo),
			TaskID:          req.TaskID,
			IsApproved:      derefInt(req.IsTaskApproved),
			IsReturn:        derefInt(req.IsTaskReturn),
			Conditions:      conditions,
		}

		out, err := micro.NextActivity(input, pid, key)
		if err != nil {
			sendSadError(w, pid, key, err.Error())
			return
		}

		//Apply workflow outputs
		req.ActivitySeqNo = &out.NextActivitySeq
		req.AssignTo = &out.AssignTo
		req.AssignedRole = &out.AssignedRole
		req.EmailFlag = &out.EmailFlag
		req.RejectFlag = &out.RejectFlag

		if out.TemplateID != nil {
			req.TemplateID = out.TemplateID
		}

		//Email only for submit
		if out.EmailFlag > 0 && out.TemplateID != nil && req.TaskID != nil {
			_ = micro.EmailQueue(
				token,
				req.ProcessID,
				*req.TaskID,
				*out.TemplateID,
				pid,
				key,
			)
		}
	}

	/* =====================================================
	   DB CALL UNCHANGED
	   ===================================================== */

	result, err := databasesad.ExecuteMasterSad(req)
	if err != nil {
		sendSadError(w, pid, key, err.Error())
		return
	}

	resp := APIResponseMasterSad{
		Status:  200,
		Message: "SAD processed successfully",
		Data:    result,
	}

	sendSadSuccess(w, pid, key, resp)
}

/* ===============================
   HELPERS
   =============================== */

func extractToken(jsonStr string) string {
	var m map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &m)
	if t, ok := m["token"].(string); ok {
		return t
	}
	return ""
}

func derefInt(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

func sendSadError(w http.ResponseWriter, pid, key, msg string) {
	resp := APIResponseMasterSad{
		Status:  400,
		Message: msg,
		Data:    nil,
	}
	j, _ := json.Marshal(resp)
	enc, _ := utils.EncryptAES(string(j), key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Data": fmt.Sprintf("%s||%s", pid, enc),
	})
}

func sendSadSuccess(w http.ResponseWriter, pid, key string, payload APIResponseMasterSad) {
	j, _ := json.Marshal(payload)
	enc, _ := utils.EncryptAES(string(j), key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Data": fmt.Sprintf("%s||%s", pid, enc),
	})
}
