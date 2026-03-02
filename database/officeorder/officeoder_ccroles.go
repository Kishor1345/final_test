// Package databaseofficeorder handles database access for the CC Roles Master.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 20-11-2025
// Last Modified By: Sridharan
// Last Modified Date: 20-11-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"fmt"

	_ "github.com/lib/pq"
)

// GetCcRolesFromDB retrieves Carbon Copy (CC) roles for a specific employee from the database.
//
// It extracts the "employee_id" from the decryptedData map, establishes a connection to
// the Meivan database, executes the CC roles query, and returns a slice of CcRoleStruct,
// the count of records found, and an error if the operation fails.
func GetCcRolesFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.CcRoleStruct, int, error) {

	// Database connection
	db := credentials.GetDB()

	employeeID, ok := decryptedData["employee_id"].(string)
	if !ok || employeeID == "" {
		return nil, 0, fmt.Errorf("missing 'employee_id' in request data")
	}

	rows, err := db.Query(modelsofficeorder.MyQueryCcRoles, employeeID)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsofficeorder.RetrieveCcRoles(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
