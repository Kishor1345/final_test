// // Package modelsofficeorder contains structs and queries for approval page OfficeOrder_approval_remarks API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/officeorder
//
//--- Creator's Info ---
//Creator: Sridharan
//
//Created On: 15-09-2025
//
//Last Modified By: Sridharan
// 
// Last Modified Date: 15-09-2025
package modelsofficeorder

import (
	"database/sql"
	"fmt"
)

// --- Query for fetching comments using process_id + task_id ---
var QueryOfficeComments = `
    SELECT 
        user_display,
        user_role,
        remarks,
        updated_on
    FROM meivan.getcomments($1, $2)
    ORDER BY updated_on DESC;
`

// --- Struct for Result Data ---
type OfficeCommentStructure struct {
	UserDisplay *string `json:"UserID"` // user_display (id - name)
	UserRole    *string `json:"UserRole"`    // user_role
	Remarks     *string `json:"Remarks"`     // remarks
	UpdatedOn   *string `json:"UpdatedOn"`   // updated_on (formatted: YYYY-MM-DD HH:MI)
}

// --- Function to read from rows ---
func RetrieveOfficeComments(rows *sql.Rows) ([]OfficeCommentStructure, error) {
	var comments []OfficeCommentStructure

	for rows.Next() {
		var comment OfficeCommentStructure

		err := rows.Scan(
			&comment.UserDisplay,
			&comment.UserRole,
			&comment.Remarks,
			&comment.UpdatedOn,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning OfficeCommentStructure row: %v", err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
