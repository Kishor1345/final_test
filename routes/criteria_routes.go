// Package routes defines the HTTP route registrations for the Criteria Master module.
// It maps endpoints to controllers responsible for managing criteria definitions,
// data fetching, and approval workflows.
//
// --- Creator's Info ---
// Creator: Kishorekumar
// Created On: 09-12-2025
package routes

import (
	"Hrmodule/auth"
	controllerscriteria "Hrmodule/controllers/criteria"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Criteria initializes and registers all Criteria Master-related API endpoints on the provided Gin engine.
//
// This function sets up routes for:
//   - Dropdown data retrieval (Great Pay).
//   - Insertion of new criteria definitions.
//   - General and specific data fetching for existing records.
//   - Approval-specific data retrieval for the workflow process.
//
// All routes registered in this function are protected by auth.JwtMiddleware to ensure
// that only authenticated users with valid tokens can modify or view criteria configurations.
func Criteria(router *gin.Engine) {

	router.Any("/greatpaydropdown", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscriteria.GreatPayDropdown))))
	router.Any("/insertcriteria", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscriteria.CriteriaInsert))))
	router.Any("/criteriamasterdatafetch", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscriteria.CriteriaMasterDataFetch))))
	router.Any("/criteriamasterexistingdatafetch", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscriteria.CriteriaMasterExistingDataFetch))))
	router.Any("/criteriamasterdatafetchforapproval", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscriteria.CriteriaMasterDataFetchForApproval))))
}
