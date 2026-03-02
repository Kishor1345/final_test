// Package modelsqms contains structs and queries for floor Dropdown API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:22/01/2026
package quartersmasterestate

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// FloorDetailsStructure represents
type FloorDetailsStructure struct {
	ID          *int    `json:"id"`
	FloorName   *string `json:"floor_name"`
	FloorNumber *string `json:"floor_number"`
}

// SQL query to fetch quarters master floor drop down data
var MyQueryFloorDropdown = `
SELECT id,floor_name,floor_number
FROM humanresources.floormaster
`

// RetrieveFloorDropdown maps SQL rows
// into FloorDetailsStructure
func RetrieveFloorDropdown(rows *sql.Rows) ([]FloorDetailsStructure, error) {
	var list []FloorDetailsStructure

	for rows.Next() {
		var s FloorDetailsStructure

		// Scan database row into struct fields
		err := rows.Scan(
			&s.ID,
			&s.FloorName,
			&s.FloorNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning Floor dropdown: %v", err)
		}

		list = append(list, s)
	}

	return list, nil
}
