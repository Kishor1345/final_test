// Package routes defines the HTTP route registrations and background job initializations
// for the Office Order module, specifically handling data synchronization and cron scheduling.
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
	"Hrmodule/controllers/cronjob"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// InitCronforofficeorder initializes and starts a background scheduler for the Office Order module.
//
// It schedules the SyncData function to run once every hour. Additionally, it triggers
// an immediate synchronization in a separate goroutine upon initialization to ensure
// data consistency at startup.
//
// Returns:
//   - A pointer to the cron.Cron instance if successful.
//   - nil if the cron job failed to be added to the scheduler.
func InitCronforofficeorder() *cron.Cron {
	c := cron.New()
	_, err := c.AddFunc("@every 1h", cronjob.SyncData)
	if err != nil {
		log.Println("Failed to start hourly cron:", err)
		return nil
	}
	c.Start()
	log.Println("Hourly cron started for SyncData()")
	go cronjob.SyncData()
	return c
}

// Cronjobs sets up the background schedulers and registers HTTP endpoints for
// manual synchronization triggers.
//
// This function initializes the hourly office order cron job and maps the
// "/OfficeOrder_Sync" endpoint to the HandleSync controller, allowing for
// on-demand data synchronization.
func Cronjobs(router *gin.Engine) {

	// --- OFFICE ORDER CRON JOBS ---
	InitCronforofficeorder()

	router.Any("/OfficeOrder_Sync", gin.WrapH(http.HandlerFunc(cronjob.HandleSync)))
}
