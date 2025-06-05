package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // driver MySQL untuk database/sql
	"github.com/joho/godotenv"          // package untuk load file .env
)

// InitDB melakukan inisialisasi koneksi ke database MySQL dan mengembalikan *sql.DB
func InitDB() *sql.DB {
	// Load konfigurasi dari file .env, jika gagal maka program langsung berhenti
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ambil variabel environment yang diperlukan dari file .env
	user := os.Getenv("DB_USER")   // username database
	pass := os.Getenv("DB_PASS")   // password database
	host := os.Getenv("DB_HOST")   // alamat host database, misal localhost atau IP
	port := os.Getenv("DB_PORT")   // port database, misal 3306
	dbname := os.Getenv("DB_NAME") // nama database

	// Buat Data Source Name (DSN) untuk koneksi MySQL
	// Tambahkan opsi parseTime=true agar driver dapat membaca tipe waktu
	// Tambahkan lokasi waktu (timezone) Asia/Jakarta
	loc := "&loc=Asia%2FJakarta"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true%s", user, pass, host, port, dbname, loc)

	// Buka koneksi ke database dengan driver MySQL dan DSN yang sudah dibuat
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke db", err) // hentikan program jika gagal koneksi
	}

	// Test koneksi ke database agar yakin koneksi berhasil
	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak bisa diakses", err) // hentikan program jika koneksi gagal
	}

	// Kembalikan objek *sql.DB yang siap digunakan untuk query database
	return db
}