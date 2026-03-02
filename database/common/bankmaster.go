// Package databasecommon contains structs and queries for Bank.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all Bank Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetBankMaster fetches the list of Bank values from Postgres
func GetBankMaster() ([]modelssad.BankMaster, error) {
	var result []modelssad.BankMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.BankMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveBankMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving Bank master data: %v", err)
	}

	result = records
	return result, nil
}
