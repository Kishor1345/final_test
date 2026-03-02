// Package quartersmasterestate contains structs and queries for Estate Quarters Number Dropdown API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Ramya M R
// Created On: 12-01-2026
//
// Last Modified By:  
//
// Last Modified Date: 24-01-2026
package quartersmasterestate

import (
	"database/sql"
	"fmt"
)

// SQL Query for Estate Quarters Number Dropdown
// Supports:
// 1. Single category selection
// 2. Multi-select building IDs (comma-separated)
// Modified to filter out null displayname values
var MyQueryEstateQuartersNumberDropdown = (`
SELECT DISTINCT
    qm.id,
    qm.displayname
FROM humanresources.quarterscategory qc
JOIN humanresources.buildingmaster bm
    ON bm.quarters_category = qc.id
JOIN humanresources.quartersmaster qm
    ON qm.building_id = bm.id
WHERE qc.id = $1
  AND qm.displayname IS NOT NULL
  AND (
        $2::text IS NULL
        OR bm.id = ANY (string_to_array($2::text, ',')::bigint[])
      )
ORDER BY qm.displayname;
`)

// Struct for Estate Quarters Number Dropdown
type EstateQuartersNumberDropdownStruct struct {
	Quarters_Id     *int    `json:"Quarters_Id"`     // dropdown value
	Quarters_Number *string `json:"Quarters_Number"` // dropdown label
}

// RetrieveEstateQuartersNumberDropdown retrieves quarters number dropdown data
// Now filters out records with null Quarters_Number
func RetrieveEstateQuartersNumberDropdown(rows *sql.Rows) ([]EstateQuartersNumberDropdownStruct, error) {
	var list []EstateQuartersNumberDropdownStruct

	for rows.Next() {
		var s EstateQuartersNumberDropdownStruct

		err := rows.Scan(
			&s.Quarters_Id,
			&s.Quarters_Number,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning estate quarters number dropdown: %v", err)
		}

		// Only append if Quarters_Number is not nil
		if s.Quarters_Number != nil && *s.Quarters_Number != "" {
			list = append(list, s)
		}
	}

	return list, nil
}