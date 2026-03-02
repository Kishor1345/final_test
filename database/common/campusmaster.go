// Package databasecommon contains data structures and database access logic for the Campus Master page.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/common
// --- Creator's Info ---
// Creator: Ramya M R
//
// Created On:10-02-2026
//
// Last Modified By:
//
// Last Modified Date:
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"fmt"

	_ "github.com/lib/pq"
)

func CampusMasterdatabase(decryptedData map[string]interface{}) ([]modelscommon.CampusMasterStructure, int, error) {

	// Database connection
	db := credentials.GetDB()

	rows, err := db.Query(modelscommon.MyQueryCampusMaster)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	data, err := modelscommon.RetrieveCampusMaster(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return data, len(data), nil
}
