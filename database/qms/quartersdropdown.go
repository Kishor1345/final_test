// Package databasequarters handles DB access for Quarters Dropdown API.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/qms
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 06-12-2025
//
// Last Modified By:
//
// Last Modified Date:
package databaseqms

import (
	credentials "Hrmodule/dbconfig"
	modelsquarters "Hrmodule/models/qms"
	"fmt"
)

// GetQuartersDropdownFromDB fetches quarters number dropdown data from DB
func GetQuartersDropdownFromDB(decryptedData map[string]interface{}) ([]modelsquarters.QuartersDropdownStruct, int, error) {

	// Database connection
	db := credentials.GetDB()

	rows, err := db.Query(modelsquarters.MyQueryQuartersDropdown)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsquarters.RetrieveQuartersDropdown(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
