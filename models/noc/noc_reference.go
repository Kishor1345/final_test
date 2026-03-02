// Package modelsnoc contains data structures and database access logic  for retrieving NOC Reference order numbers.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/noc
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 08-01-2026
//
// Last Modified By:
//
// Last Modified Date:
 package modelsnoc

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Query for retrieving NOC Reference order numbers for a given employee
// Fixed: purpose uses integer code instead of string
var MyQueryNocReference = `
	SELECT order_no
	FROM meivan.noc_m
	WHERE employee_id =$1
	  AND purpose = 1  -- replace 1 with actual integer code for "Intimation"
	  and  certificate_type ='4'
	  AND task_status_id = 3

`

// Struct for NOC Reference
type NocReferenceStructure struct {
	OrderNo string `json:"reference_number"`
}

// Row Mapper for NOC Reference
func RetrieveNocReference(rows *sql.Rows) ([]NocReferenceStructure, error) {
	var list []NocReferenceStructure

	for rows.Next() {
		var noc NocReferenceStructure

		err := rows.Scan(&noc.OrderNo)
		if err != nil {
			return nil, fmt.Errorf("error scanning noc reference data: %v", err)
		}

		list = append(list, noc)
	}

	return list, nil
}

