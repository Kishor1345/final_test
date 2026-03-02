// Package controllerscommon contains APIs for fetching employee signatures.
//
// These APIs retrieve signature details from the employee_signatures_new table.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/common
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 27-10-2025
package controllerscommon

import (
	credentials "Hrmodule/dbconfig"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var EncryptionKeyforsignature string

// -------------------- INIT FUNCTION --------------------
func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Fetch encryption key from .env
	EncryptionKeyforsignature = os.Getenv("EncryptionKeyforsignature")
	if EncryptionKeyforsignature == "" {
		log.Println("ERROR: EncryptionKeyforsignature missing in .env")
	}

}

func DownloadSignatureHandler(w http.ResponseWriter, r *http.Request) {

	db := credentials.GetDB() // ✅ safe, pooled

	employeeid := r.URL.Query().Get("role")
	if employeeid == "" {
		http.Error(w, "employeeid is required", http.StatusBadRequest)
		return
	}

	var decrypted []byte
	query := `
		SELECT public.pgp_sym_decrypt_bytea(signature, $1)
		FROM meivan.employee_signatures_new
		WHERE role=$2
		ORDER BY createdon DESC
		LIMIT 1
	`

	err := db.QueryRow(query, EncryptionKeyforsignature, employeeid).Scan(&decrypted)
	if err != nil {
		http.Error(w, "Decryption failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=signature.png")
	w.Header().Set("Content-Type", "image/png")
	w.Write(decrypted)
}
