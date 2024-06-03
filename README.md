## PCS HOME TEST

Project ini menggunakan bahasa program golang & database postgres

## Installation

1. Replace dan ganti nama config.yml.example menjadi config.yml
2. Setup konfigurasi database
3. Jalankan migrasi database menggunakan command `go run main.go migration up`
4. Jalankan service menggunakan command `go run main.go start`, akses menggunakan port `8081`

## Docker

1. Jalankan command `docker-compose up -d --build`
2. Jalankan command `docker exec pcs go run main.go migration up`
3. Run command `sh bash.sh`
3. Akses menggunakan port `8081`

## API Docs

1. Import collection `Pcs.postman_collection.json` pada aplikasi Postman
2. Akses user tersedia menggunakan username dengan password `secret` :
    - Seller / Role seller
    - Buyer / Role buyer
