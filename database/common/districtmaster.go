// Package databasecommon contains structs and queries for District.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all District Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetDistrictMaster fetches the list of District values from Postgres
func GetDistrictMaster(countryCode string, stateID int) ([]modelssad.DistrictMaster, error) {
	var result []modelssad.DistrictMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.DistrictMasterQuery, countryCode, stateID)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveDistrictMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving District master data: %v", err)
	}

	result = records
	return result, nil
}
