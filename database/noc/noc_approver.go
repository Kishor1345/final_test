// Package databasenoc contains data structures and database access logic for the NOC Approver page.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/noc
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 07-01-2026
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

// NocApproverDatabase fetches NOC approver details from the database
func NocApproverDatabase(decryptedData map[string]interface{}) ([]modelsnoc.NocApproverStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract task_id from request data
	taskID, ok := decryptedData["task_id"].(string)
	if !ok || taskID == "" {
		return nil, 0, fmt.Errorf("missing 'task_id' in request data")
	}

	// Execute query
	rows, err := db.Query(modelsnoc.MyQueryNocApprover, taskID)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map rows to struct
	data, err := modelsnoc.RetrieveNocApprover(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Return data, count, and nil error
	return data, len(data), nil
}
