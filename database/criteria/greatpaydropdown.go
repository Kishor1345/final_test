// Package databasecriteria contains data structures and database access logic for greatpaydropdown.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/criteria
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:08/01/2026
package databasecriteria

import (
	credentials "Hrmodule/dbconfig"
	modelscriteria "Hrmodule/models/criteria"
	"fmt"

	_ "github.com/lib/pq"
)

// Greatpaydatabase fetches Criteria Master data from the database
func Greatpaydatabase(decryptedData map[string]interface{}) ([]modelscriteria.GreatpayDetailsStructure, int, error) {

	// Extract cpc from decrypted data
	cpc, ok := decryptedData["cpc"].(string)
	if !ok || cpc == "" {
		return nil, 0, fmt.Errorf("Cpc is required")
	}

	// Database connection
	db := credentials.GetDB()

	// Execute criteria master data fetch query
	rows, err := db.Query(modelscriteria.MyQueryGreatpayDropdown, cpc)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Map query result to response structure
	data, err := modelscriteria.RetrieveGreatpayDropdown(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	// Return result and record count
	return data, len(data), nil
}
