// Package databaseofficeorder handles database access for the Office Order Approval page.
// It provides functionality to retrieve and process office order master records.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/officeorder
//
// --- Creator's Info ---
// Creator: Ramya
// Created On: 09-10-2025
// Last Modified By:
// Last Modified Date:
package databaseofficeorder

import (
	// Assuming the package path structure based on your sample code
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"fmt"
	"net/http"
	"strconv" // Need this for string-to-int conversion
)

// GetOfficeOrderMasterFromDB retrieves office order master records based on task status, employee ID, and cover page number.
//
// It performs the following steps:
// 1. Converts the taskStatusID from a string to an integer.
// 2. Establishes a connection to the Meivan database.
// 3. Calls the model layer to fetch office order master data.
//
// Parameters:
//   - w: http.ResponseWriter to handle and send error responses back to the client.
//   - r: *http.Request representing the current HTTP request.
//   - taskStatusID: The string representation of the status ID to filter by.
//   - employeeID: The unique identifier of the employee.
//   - coverPageNo: The specific cover page number for the query filter.
//
// Returns:
//   - A slice of OfficeOrderMasterStructure containing the results.
//   - An integer representing the total count of records found.
//   - An error object if any step in the process fails.
func GetOfficeOrderMasterFromDB(w http.ResponseWriter, r *http.Request, taskStatusID string, employeeID string, coverPageNo string) ([]modelsofficeorder.OfficeOrderMasterStructure, int, error) {

	statusID, err := strconv.Atoi(taskStatusID)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid taskStatusID provided: %s. Must be an integer.", taskStatusID)
		http.Error(w, errMsg, http.StatusBadRequest)
		return nil, 0, fmt.Errorf(errMsg)
	}
	// Database connection
	db := credentials.GetDB()

	data, err := modelsofficeorder.GetOfficeOrderMasters(db, statusID, employeeID, coverPageNo)
	if err != nil {

		return nil, 0, fmt.Errorf("retrieving office order master failed: %v", err)
	}

	return data, len(data), nil
}
