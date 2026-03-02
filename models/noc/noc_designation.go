// Package modelsnoc contains data structures and database access logic for the NOC Designation details.
//
// Path : /var/www/html/go_projects/HRMODULE/Ramya/Hr_test7007/models/noc
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 07-01-2026
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


// Query for retrieving NOC Designation details
var MyQueryNocDesignation = `
	SELECT
		dm.designationname,
		dh.effectivefrom,
		dh.effectiveto
	FROM humanresources.departmentdesignationmappinghistory dh
	JOIN humanresources.designationmaster dm
		ON dh.designationid = dm.designationid
	WHERE dh.employeeid = $1
`


// Struct for NOC Designation (Formatted Dates)
type NocDesignationStructure struct {
	DesignationName string  `json:"post_name"`
	EffectiveFrom   string  `json:"from_date"`
	EffectiveTo     *string `json:"to_date"`
}


// Row Mapper for NOC Designation with date formatting
func RetrieveNocDesignation(rows *sql.Rows) ([]NocDesignationStructure, error) {
	var list []NocDesignationStructure

	for rows.Next() {

		var (
			noc             NocDesignationStructure
			effectiveFrom   time.Time
			effectiveTo     sql.NullTime
		)

		err := rows.Scan(
			&noc.DesignationName,
			&effectiveFrom,
			&effectiveTo,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning noc designation data: %v", err)
		}

		// Format dates as DD-MM-YYYY
		noc.EffectiveFrom = effectiveFrom.Format("02-01-2006")

		if effectiveTo.Valid {
			formatted := effectiveTo.Time.Format("02-01-2006")
			noc.EffectiveTo = &formatted
		} else {
			noc.EffectiveTo = nil
		}

		list = append(list, noc)
	}

	return list, nil
}

