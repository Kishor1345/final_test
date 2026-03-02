// Package databasecommon handles database connections and queries related to DefaultRoleName data.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/common
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:30-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 30-07-2025
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

func DefaultRoleNamedatabase(decryptedData map[string]interface{}) ([]modelscommon.DefaultRoleNamestructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	// Extract order_type_id from decrypted data
	UserName, ok := decryptedData["UserName"].(string)
	if !ok || UserName == "" {
		return nil, 0, fmt.Errorf("missing 'UserName' in request data")
	}

	// Execute the query
	rows, err := db.Query(modelscommon.MyQueryDefaultRoleName, UserName)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map results
	DefaultRoleNameapi, err := modelscommon.RetrieveDefaultRoleName(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return DefaultRoleNameapi, len(DefaultRoleNameapi), nil
}
