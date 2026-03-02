// Package modelsnoc contains data structures and database access logic for the NOC Intimation Details.
//
// Path : /var/www/html/go_projects/HRMODULE/Ramya/Hr_test7007/models/noc
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 12-01-2026
//
// Last Modified By:
//
// Last Modified Date:
package modelsnoc

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Query for retrieving NOC Intimation Details
var MyQueryNocIntimationDetails = `
	SELECT 
		m.order_no,
		f.advertisement_no,
		f.advertisement_date,
		f.post_name,
		f.institution_name,
		f.last_date_for_application,
		f.nature_of_post_applied,
		f.post_applied_under,
		f.position_details
	FROM 
		meivan.noc_m m
	JOIN 
		meivan.noc_fa f
		ON m.task_id = f.task_id
	WHERE 
		m.order_no = $1
`

// Struct for NOC Intimation Details (Formatted Dates)
type NocIntimationDetailsStructure struct {
	ReferenceNo                string  `json:"reference_no"`
	AdvertisementNo            string  `json:"advertisement_no"`
	AdvertisementDate          string  `json:"advertisement_date"`
	PostName                   string  `json:"post_name"`
	InstitutionName            string  `json:"institution_name"`
	LastDateForApplication     *string `json:"last_date_for_application"`
	NatureOfPostApplied        string  `json:"nature_of_post_applied"`
	PostAppliedUnder           string  `json:"post_applied_under"`
	PositionDetails            string  `json:"position_details"`
}

// Row Mapper for NOC Intimation Details with date formatting
func RetrieveNocIntimationDetails(rows *sql.Rows) ([]NocIntimationDetailsStructure, error) {
	var list []NocIntimationDetailsStructure

	for rows.Next() {
		var (
			noc                  NocIntimationDetailsStructure
			advertisementDate     time.Time
			lastDateForApplication sql.NullTime
		)

		err := rows.Scan(
			&noc.ReferenceNo,
			&noc.AdvertisementNo,
			&advertisementDate,
			&noc.PostName,
			&noc.InstitutionName,
			&lastDateForApplication,
			&noc.NatureOfPostApplied,
			&noc.PostAppliedUnder,
			&noc.PositionDetails,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning noc intimation data: %v", err)
		}

		// Format dates as DD-MM-YYYY
		noc.AdvertisementDate = advertisementDate.Format("02-01-2006")

		if lastDateForApplication.Valid {
			formatted := lastDateForApplication.Time.Format("02-01-2006")
			noc.LastDateForApplication = &formatted
		} else {
			noc.LastDateForApplication = nil
		}

		list = append(list, noc)
	}

	return list, nil
}
