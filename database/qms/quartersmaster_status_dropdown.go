// Package databaseqms handles database operations for the Quarters Management System (QMS).
// It manages the retrieval and processing of quarters-related data from the database.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/qms
//
// --- Creator's Info ---
// Creator: Elakiya
// Created On:
// Last Modified By:
// Last Modified Date:
package databaseqms

import (
	credentials "Hrmodule/dbconfig"
	models "Hrmodule/models/qms"
	"fmt"
)

// GetQuartersStatusFromDB retrieves the current status details of quarters filtered by section ID and status.
//
// It extracts "section_id" and "status" from the provided decryptedData map. Note that these
// values are expected to be float64 (standard for numbers in a map[string]interface{} unmarshaled from JSON).
// The function establishes a connection to the Meivan database and executes the predefined QuartersStatus query.
//
// Returns:
//   - A slice of QuartersStatusStruct containing the matching quarters records.
//   - An integer representing the total count of records found.
//   - An error if the required parameters are missing or if any part of the database operation fails.
func GetQuartersStatusFromDB(
	decryptedData map[string]interface{},
) ([]models.QuartersStatusStruct, int, error) {

	// =====================
	// INPUT PARAMETERS
	// =====================
	sectionIDFloat, ok := decryptedData["section_id"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("section_id is required")
	}
	statusFloat, ok := decryptedData["status"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("status is required")
	}

	sectionID := int(sectionIDFloat)
	status := int(statusFloat)

	// Database connection
	db := credentials.GetDB()

	rows, err := db.Query(
		models.MyQueryQuartersStatus,
		sectionID,
		status,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	data, err := models.RetrieveQuartersStatus(rows)
	if err != nil {
		return nil, 0, err
	}

	return data, len(data), nil
}
