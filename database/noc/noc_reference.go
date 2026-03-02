// Package databasenoc contains data structures and database access logic for the NOC Reference page.
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
	"log"

	_ "github.com/lib/pq"
)

// NocReferenceDatabase fetches NOC reference order numbers for a given employee
func NocReferenceDatabase(decryptedData map[string]interface{}) ([]modelsnoc.NocReferenceStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract employeeid from request data
	employeeIDRaw, ok := decryptedData["employeeid"]
	if !ok {
		return nil, 0, fmt.Errorf("missing 'employeeid' in request data")
	}

	employeeID, ok := employeeIDRaw.(string)
	if !ok || employeeID == "" {
		return nil, 0, fmt.Errorf("employeeid must be a non-empty string")
	}

	// Ensure employeeID contains only digits (preserve leading zeros)
	for _, ch := range employeeID {
		if ch < '0' || ch > '9' {
			return nil, 0, fmt.Errorf("employeeid must contain only digits")
		}
	}

	// Log the value being used for DB query
	log.Printf("Querying NOC for employee_id: %s", employeeID)

	// Execute query safely
	rows, err := db.Query(modelsnoc.MyQueryNocReference, employeeID)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map rows to struct
	data, err := modelsnoc.RetrieveNocReference(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Return data, count, and nil error
	return data, len(data), nil
}
