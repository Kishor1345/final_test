// Package routes defines the HTTP route registrations for the Quarters Management System (QMS).
// It maps internal QMS endpoints to their respective controllers for managing housing and quarters data.
//
// --- Creator's Info ---
// Creator: Ramya M R
// Created On: 06-12-2025
//
// Last Modified By:
// Last Modified Date:
package routes

import (
	"Hrmodule/auth"
	controllersquarters "Hrmodule/controllers/qms"
	"net/http"

	"github.com/gin-gonic/gin"
)

// QMS initializes and registers all Quarters Management System-related API endpoints on the provided Gin engine.
//
// This function sets up routes for the following operations:
//   - /quartersdropdown: Retrieves dropdown options for quarters categorization.
//   - /quartersmaster: Fetches general quarters master details.
//   - /quartersmastereu_submit: Processes submissions and drafts for Estate Unit (EU) quarters.
//   - /quartersmastereu_data_fetch: Retrieves specific Estate Unit data for existing tasks.
//   - /quartersmaster_status_dropdown: Fetches status-specific dropdown values for quarters.
//
// All routes registered in this function are protected by auth.JwtMiddleware to ensure
// that only authenticated users with valid tokens can access or modify QMS data.
func QMS(router *gin.Engine) {

	router.Any("/quartersdropdown", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquarters.QuartersDropdown))))
	router.Any("/quartersmaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquarters.QuartersDetails))))
	router.Any("/quartersmastereu_submit", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquarters.QmseuSubmit))))
	router.Any("/quartersmastereu_data_fetch", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquarters.QMSEUDetails))))
	router.Any("/quartersmaster_status_dropdown", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersquarters.QuartersStatus))))
	
}
