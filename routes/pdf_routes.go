// Package routes defines  HTTP routes for the PDF.
//
// --- Creator's Info ---
//
// Creator: Sridharan
// Created On: 31-10-2025

// Last Modified By: Sridharan

// Last Modified Date: 31-10-2025
package routes

import (
	controllerspdf "Hrmodule/controllers/pdf"

	"net/http"

	"github.com/gin-gonic/gin"
)

func PDF(router *gin.Engine) {
	router.Any("/pdfapi", gin.WrapH(http.HandlerFunc(controllerspdf.GeneratePDFHandler)))
}
