// Package databasecommon contains structs and queries for State.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all State Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetStateMaster fetches the list of State values from Postgres
func GetStateMaster(countryCode string) ([]modelssad.StateMaster, error) {
	var result []modelssad.StateMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.StateMasterQuery, countryCode)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveStateMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving State master data: %v", err)
	}

	result = records
	return result, nil
}
