// Package modelssad contains structs and queries for Modified feilds.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/Staffadditionaldetails
// --- Creator's Info ---
// Creator: Rovita
//
// Created On: 29-01-2026
//
// Last Modified By:
//  
// Last Modified Date:
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/Staffadditionaldetails"
	"fmt"
	"database/sql"


	_ "github.com/lib/pq"
)

func GetSadPersonalDetailsFromDB(
	decryptedData map[string]interface{},
) ([]modelssad.SadPersonalDetails, int, error) {

	taskID, ok := decryptedData["task_id"].(string)
	if !ok || taskID == "" {
		return nil, 0, fmt.Errorf("task_id is required")
	}



	db, err := sql.Open("postgres", credentials.Getdatabasemeivan())
	if err != nil {
		return nil, 0, err
	}
	defer db.Close()

	rows, err := db.Query(
		modelssad.MyQuerySadPersonalDetails,
		taskID,
		
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	data, err := modelssad.RetrieveSadPersonalDetails(rows)
	if err != nil {
		return nil, 0, err
	}

	return data, len(data), nil
}
