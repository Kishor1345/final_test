// Package modelssad contains structs and queries for Modified feilds.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/Staffadditionaldetails
// --- Creator's Info ---
// Creator: Rovita
//
// Created On: 29-01-2026
//
// Last Modified By:
//  
// Last Modified Date:
package modelssad

import (
	"database/sql"
	"encoding/json"
	"fmt"
)


var MyQuerySadPersonalDetails = `

WITH params AS (
    SELECT $1::uuid AS task_id
)


-- =====================================================
-- PERSONAL DETAILS
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT
        SB.task_id,
        SM.module_name,
        SM.modified_field,
        to_jsonb(SB.*) ->> SM.modified_field AS current_value,
        SM.row_id,
        NULL::jsonb AS is_deleted_json
    FROM meivan.sad_basic SB
    JOIN meivan.sad_modified SM
        ON SB.task_id = SM.task_id
    JOIN params p
        ON SB.task_id = p.task_id
) t
WHERE t.current_value IS NOT NULL

UNION ALL

-- =====================================================
-- EDUCATION DETAILS (NORMAL)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT
        SM.task_id,
        SM.module_name,
        SM.modified_field,
        to_jsonb(SE.*) ->> SM.modified_field AS current_value,
        SM.row_id,
        NULL::jsonb AS is_deleted_json
    FROM meivan.sad_modified SM
    JOIN meivan.sad_education SE
        ON SE.sad_id = SM.row_id
    JOIN params p
        ON SM.task_id = p.task_id
    WHERE SM.module_name = 'EDUCATION_DETAILS'
      AND SE.status = 1
      AND SE.is_deleted = 0
) t
WHERE t.current_value IS NOT NULL

UNION ALL

-- =====================================================
-- EDUCATION DETAILS (DELETED – LATEST ONLY)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT ON (se.sad_id)
        se.task_id          AS task_id,
        'EDUCATION_DETAILS' AS module_name,
        NULL                AS modified_field,
        NULL                AS current_value,
        se.sad_id           AS row_id,
        to_jsonb(se)        AS is_deleted_json
    FROM meivan.sad_education se
    JOIN params p
        ON se.task_id = p.task_id
    WHERE se.is_deleted = 1
      AND NOT EXISTS (
          SELECT 1
          FROM meivan.sad_education a
          WHERE a.task_id = se.task_id
            AND a.sad_id  = se.sad_id
            AND a.status = 1
            AND a.is_deleted = 0
      )
    ORDER BY se.sad_id, se.updated_on DESC
) edu_deleted

UNION ALL

-- =====================================================
-- EXPERIENCE DETAILS (NORMAL)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT
        SM.task_id,
        SM.module_name,
        SM.modified_field,
        to_jsonb(EXP.*) ->> SM.modified_field AS current_value,
        SM.row_id,
        NULL::jsonb AS is_deleted_json
    FROM meivan.sad_modified SM
    JOIN meivan.sad_experience EXP
        ON EXP.sad_id = SM.row_id
    JOIN params p
        ON SM.task_id = p.task_id
    WHERE SM.module_name = 'EXPERIENCE_DETAILS'
      AND EXP.status = 1
      AND EXP.is_deleted = 0
) t
WHERE t.current_value IS NOT NULL

UNION ALL

-- =====================================================
-- EXPERIENCE DETAILS (DELETED – LATEST ONLY)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT ON (exp.sad_id)
        exp.task_id           AS task_id,
        'EXPERIENCE_DETAILS'  AS module_name,
        NULL                  AS modified_field,
        NULL                  AS current_value,
        exp.sad_id            AS row_id,
        to_jsonb(exp)         AS is_deleted_json
    FROM meivan.sad_experience exp
    JOIN params p
        ON exp.task_id = p.task_id
    WHERE exp.is_deleted = 1
      AND NOT EXISTS (
          SELECT 1
          FROM meivan.sad_experience a
          WHERE a.task_id = exp.task_id
            AND a.sad_id  = exp.sad_id
            AND a.status = 1
            AND a.is_deleted = 0
      )
    ORDER BY exp.sad_id, exp.updated_on DESC
) exp_deleted

UNION ALL

-- =====================================================
-- NOMINEE DETAILS (NORMAL)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT
        SM.task_id,
        SM.module_name,
        SM.modified_field,
        to_jsonb(NM.*) ->> SM.modified_field AS current_value,
        SM.row_id,
        /* ONLY nomineetype */
        jsonb_build_object(
            'nomineetype', NM.nomineetype
        ) AS is_deleted_json
    FROM meivan.sad_modified SM
    JOIN meivan.sad_nominee NM
        ON NM.sad_id = SM.row_id
    JOIN params p
        ON SM.task_id = p.task_id
    WHERE SM.module_name = 'NOMINEE_DETAILS'
      AND NM.status = 1
      AND NM.is_deleted = 0
) t
WHERE t.current_value IS NOT NULL

UNION ALL

-- =====================================================
-- NOMINEE DETAILS (DELETED – LATEST ONLY)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT ON (nm.sad_id)
        nm.task_id          AS task_id,
        'NOMINEE_DETAILS'   AS module_name,
        NULL                AS modified_field,
        NULL                AS current_value,
        nm.sad_id           AS row_id,
        to_jsonb(nm)        AS is_deleted_json
    FROM meivan.sad_nominee nm
    JOIN params p
        ON nm.task_id = p.task_id
    WHERE nm.is_deleted = 1
      AND NOT EXISTS (
          SELECT 1
          FROM meivan.sad_nominee a
          WHERE a.task_id = nm.task_id
            AND a.sad_id  = nm.sad_id
            AND a.status = 1
            AND a.is_deleted = 0
      )
    ORDER BY nm.sad_id, nm.updated_on DESC
) nom_deleted

UNION ALL

-- =====================================================
-- LANGUAGE DETAILS (NORMAL)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT
        SM.task_id,
        SM.module_name,
        SM.modified_field,
        to_jsonb(LG.*) ->> SM.modified_field AS current_value,
        SM.row_id,
        NULL::jsonb AS is_deleted_json
    FROM meivan.sad_modified SM
    JOIN meivan.sad_language LG
        ON LG.sad_id = SM.row_id
    JOIN params p
        ON SM.task_id = p.task_id
    WHERE SM.module_name = 'LANGUAGE_DETAILS'
      AND LG.status = 1
      AND LG.is_deleted = 0
) t
WHERE t.current_value IS NOT NULL

UNION ALL

-- =====================================================
-- LANGUAGE DETAILS (DELETED – LATEST ONLY)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT ON (lg.sad_id)
        lg.task_id           AS task_id,
        'LANGUAGE_DETAILS'   AS module_name,
        NULL                 AS modified_field,
        NULL                 AS current_value,
        lg.sad_id            AS row_id,
        to_jsonb(lg)         AS is_deleted_json
    FROM meivan.sad_language lg
    JOIN params p
        ON lg.task_id = p.task_id
    WHERE lg.is_deleted = 1
      AND NOT EXISTS (
          SELECT 1
          FROM meivan.sad_language a
          WHERE a.task_id = lg.task_id
            AND a.sad_id  = lg.sad_id
            AND a.status = 1
            AND a.is_deleted = 0
      )
    ORDER BY lg.sad_id, lg.updated_on DESC
) lang_deleted

UNION ALL

-- =====================================================
-- DEPENDENTS DETAILS (NORMAL)
-- =====================================================

SELECT *
FROM (
    SELECT DISTINCT
        SM.task_id,
        SM.module_name,
        SM.modified_field,
        to_jsonb(DP.*) ->> SM.modified_field AS current_value,
        SM.row_id,
        NULL::jsonb AS is_deleted_json
    FROM meivan.sad_modified SM
    JOIN meivan.sad_dependents DP
        ON DP.sad_id = SM.row_id
    JOIN params p
        ON SM.task_id = p.task_id
    WHERE SM.module_name = 'DEPENDENTS_DETAILS'
      AND DP.status = 1 and DP.is_active = 1
	and not exists (select * from meivan.sad_modified where DP.is_active = 0 )
   ) t
WHERE t.current_value IS NOT NULL

UNION ALL

-- -- =====================================================
-- -- DEPENDENTS DETAILS (DELETED – LATEST ONLY)
-- -- =====================================================
SELECT *
FROM (
    SELECT DISTINCT ON (dp.sad_id)
        dp.task_id            AS task_id,
        'DEPENDENTS_DETAILS'  AS module_name,
        NULL                  AS modified_field,
        NULL                  AS current_value,
        dp.sad_id             AS row_id,
        to_jsonb(dp)          AS is_deleted_json
		
    FROM meivan.sad_dependents dp
    JOIN params p
        ON dp.task_id = p.task_id
    WHERE dp.is_active = 0 and dp.status = 1
	
     
    ORDER BY dp.sad_id, dp.updated_on DESC
) dep_deleted

 
UNION ALL

SELECT *
FROM (
    SELECT DISTINCT
        SM.task_id,
        SM.module_name,
        SM.modified_field,
        to_jsonb(C.*) ->> SM.modified_field AS current_value,
        SM.row_id,
        NULL::jsonb AS is_deleted_json
    FROM meivan.sad_modified SM
    JOIN meivan.sad_contact C
        ON C.sad_id = SM.row_id
    JOIN params p
        ON SM.task_id = p.task_id
    WHERE SM.module_name = 'CONTACT_DETAILS'
      AND C.status = 1                 -- ACTIVE CONTACT
) t
WHERE t.current_value IS NOT NULL

UNION ALL

SELECT *
FROM (
    SELECT DISTINCT ON (c.sad_id)
        c.task_id          AS task_id,
        'CONTACT_DETAILS'  AS module_name,
        NULL               AS modified_field,
        NULL               AS current_value,
        c.sad_id           AS row_id,
        to_jsonb(c)        AS is_deleted_json
    FROM meivan.sad_contact c
    JOIN params p
        ON c.task_id = p.task_id
    WHERE c.status = 0                 -- OLD / REPLACED CONTACT
      AND NOT EXISTS (
          SELECT 1
          FROM meivan.sad_contact a
          WHERE a.task_id = c.task_id
            AND a.sad_id  = c.sad_id
            AND a.status  = 1          -- BLOCK IF ACTIVE EXISTS
      )
    ORDER BY c.sad_id, c.updated_on DESC
) contact_deleted
 
UNION ALL

-- =====================================================
-- HINDI PROFICIENCY (NORMAL)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT
        SM.task_id,
        SM.module_name,
        SM.modified_field,
        to_jsonb(HP.*) ->> SM.modified_field AS current_value,
        SM.row_id,
        NULL::jsonb AS is_deleted_json
    FROM meivan.sad_modified SM
    JOIN meivan.sad_hindiproficiency HP
        ON HP.sad_id = SM.row_id
    JOIN params p
        ON SM.task_id = p.task_id
    WHERE SM.module_name = 'HINDI_PROFICIENCY'
      AND HP.status = 1
      AND HP.is_deleted = 0
) t
WHERE t.current_value IS NOT NULL

UNION ALL

-- =====================================================
-- HINDI PROFICIENCY (DELETED – LATEST ONLY)
-- =====================================================
SELECT *
FROM (
    SELECT DISTINCT ON (hp.sad_id)
        hp.task_id            AS task_id,
        'HINDI_PROFICIENCY'   AS module_name,
        NULL                  AS modified_field,
        NULL                  AS current_value,
        hp.sad_id             AS row_id,
        to_jsonb(hp)          AS is_deleted_json
    FROM meivan.sad_hindiproficiency hp
    JOIN params p
        ON hp.task_id = p.task_id
    WHERE hp.is_deleted = 1
      AND NOT EXISTS (
          SELECT 1
          FROM meivan.sad_hindiproficiency a
          WHERE a.task_id = hp.task_id
            AND a.sad_id  = hp.sad_id
            AND a.status  = 1
            AND a.is_deleted = 0
      )
    ORDER BY hp.sad_id, hp.updated_on DESC
) hindi_deleted

`

// =====================
// RESPONSE STRUCT
// =====================
type SadPersonalDetails struct {
	TaskID        *string `json:"task_id"`
	ModuleName    *string `json:"module_name"`
	ModifiedField *string `json:"modified_field"`
	CurrentValue  *string `json:"current_value"`
	Rowid         *int    `json:"row_id"`

	// IMPORTANT: scan-safe holder for NULL jsonb
	IsDeletedJSON sql.NullString `json:"-"`
}

// =====================
// ROW SCANNER
// =====================
func RetrieveSadPersonalDetails(rows *sql.Rows) ([]SadPersonalDetails, error) {

	var list []SadPersonalDetails

	for rows.Next() {
		var s SadPersonalDetails

		err := rows.Scan(
			&s.TaskID,
			&s.ModuleName,
			&s.ModifiedField,
			&s.CurrentValue,
			&s.Rowid,
			&s.IsDeletedJSON, // SAFE FOR NULL
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		list = append(list, s)
	}

	return list, nil
}

// =====================
// JSON OUTPUT HANDLER
// =====================
func (s SadPersonalDetails) MarshalJSON() ([]byte, error) {

	type Alias SadPersonalDetails

	var deleted json.RawMessage
	if s.IsDeletedJSON.Valid {
		deleted = json.RawMessage(s.IsDeletedJSON.String)
	}

	return json.Marshal(&struct {
		Alias
		IsDeletedJSON json.RawMessage `json:"is_deleted_json"`
	}{
		Alias:         (Alias)(s),
		IsDeletedJSON: deleted,
	})
}
