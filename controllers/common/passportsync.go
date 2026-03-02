// Package controllerscommon contains APIs for passport sync
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/controllers/common
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 27-10-2025
package controllerscommon

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	credentials "Hrmodule/dbconfig"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

/*
==========================

	RESPONSE STRUCT

==========================
*/
type SyncResponse struct {
	Status       string `json:"status"`
	UpdatedCount int    `json:"updated_count"`
	SkippedEmpty int    `json:"skipped_empty"`
	TotalRead    int    `json:"total_read"`
}

/*
==========================

	API HANDLER

==========================
*/
func SyncPassportHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	/* ==========================
	   MSSQL CONNECTION
	========================== */

	// Get MSSQL connection string
	mssqlConn := credentials.GetLivedatabase10IITM()

	msdb, err := sql.Open("sqlserver", mssqlConn)
	if err != nil {
		http.Error(w, "MSSQL connection failed", http.StatusInternalServerError)
		return
	}
	defer msdb.Close()

	/* ==========================
	   POSTGRESQL CONNECTION
	========================== */

	// Database connection
	db := credentials.GetDB()

	/* ==========================
	   MSSQL QUERY (LATEST VALID PASSPORT)
	========================== */

	rows, err := msdb.Query(`
		SELECT employeeid, Number
		FROM (
			SELECT employeeid, Number,
			       ROW_NUMBER() OVER (PARTITION BY employeeid ORDER BY ModifiedOn DESC) rn
			FROM IITM..EmployeePassportDetails
			WHERE Number IS NOT NULL
			  AND LTRIM(RTRIM(Number)) <> ''
		) x
		WHERE rn = 1
	`)
	if err != nil {
		http.Error(w, "MSSQL query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	updated := 0
	skipped := 0
	total := 0

	/* ==========================
	   UPDATE POSTGRESQL
	========================== */

	for rows.Next() {
		total++

		var empID string
		var passport sql.NullString

		if err := rows.Scan(&empID, &passport); err != nil {
			continue
		}

		// Extra safety
		if !passport.Valid || strings.TrimSpace(passport.String) == "" {
			skipped++
			continue
		}

		res, err := db.Exec(`
			UPDATE humanresources.employeebasicinfo
			SET passportnumber = $1
			WHERE employeeid = $2
			  AND (passportnumber IS NULL OR passportnumber = '')
		`, passport.String, empID)

		if err == nil {
			affected, _ := res.RowsAffected()
			if affected > 0 {
				updated++
			}
		}
	}

	/* ==========================
	   RESPONSE
	========================== */

	resp := SyncResponse{
		Status:       "success",
		UpdatedCount: updated,
		SkippedEmpty: skipped,
		TotalRead:    total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
