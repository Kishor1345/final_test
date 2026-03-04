// Package mainroutes handles the central initialization and orchestration of application routing.
// It configures global middleware, aggregates sub-routes from various modules, and starts the secure server.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/mainroutes
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 30-07-2025
// Last Modified By: Sivabala
// Last Modified Date: 30-07-2025
package mainroutes

import (
	"Hrmodule/routes"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Registerroutes initializes the Gin engine, configures global CORS settings,
// and registers all module-specific API routes (Common, Login, Office Order, NOC, and QMS).
//
// It also starts the HTTPS server on port 2000 using the provided SSL certificates.
func Registerroutes() {
	// Create Gin router
	router := gin.Default()

	// --- CORS Configuration ---
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Change this in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	/*************************************************Cron Jobs*****************************************************/
	//routes.Cronjobs(router)
	/***************************************************************************************************************/
	/*************************************************Common API's**************************************************/
	routes.Common(router)
	/***************************************************************************************************************/
	/*************************************************LOGIN  API's**************************************************/
	routes.Login(router)
	/***************************************************************************************************************/
	/*************************************************OFFICE ORDER API's********************************************/
	routes.Officeorder(router)
	/***************************************************************************************************************/
	/************************************************* NOC API's ***************************************************/
	//routes.NOC(router)
	/*********************************************QMS API's*********************************************************/
	routes.QMS(router)
	/***************************************************************************************************************/
	/*************************************************Employee Efile API's******************************************/
	routes.EmployeeEfile(router)
	/***************************************************************************************************************/
	/*************************************************Staffadditionaldetails API's**********************************/
	routes.Staffadditionaldetails(router)
	/***************************************************************************************************************/
	/*************************************************PDF API's*****************************************************/
	routes.PDF(router)
	/***************************************************************************************************************/
	/*************************************************NOC API's*****************************************************/
	routes.NOC(router)
	/***************************************************Quartersmasterestate********************************************/
	routes.Quartersmasterestate(router)
	/***************************************************Criteria***************************************************/
	routes.Criteria(router)
	/***************************************************Circular***************************************************/
	routes.Circular(router)
	/***************************************************Quarters****************************************************/
	routes.Quarters(router)
	/***************************************************************************************************************/
	// --- HTTPS SERVER START ---
	fmt.Println("Server starting on port 2580 (HTTPS Enabled)")

	certFile := "certificate.pem"
	keyFile := "key.pem"

	if err := router.RunTLS(":2580", certFile, keyFile); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}



