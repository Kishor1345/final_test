// Package databasecommon contains structs and queries for DepartmentMaster.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all DepartmentMaster Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetDepartmentMaster fetches the list of DepartmentMaster values from Postgres
func GetDepartmentMaster() ([]modelssad.DepartmentMaster, error) {
	var result []modelssad.DepartmentMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.DepartmentMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveDepartmentMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving DepartmentMaster data: %v", err)
	}

	result = records
	return result, nil
}
