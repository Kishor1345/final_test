// Package databasecommon contains structs and queries for Designation.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all Designation Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetDesignationMaster fetches the list of Designation values from Postgres
func GetDesignationMaster() ([]modelssad.DesignationMaster, error) {
	var result []modelssad.DesignationMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.DesignationMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveDesignationMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving Designation master data: %v", err)
	}

	result = records
	return result, nil
}
