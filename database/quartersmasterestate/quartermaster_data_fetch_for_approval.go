// Package databaseqms contains database connection and query execution data fetch for approval.
//
// // Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:22/01/2026
package databasequartersmasterestate

import (
	credentials "Hrmodule/dbconfig"
	modelsquartersmasterestate "Hrmodule/models/quartersmasterestate"
	"fmt"
)

// GetQuartersmasterFromDB fetches Quarters Master data from the database
func GetQuartersmasterFromDB(decryptedData map[string]interface{}) ([]modelsquartersmasterestate.QuarterMasterDataFetchForApproval, int, error) {

	// Extract task_id from decrypted data
	taskID, ok := decryptedData["task_id"].(string)
	if !ok || taskID == "" {
		return nil, 0, fmt.Errorf("task_id is required")
	}

	// Database connection
	db := credentials.GetDB()

	// Execute main Quarters master data fetch query
	rows, err := db.Query(modelsquartersmasterestate.MyQueryQMSApproval, taskID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Map query result to response structure
	data, err := modelsquartersmasterestate.RetrieveCriteriaMasterDataFetchForApproval(rows)
	if err != nil {
		return nil, 0, err
	}

	// Return result and record count
	return data, len(data), nil
}
