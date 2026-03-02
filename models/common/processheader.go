// Package modelscommon contains data structures and DB scan logic for ProcessHeader API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/common
//
// --- Creator's Info ---
//
// Creator: Rovita
//
// Created On:30-12-2025
//
// Last Modified By: Rovita
//
// Last Modified Date: 30-12-2025
package modelscommon

import (
    "database/sql"
    "fmt"
)

var MyQueryProcessHeader = `
    SELECT id, process_code, process_name, status 
    FROM meivan.process_master 
    WHERE status = '1' AND id = $1
`

// ProcessHeader defines structure for process_master table
type ProcessHeader struct {
    ID          *int    `json:"id"`
    ProcessCode *string `json:"process_code"`
    ProcessName *string `json:"process_name"`
    Status      *string `json:"status"`
}

// RetrieveProcessHeader scans rows into []ProcessHeader
func RetrieveProcessHeader(rows *sql.Rows) ([]ProcessHeader, error) {
    var result []ProcessHeader
    for rows.Next() {
        var p ProcessHeader
        err := rows.Scan(
            &p.ID,
            &p.ProcessCode,
            &p.ProcessName,
            &p.Status,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning row: %v", err)
        }
        result = append(result, p)
    }
    return result, nil
}