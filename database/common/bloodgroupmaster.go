// Package databasecommon contains structs and queries for BloodGroup.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all BloodGroup Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetBloodGroupMaster fetches the list of BloodGroup values from Postgres
func GetBloodGroupMaster() ([]modelssad.BloodGroupMaster, error) {
	var result []modelssad.BloodGroupMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.BloodGroupMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveBloodGroupMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving BloodGroup master data: %v", err)
	}

	result = records
	return result, nil
}
