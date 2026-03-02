// Package routes defines the HTTP route registrations for the Circular Master module.
//
//path:/var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/routes
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:13/02/2026
package routes

import (
	"Hrmodule/auth"
	controllerscircular "Hrmodule/controllers/HR_008"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Circular(router *gin.Engine) {

	router.Any("/HR-CIR-003", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscircular.CircularEligibilityChoice))))
	router.Any("/HR-CIR-013", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscircular.CircularDetailFetch))))
	router.Any("/HR-CIR-002", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscircular.CircularDetailFetchForApproval))))
	router.Any("/HR-CIR-001", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscircular.CircularInsert))))
	router.Any("/HR-CIR-014", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscircular.CircularQuartersNumber))))
}
