// Package routes defines the HTTP route registrations for user authentication,
// session management, and LDAP integration.
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
	controllerslogin "Hrmodule/controllers/login"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login initializes and registers authentication-related routes on the provided Gin engine.
//
// This function sets up endpoints for the complete login lifecycle, including:
//   - Initial authentication via LDAP and credential exchange.
//   - Multi-factor authentication through OTP generation and validation.
//   - Session management (timeout tracking and session data retrieval).
//   - Secure utility services such as user activity logging and data decryption.
//
// Endpoints requiring valid JWT tokens are protected by the JwtMiddleware.
func Login(router *gin.Engine) {

	// Public Authentication Endpoints
	router.Any("/login", gin.WrapH(http.HandlerFunc(controllerslogin.Getkey)))
	router.Any("/HRldap", gin.WrapH(http.HandlerFunc(controllerslogin.HandleLDAPAuth)))
	router.Any("/HRldapfailure", gin.WrapH(http.HandlerFunc(controllerslogin.HandleLDAPAuthf)))
	router.Any("/Loginotp", gin.WrapH(http.HandlerFunc(controllerslogin.InsertOTPHandler)))
	router.Any("/Loginotpupdate", gin.WrapH(http.HandlerFunc(controllerslogin.ValidateOTPHandler)))

	// Protected Session & Utility Endpoints (JWT Required)
	router.Any("/SessionTimeout", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.SessionTimeoutHandler))))
	router.Any("/Sessiondata", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.SessionData))))
	router.Any("/InsertUserActivityLog", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.InsertUserActivityLog))))
	router.Any("/SendOTP", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.SendOTPHandler))))
	router.Any("/Datadecrypt", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.DatadecryptHandler))))
	router.Any("/DatadecryptKey", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.Datadecryptsessionkey))))
	router.Any("/Menuvalidation", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.MenuValidationHandler))))
}
