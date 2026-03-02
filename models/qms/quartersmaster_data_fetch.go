// Package modelsquarters contains structs and queries for Quarters details API.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/qms
// --- Creator's Info ---
// Creator: Elakiya
//
// Created On: 
//
// Last Modified By:
//
// Last Modified Date:
package modelsqms

import (
	"database/sql"
	"fmt"
)

// =====================
// SQL QUERY
// =====================
var MyQueryQMSEUDetails = `
SELECT
	QM.reference_no,
	QM.task_id,
	QM.process_id,
	QM.quartersnumber,
	QM.address,
	QM.floor,
	FM.floor_name,
	QM.street,
	QM.quarterscategoryid,
	QC.name AS quarters_category_name,
	QM.buildingtypeid,
	BM.building_name,
	QM.plintharea,
	QM.quartersstatus,
	QM.employee_id,
	QM.employee_name,
	QM.department,
	QM.designation,
	TO_CHAR(QM.ocupieddate,   'DD-MM-YYYY') AS ocupieddate,
	QM.contactno,
	QM.licencefee,
	QM.swdcharges,
	QM.servicecharges,
	QM.garagecharges,
	QM.cautiondeposit,
	TO_CHAR(QM.effectivefrom, 'DD-MM-YYYY') AS effectivefrom,
	TO_CHAR(QM.validtill,     'DD-MM-YYYY') AS validtill,
	QM.assign_to,
	QM.assigned_role,
	QM.task_status_id,
	QM.activity_seq_no,
	QM.is_task_return,
	QM.is_task_approved,
	QM.email_flag,
	QM.template_id,
	QM.reject_flag,
	QM.reject_role,
	QM.initiated_by,
	TO_CHAR(QM.initiated_on, 'DD-MM-YYYY') AS initiated_on,
	QM.updated_by,
	TO_CHAR(QM.updated_on,  'DD-MM-YYYY') AS updated_on,
	QM.status,
	QM.badge,
	QM.priority,
	QM.starred, 
	DQM.quartersstatus AS domain_quarters_status
FROM meivan.qmseu_m QM
JOIN humanresources.buildingmaster BM
	ON QM.buildingtypeid = BM.id
JOIN humanresources.quarterscategory QC
	ON BM.quarters_category = QC.id
JOIN humanresources.floormaster FM
	ON QM.floor = FM.id
LEFT JOIN humanresources.quartersmaster DQM
	ON DQM.displayname = QM.quartersnumber
WHERE QM.task_id = $1;
`

// NOTE: PropertyInfo, CurrentResidingEmployee, and ChargesAndFees structs
// are already defined in quartersmaster.go
// You need to add DomainQuartersStatus field to the existing PropertyInfo struct there

// =====================
// TaskInfo Struct
// =====================
type TaskInfo struct {
	ReferenceNo    *string `json:"reference_no"`
	TaskID         *string `json:"task_id"`
	ProcessID      *int    `json:"process_id"`
	AssignTo       *string `json:"assign_to"`
	AssignedRole   *string `json:"assigned_role"`
	TaskStatusID   *int    `json:"task_status_id"`
	ActivitySeqNo  *int    `json:"activity_seq_no"`
	IsTaskReturn   *bool   `json:"is_task_return"`
	IsTaskApproved *bool   `json:"is_task_approved"`
	EmailFlag      *bool   `json:"email_flag"`
	TemplateID     *int    `json:"template_id"`
	RejectFlag     *bool   `json:"reject_flag"`
	RejectRole     *string `json:"reject_role"`
	InitiatedBy    *string `json:"initiatedby"`
	InitiatedOn    *string `json:"initiatedon"`
	UpdatedBy      *string `json:"updatedby"`
	UpdatedOn      *string `json:"updatedon"`
	Status         *int    `json:"status"`
	Badge          *string `json:"badge"`
	Priority       *int    `json:"priority"`
	Starred        *bool   `json:"starred"`
}

// =====================
// Status Struct
// =====================
type Status struct {
	QuartersStatus *string `json:"Client_quarters_status"`
	EffectiveFrom  *string `json:"effectivefrom"`
}

// =====================
// MAIN RESPONSE STRUCT
// =====================
type QMSEUDetailsStruct struct {
	PropertyInfo            PropertyInfo            `json:"property_info"`
	CurrentResidingEmployee CurrentResidingEmployee `json:"current_residing_employee"`
	ChargesAndFees          ChargesAndFees          `json:"charges_and_fees"`
	TaskInfo                TaskInfo                `json:"task_info"`
	Status                  Status                  `json:"status"`
}

// =====================
// ROW SCANNER
// =====================
func RetrieveQMSEUDetails(rows *sql.Rows) ([]QMSEUDetailsStruct, error) {

	var list []QMSEUDetailsStruct

	for rows.Next() {

		var s QMSEUDetailsStruct

		// Scan variables
		var (
			referenceNo, taskID, quartersNumber, address, floorName, street,
			quartersCategoryName, buildingName, plinthArea, quartersStatus,
			residentID, residentName, department, designation, ocupiedDate,
			contactNo, licenceFee, swdCharges, serviceCharges, garageCharges,
			cautionDeposit, effectiveFrom, validTill, assignTo, assignedRole,
			rejectRole, initiatedBy, initiatedOn, updatedBy, updatedOn,
			badge, domainQuartersStatus *string

			processID, floor, quartersCategoryID, buildingTypeID,
			taskStatusID, activitySeqNo, templateID,
			status, priority *int

			isTaskReturn, isTaskApproved, emailFlag,
			rejectFlag, starred *bool
		)

		err := rows.Scan(
			&referenceNo,
			&taskID,
			&processID,
			&quartersNumber,
			&address,
			&floor,
			&floorName,
			&street,
			&quartersCategoryID,
			&quartersCategoryName,
			&buildingTypeID,
			&buildingName,
			&plinthArea,
			&quartersStatus,
			&residentID,
			&residentName,
			&department,
			&designation,
			&ocupiedDate,
			&contactNo,
			&licenceFee,
			&swdCharges,
			&serviceCharges,
			&garageCharges,
			&cautionDeposit,
			&effectiveFrom,
			&validTill,
			&assignTo,
			&assignedRole,
			&taskStatusID,
			&activitySeqNo,
			&isTaskReturn,
			&isTaskApproved,
			&emailFlag,
			&templateID,
			&rejectFlag,
			&rejectRole,
			&initiatedBy,
			&initiatedOn,
			&updatedBy,
			&updatedOn,
			&status,
			&badge,
			&priority,
			&starred,
			&domainQuartersStatus,
		)

		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		// Property Info
		s.PropertyInfo.QuartersNumber = quartersNumber
		s.PropertyInfo.Address = address
		s.PropertyInfo.FloorID = floor
		s.PropertyInfo.FloorName = floorName
		s.PropertyInfo.Street = street
		s.PropertyInfo.QuartersCategory = quartersCategoryID
		s.PropertyInfo.CategoryName = quartersCategoryName
		s.PropertyInfo.BuildingID = buildingTypeID
		s.PropertyInfo.BuildingName = buildingName
		s.PropertyInfo.PlinthArea = plinthArea
		s.PropertyInfo.QuartersStatus = domainQuartersStatus // Domain quarters status from DQM table

		// Current Resident
		s.CurrentResidingEmployee.ResidentID = residentID
		s.CurrentResidingEmployee.ResidentName = residentName
		s.CurrentResidingEmployee.DepartmentName = department
		s.CurrentResidingEmployee.DesignationName = designation
		s.CurrentResidingEmployee.DateOfOccupied = ocupiedDate
		s.CurrentResidingEmployee.MobileNumber = contactNo
		s.CurrentResidingEmployee.ValidTill = validTill

		// Charges
		s.ChargesAndFees.LicenceFee = licenceFee
		s.ChargesAndFees.SwdCharges = swdCharges
		s.ChargesAndFees.ServiceCharges = serviceCharges
		s.ChargesAndFees.GarageCharges = garageCharges
		s.ChargesAndFees.CautionDeposit = cautionDeposit

		// Task Info
		s.TaskInfo.ReferenceNo = referenceNo
		s.TaskInfo.TaskID = taskID
		s.TaskInfo.ProcessID = processID
		s.TaskInfo.AssignTo = assignTo
		s.TaskInfo.AssignedRole = assignedRole
		s.TaskInfo.TaskStatusID = taskStatusID
		s.TaskInfo.ActivitySeqNo = activitySeqNo
		s.TaskInfo.IsTaskReturn = isTaskReturn
		s.TaskInfo.IsTaskApproved = isTaskApproved
		s.TaskInfo.EmailFlag = emailFlag
		s.TaskInfo.TemplateID = templateID
		s.TaskInfo.RejectFlag = rejectFlag
		s.TaskInfo.RejectRole = rejectRole
		s.TaskInfo.InitiatedBy = initiatedBy
		s.TaskInfo.InitiatedOn = initiatedOn
		s.TaskInfo.UpdatedBy = updatedBy
		s.TaskInfo.UpdatedOn = updatedOn
		s.TaskInfo.Status = status
		s.TaskInfo.Badge = badge
		s.TaskInfo.Priority = priority
		s.TaskInfo.Starred = starred

		// Status
		s.Status.QuartersStatus = quartersStatus // Changed: from QM table
		s.Status.EffectiveFrom = effectiveFrom

		list = append(list, s)
	}

	return list, nil
}
