// Package databaseofficeorder handles database operations for the Office Order module,
// focusing on task details and summary information.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 21-11-2025
// Last Modified By: Sridharan
// Last Modified Date: 21-11-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"fmt"

	_ "github.com/lib/pq"
)

// GetPCRTaskDetailsFromDB retrieves specific task details based on a process ID and task ID.
//
// It extracts "process_id" and "task_id" from the decryptedData map, establishes a
// connection to the database, and executes the PCR task details query.
//
// Returns:
//   - A slice of PCRTaskDetailsStruct containing the retrieved task information.
//   - An integer representing the number of records found.
//   - An error if mandatory fields are missing or if the database operation fails.
func GetPCRTaskDetailsFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.PCRTaskDetailsStruct, int, error) {

	// Database connection
	db := credentials.GetDB()

	processID := fmt.Sprintf("%v", decryptedData["process_id"])
	taskID := fmt.Sprintf("%v", decryptedData["task_id"])

	if processID == "" || taskID == "" {
		return nil, 0, fmt.Errorf("process_id and task_id are required")
	}

	rows, err := db.Query(modelsofficeorder.MyQueryPCRTaskDetails, processID, taskID)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	result, err := modelsofficeorder.RetrievePCRTaskDetails(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return result, len(result), nil
}
