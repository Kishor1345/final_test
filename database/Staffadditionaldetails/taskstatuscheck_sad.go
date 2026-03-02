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
	"log"

	_ "github.com/lib/pq"
)

// =============================
// DB ACCESS FUNCTION
// =============================

// GetEmployeeTaskStatus fetches the list of employee IDs from Postgres
func GetEmployeeTaskStatus(employeeID string) ([]modelssad.TaskStatusCheck, error) {

	var results []modelssad.TaskStatusCheck

	if employeeID == "" {
		return results, fmt.Errorf("employeeid cannot be empty")
	}

	// Database connection
	db := credentials.GetDB()

	if err := db.Ping(); err != nil {
		return results, fmt.Errorf("database ping failed: %v", err)
	}

	log.Printf("Checking task status for employee: %s", employeeID)

	rows, err := db.Query(modelssad.TaskStatusCheckQuery, employeeID)
	if err != nil {
		return results, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	results, err = modelssad.RetrieveTaskStatusChecks(rows)
	if err != nil {
		return results, fmt.Errorf("error retrieving task status: %v", err)
	}

	log.Printf("Task status check completed for employee: %s", employeeID)
	return results, nil
}
