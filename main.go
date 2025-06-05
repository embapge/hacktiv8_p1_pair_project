package main

import (
	"context"
	"pairproject/cli"    // Package untuk menangani command line interface (CLI)
	"pairproject/config" // Package untuk konfigurasi aplikasi, termasuk inisialisasi database
)

func main() {
	// Inisialisasi koneksi database menggunakan konfigurasi dari package config
	db := config.InitDB()
	// Pastikan koneksi database ditutup ketika aplikasi selesai dijalankan
	defer db.Close()

	// Buat context dasar untuk digunakan di seluruh aplikasi (misal untuk request scope)
	ctx := context.Background()

	// Buat handler CLI dengan memasukkan koneksi db dan context
	cli := cli.NewCLIHandler(db, ctx)

	// Jalankan menu utama CLI agar user bisa mulai berinteraksi dengan aplikasi
	cli.Menu()
}