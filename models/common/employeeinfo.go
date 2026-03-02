// Package modelscommon contains data structures and database access logic for the EmployeeBasicInfo page.
//
// Path : /var/www/html/go_projects/HRMODULE/Ramya/Hr_test7007/models/common
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 05-01-2026
//
// Last Modified By: Kishorekumar
//
// Last Modified Date: 29-01-2026
/*
package modelscommon

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Query for retrieving employee basic info
var MyQueryEmployeeBasicInfo = `
  SELECT
    e.employeeid AS "EmployeeId",
    e.displayname AS "Employee Name",
    TO_CHAR(e.dob, 'DD-MM-YYYY') AS "Date of Birth",
    e.fatherorhusbandname AS "Father's Name",
       e.spousename as "Spouse",
    c.name AS "Category",

    TO_CHAR(a.doj, 'DD-MM-YYYY') AS "Date of Joining",
    TO_CHAR(a.retirementdate, 'DD-MM-YYYY') AS "Date of Retirement",

    cv.combovalue AS "Employee Type",

    dg.designationname AS "Designation",
    dm.departmentname AS "Department",

    ep.paybandscale || ' / ' || ep.grade AS "Grade Pay & Pay Info",

    TO_CHAR(
        TO_DATE(ec.dateofappointment, 'DD/MM/YYYY'),
        'DD-MM-YYYY'
    ) AS "Date of Appointment",

    TO_CHAR(p.effectivefrom, 'DD-MM-YYYY')
        AS "Date of Appointment in Present Post",

    TO_CHAR(ec.confirmationeffectivefrom, 'DD-MM-YYYY')
        AS "Date of Confirmation",

    TO_CHAR(
        TO_DATE(ec.probationperiodfrom, 'DD/MM/YYYY'),
        'DD-MM-YYYY'
    )
    || ' - ' ||
    TO_CHAR(
        TO_DATE(ec.probationperiodto, 'DD/MM/YYYY'),
        'DD-MM-YYYY'
    )
    AS "Probation period from - to",

    -- ✅ Service at IITM
    TRIM(
        EXTRACT(YEAR FROM AGE(CURRENT_DATE, a.doj))::INT || ' Years ' ||
        EXTRACT(MONTH FROM AGE(CURRENT_DATE, a.doj))::INT || ' Months ' ||
        EXTRACT(DAY FROM AGE(CURRENT_DATE, a.doj))::INT || ' Days'
    ) AS "Service at IITM",

    -- ✅ New Column Added
    a.employeegroup AS "Employee Group",

	TO_CHAR(
    TO_DATE(NULLIF(ec.poshtrainingcompleted,''), 'DD/MM/YYYY'),
    'DD-MM-YYYY'
    ) AS "Posh Training Completed On",

    TO_CHAR(
    TO_TIMESTAMP(NULLIF(ec.iprcompletedon,''), 'Mon DD YYYY HH:MIAM'),
    'DD-MM-YYYY'
    ) AS "IPR Completed On"

FROM humanresources.employeebasicinfo e

LEFT JOIN humanresources.castecategory c
       ON e.caste = c.id

LEFT JOIN humanresources.employeeappointmentdetails a
       ON e.employeeid = a.employeeid

LEFT JOIN humanresources.combovaluesmaster cv
       ON a.employeetype = cv.displayseq
      AND cv.comboname = 'EmployeeType'

LEFT JOIN humanresources.employeepresentscalemaster ep
       ON a.presentscaleid = ep.id

LEFT JOIN humanresources.departmentdesignationmapping ddm
       ON e.employeeid = ddm.employeeid

LEFT JOIN humanresources.departmentmaster dm
       ON ddm.departmentcode = dm.departmentcode

LEFT JOIN humanresources.designationmaster dg
       ON ddm.designationid = dg.designationid

LEFT JOIN humanresources.employeeposting p
       ON e.employeeid = p.employeeid::VARCHAR

LEFT JOIN humanresources.employeeconfirmation ec
       ON e.employeeid = ec.employeeid

WHERE e.employeeid = $1
`

// ✅ Struct for Employee Basic Info (16 fields now)
type EmployeeBasicInfoStructure struct {
	EmployeeID                     *string `json:"EmployeeID"`
	EmployeeName                   *string `json:"EmployeeName"`
	DateOfBirth                    *string `json:"DateOfBirth"`
	FatherName                     *string `json:"FatherName"`
	Spouse                         *string `json:"Spouse"`
	Category                       *string `json:"Category"`
	DateOfJoining                  *string `json:"DateOfJoining"`
	DateOfRetirement               *string `json:"DateOfRetirement"`
	EmployeeType                   *string `json:"EmployeeType"`
	Designation                    *string `json:"Designation"`
	Department                     *string `json:"Department"`
	GradePayAndPayInfo             *string `json:"GradePayAndPayInfo"`
	DateOfAppointment              *string `json:"DateOfAppointment"`
	DateOfAppointmentInPresentPost *string `json:"DateOfAppointmentInPresentPost"`
	DateOfConfirmation             *string `json:"DateOfConfirmation"`
	ProbationPeriod                *string `json:"ProbationPeriod"`
	ServiceAtIITM                  *string `json:"ServiceAtIITM"`
	EmployeeGroup                  *string `json:"EmployeeGroup"`
	PoshTrainingCompletedOn        *string `json:"PoshTrainingCompletedOn"`
	IPRCompletedOn                 *string `json:"IPRCompletedOn"`
}

// ✅ Row Mapper
func RetrieveEmployeeBasicInfo(rows *sql.Rows) ([]EmployeeBasicInfoStructure, error) {
	var list []EmployeeBasicInfoStructure

	for rows.Next() {
		var emp EmployeeBasicInfoStructure

		err := rows.Scan(
			&emp.EmployeeID,
			&emp.EmployeeName,
			&emp.DateOfBirth,
			&emp.FatherName,
			&emp.Spouse,
			&emp.Category,
			&emp.DateOfJoining,
			&emp.DateOfRetirement,
			&emp.EmployeeType,
			&emp.Designation,
			&emp.Department,
			&emp.GradePayAndPayInfo,
			&emp.DateOfAppointment,
			&emp.DateOfAppointmentInPresentPost,
			&emp.DateOfConfirmation,
			&emp.ProbationPeriod,
			&emp.ServiceAtIITM,
			&emp.EmployeeGroup,
			&emp.PoshTrainingCompletedOn,
			&emp.IPRCompletedOn,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		list = append(list, emp)
	}

	return list, nil
}
*/

package modelscommon

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Query for retrieving employee basic info
var MyQueryEmployeeBasicInfo = `
  SELECT
    e.employeeid AS "EmployeeId",
    e.displayname AS "Employee Name",
    TO_CHAR(e.dob, 'DD-MM-YYYY') AS "Date of Birth",
    e.fatherorhusbandname AS "Father's Name",
    e.spousename as "Spouse",
    c.name AS "Category",

    TO_CHAR(a.doj, 'DD-MM-YYYY') AS "Date of Joining",
    TO_CHAR(a.retirementdate, 'DD-MM-YYYY') AS "Date of Retirement",

    cv.combovalue AS "Employee Type",

    dg.designationname AS "Designation",
    dm.departmentname AS "Department",

    ep.paybandscale || ' / ' || ep.grade AS "Grade Pay & Pay Info",

    TO_CHAR(
        TO_DATE(ec.dateofappointment, 'DD/MM/YYYY'),
        'DD-MM-YYYY'
    ) AS "Date of Appointment",

    TO_CHAR(p.effectivefrom, 'DD-MM-YYYY')
        AS "Date of Appointment in Present Post",

    TO_CHAR(ec.confirmationeffectivefrom, 'DD-MM-YYYY')
        AS "Date of Confirmation",

    TO_CHAR(
        TO_DATE(ec.probationperiodfrom, 'DD/MM/YYYY'),
        'DD-MM-YYYY'
    )
    || ' - ' ||
    TO_CHAR(
        TO_DATE(ec.probationperiodto, 'DD/MM/YYYY'),
        'DD-MM-YYYY'
    )
    AS "Probation period from - to",

    -- ✅ Service at IITM
    TRIM(
        EXTRACT(YEAR FROM AGE(CURRENT_DATE, a.doj))::INT || ' Years ' ||
        EXTRACT(MONTH FROM AGE(CURRENT_DATE, a.doj))::INT || ' Months ' ||
        EXTRACT(DAY FROM AGE(CURRENT_DATE, a.doj))::INT || ' Days'
    ) AS "Service at IITM",

    -- ✅ Employee Group
    a.employeegroup AS "Employee Group",

    TO_CHAR(
        TO_DATE(NULLIF(ec.poshtrainingcompleted,''), 'DD/MM/YYYY'),
        'DD-MM-YYYY'
    ) AS "Posh Training Completed On",

    TO_CHAR(
        TO_TIMESTAMP(NULLIF(ec.iprcompletedon,''), 'Mon DD YYYY HH:MIAM'),
        'DD-MM-YYYY'
    ) AS "IPR Completed On",

    -- ✅ Role Name (New Column)
    (
        SELECT RoleName
        FROM (
            -- Check for default role first
            SELECT 
                edr.defaultrole AS RoleName,
                1 as priority
            FROM humanresources.employeedefaultroles edr
            WHERE edr.user_id = e.employeeid
            
            UNION ALL
            
            -- Get alphabetically first role if no default exists
            SELECT DISTINCT 
                B.campuscode || ' ' ||
                CASE WHEN A.sectionid IS NOT NULL THEN D.sectioncode ELSE A.departmentcode END || ' ' || 
                C.rolename AS RoleName,
                2 as priority
            FROM humanresources.employeerolemapping A
            JOIN humanresources.campus B ON A.campusid = B.id
            JOIN meivan.rolemaster C ON A.roleid = C.id
            LEFT JOIN humanresources.section D ON A.sectionid = D.id
            WHERE A.employeeid = e.employeeid
            AND NOT EXISTS (
                SELECT 1 FROM humanresources.employeedefaultroles edr2
                WHERE edr2.user_id = e.employeeid
            )
            ORDER BY RoleName
            LIMIT 1
        ) role_query
        ORDER BY priority
        LIMIT 1
    ) AS "Role Name"

FROM humanresources.employeebasicinfo e

LEFT JOIN humanresources.castecategory c
       ON e.caste = c.id

LEFT JOIN humanresources.employeeappointmentdetails a
       ON e.employeeid = a.employeeid

LEFT JOIN humanresources.combovaluesmaster cv
       ON a.employeetype = cv.displayseq
      AND cv.comboname = 'EmployeeType'

LEFT JOIN humanresources.employeepresentscalemaster ep
       ON a.presentscaleid = ep.id

LEFT JOIN humanresources.departmentdesignationmapping ddm
       ON e.employeeid = ddm.employeeid

LEFT JOIN humanresources.departmentmaster dm
       ON ddm.departmentcode = dm.departmentcode

LEFT JOIN humanresources.designationmaster dg
       ON ddm.designationid = dg.designationid

LEFT JOIN humanresources.employeeposting p
       ON e.employeeid = p.employeeid::VARCHAR

LEFT JOIN humanresources.employeeconfirmation ec
       ON e.employeeid = ec.employeeid

WHERE e.employeeid = $1
`

// ✅ Struct for Employee Basic Info (16 fields now)
type EmployeeBasicInfoStructure struct {
	EmployeeID                     *string `json:"EmployeeID"`
	EmployeeName                   *string `json:"EmployeeName"`
	DateOfBirth                    *string `json:"DateOfBirth"`
	FatherName                     *string `json:"FatherName"`
	Spouse                         *string `json:"Spouse"`
	Category                       *string `json:"Category"`
	DateOfJoining                  *string `json:"DateOfJoining"`
	DateOfRetirement               *string `json:"DateOfRetirement"`
	EmployeeType                   *string `json:"EmployeeType"`
	Designation                    *string `json:"Designation"`
	Department                     *string `json:"Department"`
	GradePayAndPayInfo             *string `json:"GradePayAndPayInfo"`
	DateOfAppointment              *string `json:"DateOfAppointment"`
	DateOfAppointmentInPresentPost *string `json:"DateOfAppointmentInPresentPost"`
	DateOfConfirmation             *string `json:"DateOfConfirmation"`
	ProbationPeriod                *string `json:"ProbationPeriod"`
	ServiceAtIITM                  *string `json:"ServiceAtIITM"`
	EmployeeGroup                  *string `json:"EmployeeGroup"`
	PoshTrainingCompletedOn        *string `json:"PoshTrainingCompletedOn"`
	IPRCompletedOn                 *string `json:"IPRCompletedOn"`
	RoleName                       *string `json:"RoleName"` // ✅ New field
}

// ✅ Row Mapper
func RetrieveEmployeeBasicInfo(rows *sql.Rows) ([]EmployeeBasicInfoStructure, error) {
	var list []EmployeeBasicInfoStructure

	for rows.Next() {
		var emp EmployeeBasicInfoStructure

		err := rows.Scan(
			&emp.EmployeeID,
			&emp.EmployeeName,
			&emp.DateOfBirth,
			&emp.FatherName,
			&emp.Spouse,
			&emp.Category,
			&emp.DateOfJoining,
			&emp.DateOfRetirement,
			&emp.EmployeeType,
			&emp.Designation,
			&emp.Department,
			&emp.GradePayAndPayInfo,
			&emp.DateOfAppointment,
			&emp.DateOfAppointmentInPresentPost,
			&emp.DateOfConfirmation,
			&emp.ProbationPeriod,
			&emp.ServiceAtIITM,
			&emp.EmployeeGroup,
			&emp.PoshTrainingCompletedOn,
			&emp.IPRCompletedOn,
			&emp.RoleName, // ✅ New field
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		list = append(list, emp)
	}

	return list, nil
}
