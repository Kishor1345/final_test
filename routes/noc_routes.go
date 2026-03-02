// Package routes defines the HTTP route registrations for the No Objection Certificate (NOC) system.
// It maps various endpoints to their respective controllers, handling applications, approvals, and dynamic content.
//
// Path : /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/routes
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 03-12-2025
// Last Modified By:
// Last Modified Date:
package routes

import (
	"Hrmodule/auth"
	controllersnoc "Hrmodule/controllers/noc"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NOC initializes and registers all NOC-related API endpoints on the provided Gin engine.
//
// This function sets up routes for the complete NOC lifecycle, including:
//   - Application submission and drafting (/nocsubmit).
//   - Role-based approvals and designation management (/noc_approver, /noc_designation).
//   - Template rendering and data fetching (/nocdatatemplate, /noc_template).
//   - Specialized sub-modules like Intimations and Questionnaires.
//
// All routes registered in this function are protected by auth.JwtMiddleware to ensure
// that only authenticated users with valid tokens can access NOC services.
func NOC(router *gin.Engine) {

	router.Any("/nocsubmit", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersnoc.NocSubmit))))
	router.Any("/noc_approver", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersnoc.NocApproverHandler))))
	router.Any("/noc_designation", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersnoc.NocDesignationHandler))))
	router.Any("/noc_reference", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersnoc.NocReferenceHandler))))
	router.Any("/noc_template", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersnoc.NocTemplateHandler))))
	router.Any("/noc_certificate", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersnoc.NocCertificateHandler))))
	router.Any("/noc_intimationdetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersnoc.NocIntimationDetailsHandler))))
	router.Any("/noc_questionnaire", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersnoc.QuestionnaireCertificateHandler))))
}
