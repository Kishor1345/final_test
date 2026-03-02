// Package modelssad contains structs and retriever logic for Employee Contact Details, including current and permanent address mappings for SAD workflow APIs.
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
	"log"
)

/* =============================
   QUERY (SP CALL)
============================= */

const EmployeeContactDetailsSP = `
SELECT *
FROM humanresources.get_employee_contact_details_sad($1)
`

/* =============================
   RESPONSE STRUCT
============================= */

type EmployeeContactDetails struct {
	EmployeeID         string `json:"employeeid"`
	PresentAddress1    string `json:"present_address1"`
	PresentAddress2    string `json:"present_address2"`
	PresentCountry     string `json:"present_country"`
	PresentState       string `json:"present_state"`
	PresentDistrict    string `json:"present_district"`
	PresentCity        string `json:"present_city"`
	PresentPincode     string `json:"present_pincode"`
	PresentCountryCode string `json:"present_countrycode"`
	PresentAreaCode    string `json:"present_areacode"`

	PermanentAddress1 string `json:"permanent_address1"`
	PermanentAddress2 string `json:"permanent_address2"`
	PermanentCountry  string `json:"permanent_country"`
	PermanentState    string `json:"permanent_state"`
	PermanentDistrict string `json:"permanent_district"`
	PermanentCity     string `json:"permanent_city"`
	PermanentPincode  string `json:"permanent_pincode"`

	IsQuarters string `json:"is_quarters"` // "1" or "0"
	QuartersNo string `json:"quarters_no"`
}

/* =============================
   ROW SCANNER
============================= */

func RetrieveEmployeeContactDetails(rows *sql.Rows) ([]EmployeeContactDetails, error) {

	var results []EmployeeContactDetails
	count := 0

	for rows.Next() {
		var rec EmployeeContactDetails

		if err := rows.Scan(
			&rec.EmployeeID,
			&rec.PresentAddress1,
			&rec.PresentAddress2,
			&rec.PresentCountry,
			&rec.PresentState,
			&rec.PresentDistrict,
			&rec.PresentCity,
			&rec.PresentPincode,
			&rec.PresentCountryCode,
			&rec.PresentAreaCode,
			&rec.PermanentAddress1,
			&rec.PermanentAddress2,
			&rec.PermanentCountry,
			&rec.PermanentState,
			&rec.PermanentDistrict,
			&rec.PermanentCity,
			&rec.PermanentPincode,
			&rec.IsQuarters,
			&rec.QuartersNo,
		); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		results = append(results, rec)
		count++
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("Fetched %d contact detail records", count)
	return results, nil
}
