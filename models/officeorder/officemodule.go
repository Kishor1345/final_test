// Package modelsofficeorder provides data structures and database retrieval logic 
// for office order modules and sub-modules.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/officeorder
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 15-09-2025
// Last Modified By: Sridharan
// Last Modified Date: 15-09-2025
package modelsofficeorder

import (
	"database/sql"
	"fmt"
)

// MyQueryOrderSubModule is the SQL query used to fetch process master details 
// including module ID, process code, and process name based on a specific ID.
var MyQueryOrderSubModule = (`
SELECT module_id, process_code, process_name, description
FROM meivan.process_master
WHERE id = $1
`)

// OrderSubModuleStruct defines the structure for an office order sub-module record,
// mapping the database columns from the process_master table.
type OrderSubModuleStruct struct {
	Module_id    *int    `json:"module_id"`
	Process_code *string `json:"process_code"`
	Process_name *string `json:"process_name"`
	Description  *string `json:"description"`
}

// RetrieveOrderSubModule scans database rows into a slice of OrderSubModuleStruct.
//
// It iterates through the provided *sql.Rows, maps each column to the corresponding 
// field in the struct, and returns a compiled list. Returns an error if the 
// scanning process fails.
func RetrieveOrderSubModule(rows *sql.Rows) ([]OrderSubModuleStruct, error) {
	var list []OrderSubModuleStruct
	for rows.Next() {
		var s OrderSubModuleStruct
		err := rows.Scan(
			&s.Module_id,
			&s.Process_code,
			&s.Process_name,
			&s.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		list = append(list, s)
	}
	return list, nil
}
