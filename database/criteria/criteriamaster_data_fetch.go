// Package databasecriteria contains data structures and database access logic for criteria master data fetch.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/criteria
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:13/01/2026
package databasecriteria

import (
	credentials "Hrmodule/dbconfig"
	modelscriteria "Hrmodule/models/criteria"
	"fmt"

	_ "github.com/lib/pq"
)

// CriteriaMasterDataFetchdatabase fetches Criteria Master data from the database
func CriteriaMasterDataFetchdatabase(decryptedData map[string]interface{}) ([]modelscriteria.CriteriaMasterDataFetchStructure, int, int, error) {

	// Extract task_id from decrypted data
	var Taskid *string
	if t, ok := decryptedData["task_id"].(string); ok && t != "" {
		Taskid = &t
	} else {
		Taskid = nil
	}
	// Database connection
	db := credentials.GetDB()

	// Execute main criteria master data fetch query
	rows, err := db.Query(modelscriteria.MyQueryCriteriaMasterDataFetch, Taskid)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Map query result to response structure
	data, err := modelscriteria.RetrieveCriteriaMasterDataFetch(rows)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	// If process is ongoing, return message only
	count, err := modelscriteria.RetrieveOngoingCountForDataFetch(db)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("count query failed: %v", err)
	}

	// Ensure empty slice is returned instead of nil
	if len(data) == 0 {
		data = []modelscriteria.CriteriaMasterDataFetchStructure{} // return empty slice instead of nil
	}
	// Return result and record count
	return data, len(data), count, nil
}
