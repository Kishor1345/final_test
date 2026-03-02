// Package modelscircular contains structs and queries for  circular data fetch for Eligibility Choice.
//
//Path : /var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/models/hr_008
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:13/02/2026
package modelcircular

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type EligibilityChoiceDataStructure struct{
	Id 			int     `json:"id"`
	CriteriaID  string  `json:"criteria_id"`
	Status      int     `json:"status"`
}


var MyQueryForEligibilityChoiceData = 
`
select id,criteria_id,status
from humanresources.criteria_master
where status = '1'
order by criteria_id
`

func RetrieveCircularDataFetchForEligibilityChoiceData(rows *sql.Rows) ([]EligibilityChoiceDataStructure, error) {

	var results []EligibilityChoiceDataStructure

	for rows.Next() {
		var r EligibilityChoiceDataStructure

        // Scan database row into struct fields
		err := rows.Scan(
			&r.Id,
			&r.CriteriaID,
			&r.Status,
		)
		if err != nil {
			return nil, err
		}


		results = append(results, r)
	}

	return results, nil
}
