// Package databasecommon contains structs and queries for City.
// path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/common
// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to fetch the all City Value Master.
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

// GetCityMaster fetches the list of City values from Postgres
func GetCityMaster(countryCode string, stateID int) ([]modelssad.CityMaster, error) {
	var result []modelssad.CityMaster

	// Database connection
	db := credentials.GetDB()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.CityMasterQuery, countryCode, stateID)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveCityMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving City master data: %v", err)
	}

	result = records
	return result, nil
}
