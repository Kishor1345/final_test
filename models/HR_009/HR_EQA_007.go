// Package modelQuarters for fetching preference details.
//
//Path : /var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/models/HR_009
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:23/02/2026
package modelQuarters

import (
	"database/sql"
	_ "github.com/lib/pq"
)


type QuartersPreferenceFetchStructure struct{
	Preference    string  `json:"preference"`
	OrderNo       *string  `json:"order_no"`
	QuartersNo    *string  `json:"quarters_no"`
}


var MyQueryForQuartersPreference = 
`
SELECT * FROM meivan.quarters_application_preference($1)
`

func RetrieveQuartersPreferenceDetails(rows *sql.Rows) ([]QuartersPreferenceFetchStructure, error) {

	var results []QuartersPreferenceFetchStructure
	
	for rows.Next() {
		var r QuartersPreferenceFetchStructure
        // Scan database row into struct fields
	err := rows.Scan(
		&r.Preference,
		&r.OrderNo,
		&r.QuartersNo,
	)

		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}
