// Package modelsefile contains structs and queries for InsertUpdateModules.
//
// Path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/models/Efile
// --- Creator's Info ---
// Creator: Rovita
//
//  Created On: 24-11-2025
// 
// Last Modified By: AI Assistant
// 
// Last Modified Date: 27-11-2025
// 
// This file contains queries for category_role_map operations with replace functionality.
package modelsefile

// Legacy single insert query (kept for backward compatibility)
const InsertupdateCategoryRoleMapQuery = `
INSERT INTO meivan.category_role_map (module_id, role_name, created_at, updated_at, status, created_by)
VALUES ($1, $2, NOW(), NOW(), $3, $4)`

const InsertUpdateCategoryRoleMapQuery = InsertupdateCategoryRoleMapQuery // Alias for consistency

// 1️⃣ Insert NEW modules (ignore if already exists)
const BulkInsertCategoryRoleMapQuery = `
INSERT INTO meivan.category_role_map (module_id, role_name, created_at, updated_at, status , created_by)
SELECT unnest($2::int[]), $1, NOW(), NOW(), 1 , $3
ON CONFLICT (module_id, role_name) DO NOTHING;
`

// 2️⃣ Activate selected modules (set status=1 for modules in the list)
const BulkActivateCategoryRoleMapQuery = `
UPDATE meivan.category_role_map
SET status = 1, updated_at = NOW(),
updated_by = $3
WHERE role_name = $1
  AND module_id = ANY($2::int[]);
`

// 3️⃣ Deactivate unselected modules (set status=0 for modules NOT in the list)
const BulkDeactivateCategoryRoleMapQuery = `
UPDATE meivan.category_role_map
SET status = 0, updated_at = NOW(), updated_by = $3
WHERE role_name = $1
  AND module_id != ALL($2::int[]);
`

// 4️⃣ Deactivate ALL modules for a role (used when module_id is empty string)
const DeactivateAllModulesForRoleQuery = `
UPDATE meivan.category_role_map
SET status = 0, updated_at = NOW(), updated_by = $2
WHERE role_name = $1;
`