// Package utils provides cryptographic utility functions, specifically focusing on
// session-based AES-GCM encryption for secure data handling.
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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	credentials "Hrmodule/dbconfig"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// dbAESKey stores the AES key retrieved from the database for the current session.
var dbAESKey []byte

// LoadKeyFromDB retrieves an active AES decryption key from the database for a
// specific session ID.
//
// It connects to the Meivan database, queries the session_data table for an
// active record matching the sessionID, and stores the resulting 32-byte key
// in the package-level variable for use by encryption functions.
//
// Returns an error if the database connection fails, the key is not found,
// or the key length is not exactly 32 bytes.
func LoadKeyFromDB(sessionID string) error {
	// Database connection
	db := credentials.GetDB()

	query := `
		SELECT decryptkey
		FROM session_data
		WHERE is_active = 1 AND session_id = $1
	`

	var key string
	var err error
	err = db.QueryRow(query, sessionID).Scan(&key)
	if err != nil {
		return fmt.Errorf("failed to fetch decryptkey: %w", err)
	}
	// Print the fetched key (for debugging, remove in production)
	fmt.Printf("Fetched AES Key for session %s: %s\n", sessionID, key)

	if len(key) != 32 {
		return fmt.Errorf("decryptkey must be 32 bytes, got %d", len(key))
	}

	dbAESKey = []byte(key)
	return nil
}

// Encryptnew encrypts a plaintext byte slice using the previously loaded AES key.
//
// It uses the AES-GCM (Galois/Counter Mode) encryption algorithm, generating
// a unique random nonce for each operation. The resulting string is a
// base64-encoded representation of the concatenated nonce and ciphertext.
//
// LoadKeyFromDB must be called successfully before using this function.
//
// Returns the base64-encoded encrypted string or an error if the key is not
// loaded or the encryption process fails.
func Encryptnew(plainText []byte) (string, error) {
	if dbAESKey == nil {
		return "", errors.New("AES key is not loaded")
	}

	block, err := aes.NewCipher(dbAESKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
