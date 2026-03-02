// Package databasecommon contains structs and queries for castecagetory.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all castecagetory Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetCasteCategoryMaster fetches the list of caste category values from Postgres
func GetCasteCategoryMaster() ([]modelssad.CasteCategoryMaster, error) {
	var result []modelssad.CasteCategoryMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.CasteCategoryMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveCasteCategoryMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving caste category master data: %v", err)
	}

	result = records
	return result, nil
}
