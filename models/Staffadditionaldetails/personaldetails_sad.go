// Package modelssad contains structs and builders for Employee Personal Details (SAD).
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
	"fmt"
)

// =============================
// FLAT DB STRUCT (NULL SAFE)
// =============================
type PersonalDetails_sad struct {
	EmployeeID             string
	TaskID                 sql.NullString
	Firstname              sql.NullString
	Middlename             sql.NullString
	Lastname               sql.NullString
	Gender                 sql.NullString
	MaritalStatus          sql.NullString
	DOB                    sql.NullString
	Age                    sql.NullString
	Nationality            sql.NullString
	BirthState             sql.NullString
	BirthDistrict          sql.NullString
	BirthPlace             sql.NullString
	Hometown               sql.NullString
	Religion               sql.NullString
	CasteCategory          sql.NullString
	EmergencyContactNo     sql.NullString
	MobileNo               sql.NullString
	IsPhysicallyChallenged sql.NullString
	PercentageOfDisability sql.NullString
	NatureOfDisability     sql.NullString
	PersonalEmail          sql.NullString
	AadhaarNo              sql.NullString
	MotherTongue           sql.NullString
	PanCardNo              sql.NullString
	FatherName             sql.NullString
	MotherName             sql.NullString
	SpouseName             sql.NullString
	BankName               sql.NullString
	IFSCCode               sql.NullString
	BankAccountNo          sql.NullString
}

// =============================
// STRUCTURED RESPONSE
// =============================
type PersonalDetailsStructuredResponse_sad struct {
	EmployeeID      string                 `json:"employeeid"`
	TaskID          string                 `json:"task_id,omitempty"`
	PersonalDetails PersonalDetailsSec_sad `json:"personal_details"`
	FamilyDetails   FamilyDetailsSec_sad   `json:"family_details"`
	BankDetails     BankDetailsSec_sad     `json:"bank_details"`
}

type PersonalDetailsSec_sad struct {
	FirstName              string `json:"first_name"`
	MiddleName             string `json:"middle_name"`
	LastName               string `json:"last_name"`
	Gender                 string `json:"gender"`
	MaritalStatus          string `json:"marital_status"`
	DOB                    string `json:"dob"`
	Age                    string `json:"age"`
	Nationality            string `json:"nationality"`
	BirthState             string `json:"birth_state"`
	BirthDistrict          string `json:"birth_district"`
	BirthPlace             string `json:"birth_place"`
	Hometown               string `json:"hometown"`
	Religion               string `json:"religion"`
	CasteCategory          string `json:"caste_category"`
	EmergencyContactNo     string `json:"emergency_contact_no"`
	MobileNo               string `json:"mobile_no"`
	IsPhysicallyChallenged string `json:"is_physically_challenged"`
	PercentageOfDisability string `json:"percentage_of_disability"`
	NatureOfDisability     string `json:"nature_of_disability"`
	PersonalEmail          string `json:"personal_email"`
	AadhaarNo              string `json:"aadhaar_no"`
	MotherTongue           string `json:"mother_tongue"`
	PanCardNo              string `json:"pan_card_no"`
}

type FamilyDetailsSec_sad struct {
	FatherName string `json:"father_name"`
	MotherName string `json:"mother_name"`
	SpouseName string `json:"spouse_name"`
}

type BankDetailsSec_sad struct {
	BankName      string `json:"bank_name"`
	IFSCCode      string `json:"ifsc_code"`
	BankAccountNo string `json:"bank_account_no"`
}

// =============================
// NULL STRING HELPER
// =============================
func ns(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

// =============================
// BUILDER FOR NEW
// =============================
func BuildPersonalDetailsResponseNew_sad(data PersonalDetails_sad) PersonalDetailsStructuredResponse_sad {
	return buildCommonResponse(data, false)
}

// =============================
// BUILDER FOR CONTINUE
// =============================
func BuildPersonalDetailsResponseContinue_sad(data PersonalDetails_sad) PersonalDetailsStructuredResponse_sad {
	return buildCommonResponse(data, true)
}

// =============================
// INTERNAL COMMON BUILDER
// =============================
func buildCommonResponse(data PersonalDetails_sad, includeTask bool) PersonalDetailsStructuredResponse_sad {

	resp := PersonalDetailsStructuredResponse_sad{
		EmployeeID: data.EmployeeID,
		PersonalDetails: PersonalDetailsSec_sad{
			FirstName:              ns(data.Firstname),
			MiddleName:             ns(data.Middlename),
			LastName:               ns(data.Lastname),
			Gender:                 ns(data.Gender),
			MaritalStatus:          ns(data.MaritalStatus),
			DOB:                    ns(data.DOB),
			Age:                    ns(data.Age),
			Nationality:            ns(data.Nationality),
			BirthState:             ns(data.BirthState),
			BirthDistrict:          ns(data.BirthDistrict),
			BirthPlace:             ns(data.BirthPlace),
			Hometown:               ns(data.Hometown),
			Religion:               ns(data.Religion),
			CasteCategory:          ns(data.CasteCategory),
			EmergencyContactNo:     ns(data.EmergencyContactNo),
			MobileNo:               ns(data.MobileNo),
			IsPhysicallyChallenged: ns(data.IsPhysicallyChallenged),
			PercentageOfDisability: ns(data.PercentageOfDisability),
			NatureOfDisability:     ns(data.NatureOfDisability),
			PersonalEmail:          ns(data.PersonalEmail),
			AadhaarNo:              ns(data.AadhaarNo),
			MotherTongue:           ns(data.MotherTongue),
			PanCardNo:              ns(data.PanCardNo),
		},
		FamilyDetails: FamilyDetailsSec_sad{
			FatherName: ns(data.FatherName),
			MotherName: ns(data.MotherName),
			SpouseName: ns(data.SpouseName),
		},
		BankDetails: BankDetailsSec_sad{
			BankName:      ns(data.BankName),
			IFSCCode:      ns(data.IFSCCode),
			BankAccountNo: ns(data.BankAccountNo),
		},
	}

	if includeTask {
		resp.TaskID = ns(data.TaskID)
	}

	return resp
}

// =============================
// NEW → FROM FUNCTION
// =============================
func GenericPersonalDetailsRetriever_sad(db *sql.DB, employeeID string) (interface{}, error) {

	row := db.QueryRow(`
		SELECT * 
		FROM humanresources.get_employee_personal_details($1)
	`, employeeID)

	var data PersonalDetails_sad

	err := row.Scan(
		&data.EmployeeID,
		&data.Firstname,
		&data.Middlename,
		&data.Lastname,
		&data.Gender,
		&data.MaritalStatus,
		&data.FatherName,
		&data.MotherName,
		&data.SpouseName,
		&data.DOB,
		&data.Age,
		&data.Nationality,
		&data.BirthState,
		&data.BirthDistrict,
		&data.BirthPlace,
		&data.Hometown,
		&data.Religion,
		&data.CasteCategory,
		&data.EmergencyContactNo,
		&data.MobileNo,
		&data.IsPhysicallyChallenged,
		&data.PercentageOfDisability,
		&data.NatureOfDisability,
		&data.PersonalEmail,
		&data.AadhaarNo,
		&data.MotherTongue,
		&data.BankName,
		&data.IFSCCode,
		&data.BankAccountNo,
		&data.PanCardNo,
	)
	if err != nil {
		return nil, err
	}

	return BuildPersonalDetailsResponseContinue_sad(data), nil

}

// =============================
// CONTINUE → FROM SAD TABLES
// =============================
func ContinuePersonalDetailsRetriever_sad(db *sql.DB, employeeID string) (interface{}, error) {

	query := `
	SELECT
	    m.employee_id,
	    m.task_id,
	    b.first_name,
	    b.middle_name,
	    b.last_name,
	    b.gender,
	    b.marital_status,
	    b.father_name,
	    b.mother_name,
	    b.spouse_name,
	    b.dob,
	    b.age,
	    b.nationality,
	    b.religion,
	    b.caste_category,
	    b.emergency_contact_no,
	    b.mobile_no,
	    b.is_physically_challenged,
	    b.percentage_of_disability,
	    b.nature_of_disability,
	    b.personal_email,
	    b.aadhaar_no,
	    b.mother_tongue,
	    b.bank_name,
	    b.ifsc_code,
	    b.bank_acct_no,
	    b.pan_card_no
	FROM meivan.sad_m m
	JOIN meivan.sad_basic b
	    ON b.task_id = m.task_id
	WHERE m.employee_id = $1
	  AND m.task_status_id = 6
	`

	rows, err := db.Query(query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response []PersonalDetailsStructuredResponse_sad

	for rows.Next() {
		var data PersonalDetails_sad

		err := rows.Scan(
			&data.EmployeeID,
			&data.TaskID,
			&data.Firstname,
			&data.Middlename,
			&data.Lastname,
			&data.Gender,
			&data.MaritalStatus,
			&data.FatherName,
			&data.MotherName,
			&data.SpouseName,
			&data.DOB,
			&data.Age,
			&data.Nationality,
			&data.Religion,
			&data.CasteCategory,
			&data.EmergencyContactNo,
			&data.MobileNo,
			&data.IsPhysicallyChallenged,
			&data.PercentageOfDisability,
			&data.NatureOfDisability,
			&data.PersonalEmail,
			&data.AadhaarNo,
			&data.MotherTongue,
			&data.BankName,
			&data.IFSCCode,
			&data.BankAccountNo,
			&data.PanCardNo,
		)
		if err != nil {
			return nil, err
		}

		response = append(response, BuildPersonalDetailsResponseContinue_sad(data))
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("No records found")
	}

	return response, nil
}
