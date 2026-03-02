// Package databasecriteria contains data structures and database access logic for floor dropdown.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:22/01/2026
package databasequartersmasterestate

import (
	credentials "Hrmodule/dbconfig"
	modelsquartersmasterestate "Hrmodule/models/quartersmasterestate"
	"fmt"

	_ "github.com/lib/pq"
)

// QmesFloordatabase fetches Quarters Master data from the database
func QmesFloordatabase(decryptedData map[string]interface{}) ([]modelsquartersmasterestate.FloorDetailsStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Execute main Quarters master data fetch query
	rows, err := db.Query(modelsquartersmasterestate.MyQueryFloorDropdown)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Map query result to response structure
	data, err := modelsquartersmasterestate.RetrieveFloorDropdown(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	// Return result and record count
	return data, len(data), nil
}
