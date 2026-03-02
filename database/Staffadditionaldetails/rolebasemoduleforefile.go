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

// GetRoleBasedModules fetches the list of modules based on role_name
func GetRoleBasedModules(roleName string) ([]modelssad.RoleBasedModule, error) {
	var modules []modelssad.RoleBasedModule

	// Validate input
	if roleName == "" {
		return modules, fmt.Errorf("role_name cannot be empty")
	}

	// Database connection
	db := credentials.GetDB()

	// Test database connection
	if err := db.Ping(); err != nil {
		return modules, fmt.Errorf("database ping failed: %v", err)
	}

	log.Printf("Fetching modules for role: %s", roleName)

	// Execute parameterized query
	rows, err := db.Query(modelssad.RoleBasedModuleQuery, roleName)
	if err != nil {
		return modules, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	// Retrieve and process results
	modules, err = modelssad.RetrieveRoleBasedModules(rows)
	if err != nil {
		return modules, fmt.Errorf("error retrieving modules: %v", err)
	}

	log.Printf("Successfully retrieved %d modules for role: %s", len(modules), roleName)
	return modules, nil
}
