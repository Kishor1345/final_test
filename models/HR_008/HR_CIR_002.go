// Package modelscircular contains structs and queries for  circular details for approval.
//
// Path : /var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/models/HR_008
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:16/02/2026
package modelcircular

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type GridDataForCircularApproval struct {
	QuartersCategory string `json:"quarters_category"`
	QuartersNo       string `json:"quarters_no"`
	Floor            string `json:"quarters_floor"`
	Location         string `json:"location"`
	Password         string `json:"password"`
	KeyBox           string `json:"key_box"`
	FirstChoice      string `json:"first_choice"`
	SecondChoice     string `json:"second_choice"`
	ThirdChoice      string `json:"third_choice"`
	Campus           string `json:"campus"`
}

type CircularDetailForApprovalStructure struct {
	//Circular Master
	TaskId                   uuid.UUID `json:"task_id"`
	CriteriaType             string    `json:"criteria_type"`
	CircularFor              string    `json:"circular_for"`
	RegisterationOpeningDate string    `json:"open_date_registration"`
	CancellationDate         string    `json:"last_date_cancellation"`
	ProcessID                int       `json:"process_id"`
	ActivitySeqNo            int       `json:"activity_seq_no"`
	AssignTo                 string    `json:"assign_to"`
	AssignedRole             string    `json:"assigned_role"`
	IsTaskReturn             int       `json:"is_task_return"`
	IsTaskApproved           int       `json:"is_task_approved"`
	RejectFlag               int       `json:"reject_flag"`
	RejectRole               *string   `json:"reject_role"`
	Priority                 int       `json:"priority"`
	Starred                  int       `json:"starred"`
	Badge                    int       `json:"badge"`
	EmailFlag                int       `json:"email_flag"`
	TaskStatusID             int       `json:"task_status_id"`
	Status                   int       `json:"status"`
	InitiatedBy              string    `json:"initiated_by"`
	InitiatedOn              string    `json:"initiated_on"`
	UpdatedBy                string    `json:"updated_by"`
	UpdatedOn                string    `json:"updated_on"`
	NoOpenDay                string    `json:"no_of_open_days"`
	NoCancellationDate       string    `json:"no_of_cancellation_days"`
	ReferenceNO              string    `json:"reference_no"`

	GridDataApproval []GridDataForCircularApproval `json:"circulargriddataapproval"`

	//Circular Template
	OrderNo       *string `json:"order_no"`
	OrderDate     *string `json:"order_date"`
	HeaderHTML    string  `json:"header_html"`
	ToColumn      string  `json:"to_column"`
	Subject       string  `json:"subject"`
	Reference     string  `json:"reference"`
	BodyHTML      string  `json:"body_html"`
	SignatureHTML string  `json:"signature_html"`
	CCTo          string  `json:"cc_to"`
	FooterHTML    string  `json:"footer_html"`
}

var MyQueryForCircularDataFetchForApproval = `
SELECT 
--circular_master
mcm.task_id,
mcm.criteria_type,
mcm.circular_for,
meivan.globaldate_format(open_date_registration::text)as open_date_registration,
meivan.globaldate_format(last_date_cancellation::text)as last_date_cancellation ,
mcm.process_id,
mcm.activity_seq_no,
mcm.assign_to,
mcm.assigned_role,
mcm.is_task_return,
mcm.is_task_approved,
mcm.reject_flag,
mcm.reject_role,
mcm.priority,
mcm.starred,
mcm.badge,
mcm.email_flag,
mcm.task_status_id,
mcm.status,
mcm.initiated_by,
mcm.initiated_on,
mcm.updated_by,
mcm.updated_on,
mcm.no_of_open_days,
mcm.no_of_cancellation_days,
mcm.reference_no,

--Grid Details
json_agg(
    json_build_object(
        'quarters_category', quarters_category,
        'quarters_no', quarters_no,
		'floor', floor,
        'location',  location,
        'password',password,
        'key_box',  key_box,
        'first_choice',  first_choice,
		'second_choice', second_choice,
		'third_choice',third_choice,
		'campus',campus
    )
) as circulargriddata,

--circular template
mct.order_no,
mct.order_date,
mct.header_html,
mct.to_column,
mct.subject,
mct.reference,
mct.body_html,
mct.signature_html,
mct.cc_to,
mct.footer_html
FROM meivan.clrm_m mcm 
JOIN meivan.clrm_g mcg ON mcm.task_id = mcg.task_id
JOIN meivan.clrm_templates mct ON mcg.task_id = mct.task_id
WHERE mcm.task_id = $1 and mcg.status = 1
group by mcm.task_id,
mcm.criteria_type,
mcm.circular_for,
open_date_registration,
last_date_cancellation ,
mcm.process_id,
mcm.activity_seq_no,
mcm.assign_to,
mcm.assigned_role,
mcm.is_task_return,
mcm.is_task_approved,
mcm.reject_flag,
mcm.reject_role,
mcm.priority,
mcm.starred,
mcm.badge,
mcm.email_flag,
mcm.task_status_id,
mcm.status,
mcm.initiated_by,
mcm.initiated_on,
mcm.updated_by,
mcm.updated_on,
mcm.no_of_open_days,
mcm.no_of_cancellation_days,
mcm.reference_no,
mct.order_no,
mct.order_date,
mct.header_html,
mct.to_column,
mct.subject,
mct.reference,
mct.body_html,
mct.signature_html,
mct.cc_to,
mct.footer_html
`

func RetrieveCircularDetailFetchForApproval(rows *sql.Rows) ([]CircularDetailForApprovalStructure, error) {

	var results []CircularDetailForApprovalStructure

	for rows.Next() {
		var r CircularDetailForApprovalStructure
		var cricularJSON []byte
		// Scan database row into struct fields
		err := rows.Scan(
			// circular master
			&r.TaskId,
			&r.CriteriaType,
			&r.CircularFor,
			&r.RegisterationOpeningDate,
			&r.CancellationDate,
			&r.ProcessID,
			&r.ActivitySeqNo,
			&r.AssignTo,
			&r.AssignedRole,
			&r.IsTaskReturn,
			&r.IsTaskApproved,
			&r.RejectFlag,
			&r.RejectRole,
			&r.Priority,
			&r.Starred,
			&r.Badge,
			&r.EmailFlag,
			&r.TaskStatusID,
			&r.Status,
			&r.InitiatedBy,
			&r.InitiatedOn,
			&r.UpdatedBy,
			&r.UpdatedOn,
			&r.NoOpenDay,
			&r.NoCancellationDate,
			&r.ReferenceNO,

			// ✅ JSON MUST COME HERE
			&cricularJSON,

			// circular template
			&r.OrderNo,
			&r.OrderDate,
			&r.HeaderHTML,
			&r.ToColumn,
			&r.Subject,
			&r.Reference,
			&r.BodyHTML,
			&r.SignatureHTML,
			&r.CCTo,
			&r.FooterHTML,
		)

		if err != nil {
			return nil, err
		}
		// Unmarshal JSON criteria into struct
		if err := json.Unmarshal(cricularJSON, &r.GridDataApproval); err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}
