// Package modelsnoc provides data structures and SQL queries for the NOC (No Objection Certificate) module.
// It specifically handles questionnaire data used for processes like application forwarding and vigilance clearances.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/noc
//
// --- Creator's Info ---
// Creator: Elakiya
// Created On: 21-01-2026
// Last Modified By: Vaishnavi
// Last Modified Date: 22-01-2026
package databasenoc

import (
	"fmt"

	credentials "Hrmodule/dbconfig"
	modelsnoc "Hrmodule/models/noc"

	_ "github.com/lib/pq"
)

func QuestionnaireDynamicFromDB(
	decryptedData map[string]interface{},
) ([]modelsnoc.QuestionnaireStructure, int, error) {

	// 1 Open DB connection
	db := credentials.GetDB()

	// 2 Extract QuestionType (REQUIRED)
	questionType, ok := decryptedData["QuestionType"].(string)
	if !ok || questionType == "" {
		return nil, 0, fmt.Errorf("missing or invalid 'QuestionType'")
	}

	// 4 Execute query
	rows, err := db.Query(
		modelsnoc.MyQueryQuestionnaireForwarding,
		questionType,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	// 5 Map rows  struct
	result, err := modelsnoc.RetrieveQuestionnaire(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("RetrieveQuestionnaire failed: %v", err)
	}

	return result, len(result), nil
}
