// Package databasecircular contains data structures and database access logic for Quarters NUmber fetch.
//
// Path : /var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/database/hr_008
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:16/02/2026
package databasecircular

import (
	credentials "Hrmodule/dbconfig"
	modelcircular "Hrmodule/models/HR_008"
	"fmt"

	_ "github.com/lib/pq"
)

func DatabaseQuartersNumberFetch(decryptedData map[string]interface{}) ([]modelcircular.QuartersNumberDataStructure, int, error) {

	// Extract task_id from decrypted data
	campusFloat, ok := decryptedData["campus_id"].(float64)
	if !ok || campusFloat == 0 {
	return nil, 0, fmt.Errorf("CampusId is required")
	}
	CampusId := int(campusFloat)

	categoryFloat, ok := decryptedData["category_id"].(float64)
	if !ok || categoryFloat == 0 {
	return nil, 0, fmt.Errorf("CategoryId is required")
	}
	CategoryId := int(categoryFloat)
	// Database connection
	db := credentials.GetDB()

	// Execute main criteria master data fetch query
	rows, err := db.Query(modelcircular.MyQueryForQuartersNumberData,CampusId,CategoryId)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Map query result to response structure
	data, err := modelcircular.RetrieveQuartersNumberDataFetch(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	// Ensure empty slice is returned instead of nil
	if len(data) == 0 {
		data = []modelcircular.QuartersNumberDataStructure{} // return empty slice instead of nil
	}
	// Return result and record count
	return data, len(data), nil
}