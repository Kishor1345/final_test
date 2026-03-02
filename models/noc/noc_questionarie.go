// Package modelsnoc provides data structures and SQL queries for the NOC (No Objection Certificate) module.
// It specifically handles questionnaire data used for processes like application forwarding and vigilance clearances.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/noc
//
// --- Creator's Info ---
// Creator: Elakiya
// Created On: 21-01-2026
// Last Modified By: Vaishnavi
// Last Modified Date: 22-01-2026
package modelsnoc

import (
	"database/sql"
	"fmt"
)

// MyQueryQuestionnaireForwarding retrieves active questionnaire records based on a provided 
// list of question types. It converts a comma-separated string input into an array 
// to filter the question_type column dynamically.
var MyQueryQuestionnaireForwarding = `
	SELECT 
    id,
    question_type,
    question, 
    keyvalue,
    type
FROM humanresources.questionarie
WHERE status = 1
  AND question_type = ANY(
      SELECT trim(value) 
      FROM unnest(string_to_array($1, ',')) AS value
  )
ORDER BY id ASC;

`

// QuestionnaireStructure represents the data model for a single questionnaire item,
// mapping the database fields to JSON for use in API responses.
type QuestionnaireStructure struct {
	ID           int64  `json:"id"`
	QuestionType string `json:"question_type"`
	Question     string `json:"question"`
	KeyValue     string `json:"keyvalue"`
	Type         string `json:"type"`
}

// RetrieveQuestionnaire processes database result rows into a slice of QuestionnaireStructure.
//
// It iterates through the provided *sql.Rows, mapping each column to the corresponding 
// field in the struct. Returns a list of questionnaires or an error if row scanning fails.
func RetrieveQuestionnaire(rows *sql.Rows) ([]QuestionnaireStructure, error) {

	var list []QuestionnaireStructure

	for rows.Next() {
		var q QuestionnaireStructure

		if err := rows.Scan(
			&q.ID,
			&q.QuestionType,
			&q.Question,
			&q.KeyValue,
			&q.Type,
		); err != nil {
			return nil, fmt.Errorf("error scanning questionnaire data: %v", err)
		}

		list = append(list, q)
	}

	return list, nil
}
