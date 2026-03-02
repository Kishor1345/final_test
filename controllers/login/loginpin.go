//
//Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/login
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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"Hrmodule/auth"

	"github.com/joho/godotenv"
)

/*
	============================
	  INIT

============================
*/
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
}

/*
	============================
	  REQUEST / RESPONSE

============================
*/
type Requestkey struct {
	KeyCode string `json:"key_code"`
	Token   string `json:"token"` // ADD THIS
}

type Responsekey struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

/*
	============================
	  ALIAS → ENV KEY

============================
*/
func findEnvKeyByAlias(alias string) string {

	aliasParts := strings.Split(strings.ToUpper(strings.TrimSpace(alias)), "_")

	for _, env := range os.Environ() {
		key := strings.SplitN(env, "=", 2)[0]
		keyUpper := strings.ToUpper(key)

		keyParts := strings.Split(keyUpper, "_")

		// alias parts must not be more than env parts
		if len(aliasParts) > len(keyParts) {
			continue
		}

		match := true
		for i, ap := range aliasParts {
			if !strings.HasPrefix(keyParts[i], ap) {
				match = false
				break
			}
		}

		if match {
			return key
		}
	}

	return ""
}

/*
	============================
	  AES DECRYPT

============================
*/
func decryptAESforkey(enc string) (string, error) {

	key := os.Getenv("Login_Key")
	if len(key) != 32 {
		return "", errors.New("Login_Key must be 32 bytes")
	}

	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}

	block, _ := aes.NewCipher([]byte(key))
	gcm, _ := cipher.NewGCM(block)

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid encrypted data")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}

/*
	============================
	  AES ENCRYPT

============================
*/
func encryptAESforkey(plain string) (string, error) {

	key := os.Getenv("Login_Key")
	if len(key) != 32 {
		return "", errors.New("Login_Key must be 32 bytes")
	}

	block, _ := aes.NewCipher([]byte(key))
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	cipherText := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Getkey(w http.ResponseWriter, r *http.Request) {

	// Only POST
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	// Decode request FIRST (token is in body)
	var req Requestkey
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Token from POST body
	if strings.TrimSpace(req.Token) == "" {
		http.Error(w, "Token missing", http.StatusUnauthorized)
		return
	}

	// Inject token into header (auth framework expects header)
	r.Header.Set("token", req.Token)

	// Auth
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Log + handler
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validate token again (framework style)
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN provided", http.StatusBadRequest)
			return
		}

		// Resolve env key
		envKey := findEnvKeyByAlias(req.KeyCode)
		if envKey == "" {
			http.Error(w, "Invalid key_code", http.StatusNotFound)
			return
		}

		encryptedEnvValue := os.Getenv(envKey)
		if encryptedEnvValue == "" {
			http.Error(w, "Encrypted value not found", http.StatusNotFound)
			return
		}

		// Decrypt stored value
		decryptedValue := encryptedEnvValue


		// Re-encrypt for response
		responseEncrypted, err := encryptAESforkey(decryptedValue)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"status": "success",
			"value":  responseEncrypted,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	})).ServeHTTP(w, r)
}
