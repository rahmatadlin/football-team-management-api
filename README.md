# Football Team Management API

API REST untuk manajemen tim sepak bola amatir (perusahaan) dengan **Go**, **Gin**, **GORM**, **MySQL**, dan **JWT**. Arsitektur mengikuti pemisahan **handler → service → repository**.

## Fitur utama

- Manajemen tim (CRUD + soft delete)
- Manajemen pemain (per tim, nomor punggung unik dalam tim)
- Jadwal pertandingan (`home_team_id` ≠ `away_team_id`)
- Pelaporan hasil + gol per pemain
- Endpoint laporan: skor, pemenang, top skor pertandingan, akumulasi kemenangan tim sampai pertandingan tersebut
- Autentikasi admin: **bcrypt** + **JWT** (Bearer)

## Prasyarat

- Go **1.22+**
- MySQL **8.x** (atau 5.7 dengan utf8mb4)
- Database kosong (nama DB sesuai `.env`)

## Konfigurasi

```bash
cp .env.example .env
```

Isi minimal: `JWT_SECRET`, kredensial MySQL (`DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`), opsional `PORT`.

Buat database:

```sql
CREATE DATABASE football_team_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## Menjalankan API

```bash
go run ./cmd/api
```

Server default: `http://localhost:8080` (atau `PORT` di `.env`).

Health check:

```bash
curl -s http://localhost:8080/health
```

## Migrasi & data contoh

**Auto-migrate** dijalankan saat startup `cmd/api` (membuat/ memperbarui tabel).

**Seed** (admin + data demo):

```bash
go run ./cmd/seed
```

Admin default dari `.env`: `ADMIN_EMAIL` / `ADMIN_PASSWORD` (lihat `.env.example`).

## Autentikasi

Semua endpoint di bawah `/api/v1` kecuali `POST /api/v1/auth/login` memerlukan header:

```http
Authorization: Bearer <access_token>
```

### Login

**Request**

```http
POST /api/v1/auth/login
Content-Type: application/json
```

```json
{
  "email": "admin@example.com",
  "password": "admin123"
}
```

**Response (contoh)**

```json
{
  "status": "success",
  "message": "login successful",
  "data": {
    "access_token": "<jwt>",
    "token_type": "Bearer",
    "expires_in": null,
    "admin": {
      "id": "…uuid…",
      "email": "admin@example.com"
    }
  }
}
```

## Format respons

Sukses:

```json
{
  "status": "success",
  "message": "…",
  "data": {}
}
```

Gagal:

```json
{
  "status": "error",
  "message": "pesan error",
  "data": null
}
```

## Endpoint (ringkasan)

| Metode | Path | Keterangan |
|--------|------|------------|
| POST | `/api/v1/auth/login` | Login admin |
| GET | `/api/v1/teams` | Daftar tim |
| POST | `/api/v1/teams` | Buat tim |
| GET | `/api/v1/teams/:id` | Detail tim |
| PUT | `/api/v1/teams/:id` | Update tim |
| DELETE | `/api/v1/teams/:id` | Soft delete tim |
| GET | `/api/v1/teams/:id/players` | Pemain per tim |
| POST | `/api/v1/teams/:id/players` | Buat pemain |
| PUT | `/api/v1/players/:id` | Update pemain |
| DELETE | `/api/v1/players/:id` | Soft delete pemain |
| GET | `/api/v1/matches/schedules` | Daftar jadwal |
| POST | `/api/v1/matches/schedules` | Buat jadwal |
| POST | `/api/v1/matches/:id/results` | Simpan hasil + daftar gol |
| GET | `/api/v1/matches/:id/report` | Laporan pertandingan |

**Daftar (list):** `GET /teams`, `GET /teams/:id/players`, dan `GET /matches/schedules` memakai query **Unscoped** sehingga baris yang sudah di-soft-delete tetap muncul dan field `deleted_at` terlihat. Endpoint detail tunggal (mis. `GET /teams/:id`) tetap hanya mengembalikan data **aktif** (bukan terhapus).

**GET pemain per tim:** respons berupa array objek pemain **tanpa** `team_id` dan **tanpa** objek `team` (konteks tim sudah dari URL).

## Contoh permintaan

**Buat tim** (`Authorization` wajib):

```json
POST /api/v1/teams
{
  "name": "FC Contoh",
  "logo_url": "https://example.com/logo.png",
  "founded_year": 2015,
  "address": "Jl. Contoh 1",
  "city": "Jakarta"
}
```

**Buat pemain** (posisi: `striker` | `midfielder` | `defender` | `goalkeeper`):

```json
POST /api/v1/teams/<team_id>/players
{
  "name": "Pemain Satu",
  "height": 175.5,
  "weight": 70,
  "position": "striker",
  "jersey_number": 9
}
```

**Jadwal** (`match_date`: `YYYY-MM-DD`, `match_time`: `HH:MM:SS`):

```json
POST /api/v1/matches/schedules
{
  "match_date": "2026-04-01",
  "match_time": "15:30:00",
  "home_team_id": "<uuid>",
  "away_team_id": "<uuid_lain>"
}
```

**Hasil pertandingan:** jumlah elemen `goals` **harus sama** dengan `home_score + away_score` (mis. skor 2–1 → tepat 3 entri gol). Skor 0–0 → `goals`: `[]`.

```json
POST /api/v1/matches/<match_uuid>/results
{
  "home_score": 2,
  "away_score": 1,
  "goals": [
    { "player_id": "<uuid>", "goal_time": 23 },
    { "player_id": "<uuid>", "goal_time": 55 },
    { "player_id": "<uuid>", "goal_time": 78 }
  ]
}
```

**Laporan**:

```http
GET /api/v1/matches/<match_uuid>/report
```

**Response (contoh struktur `data`)**:

```json
{
  "schedule": {
    "match_date": "2026-04-01",
    "match_time": "15:30:00"
  },
  "home_team": { "id": "…", "name": "…", "city": "…" },
  "away_team": { "id": "…", "name": "…", "city": "…" },
  "home_score": 2,
  "away_score": 1,
  "match_result": "HOME_WIN",
  "top_scorer": {
    "player_id": "…",
    "player_name": "…",
    "goals": 2
  },
  "home_team_wins_until_match": 1,
  "away_team_wins_until_match": 0
}
```

`match_result`: `HOME_WIN` | `AWAY_WIN` | `DRAW`.  
`top_scorer` bisa `null` jika tidak ada gol tercatat.

## Struktur folder

```
cmd/api/          # entrypoint API
cmd/seed/         # seed data
config/           # env & koneksi DB
models/           # entitas GORM
repository/       # akses data
service/          # logika bisnis
handler/          # HTTP Gin
middleware/       # JWT admin
routes/           # registrasi route
utils/            # JWT, bcrypt, response, apperror
```

## Pengembangan

Build biner:

```bash
go build -o bin/api ./cmd/api
go build -o bin/seed ./cmd/seed
```

## Catatan produksi

- Ganti `JWT_SECRET` dengan rahasia panjang dan acak.
- Gunakan HTTPS di depan reverse proxy.
- Sesuaikan level log GORM / mode Gin (`APP_ENV=production`).
- Pertimbangkan rate limiting dan rotasi token sesuai kebutuhan.
