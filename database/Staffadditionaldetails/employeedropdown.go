// Package modelssad contains structs and queries for Staff Additional details API.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/Staffadditionaldetails
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all active employees.
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/Staffadditionaldetails"
	"fmt"

	_ "github.com/lib/pq"
)

// GetEmployeeDropdown fetches the list of employee IDs from Postgres
func GetEmployeeDropdown() ([]modelssad.EmployeeDropdown, error) {
	var result []modelssad.EmployeeDropdown

	// Database connection
	db := credentials.GetDB()

	rows, err := db.Query(modelssad.MyQueryEmployeeDropdown)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveEmployeeDropdown(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving employee dropdown data: %v", err)
	}

	result = records
	return result, nil
}
