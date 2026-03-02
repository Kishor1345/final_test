// Package databasenoc contains data structures and database access logic for the NOC Template page.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/noc
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 08-01-2026
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

// NocTemplateDatabase fetches NOC template details
func NocTemplateDatabase(decryptedData map[string]interface{}) ([]modelsnoc.NocTemplateStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	employeeID, ok := decryptedData["employeeid"].(string)
	if !ok || employeeID == "" {
		return nil, 0, fmt.Errorf("missing 'employeeid' in request data")
	}

	// FIXED PART
	processIDFloat, ok := decryptedData["processid"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("missing 'processid' in request data")
	}
	processID := int(processIDFloat)

	templateType, ok := decryptedData["templatetype"].(string)
	if !ok || templateType == "" {
		return nil, 0, fmt.Errorf("missing 'templatetype' in request data")
	}

	taskID, ok := decryptedData["taskid"].(string)
	if !ok || taskID == "" {
		return nil, 0, fmt.Errorf("missing 'taskid' in request data")
	}

	rows, err := db.Query(
		modelsnoc.MyQueryNocTemplate,
		employeeID,
		processID,
		templateType,
		taskID,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	data, err := modelsnoc.RetrieveNocTemplate(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return data, len(data), nil
}
