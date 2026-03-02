// Package databasequartersmasterestate handles DB access for Estate Quarters Number Dropdown API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/quartersmasterestate
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 12-01-2026
//
// Last Modified By: Ramya M R
//
// Last Modified Date: 24-01-2026
package databasequartersmasterestate

import (
	credentials "Hrmodule/dbconfig"
	modelsquartersmasterestate "Hrmodule/models/quartersmasterestate"
	"fmt"
	"strings"
)

// GetEstateQuartersNumberDropdownFromDB fetches estate quarters number dropdown data from DB
// decryptedData should contain:
// 1. "category_id"    -> number (decoded as float64)
// 2. "building_ids"   -> comma-separated string (optional)
// Returns only records where Quarters_Number is not null
func GetEstateQuartersNumberDropdownFromDB(
	decryptedData map[string]interface{},
) ([]modelsquartersmasterestate.EstateQuartersNumberDropdownStruct, int, error) {

	// Database connection
	db := credentials.GetDB()

	//  Read category_id (JSON numbers are float64)
	categoryVal, exists := decryptedData["category_id"]
	if !exists {
		return nil, 0, fmt.Errorf("category_id not provided")
	}

	categoryFloat, ok := categoryVal.(float64)
	if !ok {
		return nil, 0, fmt.Errorf("category_id invalid type")
	}

	categoryID := int(categoryFloat)

	// Read building_ids (optional, comma-separated, trim-safe)
	var buildingIDs *string
	if val, exists := decryptedData["building_ids"]; exists {
		if strVal, ok := val.(string); ok {
			trimmed := strings.TrimSpace(strVal)
			if trimmed != "" {
				buildingIDs = &trimmed
			}
		}
	}

	// Execute query
	rows, err := db.Query(
		modelsquartersmasterestate.MyQueryEstateQuartersNumberDropdown,
		categoryID,
		buildingIDs,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Retrieve results (null quarters numbers already filtered)
	data, err := modelsquartersmasterestate.RetrieveEstateQuartersNumberDropdown(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
