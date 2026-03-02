// Package databaseofficeorder handles database operations for status dropdowns and office order processing.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 15-09-2025
// Last Modified By: Sridharan
// Last Modified Date: 30-09-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"fmt"
)

// GetDropdownValuesFromDB retrieves the status dropdown options based on a cover page number and employee ID.
//
// It extracts "coverpageno" and "employeeid" from the decryptedData map, establishes a connection
// to the "meivan" schema via the GnanaThalam connection utility, and executes the dropdown query.
//
// Returns:
//   - A slice of DropdownValueStruct containing the retrieved options.
//   - An integer representing the count of items found.
//   - An error if required parameters are missing or database operations fail.
func GetDropdownValuesFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.DropdownValueStruct, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract order_type_id from decrypted data
	CoverPageNo, ok := decryptedData["coverpageno"].(string)
	if !ok || CoverPageNo == "" {
		return nil, 0, fmt.Errorf("missing 'coverpageno' in request data")
	}
	EmployeeID, ok := decryptedData["employeeid"].(string)
	if !ok || EmployeeID == "" {
		return nil, 0, fmt.Errorf("missing 'employeeid' in request data")
	}

	rows, err := db.Query(modelsofficeorder.MyQueryDropdownValues, CoverPageNo, EmployeeID)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsofficeorder.RetrieveDropdownValues(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
