// Package databaseofficeorder handles database access for the ReturnDropdown API.
// It manages the retrieval of dropdown options used in the return process of office orders.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 29-10-2025
// Last Modified By: Sridharan
// Last Modified Date: 29-10-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"fmt"
)

// ReturnDropdownRequest represents the structure for a return dropdown request,
// specifically capturing the unique identifier for a task.
type ReturnDropdownRequest struct {
	TaskID string `json:"task_id"`
}

// GetReturnDropdownFromDB fetches dropdown data from the database based on a specific Task ID.
//
// It extracts the "task_id" from the provided decryptedData map, establishes a connection
// to the Meivan database, and executes a query to retrieve the relevant dropdown records.
//
// Returns:
//   - A slice of ReturnDropdownStruct containing the query results.
//   - An integer representing the number of records retrieved.
//   - An error if the database connection or query execution fails.
func GetReturnDropdownFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.ReturnDropdownStruct, int, error) {

	// Database connection
	db := credentials.GetDB()

	TaskID, _ := decryptedData["task_id"].(string)

	rows, err := db.Query(modelsofficeorder.MyQueryReturnDropdown, TaskID)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsofficeorder.RetrieveReturnDropdown(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
