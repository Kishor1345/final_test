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

// GetEmployeeContactDetails fetches the list of employee IDs from Postgres
func GetEmployeeContactDetails(employeeID string) ([]modelssad.EmployeeContactDetails, error) {

	var result []modelssad.EmployeeContactDetails

	if employeeID == "" {
		return result, fmt.Errorf("employeeid cannot be empty")
	}

	// Database connection
	db := credentials.GetDB()

	rows, err := db.Query(modelssad.EmployeeContactDetailsSP, employeeID)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	return modelssad.RetrieveEmployeeContactDetails(rows)
}
