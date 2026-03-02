// Package databasecircular contains data structures and database access logic for Quarters Preference Details.
//
// Path : /var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/database/hr_008
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:23/02/2026
package databaseQuarters

import (
	credentials "Hrmodule/dbconfig"
	modelQuarters "Hrmodule/models/HR_009"
	"fmt"

	_ "github.com/lib/pq"
)

func DatabaseQuartersPreference(decryptedData map[string]interface{}) ([]modelQuarters.QuartersPreferenceFetchStructure, int, error) {

	// Extract task_id from decrypted data
	OrderNo, ok := decryptedData["order_no"].(string)
	if !ok || OrderNo == "" {
		return nil, 0, fmt.Errorf("OrderNo is required")
	}
	// Database connection
	db := credentials.GetDB()

	// Execute main criteria master data fetch query
	rows, err := db.Query(modelQuarters.MyQueryForQuartersPreference,OrderNo)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Map query result to response structure
	data, err := modelQuarters.RetrieveQuartersPreferenceDetails(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	// Ensure empty slice is returned instead of nil
	if len(data) == 0 {
		data = []modelQuarters.QuartersPreferenceFetchStructure{} // return empty slice instead of nil
	}
	// Return result and record count
	return data, len(data), nil
}