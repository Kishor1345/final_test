// Package modelscircular contains structs and queries for  circular data fetch for  Quarters Number.
//
//Path : /var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/models/hr_008
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:25/02/2026
package modelcircular

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type QuartersNumberDataStructure struct{
	QuartersId      int         `json:"quarters_id"`
	DisplayName  string      `json:"quarters_number""`
	QuartersStatus  string      `json:"quarters_status"`
	CategoryName    string      `json:"category_name"`
}


var MyQueryForQuartersNumberData = 
`
SELECT hqm.id,hqm.displayname,hqm.quartersstatus,hqc.name
FROM humanresources.quartersmaster hqm
JOIN humanresources.buildingmaster hbm ON hqm.building_id = hbm.id
JOIN humanresources.quarterscategory hqc ON hbm.quarters_category = hqc.id
WHERE hqm.campus_id = $1 AND hqc.id = $2 AND hqm.quartersstatus = 'Fit for occupation'
`

func RetrieveQuartersNumberDataFetch(rows *sql.Rows) ([]QuartersNumberDataStructure, error) {

	var results []QuartersNumberDataStructure

	for rows.Next() {
		var r QuartersNumberDataStructure

        // Scan database row into struct fields
		err := rows.Scan(
			&r.QuartersId,
			&r.DisplayName,
			&r.QuartersStatus,
			&r.CategoryName,
		)
		if err != nil {
			return nil, err
		}


		results = append(results, r)
	}

	return results, nil
}
