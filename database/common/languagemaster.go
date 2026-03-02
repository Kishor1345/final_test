// Package databasecommon contains structs and queries for languagemaster.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all languagemaster.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetLanguageMaster fetches the list of language master values from Postgres
func GetLanguageMaster() ([]modelssad.LanguageMaster, error) {
	var result []modelssad.LanguageMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.LanguageMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveLanguageMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving language master data: %v", err)
	}

	result = records
	return result, nil
}
