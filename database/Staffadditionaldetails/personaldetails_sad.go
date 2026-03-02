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
	"strings"

	_ "github.com/lib/pq"
)

// FetchEmployeePersonalDetails_sad fetches the list of employee IDs from Postgres
func FetchEmployeePersonalDetails_sad(employeeID string, reqType string) (interface{}, error) {

	// Database connection
	db := credentials.GetDB()

	switch strings.ToLower(reqType) {
	case "new":
		return modelssad.GenericPersonalDetailsRetriever_sad(db, employeeID)

	case "continue":
		return modelssad.ContinuePersonalDetailsRetriever_sad(db, employeeID)

	default:
		return nil, fmt.Errorf("Invalid Type value: %s", reqType)
	}
}
