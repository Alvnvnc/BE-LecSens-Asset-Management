# Database Management Guide

## Seeding Control

Aplikasi ini menggunakan sistem seeding otomatis yang dapat dikontrol melalui environment variables.

### Environment Variables

- `AUTO_SEED`: Mengontrol apakah seeding otomatis dijalankan
  - `true`: Seeding akan dijalankan
  - `false`: Seeding tidak akan dijalankan (default)

- `FORCE_RESEED`: Mengontrol apakah akan melakukan reseed meskipun data sudah ada
  - `true`: Akan menghapus data lama dan seeding ulang
  - `false`: Akan skip seeding jika data sudah ada (default)

- `ENVIRONMENT`: Environment aplikasi
  - `development`: Mode development
  - `production`: Mode production

### Cara Penggunaan

#### 1. Setup Awal (First Run)
```bash
# Copy file environment
cp .env.example .env

# Edit file .env dan set:
AUTO_SEED=true
FORCE_RESEED=false

# Jalankan docker compose
docker-compose up -d
```

#### 2. Running Normal (Data sudah ada)
```bash
# File .env dengan setting:
AUTO_SEED=false
FORCE_RESEED=false

# Jalankan docker compose
docker-compose up -d
```

#### 3. Force Reseed (Reset Database)
```bash
# File .env dengan setting:
AUTO_SEED=true
FORCE_RESEED=true

# Jalankan docker compose
docker-compose up -d --force-recreate
```

### Perilaku Seeding

1. **Database Kosong**: Akan menjalankan migration dan seeding lengkap
2. **Database Sudah Ada Data**: 
   - Jika `FORCE_RESEED=false`: Skip seeding
   - Jika `FORCE_RESEED=true`: Bersihkan data lama dan seeding ulang
3. **Tables Ada tapi Kosong**: Akan menjalankan seeding tanpa cleanup

### Persistent Volume

Database menggunakan persistent volume `pgdata` yang akan mempertahankan data meskipun container dihapus.

Untuk benar-benar reset database:
```bash
# Stop containers
docker-compose down

# Remove volumes
docker volume rm asset-management_pgdata

# Start fresh
docker-compose up -d
```

### Urutan Seeding

Seeding dilakukan dalam urutan yang sesuai dengan foreign key constraints:
1. Locations
2. Asset Types
3. Sensor Types
4. Measurement Types
5. Measurement Fields
6. Assets
7. Asset Sensors
8. Sensor Thresholds
9. Sensor Status
10. Sensor Readings (7 hari terakhir)
11. Sensor Logs
12. Asset Alerts

### Troubleshooting

#### Database tidak terhubung
- Pastikan container database sudah running
- Check healthcheck status: `docker-compose ps`

#### Seeding error
- Check logs: `docker-compose logs app`
- Pastikan urutan seeding sudah benar
- Check foreign key constraints

#### Data tidak persistent
- Pastikan volume `pgdata` tidak dihapus
- Check volume: `docker volume ls`
