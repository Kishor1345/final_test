// Package databasequartersmasterestate handles DB access for Building Master Dropdown API.
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
	"strings"
)

// GetBuildingMasterDropdownFromDB fetches building master dropdown data from DB
// `decryptedData` should contain "quarter_category_ids" as comma-separated string, e.g. "1,2,3"
func GetBuildingMasterDropdownFromDB(decryptedData map[string]interface{}) ([]modelsquartersmasterestate.BuildingMasterDropdownStruct, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Read comma-separated IDs from decryptedData
	idsStr, ok := decryptedData["quarter_category_ids"].(string)
	if !ok || idsStr == "" {
		return nil, 0, fmt.Errorf("quarter_category_ids not provided")
	}

	// Convert string to slice of IDs
	idStrArr := strings.Split(idsStr, ",")
	ids := make([]interface{}, len(idStrArr))
	for i, v := range idStrArr {
		ids[i] = strings.TrimSpace(v)
	}

	// Generate Postgres placeholders: $1,$2,$3
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	query := fmt.Sprintf(modelsquartersmasterestate.MyQueryBuildingMasterDropdown, strings.Join(placeholders, ","))

	// Execute query
	rows, err := db.Query(query, ids...)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Retrieve results
	data, err := modelsquartersmasterestate.RetrieveBuildingMasterDropdown(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
