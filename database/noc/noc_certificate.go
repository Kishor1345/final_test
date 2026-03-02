// Package databasenoc contains data structures and database access logic for the NOC Certificate master page.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/noc
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On: 07-01-2026
//
// Last Modified By:
//
// Last Modified Date:
package databasenoc

import (
	credentials "Hrmodule/dbconfig"
	modelsnoc "Hrmodule/models/noc"
	"fmt"

	_ "github.com/lib/pq"
)

// NocCertificateDatabase fetches NOC Certificate master details from the database
func NocCertificateDatabase(decryptedData map[string]interface{}) (interface{}, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract inputs (optional based on API call)
	certificateType, _ := decryptedData["certificate_type"].(string)
	purpose, _ := decryptedData["purpose"].(string)

	// Extract employee_id as int
	var employeeID int
	if empIDFloat, ok := decryptedData["employee_id"].(float64); ok {
		employeeID = int(empIDFloat)
	} else if empIDStr, ok := decryptedData["employee_id"].(string); ok {
		fmt.Sscanf(empIDStr, "%d", &employeeID)
	}
	// -------------------------------------------------------------------------
	//  CASE 1: Additional Options (highest priority)
	// -------------------------------------------------------------------------
	if purpose != "" {

		rows, err := db.Query(modelsnoc.MyQueryNocCertificateAdditionalOption, purpose)
		if err != nil {
			return nil, 0, fmt.Errorf("error querying additional options: %v", err)
		}
		defer rows.Close()

		data, err := modelsnoc.RetrieveNocCertificateAdditionalOptions(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("error retrieving additional options: %v", err)
		}

		return data, len(data), nil
	}

	// -------------------------------------------------------------------------
	// CASE 2: Purpose based on Certificate Type and Employee ID
	// -------------------------------------------------------------------------
	if certificateType != "" {

		// Convert employeeID to string with leading zeros (6 digits)
		employeeIDStr := fmt.Sprintf("%06d", employeeID)

		rows, err := db.Query(modelsnoc.MyQueryNocCertificatePurpose, employeeIDStr, certificateType)
		if err != nil {
			return nil, 0, fmt.Errorf("error querying certificate purposes: %v", err)
		}
		defer rows.Close()

		data, err := modelsnoc.RetrieveNocCertificatePurposes(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("error retrieving certificate purposes: %v", err)
		}

		return data, len(data), nil
	}

	// -------------------------------------------------------------------------
	// CASE 3: Certificate Types (default)
	// -------------------------------------------------------------------------
	rows, err := db.Query(modelsnoc.MyQueryNocCertificateType)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying certificate types: %v", err)
	}
	defer rows.Close()

	data, err := modelsnoc.RetrieveNocCertificateTypes(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving certificate types: %v", err)
	}

	return data, len(data), nil
}
