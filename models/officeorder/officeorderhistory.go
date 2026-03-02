// Package modelsofficeorder provides data structures and SQL queries for tracking 
// and retrieving the history of office orders within the Meivan module.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 14-11-2025
// Last Modified By:
// Last Modified Date:
package modelsofficeorder

import (
	"database/sql"
	"fmt"
)


// MyQueryOrderHistory retrieves the history of completed orders related to a base order number.
// It uses regex to strip amendment suffixes from the input order number and finds all 
// matching records in the pcr_m table that have a task_status_id of 3 (Completed).
var MyQueryOrderHistory = (`
WITH base AS (
    SELECT regexp_replace($1, '/A[0-9]+(/CAN)?$', '') AS base_order
)
SELECT
    order_no,
    task_status_id
FROM meivan.pcr_m m, base b
WHERE m.order_no LIKE b.base_order || '%'
AND m.task_status_id = 3   -- only completed
ORDER BY length(m.order_no), m.order_no
`)



// MyQueryOrderHistory retrieves the history of completed orders related to a base order number.
// It uses regex to strip amendment suffixes from the input order number and finds all 
// matching records in the pcr_m table that have a task_status_id of 3 (Completed).
type OrderHistoryStruct struct {
	OrderNo      *string `json:"order_no"`
	TaskStatusID *int    `json:"task_status_id"`
}


// RetrieveOrderHistory processes the database result set into a slice of OrderHistoryStruct.
//
// It iterates through the provided *sql.Rows and scans the order_no and task_status_id 
// into the corresponding struct fields. Returns a list of history records or an error 
// if scanning fails.
func RetrieveOrderHistory(rows *sql.Rows) ([]OrderHistoryStruct, error) {
	var list []OrderHistoryStruct
	for rows.Next() {
		var s OrderHistoryStruct
		err := rows.Scan(
			&s.OrderNo,
			&s.TaskStatusID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		list = append(list, s)
	}
	return list, nil
}
