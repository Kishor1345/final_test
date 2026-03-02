// Package modelssad contains structs and retriever logic for Employee Document Details.
//
// Path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/models/Staffadditionaldetails
// --- Creator's Info ---
// Creator: Rovita
//
// Created On: 29-01-2026
//
// Last Modified By:
//  
// Last Modified Date: 29-01-2026
package modelssad

import (
	"database/sql"
	"fmt"
)

// DocumentDetail represents a single document record
type DocumentDetail struct {
	EmployeeID   string `json:"employeeid"`
	DocumentName string `json:"document_name"`
}

// GenericDocumentDetailsRetriever retrieves and processes document details
func GenericDocumentDetailsRetriever(db *sql.DB, employeeID string) (interface{}, error) {
	rows, err := db.Query(`SELECT * FROM humanresources.get_employee_document_details($1)`, employeeID)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var documents []DocumentDetail

	for rows.Next() {
		var doc DocumentDetail
		if err := rows.Scan(&doc.EmployeeID, &doc.DocumentName); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		documents = append(documents, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return documents, nil
}
