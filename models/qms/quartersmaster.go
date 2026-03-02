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
var MyQueryQuartersDetails = ` SELECT 
	QM.quartersnumber,
	QM.displayname,
	QM.address,
	QM.street,

	QM.floor_id,
	FM.floor_name,

	BM.quarters_category,
	QC.name AS category_name,

	QM.quartersstatus,
	QM.plintharea,

	QM.building_id,
	BM.building_name,

	QM.licencefee,
	QM.watercharges,
	QM.swdcharges,
	QM.servicecharges,
	QM.cautiondeposit,
	QM.ebcharges,
	QM.garagecharges,

	QA.resident_id,
	QA.resident_name,
	dm.departmentname,
	dg.designationname,
	a.mobilenumber,

	TO_CHAR(QA.dateofoccupied, 'DD-MM-YYYY') AS dateofoccupied,
	TO_CHAR(QA.dateofalloted,  'DD-MM-YYYY') AS dateofalloted,
	TO_CHAR(QA.validtill,      'DD-MM-YYYY') AS validtill
FROM humanresources.quartersmaster QM
RIGHT JOIN humanresources.floormaster FM ON QM.floor_id = FM.id
RIGHT JOIN humanresources.buildingmaster BM ON QM.building_id = BM.id
 JOIN humanresources.quarterscategory QC ON BM.quarters_category = QC.id
RIGHT JOIN humanresources.quartersallotedmaster QA ON QM.id = QA.quartersid
LEFT JOIN humanresources.employeebasicinfo a ON QA.resident_id = a.employeeid
LEFT JOIN humanresources.departmentdesignationmapping m ON a.employeeid = m.employeeid
LEFT JOIN humanresources.departmentmaster dm ON m.departmentcode = dm.departmentcode
LEFT JOIN humanresources.designationmaster dg ON m.designationid = dg.designationid
WHERE QM.displayname =$1  
AND NOT EXISTS (
    SELECT 1
    FROM meivan.qmseu_m q
    WHERE q.task_status_id IN (4,22) AND q.quartersnumber=QM.displayname 
)
ORDER BY QM.quartersnumber;
`

// =====================
// PROPERTY INFO
// =====================
type PropertyInfo struct {
	QuartersNumber   *string `json:"quartersnumber"`
	DisplayName      *string `json:"displayname"`
	Address          *string `json:"address"`
	Street           *string `json:"street"`
	FloorID          *int    `json:"floor_id"`
	FloorName        *string `json:"floor_name"`
	QuartersCategory *int    `json:"quarters_category"`
	CategoryName     *string `json:"quarters_category_name"`
	QuartersStatus   *string `json:"quartersstatus"`
	PlinthArea       *string `json:"plintharea"`
	BuildingID       *int    `json:"building_id"`
	BuildingName     *string `json:"building_name"`
}

// =====================
// CHARGES & FEES
// =====================
type ChargesAndFees struct {
	LicenceFee     *string `json:"licencefee"`
	WaterCharges   *string `json:"watercharges"`
	SwdCharges     *string `json:"swdcharges"`
	ServiceCharges *string `json:"servicecharges"`
	CautionDeposit *string `json:"cautiondeposit"`
	EbCharges      *string `json:"ebcharges"`
	GarageCharges  *string `json:"garagecharges"`
}

// =====================
// CURRENT RESIDING EMPLOYEE
// =====================
type CurrentResidingEmployee struct {
	ResidentID      *string `json:"resident_id"`
	ResidentName    *string `json:"resident_name"`
	DepartmentName  *string `json:"departmentname"`
	DesignationName *string `json:"designationname"`
	MobileNumber    *string `json:"mobilenumber"`
	DateOfOccupied  *string `json:"dateofoccupied"`
	DateOfAlloted   *string `json:"dateofalloted"`
	ValidTill       *string `json:"validtill"`
}

// =====================
// FINAL RESPONSE STRUCT
// =====================
type QuartersDetailsStruct struct {
	PropertyInfo            PropertyInfo            `json:"property_info"`
	ChargesAndFees          ChargesAndFees          `json:"charges_and_fees"`
	CurrentResidingEmployee CurrentResidingEmployee `json:"current_residing_employee"`
}

// =====================
// ROW SCANNER
// =====================
func RetrieveQuartersDetails(rows *sql.Rows) ([]QuartersDetailsStruct, error) {

	var list []QuartersDetailsStruct

	for rows.Next() {

		var s QuartersDetailsStruct

		err := rows.Scan(
			// PROPERTY INFO
			&s.PropertyInfo.QuartersNumber,
			&s.PropertyInfo.DisplayName,
			&s.PropertyInfo.Address,
			&s.PropertyInfo.Street,
			&s.PropertyInfo.FloorID,
			&s.PropertyInfo.FloorName,
			&s.PropertyInfo.QuartersCategory,
			&s.PropertyInfo.CategoryName,
			&s.PropertyInfo.QuartersStatus,
			&s.PropertyInfo.PlinthArea,
			&s.PropertyInfo.BuildingID,
			&s.PropertyInfo.BuildingName,

			// CHARGES & FEES
			&s.ChargesAndFees.LicenceFee,
			&s.ChargesAndFees.WaterCharges,
			&s.ChargesAndFees.SwdCharges,
			&s.ChargesAndFees.ServiceCharges,
			&s.ChargesAndFees.CautionDeposit,
			&s.ChargesAndFees.EbCharges,
			&s.ChargesAndFees.GarageCharges,

			// CURRENT RESIDING EMPLOYEE
			&s.CurrentResidingEmployee.ResidentID,
			&s.CurrentResidingEmployee.ResidentName,
			&s.CurrentResidingEmployee.DepartmentName,
			&s.CurrentResidingEmployee.DesignationName,
			&s.CurrentResidingEmployee.MobileNumber,
			&s.CurrentResidingEmployee.DateOfOccupied,
			&s.CurrentResidingEmployee.DateOfAlloted,
			&s.CurrentResidingEmployee.ValidTill,
		)

		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		list = append(list, s)
	}

	return list, nil
}
