// Package modelssad contains structs and retriever logic for Employee Dependents, including mapping database-dependent details to API response for SAD workflow APIs.
//
// Path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/models/Staffadditionaldetails
// --- Creator's Info ---
// Creator: Rovita
//
// Created On: 29-01-2026
//
// Last Modified By:
//  
// Last Modified Date: 
package modelssad

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

/* ============================
   DATABASE STRUCT
============================ */

type EmployeedatabaseDependent struct {
	EmployeeID                  string  `json:"employeeid"`
	DependentName               string  `json:"dependent_name"`
	DependentRelationship       string  `json:"dependent_relationship"`
	DependentDOB                string  `json:"dependent_dob"`
	DependentAge                *int    `json:"dependent_age"`
	DependentMaritalStatus      *string `json:"dependent_marital_status"`
	DependentBloodGroup         *string `json:"dependent_blood_group"`
	DependentGender             *string `json:"dependent_gender"`
	DependentEmploymentStatus   *string `json:"dependent_employment_status"`
	DependentAadhaarNo          *string `json:"dependent_aadhaar_no"`
	IsTwins                     int     `json:"is_twins"`
	IsCurrentlyDependent        string  `json:"is_currently_dependent"`
	IsActive                    int     `json:"is_active"` // NEW
	DependentMobileNo           *string `json:"dependent_mobile_no"`
	OptingForInsurance          int     `json:"opting_for_insurance"`
	OptingForLTC                int     `json:"opting_for_ltc"`
	IsPersonDisabled            string  `json:"is_person_disabled"`
	DependentNatureOfDisability *string `json:"dependent_nature_of_disability"`
}

/* ============================
   API RESPONSE STRUCT
============================ */

type EmployeeresponseDependent struct {
	DependentName               string  `json:"dependent_name"`
	DependentRelationship       string  `json:"dependent_relationship"`
	DependentDOB                string  `json:"dependent_dob"`
	DependentAge                *int    `json:"dependent_age,omitempty"`
	DependentMaritalStatus      *string `json:"dependent_marital_status,omitempty"`
	DependentBloodGroup         *string `json:"dependent_blood_group,omitempty"`
	DependentGender             *string `json:"dependent_gender,omitempty"`
	DependentEmploymentStatus   *string `json:"dependent_employment_status,omitempty"`
	DependentAadhaarNo          *string `json:"dependent_aadhaar_no,omitempty"`
	IsTwins                     bool    `json:"is_twins"`
	DependentMobileNo           *string `json:"dependent_mobile_no,omitempty"`
	IsCurrentlyDependent        string  `json:"is_currently_dependent"`
	IsActive                    bool    `json:"is_active"` // NEW
	OptingForInsurance          bool    `json:"opting_for_insurance"`
	OptingForLTC                bool    `json:"opting_for_ltc"`
	IsPersonDisabled            string  `json:"is_person_disabled"`
	DependentNatureOfDisability *string `json:"dependent_nature_of_disability,omitempty"`
}

/* ============================
   FINAL RESPONSE
============================ */

type EmployeeDependentReponse struct {
	EmployeeID string                      `json:"employeeid"`
	Dependents []EmployeeresponseDependent `json:"dependents"`
}

/* ============================
   DB → API MAPPER
============================ */

func mapDependentDBToAPI(dbDep EmployeedatabaseDependent) EmployeeresponseDependent {
	resp := EmployeeresponseDependent{
		DependentName:         dbDep.DependentName,
		DependentRelationship: dbDep.DependentRelationship,
		DependentDOB:          dbDep.DependentDOB,
		IsTwins:               dbDep.IsTwins == 1,
		IsCurrentlyDependent:  dbDep.IsCurrentlyDependent,
		IsActive:              dbDep.IsActive == 1, //  NEW
		OptingForInsurance:    dbDep.OptingForInsurance == 1,
		OptingForLTC:          dbDep.OptingForLTC == 1,
		IsPersonDisabled:      dbDep.IsPersonDisabled,
	}

	if dbDep.DependentAge != nil {
		resp.DependentAge = dbDep.DependentAge
	}
	if dbDep.DependentMaritalStatus != nil && *dbDep.DependentMaritalStatus != "" {
		resp.DependentMaritalStatus = dbDep.DependentMaritalStatus
	}
	if dbDep.DependentBloodGroup != nil && *dbDep.DependentBloodGroup != "" {
		resp.DependentBloodGroup = dbDep.DependentBloodGroup
	}
	if dbDep.DependentGender != nil && *dbDep.DependentGender != "" {
		resp.DependentGender = dbDep.DependentGender
	}
	if dbDep.DependentEmploymentStatus != nil && *dbDep.DependentEmploymentStatus != "" {
		resp.DependentEmploymentStatus = dbDep.DependentEmploymentStatus
	}
	if dbDep.DependentAadhaarNo != nil && *dbDep.DependentAadhaarNo != "" {
		resp.DependentAadhaarNo = dbDep.DependentAadhaarNo
	}
	if dbDep.DependentMobileNo != nil && *dbDep.DependentMobileNo != "" {
		resp.DependentMobileNo = dbDep.DependentMobileNo
	}
	if dbDep.DependentNatureOfDisability != nil && *dbDep.DependentNatureOfDisability != "" {
		resp.DependentNatureOfDisability = dbDep.DependentNatureOfDisability
	}

	return resp
}

/* ============================
   DB FETCH FUNCTION
============================ */

func FetchEmployeeDependentDetailsFromDB(db *sql.DB, employeeID string) (interface{}, error) {

	log.Printf("Fetching dependents for employee: %s", employeeID)

	row := db.QueryRow(
		`SELECT humanresources.get_employee_dependent_details_sad($1)`,
		employeeID,
	)

	var jsonData []byte
	if err := row.Scan(&jsonData); err != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}

	var dbDependents []EmployeedatabaseDependent
	if err := json.Unmarshal(jsonData, &dbDependents); err == nil {

		respDeps := make([]EmployeeresponseDependent, 0)
		for _, d := range dbDependents {
			respDeps = append(respDeps, mapDependentDBToAPI(d))
		}

		return EmployeeDependentReponse{
			EmployeeID: employeeID,
			Dependents: respDeps,
		}, nil
	}

	var single EmployeedatabaseDependent
	if err := json.Unmarshal(jsonData, &single); err == nil {

		return EmployeeDependentReponse{
			EmployeeID: employeeID,
			Dependents: []EmployeeresponseDependent{mapDependentDBToAPI(single)},
		}, nil
	}

	return nil, fmt.Errorf("invalid JSON returned from DB")
}
