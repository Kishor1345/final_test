// Package modelssad contains structs and queries for Role Based Module mapping.
//
// Path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/models/Staffadditionaldetails
// --- Creator's Info ---
// Creator: Rovita
//
// Created On: 29-01-2026
//
// Last Modified By:
//  
// Last Modified Date:
 package modelssad

import (
	"database/sql"
	"fmt"
	"log"
)

// RoleBasedModuleQuery fetches modules based on role_name with parameterized query
const RoleBasedModuleQuery = `
SELECT module_name
FROM meivan.category_role_map a 
JOIN meivan.category_visibility b 
    ON b.id = a.module_id
WHERE role_name = $1 and a.status='1'
ORDER BY module_name
`

// RoleBasedModule represents a module accessible by a specific role
type RoleBasedModule struct {
	ModuleName string `json:"module_name"`
}

// RetrieveRoleBasedModules scans module names from query results
func RetrieveRoleBasedModules(rows *sql.Rows) ([]RoleBasedModule, error) {
	var modules []RoleBasedModule
	var recordCount int

	for rows.Next() {
		var module RoleBasedModule
		if err := rows.Scan(&module.ModuleName); err != nil {
			return nil, fmt.Errorf("error scanning module row: %v", err)
		}
		modules = append(modules, module)
		recordCount++
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	log.Printf("Successfully scanned %d module records", recordCount)
	return modules, nil
}