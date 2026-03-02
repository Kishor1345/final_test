// Package modelscriteria contains structs and queries for Existing data fetch ion domain table.
//
//Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/criteria
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:23/01/2026
package modelscriteria

import (
	"database/sql"
	"encoding/json"

	_ "github.com/lib/pq"
	"github.com/google/uuid"
)

// CriteriaForExisting represents individual criteria details
type CriteriaForExisting struct {
	CriteriaID  string `json:"criteria_id"`
	Description string `json:"description"`
	GreatPay    string `json:"greatpay"`
	Level       string `json:"level"`
}

// CriteriaMasterExistingDataFetchStructure represents
type CriteriaMasterExistingDataFetchStructure struct {
	TaskID      *uuid.UUID            `json:"task_id"`
	ReferenceNo string                `json:"reference_no"`
	Status      int                   `json:"status"`
	CreatedBy   string                `json:"created_by"`
	CreatedAt   string                `json:"created_at"`
	UpdatedBy   string                `json:"updated_by"`
	UpdatedAt   string                `json:"updated_at"`
	Criteria    []CriteriaForExisting `json:"criteria"`
    ProcessMsg     string                `json:"process"`
}

// SQL query to fetch criteria master existing data with
// aggregated criteria using  JSON
var MyQueryCriteriaMasterExistingDataFetch = `
WITH criteria_rows AS (
    SELECT 
        hcm.task_id,
        hcm.reference_no,
        hcm.criteria_id,
        hcm.description,
        hcm.status,
        hcm.created_by,
        hcm.created_at,
        hcm.updated_by,
        hcm.updated_at,
        COALESCE(string_agg(gp.id::TEXT, ','), '') AS greatpay,
        COALESCE(string_agg(lv.id::TEXT, ','), '') AS level
    FROM humanresources.criteria_master hcm
    LEFT JOIN humanresources.cpc_master cm 
        ON cm.id = hcm.cpc_id
    LEFT JOIN humanresources.cpc_master gp 
        ON gp.id = cm.id AND gp.cpc = '6'
    LEFT JOIN humanresources.cpc_master lv 
        ON lv.id = cm.id AND lv.cpc = '7'
    WHERE hcm.status = '1'
    GROUP BY
        hcm.task_id,
        hcm.reference_no,
        hcm.criteria_id,
        hcm.description,
        hcm.status,
        hcm.created_by,
        hcm.created_at,
        hcm.updated_by,
        hcm.updated_at
)
SELECT
    task_id,
    reference_no,
    status,
    created_by,
    created_at,
    updated_by,
    updated_at,
    json_agg(
        json_build_object(
            'criteria_id', criteria_id,
            'description', description,
            'greatpay', greatpay,
            'level', level
        ) ORDER BY criteria_id
    ) AS criteria
FROM criteria_rows
GROUP BY
    task_id,
    reference_no,
    status,
    created_by,
    created_at,
    updated_by,
    updated_at;
`

// RetrieveOngoingCount checks whether
// any task is currently in ongoing status
func RetrieveOngoingCount(db *sql.DB) (int, error) {
	query := `
		SELECT COUNT(1)
		FROM meivan.cmes_m
		WHERE task_status_id = 4
	`
	var count int
	err := db.QueryRow(query).Scan(&count)
	return count, err
}

// RetrieveCriteriaMasterExistingDataFetch maps SQL rows
// into CriteriaMasterExistingDataFetchStructure
func RetrieveCriteriaMasterExistingDataFetch(
	rows *sql.Rows,
) ([]CriteriaMasterExistingDataFetchStructure, error) {
	var results []CriteriaMasterExistingDataFetchStructure

	for rows.Next() {
		var r CriteriaMasterExistingDataFetchStructure
		var criteriaJSON []byte// Holds JSON aggregated criteria

        // Scan database row into struct fields
		err := rows.Scan(
			&r.TaskID,
			&r.ReferenceNo,
			&r.Status,
			&r.CreatedBy,
			&r.CreatedAt,
			&r.UpdatedBy,
			&r.UpdatedAt,
			&criteriaJSON,
		)
		if err != nil {
			return nil, err
		}

        // Unmarshal JSON criteria into struct
		if err := json.Unmarshal(criteriaJSON, &r.Criteria); err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	//SAFETY CHECK
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
