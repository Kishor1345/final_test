// Package databaseofficeorder handles database operations for OfficeOrdervisitdetails.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 30-09-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"fmt"
)

// GetVisitDetailsFromDB fetches visit details based on employeeid and/or coverpageno
func GetVisitDetailsFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.VisitDetails, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract parameters safely
	EmployeeID, _ := decryptedData["employeeid"].(string)
	CoverPageNo, _ := decryptedData["coverpageno"].(string)

	// --- Case 1: Both EmployeeID and CoverPageNo ---
	if EmployeeID != "" && CoverPageNo != "" {
		rows, err := db.Query(modelsofficeorder.MyQueryVisitDetailsemployeecoverpage, EmployeeID, CoverPageNo)
		if err != nil {
			return nil, 0, fmt.Errorf("query error (employee+cover): %v", err)
		}
		defer rows.Close()

		visitDetails, err := modelsofficeorder.RetrieveVisitDetailsByCover(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("data retrieve error: %v", err)
		}
		return visitDetails, len(visitDetails), nil
	}

	// --- Case 2: Only EmployeeID ---
	if EmployeeID != "" {
		rows, err := db.Query(modelsofficeorder.MyQueryVisitDetails)
		if err != nil {
			return nil, 0, fmt.Errorf("query error (employee): %v", err)
		}
		defer rows.Close()

		visitDetails, err := modelsofficeorder.RetrieveVisitDetails(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("data retrieve error: %v", err)
		}
		return visitDetails, len(visitDetails), nil
	}

	// --- Case 3: No filters → return all visits ---
	rows, err := db.Query(modelsofficeorder.MyQueryVisitDetails)
	if err != nil {
		return nil, 0, fmt.Errorf("query error (all): %v", err)
	}
	defer rows.Close()

	visitDetails, err := modelsofficeorder.RetrieveVisitDetails(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("data retrieve error: %v", err)
	}
	return visitDetails, len(visitDetails), nil
}
