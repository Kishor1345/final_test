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

// FetchEmployeeEducationDetails connects to DB and retrieves employee education details
func FetchEmployeeEducationDetails(employeeID string) (interface{}, error) {

	// Database connection
	db := credentials.GetDB()

	var err error
	// Verify DB connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("DB connection failed: %v", err)
	}

	// Use the generic retriever function
	data, err := modelssad.GenericEducationDetailsRetriever(db, employeeID)
	if err != nil {
		return nil, fmt.Errorf("retrieving Education details failed: %v", err)
	}

	return data, nil
}
