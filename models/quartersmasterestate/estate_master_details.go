// Package quartersmasterestate contains structs and queries for Estate Master Details API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Ramya M R
// Created On: 19-01-2026
//
// Last Modified By: Ramya M R
//
// Last Modified Date: 11-02-2026
package quartersmasterestate

import (
	"database/sql"
	"fmt"
)

// SQL Query for Estate Master Details
var MyQueryEstateMasterDetails = `
SELECT 
--quartermaster
qm.id,
qm.quartersnumber,
qm.displayname,
CASE 
  WHEN qm.is_servant_quarters = 1 THEN 'Yes'
  WHEN qm.is_servant_quarters = 0 THEN 'No'
  ELSE NULL
END AS is_servant_quarters,
qm.servant_quartersno,
qm.effectivefrom,
qm.plintharea,
qm.street,
qm.address,
qm.licencefee,
qm.swdcharges,
qm.servicecharges,
qm.cautiondeposit,
qm.garagecharges,
qm.quartersstatus,
--buildingmaster
bm.id,
bm.building_name,

--quarterscategory
qc.id,
qc.name,

-- campus
c.id,
c.campuscode,

--floormaster
fm.id,
fm.floor_name
FROM humanresources.quartersmaster qm 
JOIN humanresources.buildingmaster bm ON qm.building_id = bm.id
JOIN humanresources.quarterscategory qc ON bm.quarters_category = qc.id
JOIN humanresources.campus c ON qc.campus_id = c.id
JOIN humanresources.floormaster fm ON qm.floor_id = fm.id
WHERE c.id = $1
	AND ($2::bigint IS NULL OR qc.id = $2)
	AND ($3::bigint[] IS NULL OR bm.id = ANY($3))
	AND ($4::bigint[] IS NULL OR qm.id = ANY($4))
ORDER BY qm.displayname
`

// Struct for Estate Master Details
type EstateMasterDetailsStruct struct {
	QuartersID        *int     `json:"Quarters_Id"`
	QuartersNumber    *string  `json:"Quarters_Number"`
	FloorID           *int     `json:"Floor_Id"`
	FloorName         *string  `json:"Floor_Name"`
	BuildingID        *int     `json:"Building_Id"`
	BuildingName      *string  `json:"Building_Name"`
	CategoryID        *int     `json:"Catergory_Id"`
	CampusID          *int     `json:"Campus_Id"`
	CampusCode        *string  `json:"Campus_Code"`
	DisplayName       *string  `json:"Display_Name"`
	Street            *string  `json:"Street"`
	PlinthArea        *float64 `json:"Plinth_Area"`
	LicenceFee        *float64 `json:"Licence_Fee"`
	SWDCharges        *float64 `json:"SWD_Charges"`
	ServiceCharges    *float64 `json:"Service_Charges"`
	EBCharges         *float64 `json:"EB_Charges"`
	GarageCharges     *float64 `json:"Garage_Charges"`
	CautionDeposit    *float64 `json:"Caution_Deposit"`
	QuartersStatus    *string  `json:"Quarters_Status"`
	Address           *string  `json:"Address"`
	IsServantQuarters *string  `json:"Is_Servant_Quarters"`
	ServantQuartersNo *string  `json:"Servant_Quarters_No"`
	EffectiveFrom     *string  `json:"Effective_From"`
}

// RetrieveEstateMasterDetails retrieves estate master details data
func RetrieveEstateMasterDetails(rows *sql.Rows) ([]EstateMasterDetailsStruct, error) {
    var list []EstateMasterDetailsStruct

    for rows.Next() {
        var s EstateMasterDetailsStruct
        var categoryName *string // qc.name — not in struct, scan to temp

        err := rows.Scan(
            &s.QuartersID,        // qm.id
            &s.QuartersNumber,    // qm.quartersnumber
            &s.DisplayName,       // qm.displayname
            &s.IsServantQuarters, // is_servant_quarters  
            &s.ServantQuartersNo, // qm.servant_quartersno 
            &s.EffectiveFrom,     // qm.effectivefrom      
            &s.PlinthArea,        // qm.plintharea         
            &s.Street,            // qm.street             
            &s.Address,           // qm.address            
            &s.LicenceFee,        // qm.licencefee         
            &s.SWDCharges,        // qm.swdcharges         
            &s.ServiceCharges,    // qm.servicecharges     
            &s.CautionDeposit,    // qm.cautiondeposit     
            &s.GarageCharges,     // qm.garagecharges   
            &s.QuartersStatus, //dm.quartersstaus    
            &s.BuildingID,        // bm.id                 
            &s.BuildingName,      // bm.building_name      
            &s.CategoryID,        // qc.id                 
            &categoryName,        // qc.name (temp var)    
            &s.CampusID,          // c.id                  
            &s.CampusCode,        // c.campuscode          
            &s.FloorID,           // fm.id                 
            &s.FloorName,         // fm.floor_name     
               
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning estate master details: %v", err)
        }

        list = append(list, s)
    }

    // Check row iteration error
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("row iteration error: %v", err)
    }

    return list, nil
}