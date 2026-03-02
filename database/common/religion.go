// Package databasecommon contains structs and queries for Religion.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all Religion Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetReligionMaster fetches the list of Religion values from Postgres
func GetReligionMaster() ([]modelssad.ReligionMaster, error) {
	var result []modelssad.ReligionMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.ReligionMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveReligionMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving Religion master data: %v", err)
	}

	result = records
	return result, nil
}
