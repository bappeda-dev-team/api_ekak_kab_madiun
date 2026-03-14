// @title           EKAK Kabupaten Madiun API
// @version         1.0
// @description     API Backend Kertas Kerja Kabupaten Madiun
// @termsOfService  http://swagger.io/terms/
// @contact.name    Tim Pengembang
// @schemes http https
// @BasePath  /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/middleware"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "ekak_kabupaten_madiun/docs"

	"github.com/joho/godotenv"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	host := os.Getenv("host")
	port := os.Getenv("port")
	addr := fmt.Sprintf("%s:%s", host, port)

	cors := helper.NewCORSMiddleware()

	if addr == ":" {
		addr = "localhost:8080"
	}

	return &http.Server{
		Addr:    addr,
		Handler: cors.Handler(authMiddleware),
	}
}

func main() {
	runSeeder := flag.Bool("seed", false, "Jalankan database seeder")
	flag.Parse()
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Cek flag seeder
	if *runSeeder {
		log.Println("Menjalankan database seeder...")
		seeder := InitializeSeeder()
		seeder.SeedAll()
		log.Println("Seeder selesai dijalankan")
		return
	}
	// Initialize dan jalankan server
	server := InitializeServer()
	log.Printf("Server berjalan di %s", server.Addr)
	err = server.ListenAndServe()
	helper.PanicIfError(err)
}
