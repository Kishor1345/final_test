// Package modelssad contains structs and queries for Task Status validation.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/Staffadditionaldetails
// Creator: Rovita
//
// Created On: 29-01-2026
//
// Last Modified By:
//  
// Last Modified Date:
package modelssad

import (
	"database/sql"
	"fmt"
	"log"
)

// =============================
// QUERY
// =============================
const TaskStatusCheckQuery = `
SELECT
    employee_id,
		task_id,
    CASE
        WHEN task_status_id = 4 THEN 'ERROR: Ongoing task already exists'
        WHEN task_status_id = 6 THEN 'ERROR: Save & Hold task exists'
    END AS message
FROM meivan.sad_m
WHERE employee_id = $1
  AND task_status_id IN (4, 6)
`

// =============================
// RESPONSE STRUCT
// =============================
type TaskStatusCheck struct {
	EmployeeID string `json:"employeeid"`
	Task_id    string `json:"task_id"`
	Message    string `json:"message"`
}

// =============================
// ROW SCANNER (REFERENCE PATTERN)
// =============================
func RetrieveTaskStatusChecks(rows *sql.Rows) ([]TaskStatusCheck, error) {

	var results []TaskStatusCheck
	var recordCount int

	for rows.Next() {
		var rec TaskStatusCheck
		if err := rows.Scan(&rec.EmployeeID, &rec.Task_id, &rec.Message); err != nil {
			return nil, fmt.Errorf("error scanning task status row: %v", err)
		}
		results = append(results, rec)
		recordCount++
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	log.Printf("Successfully scanned %d task status records", recordCount)
	return results, nil
}
