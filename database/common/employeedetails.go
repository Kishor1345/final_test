// Package databasecommon provides shared database operations and utilities
// for retrieving HR-related information.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/common
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:04-12-2025
//
// Last Modified By:
//
// Last Modified Date:
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// EmployeeDetailsdatabase fetches detailed information for a specific employee from the database.
// It extracts the "employeeid" from the provided decryptedData map, establishes a database connection,
// executes the predefined query, and returns a slice of EmployeeDetailsStructure, the record count, and any error encountered.
func EmployeeDetailsdatabase(decryptedData map[string]interface{}) ([]modelscommon.EmployeeDetailsStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	employeeID, ok := decryptedData["employeeid"].(string)
	if !ok || employeeID == "" {
		return nil, 0, fmt.Errorf("missing 'employeeid' in request data")
	}

	rows, err := db.Query(modelscommon.MyQueryEmployeeDetails, employeeID)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	data, err := modelscommon.RetrieveEmployeeDetails(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return data, len(data), nil
}
