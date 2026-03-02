// Package modelsnoc contains data structures and database access logic for the NOC Approver page.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/noc
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 07-01-2026
//
// Last Modified By:
//
// Last Modified Date:
package modelsnoc

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)


// Query for retrieving NOC Approver details (JSONB response)
var MyQueryNocApprover = `
	SELECT meivan.get_noc_details_by_id($1)::jsonb
`


// Struct for NOC Approver (Dynamic JSONB Data)
// NOTE:
// - Supports multiple types
// - Column count varies per type
// - Dates will be formatted in Service layer
type NocApproverStructure struct {
	Data map[string]interface{} `json:"data"`
}


// Row Mapper for NOC Approver JSONB response
func RetrieveNocApprover(rows *sql.Rows) ([]NocApproverStructure, error) {
	var list []NocApproverStructure

	for rows.Next() {
		var rawJSON []byte
		var noc NocApproverStructure

		err := rows.Scan(&rawJSON)
		if err != nil {
			return nil, fmt.Errorf("error scanning jsonb data: %v", err)
		}

		err = json.Unmarshal(rawJSON, &noc.Data)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling jsonb data: %v", err)
		}

		list = append(list, noc)
	}

	return list, nil
}

