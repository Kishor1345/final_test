// Package routes defines the HTTP route registrations for the Office Order and HR modules.
// It orchestrates the mapping between URL endpoints and their respective controllers.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/routes
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 31-10-2025
// Last Modified By: Sridharan
// Last Modified Date: 31-10-2025
package routes

import (
	"Hrmodule/auth"
	controllerscommon "Hrmodule/controllers/common"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Common initializes and registers shared administrative and utility routes on the provided Gin engine.
//
// Most endpoints within this function are protected by JWT authentication middleware
// to ensure that only authorized users can access task inboxes, summaries, and employee details.
//
// Routes registered include:
//   - /Defaultrole: Retrieves the default role for a user.
//   - /TaskInbox & /TaskSummary: Handles task management and overview.
//   - /Statusmaster & /Statusmasternew: Manages status configurations.
//   - /Inboxactivity: Handles updates for task activities.
//   - /download-signature: Public endpoint for signature retrieval (No JWT).
//   - /Employeedetails: Fetches comprehensive employee information.
//   - /ProcessHeader: Retrieves header data for various processes.
func Common(router *gin.Engine) {

	router.Any("/Defaultrole", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.DefaultRoleName))))
	router.Any("/TaskInbox", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.InboxTasksRole))))
	router.Any("/TaskSummary", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.TaskSummary))))
	router.Any("/Statusmaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.StatusMaster))))
	router.Any("/Inboxactivity", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.TaskUpdateHandler))))
	router.Any("/Statusmasternew", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.StatusMasternew))))
	router.Any("/download-signature", gin.WrapH(http.HandlerFunc(controllerscommon.DownloadSignatureHandler))) //without jwt
	router.Any("/Employeedetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.Employeedetails))))
	router.Any("/ProcessHeader", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.ProcessHeader))))
	router.Any("/Employeeinfo", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.EmployeeBasicInfoHandler))))
	router.Any("/ComboValueMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.ComboValueMaster))))
	router.Any("/Religion", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.Religion))))
	router.Any("/BloodGroupMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.BloodGroupMaster))))
	router.Any("/CasteCategory", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.CasteCategory))))
	router.Any("/LanguageMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.LanguageMaster))))
	router.Any("/OfficialLanguage", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.OfficialLanguage))))
	router.Any("/DesignationMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.DesignationMaster))))
	router.Any("/DepartmentMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.DepartmentMaster))))
	router.Any("/Year", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.Year))))
	router.Any("/EmployeePresentScaleMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.EmployeePresentScaleMaster))))
	router.Any("/CountryMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.CountryMaster))))
	router.Any("/StateMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.StateMaster))))
	router.Any("/DistrictMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.DistrictMaster))))
	router.Any("/BankMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.BankMaster))))
	router.Any("/CityMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.CityMaster))))
	router.Any("/Campus_master", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.Campusmaster))))
	router.Any("/nocstatusdelete", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.OfficeOrderstatusdelete))))

}
