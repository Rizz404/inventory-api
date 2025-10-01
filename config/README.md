# Config Package

Package ini berisi konfigurasi untuk inisialisasi berbagai komponen aplikasi seperti database dan external clients.

## Struktur

```
config/
├── client/           # Individual client initialization files
│   ├── cloudinary.go # Cloudinary client
│   ├── fcm.go        # Firebase Cloud Messaging client
│   └── README.md     # Client documentation
├── clients.go        # Main clients struct and initialization
├── database.go       # Database connection initialization
└── README.md         # This file
```

## Clients

Package ini menyediakan struktur `Clients` yang menampung semua external service clients:

### Clients yang Tersedia

1. **Cloudinary** - File upload dan image processing
   - Diaktifkan melalui environment variables:
     - `CLOUDINARY_URL` (format: `cloudinary://api_key:api_secret@cloud_name`)
     - Atau: `CLOUDINARY_CLOUD_NAME`, `CLOUDINARY_API_KEY`, `CLOUDINARY_API_SECRET`

2. **FCM (Firebase Cloud Messaging)** - Push notification service
   - Diaktifkan dengan setting `ENABLE_FCM=true`
   - Memerlukan Firebase credentials melalui environment variables:
     - `FIREBASE_TYPE`
     - `FIREBASE_PROJECT_ID`
     - `FIREBASE_PRIVATE_KEY_ID`
     - `FIREBASE_PRIVATE_KEY`
     - `FIREBASE_CLIENT_EMAIL`
     - `FIREBASE_CLIENT_ID`
     - `FIREBASE_AUTH_URI`
     - `FIREBASE_TOKEN_URI`
     - `FIREBASE_AUTH_PROVIDER_X509_CERT_URL`
     - `FIREBASE_CLIENT_X509_CERT_URL`
     - `FIREBASE_UNIVERSE_DOMAIN`

### Penggunaan

```go
import "github.com/Rizz404/inventory-api/config"

// Inisialisasi semua clients
clients := config.InitializeClients()

// Akses individual client
cloudinaryClient := clients.Cloudinary
fcmClient := clients.FCM

// Clients yang tidak dikonfigurasi akan bernilai nil
if clients.Cloudinary != nil {
    // Gunakan Cloudinary
}
```

## Database

Package ini menyediakan fungsi `InitializeDatabase()` untuk inisialisasi koneksi database PostgreSQL menggunakan GORM.

### Environment Variables

- `DSN` - Data Source Name untuk PostgreSQL (required)

### Penggunaan

```go
import "github.com/Rizz404/inventory-api/config"

// Inisialisasi database
db := config.InitializeDatabase()

// Dapatkan generic database object untuk connection management
sqlDB, err := db.DB()
if err != nil {
    log.Fatal(err)
}
defer sqlDB.Close()
```

## Menambahkan Client Baru

Untuk menambahkan client baru, ikuti langkah berikut:

### 1. Buat file client baru di `config/client/`

Contoh: `config/client/newclient.go`

```go
package client

import (
    "log"
    "os"

    "github.com/Rizz404/inventory-api/internal/client/newclient"
)

// InitNewClient initializes NewClient
func InitNewClient() *newclient.Client {
    // Check if enabled
    if os.Getenv("ENABLE_NEW_CLIENT") != "true" {
        log.Printf("NewClient disabled")
        return nil
    }

    // Get configuration from env
    apiKey := os.Getenv("NEW_CLIENT_API_KEY")
    if apiKey == "" {
        log.Printf("Warning: NEW_CLIENT_API_KEY not set")
        return nil
    }

    // Initialize client
    client, err := newclient.NewClient(apiKey)
    if err != nil {
        log.Printf("Warning: Failed to initialize NewClient: %v", err)
        return nil
    }

    log.Printf("NewClient initialized successfully")
    return client
}
```

### 2. Update struct `Clients` di `clients.go`

```go
// Clients holds all external service clients
type Clients struct {
    Cloudinary *cloudinary.Client
    FCM        *fcm.Client
    NewClient  *newclient.Client  // ← Field baru
}
```

### 3. Update fungsi `InitializeClients()` di `clients.go`

```go
import "github.com/Rizz404/inventory-api/config/client"

func InitializeClients() *Clients {
    return &Clients{
        Cloudinary: client.InitCloudinary(),
        FCM:        client.InitFCM(),
        NewClient:  client.InitNewClient(),  // ← Tambahkan ini
    }
}
```

### 4. Update dokumentasi di `config/client/README.md`

Tambahkan informasi tentang client baru beserta environment variables yang diperlukan.

## Best Practices

1. **Error Handling**: Semua error dalam inisialisasi client dicatat sebagai warning, bukan fatal error. Ini memungkinkan aplikasi tetap berjalan meskipun beberapa service tidak tersedia.

2. **Nil Checks**: Selalu periksa apakah client tidak nil sebelum menggunakannya:
   ```go
   if clients.Cloudinary != nil {
       // Gunakan client
   }
   ```

3. **Environment Variables**: Gunakan environment variables untuk konfigurasi sensitive seperti API keys dan credentials.

4. **Logging**: Setiap inisialisasi client harus mencatat status (success/warning) untuk memudahkan debugging.
