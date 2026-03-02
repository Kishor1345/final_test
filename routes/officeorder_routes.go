// Package routes defines the HTTP route registrations for the Office Order management system.
// It maps endpoints to controllers responsible for order lifecycles, visit details, and approval workflows.
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
	controllersofficeorder "Hrmodule/controllers/officeorder"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Officeorder registers all endpoints related to the processing and management of office orders.
//
// This function initializes routes for:
//   - Module and sub-module configurations.
//   - Visit details and PCR (Post-Commitment Report) data entry.
//   - Task status updates, approval remarks, and return dropdowns.
//   - Document generation (templates and history PDFs).
//   - Audit trails and task summaries.
//
// All routes registered in this function are protected by auth.JwtMiddleware to ensure
// that only authenticated users can access or modify office order data.
func Officeorder(router *gin.Engine) {
	router.Any("/OfficeOrder_module", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OrderSubModule))))
	router.Any("/OfficeOrder_visitdetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.Ordervisitdetails))))
	router.Any("/OfficeOrder_InsertOfficedetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.PCRInsert))))
	router.Any("/OfficeOrder_Count", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.NeedGenerateHandler))))
	router.Any("/OfficeOrder_DropdownValuesHandler", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.DropdownValuesHandler))))
	router.Any("/OfficeOrder_statusupdate", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OfficeOrderUpdateTaskStatus))))
	router.Any("/OfficeOrder_approval_remarks", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OfficeComments))))
	router.Any("/OfficeOrder_datatemplate", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.GetOfficeOrderDetailsfortemplate))))
	router.Any("/OfficeOrder_ReturnDropdown", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.ReturnDropdown))))
	router.Any("/OfficeOrder_taskvisitdetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OrderTaskVisitDetails))))
	router.Any("/OfficeOrder_History", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OfficeOrderHistory))))
	router.Any("/OfficeOrder_Historypdf", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.FetchOrderHistoryPDF))))
	router.Any("/OfficeOrder_CcRoles", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.CcRoles))))
	router.Any("/OfficeOrder_Tasksummary", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.TaskDetailsTaskSummary))))
}
