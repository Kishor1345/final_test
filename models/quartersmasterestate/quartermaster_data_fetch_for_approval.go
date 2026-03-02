// Package  modelsqms contains structs and queries for quartermaster data fetch for approval.
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:22/01/2026
package quartersmasterestate
import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/google/uuid"
)

// QuarterMasterDataFetchForApproval represents
type QuarterMasterDataFetchForApproval struct{

	//master table
	ReferenceNo      string         `json:"reference_no"`
	TaskID           *uuid.UUID        `json:"task_id"`
	ProcessID        int            `json:"process_id"`
	TaskStatusID     int            `json:"task_status_id"`
	ActivitySeqNo    int            `json:"activity_seq_no"`
	AssignTo         string         `json:"assign_to"`
	AssignedRole     string         `json:"assigned_role"`
	IsTaskReturn     int            `json:"is_task_return"`
	IsTaskApproved   int            `json:"is_task_approved"`
	RejectFlag       int            `json:"reject_flag"`
	RejectRole       *string        `json:"reject_role"`
	Priority         int            `json:"priority"`
	Starred          int            `json:"starred"`
	Badge            int            `json:"badge"`
	EmailFlag        int            `json:"email_flag"`
    TemplateId       *int            `json:"template_id"`

	//grid table
	CategoryId       int            `json:"category_id"`
	BuildingId       int            `json:"building_id"`
	QuartersType     string         `json:"quarters_type"`
	QuartersNumber   string         `json:"quarters_number"`
	QuartersStatus   string         `json:"quarters_status"`
	Street           string         `json:"street"`
	FloorId          int            `json:"floor_id"`
	PlinthArea       int        	`json:"plinth_area"`
    LicenceFee       int        	`json:"licence_fee"`
    ServiceCharge    int       		`json:"service_charge"`
    CautionDeposit   int        	`json:"caution_deposit"`
    SWDCharge        int        	`json:"swd_charge"`
    Status           int     		`json:"status"`

	
    // Building Master
    BuildingID   	int64    		`json:"id"`
    BuildingName	string 			`json:"building_name"`

	// Floor Master
    FloorID   		int64    		`json:"id"`
    FloorName 		string 			`json:"floor_name"`

	  // Quarters Category
    CategoryID   	int    			`json:"id"`
    CategoryName 	string 			`json:"category_name"`

    QuartersMasterID int            `json:"quarters_master_id`
}

// SQL query to fetch Quarters master data
var MyQueryQMSApproval=`

SELECT
--master table

mqm.task_id,mqm.reference_no,mqm.process_id,mqm.task_status_id,mqm.activity_seq_no,mqm.assign_to,
mqm.assigned_role,mqm.is_task_return,mqm.is_task_approved,mqm.reject_flag,mqm.reject_role,
mqm.priority,mqm.starred,mqm.badge,mqm.email_flag,
mqm.template_id,

--grid table

mqg.category_id,mqg.building_name,mqg.quarters_type,mqg.quarters_number,mqg.quarters_status,
mqg.street,mqg.floor_id,mqg.plinth_area,mqg.licence_fee,mqg.service_charge,mqg.caution_deposit,
mqg.swd_charge,mqg.status,

--building master

hbm.id,hbm.building_name,

--floor master

hfm.id,hfm.floor_name,

--quarterscategory

hqc.id,hqc.name,

--quartermaster

hqm.id

--table join

FROM meivan.qmes_m mqm
JOIN meivan.qmes_g mqg ON mqm.task_id = mqg.task_id
JOIN humanresources.quartersmaster hqm ON hqm.quartersnumber = mqg.quarters_number
JOIN humanresources.buildingmaster hbm ON mqg.building_id = hbm.id
JOIN humanresources.floormaster hfm ON mqg.floor_id = hfm.id
JOIN humanresources.quarterscategory hqc ON mqg.category_id = hqc.id
where mqg.status = '1' 
and mqm.task_id=$1;

`



// RetrieveCriteriaMasterDataFetchForApproval maps SQL rows
// into QuarterMasterDataFetchForApproval
func  RetrieveCriteriaMasterDataFetchForApproval(rows *sql.Rows)([]QuarterMasterDataFetchForApproval, error) {

	var results []QuarterMasterDataFetchForApproval

	for rows.Next() {
		var q QuarterMasterDataFetchForApproval
        
        // Scan database row into struct fields
		err := rows.Scan(
    // master
    &q.TaskID,
    &q.ReferenceNo,
    
    &q.ProcessID,
    &q.TaskStatusID,
    &q.ActivitySeqNo,
    &q.AssignTo,
    &q.AssignedRole,
    &q.IsTaskReturn,
    &q.IsTaskApproved,
    &q.RejectFlag,
    &q.RejectRole,
    &q.Priority,
    &q.Starred,
    &q.Badge,
    &q.EmailFlag,
    &q.TemplateId,

    // grid
    &q.CategoryID,
    &q.BuildingID,
    &q.QuartersType,
    &q.QuartersNumber,
    &q.QuartersStatus,
    &q.Street,
    &q.FloorID,
    &q.PlinthArea,
    &q.LicenceFee,
    &q.ServiceCharge,
    &q.CautionDeposit,
    &q.SWDCharge,
    &q.Status,

    // building
    &q.BuildingID,
    &q.BuildingName,

    // floor
    &q.FloorID,
    &q.FloorName,

    // category
    &q.CategoryID,
    &q.CategoryName,

    &q.QuartersMasterID,
)
	if err != nil {
			return nil, err
		}

		results = append(results, q)
	}

	return results, nil

}