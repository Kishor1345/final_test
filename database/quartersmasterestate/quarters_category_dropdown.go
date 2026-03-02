// Package databasequartersmasterestate handles DB access for Quarter Category Dropdown API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/quartersmasterestate
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 12-01-2026
//
// Last Modified By:
//
// Last Modified Date:
package databasequartersmasterestate

import (
	credentials "Hrmodule/dbconfig"
	modelsquartersmasterestate "Hrmodule/models/quartersmasterestate"
	"fmt"
)

// GetQuarterCategoryDropdownFromDB fetches quarter category dropdown data from DB
func GetQuarterCategoryDropdownFromDB(decryptedData map[string]interface{}) ([]modelsquartersmasterestate.QuarterCategoryDropdownStruct, int, error) {

	campusValue, ok := decryptedData["campus_id"]
	if !ok {
    	return nil, 0, fmt.Errorf("CampusId is required")
	}
	// Database connection
	db := credentials.GetDB()

	rows, err := db.Query(modelsquartersmasterestate.MyQueryQuarterCategoryDropdown,campusValue)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsquartersmasterestate.RetrieveQuarterCategoryDropdown(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
