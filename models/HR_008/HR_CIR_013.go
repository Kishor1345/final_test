// Package modelscircular contains structs and queries for  circular details.
//
//Path : /var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/models/hr_008
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
	_ "github.com/lib/pq"
	"github.com/google/uuid"
)


type GridDataForCircularDetails struct{
	QuartersCategory string `json:"quarters_category"`
	QuartersNo       string `json:"quarters_no"`
	Floor            string `json:"floor"`
	Location         string `json:"location"`
	Password         string `json:"password"`
	KeyBox           string `json:"key_box"`
	FirstChoice      string `json:"first_choice"`
	SecondChoice     string `json:"second_choice"`
	ThirdChoice      string `json:"third_choice"`
}


type CircularDetailStructure struct{
	//Circular Master
	TaskId  	             uuid.UUID       `json:"task_id"`
	CriteriaType             string     `json:"criteria_type"`
	CircularFor   			 string     `json:"circular_for"`
	RegisterationOpeningDate string     `json:"open_date_registration"`
	RegisterationLastDate 	 string     `json:"last_date_registration"`
	CancellationDate 		 string     `json:"last_date_cancellation"`
	ProcessID                int        `json:"process_id"`
	HeaderHTML				 string     `json:"header_html"`
	OrderNo                  string      `json:"order_no"`
	ToColumn                 string     `json:"to_column"`
	Subject                  string     `json:"subject"`
	Reference              string       `json:"reference"`
	BodyHTML                 string     `json:"body_html"`
	CCTo					 string     `json:"cc_to"`
	FooterHTML               string     `json:"footer_html"`
	Status                   int        `json:"status"`
	NoOfRegistrationOpenDays   string     `json:"no_of_registration_open_days"`
	NoOfCancellationDate         string     `json:"no_of_cancellation_days"`
	CirculaDate                    string `json:"circular_date"`
	OrderDate                    string `json:"order_date"`
	TemplateId                    string `json:"template_id"`
	Campus                    string `json:"campus"`
	GridDataCircularDetailsExisting  []GridDataForCircularDetails  `json:"circulargriddatadetails`

	
	
}


var MyQueryForCircularDataFetch = 
`
select 
task_id,
criteria_type,
circular_for,
meivan.globaldate_format (open_date_registration::text) as open_date_registration,
meivan.globaldate_format (last_date_registration::text) as last_date_registration,
meivan.globaldate_format (last_date_cancellation::text) as last_date_cancellation,
process_id,
header_html,
order_no,
to_column,
subject,
reference,
body_html,
cc_to,
footer_html,
status,
no_of_registration_open_days,
no_of_cancellation_days,
meivan.globaldate_format(circular_date::text)as circular_date,
meivan.globaldate_format(order_date::text)as order_date,
template_id,
campus,
json_agg(
    json_build_object(
        'quarters_category',   quarters_category,
        'quarters_no',   quarters_no,
		'floor', floor,
        'location',  location,
        'password',password,
        'key_box',    key_box,
        'first_choice',  first_choice,
		'second_choice', second_choice,
		'third_choice',third_choice
    )
) as circulargriddata
from humanresources.circular_master 
where order_no = $1
group by 
task_id,
criteria_type,
circular_for,
open_date_registration,
last_date_registration,
last_date_cancellation,
process_id,
header_html,
order_no,
to_column,
subject,
reference,
body_html,
cc_to,
footer_html,
status,
no_of_registration_open_days,
no_of_cancellation_days,
circular_date,
order_date,
template_id,
campus

`

func RetrieveCircularDetailFetch(rows *sql.Rows) ([]CircularDetailStructure, error) {

	var results []CircularDetailStructure

	for rows.Next() {
		var r CircularDetailStructure
		var cricularJSON []byte
        // Scan database row into struct fields
		err := rows.Scan(
			//circular master 
			&r.TaskId,
	        &r.CriteriaType,
	        &r.CircularFor,
	        &r.RegisterationOpeningDate,
	        &r.RegisterationLastDate,
	        &r.CancellationDate,
	        &r.ProcessID,
	        &r.HeaderHTML,
	        &r.OrderNo,
	        &r.ToColumn,
	        &r.Subject,
	        &r.Reference,
	        &r.BodyHTML,
	        &r.CCTo,
	        &r.FooterHTML,
	        &r.Status,
	        &r.NoOfRegistrationOpenDays,
	        &r.NoOfCancellationDate,
	        &r.CirculaDate,
	        &r.OrderDate,
	        &r.TemplateId,
	        &r.Campus,
			&cricularJSON,
		)
		if err != nil {
			return nil, err
		}
		// Unmarshal JSON criteria into struct
		if err := json.Unmarshal(cricularJSON, &r.GridDataCircularDetailsExisting); err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}
