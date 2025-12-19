package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func GetConnection() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	log.Printf("Mencoba koneksi ke database %s di %s:%s...", dbName, dbHost, dbPort)

	if dbHost == "" {
		dbHost = "localhost"
	}

	connStr := fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbHost, dbPort, dbName)

	dbUrl := os.Getenv("DB_URL")

	if dbUrl != "" {
		connStr = dbUrl
	}

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Error membuka koneksi database: %v", err)
	}

	// Set konfigurasi koneksi
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	// Cek koneksi dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Gagal terhubung ke database dalam 10 detik: %v", err)
		log.Printf("Mencoba koneksi ulang...")

		// Coba lagi dengan timeout yang lebih lama
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err != nil {
			db.Close()
			log.Fatalf("Koneksi database gagal setelah percobaan ulang: %v", err)
		}
	}

	log.Printf("Berhasil terhubung ke database %s", dbName)
	log.Printf("Max Open Connections: %d", db.Stats().MaxOpenConnections)
	log.Printf("Open Connections: %d", db.Stats().OpenConnections)
	log.Printf("In Use Connections: %d", db.Stats().InUse)
	log.Printf("Idle Connections: %d", db.Stats().Idle)

	return db
}

// package app

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/joho/godotenv"
// )

// func GetConnection() *sql.DB {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file:", err)
// 	}

// 	dbUser := os.Getenv("DB_USER")
// 	dbPort := os.Getenv("DB_PORT")
// 	dbName := os.Getenv("DB_NAME")
// 	dbHost := os.Getenv("DB_HOST")

// 	log.Printf("Mencoba koneksi ke database %s di %s:%s...", dbName, dbHost, dbPort)

// 	if dbHost == "" {
// 		dbHost = "localhost"
// 	}

// 	// Tambahkan timeout parameters untuk mengatasi lock timeout
// 	// innodb_lock_wait_timeout=120: timeout untuk menunggu lock (dalam detik)
// 	// readTimeout=60s: timeout untuk membaca data
// 	// writeTimeout=60s: timeout untuk menulis data
// 	// timeout=10s: timeout untuk koneksi awal
// 	connStr := fmt.Sprintf(
// 		"%s@tcp(%s:%s)/%s?parseTime=true&timeout=10s&readTimeout=60s&writeTimeout=60s",
// 		dbUser, dbHost, dbPort, dbName,
// 	)

// 	dbUrl := os.Getenv("DB_URL")

// 	if dbUrl != "" {
// 		// Jika DB_URL digunakan, pastikan timeout parameters ada
// 		if dbUrl[len(dbUrl)-1] != '?' {
// 			if !contains(dbUrl, "readTimeout") {
// 				dbUrl += "&readTimeout=60s&writeTimeout=60s"
// 			}
// 		} else {
// 			if !contains(dbUrl, "readTimeout") {
// 				dbUrl += "readTimeout=60s&writeTimeout=60s"
// 			} else {
// 				dbUrl = dbUrl[:len(dbUrl)-1]
// 			}
// 		}
// 		connStr = dbUrl
// 	}

// 	db, err := sql.Open("mysql", connStr)
// 	if err != nil {
// 		log.Fatalf("Error membuka koneksi database: %v", err)
// 	}

// 	// Set konfigurasi koneksi
// 	db.SetMaxOpenConns(50)
// 	db.SetMaxIdleConns(25)
// 	db.SetConnMaxIdleTime(10 * time.Minute)
// 	db.SetConnMaxLifetime(60 * time.Minute)

// 	// Cek koneksi dengan timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	err = db.PingContext(ctx)
// 	if err != nil {
// 		log.Printf("Gagal terhubung ke database dalam 10 detik: %v", err)
// 		log.Printf("Mencoba koneksi ulang...")

// 		// Coba lagi dengan timeout yang lebih lama
// 		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
// 		defer cancel()

// 		err = db.PingContext(ctx)
// 		if err != nil {
// 			db.Close()
// 			log.Fatalf("Koneksi database gagal setelah percobaan ulang: %v", err)
// 		}
// 	}

// 	log.Printf("Berhasil terhubung ke database %s", dbName)
// 	log.Printf("Max Open Connections: %d", db.Stats().MaxOpenConnections)
// 	log.Printf("Open Connections: %d", db.Stats().OpenConnections)
// 	log.Printf("In Use Connections: %d", db.Stats().InUse)
// 	log.Printf("Idle Connections: %d", db.Stats().Idle)

// 	return db
// }

// func contains(s, substr string) bool {
// 	return len(s) >= len(substr) &&
// 		(s == substr ||
// 			len(s) > len(substr) &&
// 				(s[:len(substr)] == substr ||
// 					s[len(s)-len(substr):] == substr ||
// 					findSubstring(s, substr)))
// }

// func findSubstring(s, substr string) bool {
// 	for i := 0; i <= len(s)-len(substr); i++ {
// 		if s[i:i+len(substr)] == substr {
// 			return true
// 		}
// 	}
// 	return false
// }
