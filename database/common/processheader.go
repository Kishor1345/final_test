// Package databasecommon handles DB calls for ProcessHeader API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/common
//
// --- Creator's Info ---
//
// Creator: Rovita
//
// Created On:30-12-2025
//
// Last Modified By: Rovita
//
// Last Modified Date: 30-12-2025
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"fmt"
	"strconv"
)

// ProcessHeaderDatabase executes the query
func ProcessHeaderDatabase(decryptedData map[string]interface{}) ([]modelscommon.ProcessHeader, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract id from decrypted data
	// Handle both string and numeric types
	var processID int
	var err error
	switch v := decryptedData["id"].(type) {
	case string:
		processID, err = strconv.Atoi(v)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid 'id' format: %v", err)
		}
	case float64:
		processID = int(v)
	case int:
		processID = v
	default:
		return nil, 0, fmt.Errorf("missing or invalid 'id' in request data")
	}

	// Execute query
	rows, err := db.Query(modelscommon.MyQueryProcessHeader, processID)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying DB: %v", err)
	}
	defer rows.Close()

	// Map results
	data, err := modelscommon.RetrieveProcessHeader(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return data, len(data), nil
}
