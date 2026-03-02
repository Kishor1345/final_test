// Package modelsnoc contains data structures and database access logic  for retrieving NOC Template page.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/noc
// --- Creator's Info ---
//  Creator: Ramya M R
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

// Query for retrieving NOC Template details
// (DO NOT guess columns for a function)
var MyQueryNocTemplate = `
	SELECT *
	FROM meivan.noc_template(
		$1, -- p_employee_id
		$2, -- p_process_id
		$3, -- p_template_types
		$4  -- p_task_id
	)
`

// Struct MUST match function output exactly (order + names)
type NocTemplateStructure struct {
	TemplateType  string `json:"template_type"`
	HeaderHTML    string `json:"header_html"`
	OrderNo       string `json:"order_no"`
	OrderDate     string `json:"order_date"` // format in DB if needed
	ToColumn      string `json:"to"`
	Subject       string `json:"subject"`
	Reference     string `json:"reference"`
	BodyHTML      string `json:"body_html"`
	SignatureHTML string `json:"signature_html"`
	CcTo          string `json:"cc_to"`
	FooterHTML    string `json:"footer_html"`
}

// Row Mapper — ORDER IS CRITICAL
func RetrieveNocTemplate(rows *sql.Rows) ([]NocTemplateStructure, error) {
	var list []NocTemplateStructure

	for rows.Next() {
		var noc NocTemplateStructure

		err := rows.Scan(
			&noc.TemplateType,
			&noc.HeaderHTML,
			&noc.OrderNo,
			&noc.OrderDate,
			&noc.ToColumn,
			&noc.Subject,
			&noc.Reference,
			&noc.BodyHTML,
			&noc.SignatureHTML,
			&noc.CcTo,
			&noc.FooterHTML,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning noc template data: %v", err)
		}

		list = append(list, noc)
	}

	return list, nil
}
