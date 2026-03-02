// Package modelsefile contains structs and queries for ALLMenus.
//
// Path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/models/Efile
// --- Creator's Info ---
// Creator: Rovita
// 
// Created On: 24-11-2025
// 
// Last Modified By:  
// 
// Last Modified Date: 
// 
// This api is to fetch the all ALLMenus Master.
package modelsefile

// ALLMenusMasterQuery - Query to fetch ALLMenus values for specific role
const ALLMenusMasterQuery = `SELECT 
    mcv.module_name,
    CASE 
        WHEN mcr.module_id IS NOT NULL THEN 'YES'
        ELSE 'NO'
    END AS module_status, 
    mcv.id as module_id
FROM meivan.category_visibility mcv
LEFT JOIN meivan.category_role_map mcr 
    ON mcv.id = mcr.module_id  
    AND mcr.status = '1'
    AND mcr.role_name = $1 
WHERE mcv.status = '1'
ORDER BY mcv.id`

// ALLMenusMasterGroupedByRoleQuery - Query to fetch all modules grouped by role_name
const ALLMenusMasterGroupedByRoleQuery = `-- Get all active roles
WITH active_roles AS (
    SELECT DISTINCT role_name
    FROM meivan.category_role_map
    WHERE role_name IS NOT NULL and status = '1'
),
active_modules AS (
    SELECT id AS module_id, module_name
    FROM meivan.category_visibility
    WHERE status = '1'
)
SELECT 
    r.role_name,
    m.module_id,
    m.module_name,
    CASE 
        WHEN crm.status = '1' THEN 'YES'
        ELSE 'NO'
    END AS module_status
FROM active_roles r
CROSS JOIN active_modules m
LEFT JOIN meivan.category_role_map crm
    ON crm.role_name = r.role_name
   AND crm.module_id = m.module_id
ORDER BY r.role_name, m.module_id`