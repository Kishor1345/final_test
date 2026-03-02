// Package routes defines the HTTP route registrations for the Quarters Application.
//
//path:/var/www/html/go_projects/HRMODULE/kishorenew/hr2000/Meivan/routes
//
// --- Creator's Info ---
//
// Creator: Kishorekumar
//
// Created On:23/02/2026
package routes

import (
	"Hrmodule/auth"
	controllersQuarters "Hrmodule/controllers/HR_009"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Quarters(router *gin.Engine) {

	router.Any("/HR-EQA-007", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersQuarters.QuartersFetchForPreference))))
}
