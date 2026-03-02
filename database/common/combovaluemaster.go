// Package databasesad contains structs and queries for Staff Additional details API.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all Combo Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetComboMaster fetches the list of combo values from Postgres based on comboname
func GetComboMaster(comboName string) ([]modelssad.ComboMaster, error) {
	var result []modelssad.ComboMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.ComboMasterQuery, comboName)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveComboMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving combo master data: %v", err)
	}

	result = records
	return result, nil
}
