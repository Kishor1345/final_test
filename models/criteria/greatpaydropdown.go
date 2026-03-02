// Package modelscriteria contains structs and queries for GreatPay Dropdown API.
//
//Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/criteria
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:08/01/2026
package modelscriteria

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// GreatpayDetailsStructure represents
type GreatpayDetailsStructure struct{
	ID int `json:"id"`
	Name *string `json:"name"`
}

// SQL query to fetch greatpay dropdown data 
var MyQueryGreatpayDropdown=`
SELECT id,name 
FROM humanresources.cpc_master 
WHERE cpc= $1
`


// RetrieveGreatpayDropdown maps SQL rows
// into GreatpayDetailsStructure
func RetrieveGreatpayDropdown(rows *sql.Rows) ([]GreatpayDetailsStructure, error) {
	var list []GreatpayDetailsStructure

	for rows.Next() {
		var s GreatpayDetailsStructure

		// Scan database row into struct fields
		err := rows.Scan(
			&s.ID,
			&s.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning Greatpay dropdown: %v", err)
		}

	
		list = append(list, s)
	}

	return list, nil
}