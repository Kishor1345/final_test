// Package models contains data structures and database access logic for the Campus Master page.
//
//Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/common
//--- Creator's Info ---
// Creator: Ramya M R
//
// Created On:10-02-2026
//
// Last Modified By:
//
// Last Modified Date:
package modelscommon

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Query to retrieve Campus details
var MyQueryCampusMaster = `
SELECT 
    id,
    campuscode,
    campusname,
    location
FROM humanresources.campus
`

//  Struct for Campus Master
type CampusMasterStructure struct {
	ID         *string `json:"ID"`
	CampusCode *string `json:"CampusCode"`
	CampusName *string `json:"CampusName"`
	Location   *string `json:"Location"`
}

// Row Mapper
func RetrieveCampusMaster(rows *sql.Rows) ([]CampusMasterStructure, error) {
	var list []CampusMasterStructure

	for rows.Next() {
		var campus CampusMasterStructure
		err := rows.Scan(
			&campus.ID,
			&campus.CampusCode,
			&campus.CampusName,
			&campus.Location,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		list = append(list, campus)
	}

	return list, nil
}
