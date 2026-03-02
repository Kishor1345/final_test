// Package models contains data structures and database access logic for the Employeedetails page.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/common
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:04-12-2025
//
// Last Modified By:
//
// Last Modified Date:
package modelscommon

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var MyQueryEmployeeDetails = `
SELECT 
    e.employeeid,

    CASE 
        WHEN c.combovalue IS NOT NULL 
        THEN c.combovalue || ' ' || e.firstname
        ELSE e.firstname
    END AS firstname,

    e.firstname AS firstnameoriginal,
    e.fatherorhusbandname,
    TO_CHAR(a.doj, 'DD-MM-YYYY') AS doj,
    dm.departmentname,
    dg.designationname,
    a.employeegroup,
		e.passportnumber
FROM humanresources.employeebasicinfo e

INNER JOIN humanresources.employeeappointmentdetails a
    ON e.employeeid = a.employeeid

LEFT JOIN humanresources.combovaluesmaster c
       ON c.comboname = 'Suffix'
      AND c.displayseq = e.suffix

LEFT JOIN humanresources.departmentdesignationmapping m
       ON e.employeeid = m.employeeid

LEFT JOIN humanresources.departmentmaster dm
       ON m.departmentcode = dm.departmentcode

LEFT JOIN humanresources.designationmaster dg
       ON m.designationid = dg.designationid

WHERE e.employeeid = $1
`

// Struct for Employee Details
type EmployeeDetailsStructure struct {
	EmployeeID      *string `json:"EmployeeID"`
	FacultyName     *string `json:"FacultyName"`
	FacultyOriginal *string `json:"FacultyOriginal"`
	FatherName      *string `json:"FatherName"`
	Doj             *string `json:"Doj"`
	DepartmentName  *string `json:"DepartmentName"`
	DesignationName *string `json:"DesignationName"`
	Employeegroup   *string `json:"Employeegroup"`
	Passportnumber  *string `json:"passportnumber"`
}

// Row Mapper
func RetrieveEmployeeDetails(rows *sql.Rows) ([]EmployeeDetailsStructure, error) {
	var list []EmployeeDetailsStructure

	for rows.Next() {
		var emp EmployeeDetailsStructure
		err := rows.Scan(
			&emp.EmployeeID,
			&emp.FacultyName,
			&emp.FacultyOriginal,
			&emp.FatherName,
			&emp.Doj,
			&emp.DepartmentName,
			&emp.DesignationName,
			&emp.Employeegroup,
			&emp.Passportnumber,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		list = append(list, emp)
	}

	return list, nil
}
