// Package utils provides cryptographic and database utility functions.
// It specifically handles AES-GCM encryption and decryption for secure data
// exchange and manages session-based key retrieval from the database.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/utils
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 07-07-2025
// Last Modified By: Sridharan
// Last Modified Date: 09-07-2025
package utils

import (
	credentials "Hrmodule/dbconfig"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// GetDecryptKey retrieves a 32-character AES key from the meivan.session_data table
// using the provided session ID (pid). It returns an error if the key is not found,
// inactive, or is not exactly 32 bytes.
func GetDecryptKey(pid string) (string, error) {
	db := credentials.GetDB() // ✅ safe: called after main() init

	var key string
	err := db.QueryRow(
		`SELECT decryptkey FROM meivan.session_data WHERE session_id=$1 AND is_active=1`,
		pid,
	).Scan(&key)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid or inactive id")
		}
		return "", err
	}

	if len(key) != 32 {
		return "", fmt.Errorf("decrypt key must be 32 characters long (got %d)", len(key))
	}

	return key, nil
}

// EncryptAES encrypts a plaintext string using AES-GCM with the provided 32-byte key.
// It generates a unique nonce for every operation and prepends it to the ciphertext
// before base64 encoding the final result.
func EncryptAES(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// EncryptAES encrypts a plaintext string using AES-GCM with the provided 32-byte key.
// It generates a unique nonce for every operation and prepends it to the ciphertext
// before base64 encoding the final result.
func DecryptAES(ciphertextB64, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, encrypted := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// EncryptResponse represents the JSON structure for responses containing
// encrypted data payloads.
type EncryptResponse struct {
	Data string `json:"Data"`
}

// DecryptRequest represents the JSON structure for incoming requests
// containing encrypted data to be processed
type DecryptRequest struct {
	Data string `json:"Data"`
}

// EncryptHandler is an HTTP handler that encrypts a JSON payload.
// It expects a "P_id" in the JSON to identify the session, retrieves the
// session's key, encrypts the remaining data, and returns the ciphertext
// combined with the P_id.
func EncryptHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request as generic JSON
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	pidVal, ok := req["P_id"]
	if !ok {
		http.Error(w, "Missing 'P_id' field", http.StatusBadRequest)
		return
	}

	pid := fmt.Sprintf("%v", pidVal)
	key, err := GetDecryptKey(pid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Key fetch failed: %v", err), http.StatusUnauthorized)
		return
	}

	// Remove P_id before encryption
	delete(req, "P_id")

	// Convert the remaining payload to JSON
	payloadBytes, _ := json.Marshal(req)

	// Encrypt payload using DB key
	encrypted, err := EncryptAES(string(payloadBytes), key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp := EncryptResponse{
		Data: fmt.Sprintf("%s||%s", pid, encrypted),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DecryptHandler is an HTTP handler that decrypts data received in the format "PID||Ciphertext".
// It extracts the PID, retrieves the decryption key from the database, performs
// the decryption, and returns the original JSON payload.
func DecryptHandler(w http.ResponseWriter, r *http.Request) {
	var req DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	key, err := GetDecryptKey(pid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Key fetch failed: %v", err), http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Decryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &payload); err != nil {
		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
		return
	}

	payload["P_id"] = pid

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
