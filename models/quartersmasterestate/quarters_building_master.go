// Package quartersmasterestate contains structs and queries for Building Master Dropdown API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Ramya M R
// Created On: 12-01-2026
//
// Last Modified By:
//
// Last Modified Date:
package quartersmasterestate

import (
	"database/sql"
	"fmt"
)

// SQL Query for Building Master Dropdown
// Use placeholder %s for filtering by multiple Quarter Category IDs (Postgres style: $1,$2,...)
var MyQueryBuildingMasterDropdown = (` 
SELECT 
    bm.id,
    bm.building_name
FROM humanresources.buildingmaster bm
JOIN humanresources.quarterscategory qc
    ON bm.quarters_category = qc.id
WHERE qc.id IN (%s)
ORDER BY bm.building_name;
`)

// Struct for Building Master Dropdown
type BuildingMasterDropdownStruct struct {
	Building_ID     *int    `json:"Building_ID"`      // previously "Category"
	Building_Number *string `json:"Building_Number"`  // previously "label", now "Building_Number"
}

// RetrieveBuildingMasterDropdown retrieves building master dropdown data
func RetrieveBuildingMasterDropdown(rows *sql.Rows) ([]BuildingMasterDropdownStruct, error) {
	var list []BuildingMasterDropdownStruct

	for rows.Next() {
		var s BuildingMasterDropdownStruct

		err := rows.Scan(
			&s.Building_ID,
			&s.Building_Number,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning building master dropdown: %v", err)
		}

		list = append(list, s)
	}

	return list, nil
}
