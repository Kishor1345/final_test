// Package modelssad contains request/response structures and custom date handling logic for SAD master workflow processing.
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
	"strings"
	"time"
)

// CustomDate wraps time.Time to parse "DD-MM-YYYY" from JSON
type CustomDate struct {
	time.Time
}

// UnmarshalJSON parses "DD-MM-YYYY" into CustomDate
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		cd.Time = time.Time{}
		return nil
	}

	t, err := time.Parse("02-01-2006", s)
	if err != nil {
		return err
	}

	cd.Time = t
	return nil
}

// MarshalJSON converts CustomDate to "DD-MM-YYYY" string
func (cd *CustomDate) MarshalJSON() ([]byte, error) {
	if cd == nil || cd.Time.IsZero() {
		return []byte(`null`), nil
	}
	s := cd.Format("02-01-2006")
	return []byte(`"` + s + `"`), nil
}

// MasterSadRequest represents the parameters for master_sad
// All dynamic JSON fields remain as interface{}; dates inside them will be handled at marshaling
type MasterSadRequest struct {
	ActionType string  `json:"action_type"`
	TaskID     *string `json:"task_id"`
	ProcessID  int     `json:"process_id"`
	EmployeeID string  `json:"employee_id"`
	UpdatedBy  string  `json:"updated_by"`

	EmployeeName   *string `json:"employee_name"`
	AssignTo       *string `json:"assign_to"`
	AssignedRole   *string `json:"assigned_role"`
	TaskStatusID   *int    `json:"task_status_id"`
	ActivitySeqNo  *int    `json:"activity_seq_no"`
	IsTaskReturn   *int    `json:"is_task_return"`
	IsTaskApproved *int    `json:"is_task_approved"`
	EmailFlag      *int    `json:"email_flag"`
	TemplateID     *int    `json:"template_id"`
	RejectFlag     *int    `json:"reject_flag"`
	RejectRole     *string `json:"reject_role"`
	InitiatedBy    *string `json:"initiated_by"`
	Badge          *int    `json:"badge"`
	Priority       *int    `json:"priority"`
	Starred        *int    `json:"starred"`
	// WORKFLOW INPUTS (FROM FRONTEND)
	EmployeeGroup string `json:"EmployeeGroup"`
	PBM           string `json:"PBM"`
	// Basic Info
	FirstName              *string     `json:"first_name"`
	MiddleName             *string     `json:"middle_name"`
	LastName               *string     `json:"last_name"`
	Gender                 *string     `json:"gender"`
	MaritalStatus          *string     `json:"marital_status"`
	FatherName             *string     `json:"father_name"`
	MotherName             *string     `json:"mother_name"`
	SpouseName             *string     `json:"spouse_name"`
	Dob                    *CustomDate `json:"dob"`
	Age                    *int        `json:"age"`
	Nationality            *string     `json:"nationality"`
	Religion               *string     `json:"religion"`
	CasteCategory          *string     `json:"caste_category"`
	EmergencyContactNo     *string     `json:"emergency_contact_no"`
	MobileNo               *string     `json:"mobile_no"`
	IsPhysicallyChallenged *int        `json:"is_physically_challenged"`
	PercentageOfDisability *float64    `json:"percentage_of_disability"`
	NatureOfDisability     *string     `json:"nature_of_disability"`
	PersonalEmail          *string     `json:"personal_email"`
	AadhaarNo              *string     `json:"aadhaar_no"`
	MotherTongue           *string     `json:"mother_tongue"`
	BankName               *string     `json:"bank_name"`
	IfscCode               *string     `json:"ifsc_code"`
	BankAcctNo             *string     `json:"bank_acct_no"`
	IdentificationMarks    *string     `json:"identification_marks"`
	PanCardNo              *string     `json:"pan_card_no"`
	EmployeeType           *string     `json:"employee_type"`
	Department             *string     `json:"department"`
	Designation            *string     `json:"designation"`
	Section                *string     `json:"section"`
	RouteTo                *string     `json:"route_to"`
	Grade                  *string     `json:"grade"`
	EmpGroup               *string     `json:"emp_group"`
	PayInfo                *string     `json:"pay_info"`
	BasicPay               *float64    `json:"basic_pay"`
	NonPracticePay         *float64    `json:"non_practice_pay"`
	NameOfPayBand          *string     `json:"name_of_pay_band"`
	DateOfJoining          *CustomDate `json:"date_of_joining"`
	DateOfConfirmation     *CustomDate `json:"date_of_confirmation"`
	OfficeRoomNo           *string     `json:"office_room_no"`
	OfficeExtensionNo      *string     `json:"office_extension_no"`
	IsActive               *int        `json:"is_active"`
	EmployeeStatus         *string     `json:"employee_status"`
	DateOfRetirement       *CustomDate `json:"date_of_retirement"`
	EffectiveDate          *CustomDate `json:"effective_date"`

	// Dynamic JSON Data
	ContactData    interface{} `json:"contact_data"`
	DependentsData interface{} `json:"dependents_data"`
	EducationData  interface{} `json:"education_data"`
	ExperienceData interface{} `json:"experience_data"`
	LanguageData   interface{} `json:"language_data"`
	NomineeData    interface{} `json:"nominee_data"`
	HindiData      interface{} `json:"hindi_data"`

	Comments       *string  `json:"comments"`
	UserRole       *string  `json:"user_role"`
	ModuleName     *string  `json:"module_name"`
	ModifiedFields []string `json:"modified_fields"`
	Rowid          *int     `json:"p_row_id"`
}

// MasterSadResponse represents stored procedure result
type MasterSadResponse struct {
	TaskID string  `json:"task_id"`
	SadID  int64   `json:"sad_id"`
	RowId  *string `json:"row_id"`
}
