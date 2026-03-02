// Package databaseofficeorder handles database operations for tracking and counting office orders.
// It provides functionality to interact with the WF_officeorder database.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 15-09-2025
// Last Modified By: Ramya
// Last Modified Date: 30-09-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"fmt"

	_ "github.com/lib/pq"
)

// GetCombinedNeedGenerate retrieves the office order generation counts from the PostgreSQL database.
//
// It establishes a connection to the office order database, executes a predefined query
// to count records requiring generation, and maps the first result to the Postgres
// field of the CombinedNeedGenerate structure.
//
// Returns a CombinedNeedGenerate struct containing the counts and an error if any part
// of the database operation fails.
func GetCombinedNeedGenerate() (modelsofficeorder.CombinedNeedGenerate, error) {
	var result modelsofficeorder.CombinedNeedGenerate

	// Database connection
	db := credentials.GetDB()

	rows, err := db.Query(modelsofficeorder.MyQueryNeedGeneratePostgres)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelsofficeorder.RetrieveNeedGeneratePostgres(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving Postgres data: %v", err)
	}

	if len(records) > 0 {
		result.Postgres = records[0]
	}

	return result, nil
}
