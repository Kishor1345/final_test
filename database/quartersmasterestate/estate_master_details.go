// Package databasequartersmasterestate handles DB access for Estate Master Details API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Ramya M R
//
// Created On: 19-01-2026
//
// Last Modified By:
//
// Last Modified Date:
package databasequartersmasterestate

import (
	credentials "Hrmodule/dbconfig"
	modelsquartersmasterestate "Hrmodule/models/quartersmasterestate"
	"fmt"
	"strings"
	"github.com/lib/pq" // Required for handling Go slices in Postgres
)


func parseIDList(input interface{}) []int64 {
	var ids []int64

	if input == nil {
		return nil
	}

	switch val := input.(type) {
	case float64:
		ids = append(ids, int64(val))

	case string:
		strParts := strings.Split(val, ",")
		for _, s := range strParts {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			var id int64
			if _, err := fmt.Sscanf(s, "%d", &id); err == nil {
				ids = append(ids, id)
			}
		}

	case []interface{}:
		for _, item := range val {
			if num, ok := item.(float64); ok {
				ids = append(ids, int64(num))
			} else if s, ok := item.(string); ok {
				var id int64
				if _, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &id); err == nil {
					ids = append(ids, id)
				}
			}
		}
	}

	if len(ids) == 0 {
		return nil
	}
	return ids
}

func GetEstateMasterDetailsFromDB(
    decryptedData map[string]interface{},
) ([]modelsquartersmasterestate.EstateMasterDetailsStruct, int, error) {

    db := credentials.GetDB()

    //  1. Campus ID — MANDATORY
    campusVal, ok := decryptedData["Campus_Id"]
    if !ok || campusVal == nil {
        return nil, 0, fmt.Errorf("campus_id is required")
    }
    var campusID int64
    if cv, ok := campusVal.(float64); ok {
        campusID = int64(cv)
    } else if cv, ok := campusVal.(string); ok {
        if _, err := fmt.Sscanf(cv, "%d", &campusID); err != nil {
            return nil, 0, fmt.Errorf("invalid campus_id")
        }
    } else {
        return nil, 0, fmt.Errorf("invalid campus_id format")
    }

    // 2. Category ID — OPTIONAL (pass nil if not provided)
    var categoryID *int64
    if categoryVal, ok := decryptedData["Category_Id"]; ok && categoryVal != nil {
        var cid int64
        if cv, ok := categoryVal.(float64); ok {
            cid = int64(cv)
            categoryID = &cid
        } else if cv, ok := categoryVal.(string); ok && cv != "" {
            if _, err := fmt.Sscanf(cv, "%d", &cid); err == nil {
                categoryID = &cid
            }
        }
    }

    //  3. Building IDs — OPTIONAL
    buildingIDs := parseIDList(decryptedData["Building_Id"])

    // 4. Quarters IDs — OPTIONAL
    quartersIDs := parseIDList(decryptedData["Quarters_Id"])

    // 5. Prepare array params — pass nil if empty (SQL handles NULL check)
    var bParam, qParam interface{}
    if len(buildingIDs) > 0 {
        bParam = pq.Array(buildingIDs)
    } else {
        bParam = nil  // nil so SQL $3::bigint[] IS NULL check works
    }

    if len(quartersIDs) > 0 {
        qParam = pq.Array(quartersIDs)
    } else {
        qParam = nil  //  nil so SQL $4::bigint[] IS NULL check works
    }

    // 6. Execute query
    rows, err := db.Query(
        modelsquartersmasterestate.MyQueryEstateMasterDetails,
        campusID,   // $1 MANDATORY int64
        categoryID, // $2 OPTIONAL *int64 (nil if not provided)
        bParam,     // $3 OPTIONAL (nil if not provided)
        qParam,     // $4 OPTIONAL (nil if not provided)
    )
    if err != nil {
        return nil, 0, fmt.Errorf("query execution failed: %v", err)
    }
    defer rows.Close()

    // 7. Fetch results
    data, err := modelsquartersmasterestate.RetrieveEstateMasterDetails(rows)
    if err != nil {
        return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
    }

    if len(data) == 0 {
        return []modelsquartersmasterestate.EstateMasterDetailsStruct{}, 0, nil
    }

    return data, len(data), nil
}