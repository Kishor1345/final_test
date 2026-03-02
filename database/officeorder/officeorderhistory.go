// Package databaseofficeorder provides database access layer for office order modules,
// including history tracking, approvals, and dropdown management.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 29-10-2025
// Last Modified By: Sridharan
// Last Modified Date: 29-10-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"fmt"
)

// GetOrderHistoryFromDB retrieves the chronological history of a specific office order.
//
// It extracts the "orderNo" from the provided decryptedData map, establishes a connection
// to the Meivan database, and executes a query to fetch all historical actions associated
// with that order number.
//
// Returns:
//   - A slice of OrderHistoryStruct containing the audit trail of the order.
//   - An integer representing the total number of history records retrieved.
//   - An error if the order number is missing, the DB connection fails, or the query fails.
func GetOrderHistoryFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.OrderHistoryStruct, int, error) {

	// Database connection
	db := credentials.GetDB()
	// Extract order_type_id from decrypted data
	orderNo, ok := decryptedData["orderNo"].(string)
	if !ok || orderNo == "" {
		return nil, 0, fmt.Errorf("missing 'orderNo' in request data")
	}
	rows, err := db.Query(modelsofficeorder.MyQueryOrderHistory, orderNo)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsofficeorder.RetrieveOrderHistory(rows)
	if err != nil {
		return nil, 0, err
	}

	return data, len(data), nil
}
