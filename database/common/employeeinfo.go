// Package databasecommon contains data structures and database access logic for the EmployeeBasicInfo page.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/common
// --- Creator's Info ---
//
// Creator: Ramya M R
//
// Created On: 05-01-2026
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

// EmployeeBasicInfodatabase fetches employee basic info from the database
func EmployeeBasicInfodatabase(decryptedData map[string]interface{}) ([]modelscommon.EmployeeBasicInfoStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract employeeid from request data
	employeeID, ok := decryptedData["employeeid"].(string)
	if !ok || employeeID == "" {
		return nil, 0, fmt.Errorf("missing 'employeeid' in request data")
	}

	// Execute query
	rows, err := db.Query(modelscommon.MyQueryEmployeeBasicInfo, employeeID)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map rows to struct
	data, err := modelscommon.RetrieveEmployeeBasicInfo(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Return data, count, and nil error
	return data, len(data), nil
}
