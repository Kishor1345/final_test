// Package quartersmasterestate contains structs and queries for Quarter Category Dropdown API.
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

// SQL Query for Quarter Category Dropdown
var MyQueryQuarterCategoryDropdown = (`
SELECT id,name
FROM humanresources.quarterscategory
where campus_id = $1
ORDER BY name;
`)

// Struct for Quarter Category Dropdown
type QuarterCategoryDropdownStruct struct {
	ID       int    `json:"id"`       // previously "value", now "Id"
	Name     string `json:"name"` // previously "label", now "Category"
}

// RetrieveQuarterCategoryDropdown retrieves quarter category dropdown data
func RetrieveQuarterCategoryDropdown(rows *sql.Rows) ([]QuarterCategoryDropdownStruct, error) {
	var list []QuarterCategoryDropdownStruct

	for rows.Next() {
		var s QuarterCategoryDropdownStruct

		err := rows.Scan(
			&s.ID,
			&s.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning quarter category dropdown: %v", err)
		}

		list = append(list, s)
	}

	return list, nil
}
