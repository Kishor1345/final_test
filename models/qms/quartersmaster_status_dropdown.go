// Package modelsquarters contains structs and queries for Quarters details API.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/qms
// --- Creator's Info ---
// Creator: Elakiya
//
// Created On: 
//
// Last Modified By:
//
// Last Modified Date:
package modelsqms

import (
	"database/sql"
	"fmt"
)

// =====================
// SQL QUERY
// =====================
var MyQueryQuartersStatus = `
SELECT name
FROM humanresources.quarters_status_master
WHERE section_id = $1
  AND status = $2
ORDER BY name;
`

// =====================
// RESPONSE STRUCT
// =====================
type QuartersStatusStruct struct {
	Name *string `json:"name"`
}

// =====================
// ROW SCANNER
// =====================
func RetrieveQuartersStatus(rows *sql.Rows) ([]QuartersStatusStruct, error) {

	var list []QuartersStatusStruct

	for rows.Next() {
		var s QuartersStatusStruct

		if err := rows.Scan(&s.Name); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		list = append(list, s)
	}

	return list, nil
}
