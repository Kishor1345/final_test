// Package modelscriteria contains structs and queries for  criteria data fetch for approval.
//
//Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/criteria
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:28/01/2026
package modelscriteria

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/google/uuid"
)

// CriteriaApproval represents individual criteria details
type CriteriaApproval struct {
    CriteriaID  string `json:"criteria_id"`
    Description string `json:"description"`
    GreatPay    string `json:"greatpay"`
    Level       string `json:"level"`
    GreatPayName    string `json:"greatpayname"`
    LevelName       string `json:"levelname"`
}

// CriteriaMasterDataFetchForApprovalStructure represents
type CriteriaMasterDataFetchForApprovalStructure struct{
	ReferenceNo      string         `json:"reference_no"`
	TaskID           *uuid.UUID        `json:"task_id"`
	ProcessID        int            `json:"process_id"`
	Status           int            `json:"status"`
	AssignTo         string         `json:"assign_to"`
	AssignedRole     string         `json:"assigned_role"`
	TaskStatusID     int            `json:"task_status_id"`
	ActivitySeqNo    int            `json:"activity_seq_no"`
	IsTaskReturn     int            `json:"is_task_return"`
	IsTaskApproved   int            `json:"is_task_approved"`
	EmailFlag        int            `json:"email_flag"`
	TemplateID       *int           `json:"template_id"`
	RejectFlag       int            `json:"reject_flag"`
	RejectRole       *string        `json:"reject_role"`
	InitiatedBy      string         `json:"initiated_by"`
	InitiatedOn      string         `json:"initiated_on"`
	UpdatedBy        string         `json:"updated_by"`
	UpdatedOn        string         `json:"updated_on"`
	StatusFlag       int            `json:"status_flag"`
	Badge            int            `json:"badge"`
	Priority         int            `json:"priority"`
	Starred          int            `json:"starred"`
	Criteria []CriteriaApproval `json:"criteria"` 
    ProcessMsg     string                `json:"process"`
}


// SQL query to fetch criteria master data with
// aggregated criteria using JSON
var MyQueryCriteriaMasterDataFetchForApproval = `
WITH criteria_rows AS (
       SELECT 
        cg.task_id,
        cg.criteria_id,
        cg.description,
        COALESCE(string_agg(gp.id::TEXT, ','), '') AS greatpay,
        COALESCE(string_agg(level.id::TEXT, ','), '') AS level,
        COALESCE(string_agg(gp.name::TEXT, ','), '') AS greatpayname,
        COALESCE(string_agg(level.name::TEXT, ','), '') AS levelname
    FROM meivan.cmes_g cg
    LEFT JOIN humanresources.cpc_master cm ON cm.id = cg.cpc_id
    LEFT JOIN (SELECT id,name FROM humanresources.cpc_master WHERE cpc = '6') AS gp ON cm.id = gp.id
    LEFT JOIN (SELECT id,name FROM humanresources.cpc_master WHERE cpc = '7') AS level ON cm.id = level.id
    where cg.status = 1
    GROUP BY cg.task_id, cg.criteria_id, cg.description
)

SELECT
    cm.reference_no,
    cm.task_id,
    cm.process_id,
    cm.status,
    cm.assign_to,
    cm.assigned_role,
    cm.task_status_id,
    cm.activity_seq_no,
    cm.is_task_return,
    cm.is_task_approved,
    cm.email_flag,
    cm.template_id,
    cm.reject_flag,
    cm.reject_role,
    cm.initiated_by,
    cm.initiated_on,
    cm.updated_by,
    cm.updated_on,
    cm.status_flag,
    cm.badge,
    cm.priority,
    cm.starred,
    json_agg(
        json_build_object(
            'criteria_id', cr.criteria_id,
            'description', cr.description,
            'greatpay', cr.greatpay,
            'level', cr.level,
            'greatpayname',cr.greatpayname,
            'levelname',cr.levelname
        )
    ) AS criteria
FROM meivan.cmes_m cm
JOIN criteria_rows cr ON cm.task_id = cr.task_id
WHERE 
    cm.task_id= $1 
GROUP BY
    cm.reference_no,
    cm.task_id,
    cm.process_id,
    cm.status,
    cm.assign_to,
    cm.assigned_role,
    cm.task_status_id,
    cm.activity_seq_no,
    cm.is_task_return,
    cm.is_task_approved,
    cm.email_flag,
    cm.template_id,
    cm.reject_flag,
    cm.reject_role,
    cm.initiated_by,
    cm.initiated_on,
    cm.updated_by,
    cm.updated_on,
    cm.status_flag,
    cm.badge,
    cm.priority,
    cm.starred;


`


// RetrieveCriteriaMasterDataFetchForApproval maps SQL rows
// into CriteriaMasterDataFetchStructure
func RetrieveCriteriaMasterDataFetchForApproval(rows *sql.Rows) ([]CriteriaMasterDataFetchForApprovalStructure, error) {

	var results []CriteriaMasterDataFetchForApprovalStructure

	for rows.Next() {
		var r CriteriaMasterDataFetchForApprovalStructure
		var criteriaJSON []byte

        // Scan database row into struct fields
		err := rows.Scan(
			&r.ReferenceNo,
			&r.TaskID,
			&r.ProcessID,
			&r.Status,
			&r.AssignTo,
			&r.AssignedRole,
			&r.TaskStatusID,
			&r.ActivitySeqNo,
			&r.IsTaskReturn,
			&r.IsTaskApproved,
			&r.EmailFlag,
			&r.TemplateID,
			&r.RejectFlag,
			&r.RejectRole,
			&r.InitiatedBy,
			&r.InitiatedOn,
			&r.UpdatedBy,
			&r.UpdatedOn,
			&r.StatusFlag,
			&r.Badge,
			&r.Priority,
			&r.Starred,
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

	return results, nil
}
