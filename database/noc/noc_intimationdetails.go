// Package databasenoc contains data structures and database access logic for the NOC Intimation Details page.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/noc
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 12-01-2026
//
// Last Modified By:
//
// Last Modified Date:
package databasenoc

import (
	credentials "Hrmodule/dbconfig"
	modelsnoc "Hrmodule/models/noc"
	"fmt"

	_ "github.com/lib/pq"
)

// NocIntimationDetailsDatabase fetches NOC Intimation Details from the database
func NocIntimationDetailsDatabase(decryptedData map[string]interface{}) ([]modelsnoc.NocIntimationDetailsStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract order_no from request data
	orderNo, ok := decryptedData["order_no"].(string)
	if !ok || orderNo == "" {
		return nil, 0, fmt.Errorf("missing 'order_no' in request data")
	}

	// Execute query
	rows, err := db.Query(modelsnoc.MyQueryNocIntimationDetails, orderNo)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map rows to struct
	data, err := modelsnoc.RetrieveNocIntimationDetails(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Return data, count, and nil error
	return data, len(data), nil
}
