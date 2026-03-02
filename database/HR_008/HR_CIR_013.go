// Package databasecircular contains data structures and database access logic for Circular data fetch.
//
// Path : /var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/database/hr_008
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:16/02/2026
package databasecircular

import (
	credentials "Hrmodule/dbconfig"
	modelcircular "Hrmodule/models/HR_008"
	"fmt"

	_ "github.com/lib/pq"
)

func DatabaseCircularDetailFetch(decryptedData map[string]interface{}) ([]modelcircular.CircularDetailStructure, int, error) {

	// Extract task_id from decrypted data
	OrderNo, ok := decryptedData["order_no"].(string)
	if !ok || OrderNo == "" {
		return nil, 0, fmt.Errorf("OrderNo is required")
	}
	// Database connection
	db := credentials.GetDB()

	// Execute main criteria master data fetch query
	rows, err := db.Query(modelcircular.MyQueryForCircularDataFetch,OrderNo)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Map query result to response structure
	data, err := modelcircular.RetrieveCircularDetailFetch(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	// Ensure empty slice is returned instead of nil
	if len(data) == 0 {
		data = []modelcircular.CircularDetailStructure{} // return empty slice instead of nil
	}
	// Return result and record count
	return data, len(data), nil
}