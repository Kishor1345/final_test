// Package models contains data structures and database access logic for the DefaultRoleName page.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/common
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:30-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 10-11-2025
//
// Path:Login Page
package modelscommon

import (
	//	modelstable "Hrmodule/models/tables"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// var MyQueryDefaultRoleName = `
// SELECT Distinct A.Employeeid as UserID,loginname as Username,
//     B.campuscode || ' ' ||
//     CASE WHEN A.sectionid IS NOT NULL THEN D.sectioncode ELSE A.departmentcode END || ' ' || 
//     C.rolename AS RoleName
// FROM humanresources.employeerolemapping A
// JOIN humanresources.campus B ON A.campusid = B.id
// JOIN meivan.rolemaster C ON A.roleid = C.id
// LEFT JOIN humanresources.section D ON A.sectionid = D.id
// join humanresources.employeebasicinfo e
// on A.employeeid=e.employeeid
// WHERE e.loginname = $1
// `

var MyQueryDefaultRoleName = `
SELECT DISTINCT 
    A.employeeid AS UserID,
    e.loginname AS Username,

    CONCAT_WS(
        ' ',
        TRIM(B.campuscode),

        CASE 
            WHEN A.roleid = '35' THEN TRIM(A.departmentcode)
            WHEN TRIM(A.departmentcode) <> 'ADM' THEN TRIM(A.departmentcode)
        END,

        CASE 
            WHEN A.roleid <> '35' THEN TRIM(D.sectioncode)
        END,

        TRIM(C.rolename)
    ) AS RoleName

FROM humanresources.employeerolemapping A
JOIN humanresources.campus B 
    ON A.campusid = B.id
JOIN meivan.rolemaster C 
    ON A.roleid = C.id
LEFT JOIN humanresources.section D 
    ON A.sectionid = D.id
JOIN humanresources.employeebasicinfo e
    ON A.employeeid = e.employeeid
WHERE e.loginname = $1;


`


// DefaultRoleNamestructure defines the structure of DefaultRoleName
type DefaultRoleNamestructure struct {
	UserID   *string `json:"UserID"`
	Username *string `json:"Username"`
	RoleName *string `json:"RoleName"`
}

// RetrieveDefaultRoleName scans rows into DefaultRoleNamestructure slice
func RetrieveDefaultRoleName(rows *sql.Rows) ([]DefaultRoleNamestructure, error) {
	var DefaultRoleNameapi []DefaultRoleNamestructure

	for rows.Next() {
		var DRN DefaultRoleNamestructure
		err := rows.Scan(
			&DRN.UserID,
			&DRN.Username,
			&DRN.RoleName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		DefaultRoleNameapi = append(DefaultRoleNameapi, DRN)
	}

	return DefaultRoleNameapi, nil
}
