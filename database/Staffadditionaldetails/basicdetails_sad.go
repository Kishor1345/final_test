// Package databasesad interacts with the Staff Additional Details DB.
//path :/var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/Staffadditionaldetails
// --- Creator's Info ---
// Creator: Rovita
// Created On: 11-11-2025
// Description: Insert and Update operations for employee basic details.
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/Staffadditionaldetails"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/lib/pq"
)
// convertCustomDates normalizes date fields into YYYY-MM-DD format
// and recursively processes nested structures.
func convertCustomDates(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	// Convert date strings in DD-MM-YYYY format
	if s, ok := data.(string); ok {
		if t, err := time.Parse("02-01-2006", s); err == nil {
			return t.Format("2006-01-02")
		}
		return s
	}

	val := reflect.ValueOf(data)

	switch val.Kind() {
	// Handle pointer values safely
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		return convertCustomDates(val.Elem().Interface())
		
	// Handle slices and arrays
	case reflect.Slice, reflect.Array:
		out := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			out[i] = convertCustomDates(val.Index(i).Interface())
		}
		return out
	// Handle map values
	case reflect.Map:
		out := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			out[key.String()] = convertCustomDates(val.MapIndex(key).Interface())
		}
		return out
	// Handle struct fields using JSON tags
	case reflect.Struct:
		out := make(map[string]interface{})
		typ := val.Type()

		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)
			ft := typ.Field(i)

			if !f.CanInterface() {
				continue
			}

			tag := ft.Tag.Get("json")
			if tag == "" || tag == "-" {
				continue
			}

			key := tag
			if idx := commaIndex(tag); idx != -1 {
				key = tag[:idx]
			}

			out[key] = convertCustomDates(f.Interface())
		}
		return out

	default:
		return data
	}
}
// commaIndex returns the index of the first comma in a string
func commaIndex(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			return i
		}
	}
	return -1
}

// ExecuteMasterSad calls the stored procedure meivan.master_sad
func ExecuteMasterSad(req modelssad.MasterSadRequest) (modelssad.MasterSadResponse, error) {
	var resp modelssad.MasterSadResponse

	// Database connection
	db := credentials.GetDB()

	// Prepare JSON payloads
	var contactJSON, depJSON, eduJSON, expJSON, langJSON, nomJSON, hinJSON []byte
var err error
	if req.ContactData != nil {
		contactJSON, err = json.Marshal(convertCustomDates(req.ContactData))
		if err != nil {
			return resp, err
		}
	}
	if req.DependentsData != nil {
		depJSON, err = json.Marshal(convertCustomDates(req.DependentsData))
		if err != nil {
			return resp, err
		}
	}
	if req.EducationData != nil {
		eduJSON, err = json.Marshal(convertCustomDates(req.EducationData))
		if err != nil {
			return resp, err
		}
	}
	if req.ExperienceData != nil {
		expJSON, err = json.Marshal(convertCustomDates(req.ExperienceData))
		if err != nil {
			return resp, err
		}
	}
	if req.LanguageData != nil {
		langJSON, err = json.Marshal(convertCustomDates(req.LanguageData))
		if err != nil {
			return resp, err
		}
	}
	if req.NomineeData != nil {
		nomJSON, err = json.Marshal(convertCustomDates(req.NomineeData))
		if err != nil {
			return resp, err
		}
	}
	if req.HindiData != nil {
		hinJSON, err = json.Marshal(convertCustomDates(req.HindiData))
		if err != nil {
			return resp, err
		}
	}
	// Prepare nullable date fields
	var dob, doj, doc, dor, eff *sql.NullTime

	if req.Dob != nil {
		dob = &sql.NullTime{Time: req.Dob.Time, Valid: true}
	}
	if req.DateOfJoining != nil {
		doj = &sql.NullTime{Time: req.DateOfJoining.Time, Valid: true}
	}
	if req.DateOfConfirmation != nil {
		doc = &sql.NullTime{Time: req.DateOfConfirmation.Time, Valid: true}
	}
	if req.DateOfRetirement != nil {
		dor = &sql.NullTime{Time: req.DateOfRetirement.Time, Valid: true}
	}
	if req.EffectiveDate != nil {
		eff = &sql.NullTime{Time: req.EffectiveDate.Time, Valid: true}
	}

	// Stored procedure call
	query := `
		CALL meivan.master_sad(
			$1,$2,$3,$4,$5,NULL,
			$6,$7,$8,$9,$10,
			$11,$12,$13,$14,$15,
			$16,$17,$18,$19,$20,
			$21,$22,$23,$24,$25,
			$26,$27,$28,$29,$30,
			$31,$32,$33,$34,$35,
			$36,$37,$38,$39,$40,
			$41,$42,$43,$44,$45,
			$46,$47,$48,$49,$50,
			$51,$52,$53,$54,$55,
			$56,$57,$58,$59,$60,
			$61,$62,$63,$64,$65,
			$66,$67,$68,$69,$70,
			$71,$72,$73,$74,$75,$76,$77
		)
	`

	row := db.QueryRow(
		query,
		req.ActionType,
		req.TaskID,
		req.ProcessID,
		req.EmployeeID,
		req.UpdatedBy,

		req.EmployeeName,
		req.AssignTo,
		req.AssignedRole,
		req.TaskStatusID,
		req.ActivitySeqNo,
		req.IsTaskReturn,
		req.IsTaskApproved,
		req.EmailFlag,
		req.TemplateID,
		req.RejectFlag,
		req.RejectRole,
		req.InitiatedBy,
		req.Badge,
		req.Priority,
		req.Starred,
		req.FirstName,
		req.MiddleName,
		req.LastName,
		req.Gender,
		req.MaritalStatus,
		req.FatherName,
		req.MotherName,
		req.SpouseName,
		dob,
		req.Age,
		req.Nationality,
		req.Religion,
		req.CasteCategory,
		req.EmergencyContactNo,
		req.MobileNo,
		req.IsPhysicallyChallenged,
		req.PercentageOfDisability,
		req.NatureOfDisability,
		req.PersonalEmail,
		req.AadhaarNo,
		req.MotherTongue,
		req.BankName,
		req.IfscCode,
		req.BankAcctNo,
		req.IdentificationMarks,
		req.PanCardNo,
		req.EmployeeType,
		req.Department,
		req.Designation,
		req.Section,
		req.RouteTo,
		req.Grade,
		req.EmpGroup,
		req.PayInfo,
		req.BasicPay,
		req.NonPracticePay,
		req.NameOfPayBand,
		doj,
		doc,
		req.OfficeRoomNo,
		req.OfficeExtensionNo,
		req.IsActive,
		req.EmployeeStatus,
		eff,
		dor,

		nullableBytes(contactJSON),
		nullableBytes(depJSON),
		nullableBytes(eduJSON),
		nullableBytes(expJSON),
		nullableBytes(langJSON),
		nullableBytes(nomJSON),
		nullableBytes(hinJSON),
		req.Comments,
		req.UserRole,
		req.ModuleName,
		pq.Array(req.ModifiedFields),
		req.Rowid,
	)

	// Scan OUT parameters
	err = row.Scan(&resp.TaskID, &resp.SadID, &resp.RowId)
	if err != nil {

		//DO NOT wrap postgres business errors
		if pqErr, ok := err.(*pq.Error); ok {
			return resp, fmt.Errorf(pqErr.Message)
		}

		return resp, err
	}

	return resp, nil
}

// nullableBytes helper
func nullableBytes(b []byte) interface{} {
	if len(b) == 0 {
		return nil
	}
	return string(b)
}
