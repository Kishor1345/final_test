// Package commoncontrollers exposes API for ProcessHeader.
//
//Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/common
//
// --- Creator's Info ---
//
// Creator: Rovita
//
// Created On:30-12-2025
//
// Last Modified By: Rovita
//
// Last Modified Date: 30-12-2025
package controllerscommon

import (
    "Hrmodule/auth"
    database "Hrmodule/database/common"
    "Hrmodule/utils"
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
)

// APIResponseforProcessHeader standard response
type APIResponseforProcessHeader struct {
    Status  int         `json:"Status"`
    Message string      `json:"message"`
    Data    interface{} `json:"Data"`
}

// ProcessHeaderTokenRequest wrapper for encrypted data
type ProcessHeaderTokenRequest struct {
    Data string `json:"Data"`
}

// ProcessHeader API handler
func ProcessHeader(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
        return
    }

    // Read body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Unable to read request body", http.StatusBadRequest)
        return
    }
    r.Body = io.NopCloser(bytes.NewBuffer(body))

    // Extract encrypted data
    var req ProcessHeaderTokenRequest
    if err := json.Unmarshal(body, &req); err != nil {
        log.Printf("Error unmarshalling JSON: %v", err)
        http.Error(w, "Invalid JSON body", http.StatusBadRequest)
        return
    }

    // Decrypt Data and Extract Token
    parts := strings.Split(req.Data, "||")
    if len(parts) != 2 {
        log.Printf("Invalid data format")
        http.Error(w, "Invalid Data format", http.StatusBadRequest)
        return
    }

    pid := parts[0]
    encryptedPart := parts[1]

    // Get decryption key from database
    key, err := utils.GetDecryptKey(pid)
    if err != nil {
        log.Printf("Key fetch failed: %v", err)
        http.Error(w, "Decryption failed", http.StatusUnauthorized)
        return
    }

    // Decrypt the payload
    decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
    if err != nil {
        log.Printf("Decryption error: %v", err)
        http.Error(w, "Decryption failed", http.StatusUnauthorized)
        return
    }

    // Parse decrypted JSON
    var decryptedData map[string]interface{}
    if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
        log.Printf("Invalid decrypted JSON: %v", err)
        http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
        return
    }

    // Extract and validate token
    token, ok := decryptedData["token"].(string)
    if !ok || token == "" {
        log.Printf("Token not found in decrypted data")
        http.Error(w, "Token not found", http.StatusBadRequest)
        return
    }

    // Validate id parameter exists
    if _, exists := decryptedData["id"]; !exists {
        log.Printf("Missing 'id' parameter in request")
        http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
        return
    }

    // Set Token Header for Authentication
    r.Header.Set("token", token)

    // Authenticate
    if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
        return
    }

    // Log and process
    auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if err := auth.IsValidIDFromRequest(r); err != nil {
            http.Error(w, "Invalid TOKEN provided", http.StatusBadRequest)
            return
        }

        // Query database
        data, total, err := database.ProcessHeaderDatabase(decryptedData)
        if err != nil {
            log.Printf("Database error: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Build response
        response := APIResponseforProcessHeader{
            Status:  200,
            Message: "Success",
            Data: map[string]interface{}{
                "No Of Records": total,
                "Records":       data,
            },
        }

        responseJSON, err := json.MarshalIndent(response, "", "    ")
        if err != nil {
            http.Error(w, "Response marshal failed", http.StatusInternalServerError)
            return
        }

        // Encrypt response
        encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
        if err != nil {
            http.Error(w, "Response encryption failed", http.StatusInternalServerError)
            return
        }

        finalResp := map[string]string{
            "Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
        }

        // Save response log
        auth.SaveResponseLog(
            r,
            finalResp,
            http.StatusOK,
            "application/json",
            len(responseJSON),
            string(body),
        )

        // Send to client
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(finalResp)
    })).ServeHTTP(w, r)
}