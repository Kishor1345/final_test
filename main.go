// This API serves as a core component of the Workflow-Human Resources Module system, functioning as a secure middleware layer that bridges the React-based frontend with the Microsoft SQL Server (MSSQL) backend.
//
// Once data is fetched from the database, it is encrypted, and sent to the frontend for display and user interaction.
//
// This middleware design promotes separation of concerns, enhances security through built-in authentication and encryption, and ensures smooth communication between the user interface and the underlying data infrastructure that powers Human Resources workflows and operations.
// package main

// import (
// 	"Hrmodule/mainroutes"
// )

// // new
// // main is the entry point of the application.
// // It calls Registerroutes to bind API endpoints and start the server.
// func main() {
// 	mainroutes.Registerroutes()
// }

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	credentials "Hrmodule/dbconfig"
	"Hrmodule/mainroutes"
)

func main() {

	// 1. Get PostgreSQL connection string
	connStr := credentials.Getdatabasemeivan()

	// 2. Initialize DB pool (defaults only)
	if err := credentials.InitPostgresPool(connStr); err != nil {

		log.Fatalf("Failed to initialize DB pool: %v", err)
	}
	log.Println("✅ Database pool initialized")

	// 3. Start routes (existing behavior)
	go mainroutes.Registerroutes()
	log.Println("🚀 Server started")

	// 4. Wait for shutdown signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Shutting down server...")

	// 5. Close DB pool
	if err := credentials.ClosePostgresPool(); err != nil {
		log.Printf("Error closing DB pool: %v", err)
	} else {
		log.Println("✅ Database pool closed")
	}
}
