// Package databasecriteria contains data structures and database access logic for existing data.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/criteria
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:23/01/2026
package databasecriteria

import (
	credentials "Hrmodule/dbconfig"
	modelscriteria "Hrmodule/models/criteria"
	"fmt"

	_ "github.com/lib/pq"
)

// CriteriaMasterExistingDataFetchdatabase fetches Criteria Master data from the database
func CriteriaMasterExistingDataFetchdatabase(decryptedData map[string]interface{}) ([]modelscriteria.CriteriaMasterExistingDataFetchStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Check for ongoing process before fetching data
	count, err := modelscriteria.RetrieveOngoingCount(db)
	if err != nil {
		return nil, 0, fmt.Errorf("count query failed: %v", err)
	}

	// If process is ongoing, return message only
	if count > 0 {
		return []modelscriteria.CriteriaMasterExistingDataFetchStructure{
			{
				ProcessMsg: "Process On Going",
				Criteria:   []modelscriteria.CriteriaForExisting{},
			},
		}, 0, nil
	}

	// Execute main criteria master data fetch query
	rows, err := db.Query(modelscriteria.MyQueryCriteriaMasterExistingDataFetch)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Map query result to response structure
	data, err := modelscriteria.RetrieveCriteriaMasterExistingDataFetch(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}
	// Ensure empty slice is returned instead of nil
	if len(data) == 0 {
		data = []modelscriteria.CriteriaMasterExistingDataFetchStructure{} // return empty slice instead of nil
	}
	// Return result and record count
	return data, len(data), nil
}
