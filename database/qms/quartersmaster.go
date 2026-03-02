// Package databaseqms handles database operations for the Quarters Management System (QMS).
// It provides functionality to retrieve specific quarters details based on user display names.
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

// GetQuartersDetailsFromDB fetches specific quarters information from the database.
//
// It extracts the "displayname" from the provided decryptedData map, establishes a
// connection to the Meivan database, and executes a query to retrieve the relevant
// records using the predefined QuartersDetails query.
//
// Returns:
//   - A slice of QuartersDetailsStruct containing the retrieved records.
//   - An integer representing the total count of records found.
//   - An error if the displayname is missing, the DB connection fails, or the query fails.
func GetQuartersDetailsFromDB(

	decryptedData map[string]interface{},
) ([]models.QuartersDetailsStruct, int, error) {

	displayName, ok := decryptedData["displayname"].(string)
	if !ok || displayName == "" {
		return nil, 0, fmt.Errorf("displayname is required")
	}

	// Database connection
	db := credentials.GetDB()

	rows, err := db.Query(models.MyQueryQuartersDetails, displayName)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	data, err := models.RetrieveQuartersDetails(rows)
	if err != nil {
		return nil, 0, err
	}

	return data, len(data), nil
}
