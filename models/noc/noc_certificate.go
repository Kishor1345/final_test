// Package modelsnoc contains data structures and database access logic  for NOC Certificate master details.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/noc
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

	_ "github.com/lib/pq"
)

// -----------------------------------------------------------------------------
// QUERY 1: Certificate Types (No WHERE condition)
// -----------------------------------------------------------------------------
var MyQueryNocCertificateType = `
	SELECT DISTINCT
		certificateid,
		certificate_type
	FROM humanresources.certificate_master
	ORDER BY certificateid ASC
`

// Struct for Certificate Type
type NocCertificateTypeStructure struct {
	CertificateID   int    `json:"certificate_id"`
	CertificateType string `json:"certificate_type"`
}

// Row Mapper for Certificate Type
func RetrieveNocCertificateTypes(rows *sql.Rows) ([]NocCertificateTypeStructure, error) {
	var list []NocCertificateTypeStructure

	for rows.Next() {

		var noc NocCertificateTypeStructure

		err := rows.Scan(
			&noc.CertificateID,
			&noc.CertificateType,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning noc certificate type data: %v", err)
		}

		list = append(list, noc)
	}

	return list, nil
}

// -----------------------------------------------------------------------------
// QUERY 2: Purpose based on Certificate Type and Employee ID
// -----------------------------------------------------------------------------
var MyQueryNocCertificatePurpose = `
	SELECT DISTINCT
		cm.purposeid,
		cm.purpose,
		CASE 
			WHEN nm.task_status_id = 22 THEN 'Save and Hold'
			WHEN nm.task_status_id = 4 THEN 'Ongoing'
			ELSE NULL
		END AS status_message,
		nm.task_status_id
	FROM humanresources.certificate_master cm
	LEFT JOIN meivan.noc_m nm 
		ON cm.certificateid = nm.certificate_type 
		AND cm.purposeid = nm.purpose
		AND nm.employee_id = $1
	WHERE cm.certificate_type = $2
	ORDER BY cm.purposeid ASC
`

// Struct for Certificate Purpose
type NocCertificatePurposeStructure struct {
	PurposeID     int     `json:"purpose_id"`
	Purpose       string  `json:"purpose"`
	StatusMessage *string `json:"status_message"`
	TaskStatusID  *int    `json:"task_status_id"`
}

// Row Mapper for Certificate Purpose
func RetrieveNocCertificatePurposes(rows *sql.Rows) ([]NocCertificatePurposeStructure, error) {
	var list []NocCertificatePurposeStructure

	for rows.Next() {
		var noc NocCertificatePurposeStructure

		err := rows.Scan(
			&noc.PurposeID,
			&noc.Purpose,
			&noc.StatusMessage,
			&noc.TaskStatusID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning noc certificate purpose data: %v", err)
		}

		list = append(list, noc)
	}

	return list, nil
}

// -----------------------------------------------------------------------------
// QUERY 3: Additional Options based on Purpose
// -----------------------------------------------------------------------------
var MyQueryNocCertificateAdditionalOption = `
	SELECT DISTINCT
		additionaloptionid,
		additional_option
	FROM humanresources.certificate_master
	WHERE purpose = $1
	ORDER BY additionaloptionid ASC
`

// Struct for Additional Options
type NocCertificateAdditionalOptionStructure struct {
	AdditionalOptionID int    `json:"additional_option_id"`
	AdditionalOption   string `json:"additional_option"`
}

// Row Mapper for Additional Options
func RetrieveNocCertificateAdditionalOptions(rows *sql.Rows) ([]NocCertificateAdditionalOptionStructure, error) {
	var list []NocCertificateAdditionalOptionStructure

	for rows.Next() {

		var noc NocCertificateAdditionalOptionStructure

		err := rows.Scan(
			&noc.AdditionalOptionID,
			&noc.AdditionalOption,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning noc certificate additional option data: %v", err)
		}

		list = append(list, noc)
	}

	return list, nil
}
