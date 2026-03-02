// Package databasesad interacts with the Staff Additional Details DB.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/Staffadditionaldetails
// --- Creator's Info ---
// Creator: kishorekumar
// Created On: 29-01-2026
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/Staffadditionaldetails"
	"fmt"

	_ "github.com/lib/pq"
)

// FetchEmployeeDocumentDetails fetches the list of employee IDs from Postgres
func FetchEmployeeDocumentDetails(employeeID string) (interface{}, error) {

	// Database connection
	db := credentials.GetDB()

	// Use the generic retriever function
	data, err := modelssad.GenericDocumentDetailsRetriever(db, employeeID)
	if err != nil {
		return nil, fmt.Errorf("retrieving document details failed: %v", err)
	}

	return data, nil
}
