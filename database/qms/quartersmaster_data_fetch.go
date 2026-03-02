// Package databaseqms handles database operations for the Quarters Management System (QMS).
// It provides functionality to retrieve specific quarters details based on user display names.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/qms
//
// --- Creator's Info ---
// Creator: Elakiya
// Created On:
// Last Modified By:
// Last Modified Date:
package databaseqms

import (
	credentials "Hrmodule/dbconfig"
	models "Hrmodule/models/qms"
	"fmt"
)

// GetQMSEUDetailsFromDB retrieves Estate Unit (EU) details for a specific task within the Quarters Management System.
//
// It extracts the "task_id" from the provided decryptedData map, establishes a connection to the
// Meivan database, and executes a query to fetch the EU-specific information associated with the task.
//
// Returns:
//   - A slice of QMSEUDetailsStruct containing the retrieved information.
//   - An integer representing the total count of records found.
//   - An error if the task_id is missing, the DB connection fails, or the query execution fails.
func GetQMSEUDetailsFromDB(
	decryptedData map[string]interface{},
) ([]models.QMSEUDetailsStruct, int, error) {

	taskID, ok := decryptedData["task_id"].(string)
	if !ok || taskID == "" {
		return nil, 0, fmt.Errorf("task_id is required")
	}

	// Database connection
	db := credentials.GetDB()

	rows, err := db.Query(models.MyQueryQMSEUDetails, taskID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	data, err := models.RetrieveQMSEUDetails(rows)
	if err != nil {
		return nil, 0, err
	}

	return data, len(data), nil
}
