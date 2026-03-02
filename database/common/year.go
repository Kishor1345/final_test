// Package databasecommon contains structs and queries for YearMaster.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all YearMaster Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetYearMaster fetches the list of YearMaster values from Postgres
func GetYearMaster() ([]modelssad.YearMaster, error) {
	var result []modelssad.YearMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.YearMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveYearMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving YearMaster data: %v", err)
	}

	result = records
	return result, nil
}
