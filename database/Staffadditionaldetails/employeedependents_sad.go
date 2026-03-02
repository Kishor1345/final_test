// Package databasesad interacts with the Staff Additional Details DB.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/Staffadditionaldetails
// --- Creator's Info ---
// Creator: kishorekumar
// Created On: 29-01-2026
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/Staffadditionaldetails"

	_ "github.com/lib/pq"
)

func GetEmployeeDependentDetails(employeeID string) (interface{}, error) {

	// Database connection
	db := credentials.GetDB()

	return modelssad.FetchEmployeeDependentDetailsFromDB(db, employeeID)
}
