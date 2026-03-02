// Package databasecircular contains data structures and database access logic for fetch Eligibility Choice.
//
// Path : /var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/database/hr_008
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:13/02/2026
package databasecircular

import (
	credentials "Hrmodule/dbconfig"
	modelcircular "Hrmodule/models/HR_008"
	"fmt"

	_ "github.com/lib/pq"
)

func CircularDataFetchForEligibilityChoice(decryptedData map[string]interface{}) ([]modelcircular.EligibilityChoiceDataStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Execute main criteria master data fetch query
	rows, err := db.Query(modelcircular.MyQueryForEligibilityChoiceData)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Map query result to response structure
	data, err := modelcircular.RetrieveCircularDataFetchForEligibilityChoiceData(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	// Ensure empty slice is returned instead of nil
	if len(data) == 0 {
		data = []modelcircular.EligibilityChoiceDataStructure{} // return empty slice instead of nil
	}
	// Return result and record count
	return data, len(data), nil
}