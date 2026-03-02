// Package routes defines the HTTP route registrations for the Quarters Master Estate module.
// It maps endpoints for managing quarters infrastructure, categories, and building details.
//
// --- Creator's Info ---
// Creator: Ramya M R
// Created On: 12-01-2026
// Last Modified By:
// Last Modified Date:
package routes

import (
	"Hrmodule/auth"
	controllersquartersmasterestate "Hrmodule/controllers/quartersmasterestate"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Quartersmasterestate initializes and registers API routes related to Estate Quarters infrastructure.
//
// This function sets up endpoints for:
//   - Dropdown data retrieval (Categories, Buildings, Quarters Numbers, and Floors).
//   - Fetching comprehensive master details for estate records.
//   - Submission of master data and retrieval of records specifically for approval stages.
//
// All routes registered in this function are wrapped in JwtMiddleware to ensure
// that only authenticated users with valid credentials can access or modify estate master data.
func Quartersmasterestate(router *gin.Engine) {

	router.Any("/Quarter_Category_Dropdown", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquartersmasterestate.QuarterCategoryDropdown))))
	router.Any("/Quarter_building_master", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquartersmasterestate.BuildingMasterDropdown))))
	router.Any("/Estate_quarters_number", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquartersmasterestate.EstateQuartersNumberDropdown))))
	router.Any("/Estate_master_details", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquartersmasterestate.EstateMasterDetails))))
	router.Any("/Quarters_master_submit", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquartersmasterestate.QuartersMasterSubmit))))
	router.Any("/quartersmaster_data_fetch_for_approval", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquartersmasterestate.QuartersMasterDataFetchForApproval))))
	router.Any("/qmesfloordropdown", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquartersmasterestate.QmesFloorDropdown))))
}
