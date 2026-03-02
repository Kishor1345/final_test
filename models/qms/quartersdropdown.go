// Package modelsquarters contains structs and queries for Quarters Number Dropdown API.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/qms
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 06-12-2025
//
// Last Modified By:
//
// Last Modified Date:
package modelsqms

import (
	"database/sql"
	"fmt"
)

// SQL Query for Quarters Number Dropdown
var MyQueryQuartersDropdown = (`
SELECT DISTINCT
    displayname as quartersnumber
FROM humanresources.quartersmaster
ORDER BY quartersnumber;
`)

// Struct for Quarters Dropdown
type QuartersDropdownStruct struct {
	QuartersNumber *string `json:"quartersnumber"`
}

// Function to retrieve quarters dropdown data
func RetrieveQuartersDropdown(rows *sql.Rows) ([]QuartersDropdownStruct, error) {
	var list []QuartersDropdownStruct

	for rows.Next() {
		var s QuartersDropdownStruct

		err := rows.Scan(
			&s.QuartersNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning quarters dropdown: %v", err)
		}

		list = append(list, s)
	}

	return list, nil
}
