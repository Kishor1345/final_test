// Package databasestatusmaster handles DB access for Status Master.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/common
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 29-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 29-09-2025
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"fmt"
)

func GetStatusMasternewFromDB(decryptedData map[string]interface{}) ([]modelscommon.StatusMasternewStruct, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract order_type_id from decrypted data
	StatusDescription, ok := decryptedData["statusdescription"].(string)
	if !ok || StatusDescription == "" {
		return nil, 0, fmt.Errorf("missing 'statusdescription' in request data")
	}

	rows, err := db.Query(modelscommon.MyQueryStatusMasternew, StatusDescription)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelscommon.RetrieveStatusMasternew(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
