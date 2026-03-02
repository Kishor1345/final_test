// // Package databaseefile contains structs and queries for InsertUpdateModules.
// //path : /var/www/html/go_projects/HRMODULE/Rovita/HR_test/database/Efile
// // --- Creator's Info ---
// // Creator: Rovita
// // Created On: 24-11-2025
// // Last Modified By: AI Assistant
// // Last Modified Date: 27-11-2025
// // This api is to InsertUpdate data into category_role_map with intelligent upsert logic.
package databaseefile

import (
	credentials "Hrmodule/dbconfig"
	modelsefile "Hrmodule/models/Efile"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// ReplaceModulesForRole performs a complete replacement:
// 1. Insert NEW modules (status=1, created_by)
// 2. Activate selected modules (status=1, updated_by)
// 3. Deactivate unselected modules (status=0, updated_by)
//
// Returns: insertedCount, activatedCount, deactivatedCount, error
func ReplaceModulesForRole(moduleIDs []string, roleName, userID string) (int64, int64, int64, error) {
	var insertedCount, activatedCount, deactivatedCount int64

	// Database connection
	db := credentials.GetDB()

	if err := db.Ping(); err != nil {
		return 0, 0, 0, fmt.Errorf("DB connection failed: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Convert moduleIDs from string → int64
	moduleIDInts := make([]int64, len(moduleIDs))
	for i, moduleID := range moduleIDs {
		var id int64
		if _, err := fmt.Sscanf(moduleID, "%d", &id); err != nil {
			return 0, 0, 0, fmt.Errorf("invalid module_id format: %s", moduleID)
		}
		moduleIDInts[i] = id
	}

	// -----------------------------
	// 1. Insert NEW modules (created_by)
	// -----------------------------
	result1, err := tx.Exec(
		modelsefile.BulkInsertCategoryRoleMapQuery,
		roleName,
		pq.Array(moduleIDInts),
		userID, // created_by
	)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error inserting new modules: %v", err)
	}
	insertedCount, _ = result1.RowsAffected()

	// -----------------------------
	// 2. Activate selected modules (updated_by)
	// -----------------------------
	result2, err := tx.Exec(
		modelsefile.BulkActivateCategoryRoleMapQuery,
		roleName,
		pq.Array(moduleIDInts),
		userID, // updated_by
	)
	if err != nil {
		return insertedCount, 0, 0, fmt.Errorf("error activating modules: %v", err)
	}
	activatedCount, _ = result2.RowsAffected()

	// -----------------------------
	// 3. Deactivate unselected modules (updated_by)
	// -----------------------------
	result3, err := tx.Exec(
		modelsefile.BulkDeactivateCategoryRoleMapQuery,
		roleName,
		pq.Array(moduleIDInts),
		userID, // updated_by
	)
	if err != nil {
		return insertedCount, activatedCount, 0, fmt.Errorf("error deactivating modules: %v", err)
	}
	deactivatedCount, _ = result3.RowsAffected()

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return insertedCount, activatedCount, deactivatedCount, fmt.Errorf("error committing transaction: %v", err)
	}

	return insertedCount, activatedCount, deactivatedCount, nil
}

// Deactivate ALL modules for a given role
func DeactivateAllModulesForRole(roleName, userID string) (int64, error) {
	connectionString := credentials.Getdatabasehr()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return 0, fmt.Errorf("DB connection failed: %v", err)
	}

	result, err := db.Exec(
		modelsefile.DeactivateAllModulesForRoleQuery,
		roleName,
		userID, // updated_by
	)
	if err != nil {
		return 0, fmt.Errorf("error deactivating all modules: %v", err)
	}

	deactivatedCount, _ := result.RowsAffected()
	return deactivatedCount, nil
}

// Backward compatibility
func BulkUpsertCategoryRoleMap(moduleIDs []string, roleName, userID string) (int64, int64, int64, error) {
	return ReplaceModulesForRole(moduleIDs, roleName, userID)
}

// Legacy function (backward compatible)
func InsertupdateMultipleCategoryRoleMap(moduleIDs []string, roleName, status, userID string) (int64, error) {
	inserted, activated, _, err := ReplaceModulesForRole(moduleIDs, roleName, userID)
	return inserted + activated, err
}

// Single insert legacy
func InsertupdateCategoryRoleMap(moduleID, roleName, status, userID string) (int64, error) {
	return InsertupdateMultipleCategoryRoleMap([]string{moduleID}, roleName, status, userID)
}
