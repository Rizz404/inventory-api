# Alur Detail Aplikasi Inventaris (Dengan Lokalisasi, Optimalisasi DB, Data Matrix & Google Maps)

## 1. ALUR LOGIN (Semua Role)

### Tampilan Mobile:
- **Splash Screen**: Logo PT Fujiyama Technology Solution.
- **Login Screen**:
  - Input field untuk name.
  - Input field untuk password.
  - Pilihan Bahasa (Dropdown/Icon): Indonesia (Default), English, 日本語.
  - Tombol "Login".

### Proses Backend:
1. User memasukkan credentials dan memilih bahasa. **Contoh**: `name`: "ahmad.admin", `password`: "Rahasia123", `lang`: "en-US".
2. Sistem query ke tabel `users` mencari `name` = 'ahmad.admin' dan `is_active` = TRUE.
3. Verifikasi `password_hash`. Jika cocok, sistem mendapat data user: `id`: '01J4X1...', `full_name`: 'Ahmad Susanto', `role`: 'Admin'.
4. Generate token (JWT) yang berisi `id`, `role`, dan `lang` user, lalu arahkan ke Dashboard Admin. Bahasa yang ditampilkan di seluruh aplikasi akan mengikuti `lang` ini.

---

## 2. ALUR ADMIN (Role: 'Admin')

### Dashboard Utama Admin
**Tampilan Mobile (Contoh Bahasa Inggris):**
- Header: "Welcome, Ahmad Susanto".
- Card statistik: "Total Assets: 150", "Under Maintenance: 5", "Warranty Expiring: 8".
- Quick actions: Add Asset, Add User, View Reports.
- Ikon notifikasi dengan badge merah.

### Alur 1: Mengelola Master Data (Dengan Lokalisasi)
**Cerita:** Pak Ahmad sebagai Admin perlu menyiapkan data dasar seperti kategori dan lokasi dalam tiga bahasa.

**Langkah-langkah:**
1. **Kelola Kategori Aset (Dengan Hierarki & Terjemahan)**
   - Tap menu "Master Data" → "Categories".
   - Tampil daftar kategori (nama ditampilkan sesuai bahasa yang dipilih saat login).
   - **Langkah 1: Menambah Kategori Induk**
     - Tap tombol "+".
     - Form input diisi:
       - **Data Utama**:
         - `category_code`: "ELK"
         - `parent_id`: (Dikosongkan)
       - **Terjemahan**:
         - `category_name` (ID): "Elektronik"
         - `description` (ID): "Semua aset yang berhubungan dengan elektronik."
         - `category_name` (EN): "Electronics"
         - `description` (EN): "All assets related to electronics."
         - `category_name` (JA): "電子機器"
         - `description` (JA): "電子機器に関連するすべての資産。"
     - Simpan. Data baru masuk ke tabel `categories` dan `category_translations`.
   - **Langkah 2: Menambah Kategori Anak**
     - Tap tombol "+".
     - Form input diisi:
       - **Data Utama**:
         - `category_code`: "LPT"
         - `parent_id`: Pilih "Elektronik" dari daftar.
       - **Terjemahan**:
         - `category_name` (ID): "Laptop"
         - `description` (ID): "Komputer portabel untuk karyawan."
         - `category_name` (EN): "Laptop"
         - `description` (EN): "Portable computers for employees."
         - `category_name` (JA): "ラップトップ"
         - `description` (JA): "従業員向けのポータブルコンピュータ。"
     - Simpan.

2. **Kelola Lokasi (Dengan Terjemahan)**
   - Tap menu "Master Data" → "Locations".
   - Tap tombol "+" untuk menambah lokasi baru.
   - Form input diisi:
     - **Data Utama**:
       - `location_code`: "MKT-L3"
       - `building`: "Gedung A"
       - `floor`: "3"
       - **Koordinat GPS**: Otomatis diambil dari GPS device atau input manual
         - `latitude`: -6.3734, `longitude`: 106.8284 (Jakarta)
     - **Terjemahan**:
       - `location_name` (ID): "Ruang Marketing Lantai 3"
       - `location_name` (EN): "Marketing Room 3rd Floor"
       - `location_name` (JA): "マーケティングルーム3階"
   - Simpan. Data baru masuk ke tabel `locations` dan `location_translations`.

3. **Kelola Maintenance Schedules (Dengan Terjemahan)**
   - Tap menu "Maintenance" → "Schedules".
   - Tap tombol "+" untuk menambah jadwal baru.
   - Form input diisi:
     - **Data Utama**:
       - `asset_id`: Pilih asset
       - `maintenance_type`: "Preventive"
       - `scheduled_date`: "2025-09-15"
       - `frequency_months`: 6
     - **Terjemahan**:
       - `title` (ID): "Perawatan Rutin Laptop"
       - `description` (ID): "Pembersihan dan update sistem operasi."
       - `title` (EN): "Routine Laptop Maintenance"
       - `description` (EN): "Cleaning and operating system update."
       - `title` (JA): "ラップトップの定期メンテナンス"
       - `description` (JA): "クリーニングとオペレーティングシステムの更新。"
   - Simpan. Data masuk ke `maintenance_schedules` dan `maintenance_schedule_translations`.

4. **Kelola User**
   - (Alur tidak berubah, karena data user seperti nama tidak dilokalisasi).

### Alur 2: Mengelola Aset
**Cerita:** Pak Ahmad mendaftarkan laptop Dell baru.

**Langkah-langkah:**
1. **Registrasi Aset Baru**
   - Tap "+" di dashboard.
   - Form registrasi aset diisi:
     - `asset_tag`: "LPT-2025-001"
     - `asset_name`: "Laptop Dell Latitude 5430"
     - `category_id`: Pilih "Laptop" (nama kategori ditampilkan sesuai bahasa login Pak Ahmad).
     - `brand`: "Dell"
     - `model`: "Latitude 5430"
     - `serial_number`: "DL5430-2025-001"
     - `purchase_date`: "2025-01-15"
     - `purchase_price`: 15000000
     - `vendor_name`: "PT Dell Indonesia"
     - `warranty_end`: "2027-01-15"
     - `location_id`: Pilih "Gudang IT" (nama lokasi ditampilkan sesuai bahasa).
   - Simpan.

2. **Menyiapkan Identifier Fisik (Data Matrix)**
   - Sistem generate Data Matrix code dengan value "LPT-2025-001" dan update field `data_matrix_value`.
   - Data Matrix dapat menyimpan lebih banyak data dibanding QR Code dan lebih tahan terhadap kerusakan fisik.

### Alur 3: Monitoring dan Reporting
**Langkah-langkah:**
1. **Dashboard Analytics (Contoh Bahasa Jepang)**
   - Chart menunjukkan: "カテゴリ ラップトップ: 50 ユニット", "カテゴリ モニター: 45 ユニット".
2. **Generate Laporan**
   - Menu "Laporan" → Pilih "Laporan Inventaris Lengkap".
   - Filter `location_id` = "Ruang Marketing Lantai 3" (pencarian bisa menggunakan nama dari bahasa manapun).
   - Ekspor ke PDF. Header kolom dalam PDF disesuaikan dengan bahasa yang dipilih saat ekspor.

### Alur 4: Pelacakan Aset via Google Maps
**Langkah-langkah:**
1. **Peta Lokasi Aset**
   - Menu "Maps" → "Asset Locations".
   - Google Maps menampilkan marker untuk setiap lokasi dengan koordinat dari tabel `locations`.
   - Klik marker menampilkan detail lokasi dan jumlah aset di lokasi tersebut.
2. **Riwayat Pergerakan Aset**
   - Pilih aset tertentu → "Movement History".
   - Google Maps menampilkan jalur pergerakan aset berdasarkan koordinat scan dari `scan_logs`.

---

## 3. ALUR STAFF (Role: 'Staff')

### Dashboard Staff (Contoh Bahasa Indonesia)
- Header: "Halo, Siti Aminah".
- Tombol besar "Scan Aset".
- Daftar: "Jadwal Perawatan Mendatang: Perawatan Rutin Laptop (3 hari lagi)".

### Alur 1: Pemeriksaan Fisik Aset (Stock Opname)
**Cerita:** Bu Siti melakukan stock opname.

**Langkah-langkah:**
1. **Scan Aset (Data Matrix)**
   - Tap tombol "Scan Aset".
   - Arahkan kamera ke stiker Data Matrix di laptop Budi.
   - Aplikasi membaca `data_matrix_value`: "LPT-2025-001".
   - Aplikasi query ke `assets`, `categories`, `category_translations`, `locations`, dan `location_translations` (bergantung bahasa login Bu Siti) dan menemukan data "Laptop Dell Latitude 5430".
   - Halaman detail aset ditampilkan dengan nama kategori dan lokasi sesuai bahasa Bu Siti.
   - **Data scan masuk ke `scan_logs` beserta koordinat GPS saat scan** untuk pelacakan.

### Alur 2: Mengelola Perpindahan dan Penugasan Aset
**Langkah-langkah:**
1. **Pindahkan dan Tugaskan Aset**
   - Bu Siti scan laptop LPT-2025-001.
   - Dari halaman detail, tap "Pindahkan / Tugaskan".
   - Form perubahan diisi:
     - `to_location_id`: Pilih "Marketing Room 3rd Floor" (jika Bu Siti login dengan bahasa Inggris).
     - `to_user_id`: Pilih "Budi Santoso".
     - **Terjemahan Notes**:
       - `notes` (ID): "Penyerahan untuk karyawan baru Budi Santoso"
       - `notes` (EN): "Handover to new employee Budi Santoso"
       - `notes` (JA): "新入社員ブディ・サントソへの引き継ぎ"
   - Simpan. Data masuk ke `asset_movements` dan `asset_movement_translations`.

### Alur 3: Mengelola Perawatan Aset
**Langkah-langkah:**
1. **Membuat Record Perawatan**
   - Bu Siti scan aset yang akan dirawat.
   - Tap "Add Maintenance Record".
   - Form input diisi:
     - **Data Utama**:
       - `schedule_id`: (Jika ada jadwal terkait)
       - `maintenance_date`: "2025-08-08"
       - `performed_by_user`: "Siti Aminah"
       - `actual_cost`: 500000
     - **Terjemahan**:
       - `title` (ID): "Pembersihan dan Pengecekan Hardware"
       - `notes` (ID): "Membersihkan keyboard, layar, dan mengecek kondisi baterai."
       - `title` (EN): "Hardware Cleaning and Check"
       - `notes` (EN): "Cleaned keyboard, screen, and checked battery condition."
       - `title` (JA): "ハードウェアクリーニングとチェック"
       - `notes` (JA): "キーボード、スクリーンをクリーニングし、バッテリー状態をチェックしました。"
   - Simpan. Data masuk ke `maintenance_records` dan `maintenance_record_translations`.

---

## 4. ALUR EMPLOYEE (Role: 'Employee')

### Dashboard Employee
- Header: "Halo, Budi Santoso".
- Daftar "Aset Saya".

### Alur 1: Melihat Aset yang Digunakan
**Langkah-langkah:**
1. Budi login (misal dengan bahasa Indonesia).
2. Dashboard menampilkan: "Laptop Dell Latitude 5430" yang berada di lokasi "Ruang Marketing Lantai 3".
3. Budi tap dan melihat detail aset, dengan nama kategori "Laptop".
4. **Fitur Google Maps**: Budi dapat melihat lokasi aset di peta dengan koordinat yang tersimpan.

### Alur 2: Melaporkan Masalah Aset
**Langkah-langkah:**
1. **Buat Laporan Masalah**
   - Budi tap aset bermasalah di dashboard.
   - Tap "Laporkan Masalah".
   - Form input diisi:
     - **Data Utama**:
       - `issue_type`: "Hardware"
       - `priority`: "Medium"
       - `reported_date`: "2025-08-08"
     - **Terjemahan**:
       - `title` (ID): "Keyboard Tidak Responsif"
       - `description` (ID): "Beberapa tombol keyboard tidak berfungsi dengan baik."
       - `title` (EN): "Keyboard Not Responsive"
       - `description` (EN): "Several keyboard keys are not working properly."
       - `title` (JA): "キーボードが反応しない"
       - `description` (JA): "いくつかのキーボードキーが正常に動作しません。"
   - Simpan. Data masuk ke `issue_reports` dan `issue_report_translations`.

---

## 5. FITUR BERSAMA (ADMIN, STAFF, EMPLOYEE)

### Pencarian Aset
- User bisa mengetik "Latitude" (EN), "Dell" (brand), atau "ラップトップ" (JA) di search bar. Aplikasi akan mencari di `assets.asset_name`, `assets.brand`, `assets.model` DAN di `category_translations.category_name`, `location_translations.location_name` untuk menampilkan hasil yang relevan.

### Notifikasi Multibahasa
- Sistem akan mengirim notifikasi dalam bahasa sesuai `preferred_lang` user di tabel `users`.
- Title dan message notifikasi disimpan dalam `notification_translations`.

### Google Maps Integration
- **Peta Real-time**: Menampilkan semua lokasi aset dengan marker di Google Maps.
- **Riwayat Pergerakan**: Tracking pergerakan aset berdasarkan scan history dengan koordinat GPS.
- **Geofencing**: Notifikasi jika aset keluar dari area yang ditentukan (opsional).

---

## 6. DATABASE YANG DIUSULKAN (LENGKAP DENGAN LOKALISASI, DATA MATRIX & GOOGLE MAPS)

```sql
-- Tipe data kustom
CREATE TYPE user_role AS ENUM ('Admin', 'Staff', 'Employee');
CREATE TYPE asset_status AS ENUM ('Active', 'Maintenance', 'Disposed', 'Lost');
CREATE TYPE asset_condition AS ENUM ('Good', 'Fair', 'Poor', 'Damaged');
CREATE TYPE maintenance_schedule_type AS ENUM ('Preventive', 'Corrective');
CREATE TYPE schedule_status AS ENUM ('Scheduled', 'Completed', 'Cancelled');
CREATE TYPE scan_method_type AS ENUM ('DATA_MATRIX', 'MANUAL_INPUT');
CREATE TYPE scan_result_type AS ENUM ('Success', 'Invalid ID', 'Asset Not Found');
CREATE TYPE notification_type AS ENUM ('MAINTENANCE', 'WARRANTY', 'STATUS_CHANGE', 'MOVEMENT', 'ISSUE_REPORT');
CREATE TYPE issue_priority AS ENUM ('Low', 'Medium', 'High', 'Critical');
CREATE TYPE issue_status AS ENUM ('Open', 'In Progress', 'Resolved', 'Closed');

-- 1. Tabel Users
CREATE TABLE users (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role user_role NOT NULL,
    employee_id VARCHAR(20) UNIQUE NULL,
    preferred_lang VARCHAR(5) DEFAULT 'id-ID',
    is_active BOOLEAN DEFAULT TRUE,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_role_active ON users(role, is_active);

--------------------------------------------------

-- 2. Tabel Categories (Tanpa field yang bisa diterjemahkan)
CREATE TABLE categories (
    id VARCHAR(26) PRIMARY KEY,
    parent_id VARCHAR(26) NULL,
    category_code VARCHAR(20) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_parent_category FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
);
CREATE INDEX idx_categories_parent_id ON categories(parent_id);

-- 2a. Tabel Terjemahan untuk Categories
CREATE TABLE category_translations (
    id VARCHAR(26) PRIMARY KEY,
    category_id VARCHAR(26) NOT NULL,
    lang_code VARCHAR(5) NOT NULL,
    category_name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    UNIQUE (category_id, lang_code),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
CREATE INDEX idx_category_translations_category_lang ON category_translations(category_id, lang_code);

--------------------------------------------------

-- 3. Tabel Locations (Tanpa field yang bisa diterjemahkan, DENGAN KOORDINAT GPS)
CREATE TABLE locations (
    id VARCHAR(26) PRIMARY KEY,
    location_code VARCHAR(20) UNIQUE NOT NULL,
    building VARCHAR(100) NULL,
    floor VARCHAR(20) NULL,
    latitude DECIMAL(11,8) NULL,
    longitude DECIMAL(11,8) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_locations_building_floor ON locations(building, floor);
CREATE INDEX idx_locations_coordinates ON locations(latitude, longitude);

-- 3a. Tabel Terjemahan untuk Locations
CREATE TABLE location_translations (
    id VARCHAR(26) PRIMARY KEY,
    location_id VARCHAR(26) NOT NULL,
    lang_code VARCHAR(5) NOT NULL,
    location_name VARCHAR(100) NOT NULL,
    UNIQUE (location_id, lang_code),
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);
CREATE INDEX idx_location_translations_location_lang ON location_translations(location_id, lang_code);

--------------------------------------------------

-- 4. Tabel Assets (Field yang tidak perlu diterjemahkan tetap ada, DENGAN DATA MATRIX)
CREATE TABLE assets (
    id VARCHAR(26) PRIMARY KEY,
    asset_tag VARCHAR(50) UNIQUE NOT NULL,
    data_matrix_value VARCHAR(255) UNIQUE NOT NULL,
    asset_name VARCHAR(200) NOT NULL, -- Nama teknis, tidak perlu diterjemahkan
    category_id VARCHAR(26) NOT NULL,
    brand VARCHAR(100) NULL,
    model VARCHAR(100) NULL,
    serial_number VARCHAR(100) UNIQUE NULL,
    purchase_date DATE NULL,
    purchase_price DECIMAL(15,2) NULL,
    vendor_name VARCHAR(150) NULL,
    warranty_end DATE NULL,
    status asset_status DEFAULT 'Active',
    condition_status asset_condition DEFAULT 'Good',
    location_id VARCHAR(26) NULL,
    assigned_to VARCHAR(26) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
    FOREIGN KEY (assigned_to) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_assets_location ON assets(location_id);
CREATE INDEX idx_assets_assigned_to ON assets(assigned_to);
CREATE INDEX idx_assets_category_id ON assets(category_id);
CREATE INDEX idx_assets_warranty_end ON assets(warranty_end);
CREATE INDEX idx_assets_name_brand_model ON assets(asset_name, brand, model);
CREATE INDEX idx_assets_data_matrix ON assets(data_matrix_value);

--------------------------------------------------

-- 5. Tabel Asset_Movements (Tanpa field yang bisa diterjemahkan)
CREATE TABLE asset_movements (
    id VARCHAR(26) PRIMARY KEY,
    asset_id VARCHAR(26) NOT NULL,
    from_location_id VARCHAR(26) NULL,
    to_location_id VARCHAR(26) NULL,
    from_user_id VARCHAR(26) NULL,
    to_user_id VARCHAR(26) NULL,
    movement_date TIMESTAMP WITH TIME ZONE NOT NULL,
    moved_by VARCHAR(26) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    FOREIGN KEY (from_location_id) REFERENCES locations(id) ON DELETE SET NULL,
    FOREIGN KEY (to_location_id) REFERENCES locations(id) ON DELETE SET NULL,
    FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (moved_by) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX idx_movements_asset_id ON asset_movements(asset_id);
CREATE INDEX idx_movements_movement_date ON asset_movements(movement_date);
CREATE INDEX idx_movements_to_location_user ON asset_movements(to_location_id, to_user_id);

-- 5a. Tabel Terjemahan untuk Asset_Movements
CREATE TABLE asset_movement_translations (
    id VARCHAR(26) PRIMARY KEY,
    movement_id VARCHAR(26) NOT NULL,
    lang_code VARCHAR(5) NOT NULL,
    notes TEXT NULL,
    UNIQUE (movement_id, lang_code),
    FOREIGN KEY (movement_id) REFERENCES asset_movements(id) ON DELETE CASCADE
);
CREATE INDEX idx_movement_translations_movement_lang ON asset_movement_translations(movement_id, lang_code);

--------------------------------------------------

-- 6. Tabel Maintenance_Schedules (Tanpa field yang bisa diterjemahkan)
CREATE TABLE maintenance_schedules (
    id VARCHAR(26) PRIMARY KEY,
    asset_id VARCHAR(26) NOT NULL,
    maintenance_type maintenance_schedule_type NOT NULL,
    scheduled_date DATE NOT NULL,
    frequency_months INT NULL,
    status schedule_status DEFAULT 'Scheduled',
    created_by VARCHAR(26) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX idx_maintenance_schedules_asset_id ON maintenance_schedules(asset_id);
CREATE INDEX idx_maintenance_schedules_status ON maintenance_schedules(status);
CREATE INDEX idx_maintenance_schedules_scheduled_date ON maintenance_schedules(scheduled_date);

-- 6a. Tabel Terjemahan untuk Maintenance_Schedules
CREATE TABLE maintenance_schedule_translations (
    id VARCHAR(26) PRIMARY KEY,
    schedule_id VARCHAR(26) NOT NULL,
    lang_code VARCHAR(5) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NULL,
    UNIQUE (schedule_id, lang_code),
    FOREIGN KEY (schedule_id) REFERENCES maintenance_schedules(id) ON DELETE CASCADE
);
CREATE INDEX idx_schedule_translations_schedule_lang ON maintenance_schedule_translations(schedule_id, lang_code);

--------------------------------------------------

-- 7. Tabel Maintenance_Records (Tanpa field yang bisa diterjemahkan)
CREATE TABLE maintenance_records (
    id VARCHAR(26) PRIMARY KEY,
    schedule_id VARCHAR(26) NULL,
    asset_id VARCHAR(26) NOT NULL,
    maintenance_date DATE NOT NULL,
    performed_by_user VARCHAR(26) NULL,
    performed_by_vendor VARCHAR(150) NULL,
    actual_cost DECIMAL(12,2) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (schedule_id) REFERENCES maintenance_schedules(id) ON DELETE SET NULL,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    FOREIGN KEY (performed_by_user) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_maintenance_records_asset_id ON maintenance_records(asset_id);
CREATE INDEX idx_maintenance_records_date ON maintenance_records(maintenance_date);

-- 7a. Tabel Terjemahan untuk Maintenance_Records
CREATE TABLE maintenance_record_translations (
    id VARCHAR(26) PRIMARY KEY,
    record_id VARCHAR(26) NOT NULL,
    lang_code VARCHAR(5) NOT NULL,
    title VARCHAR(200) NOT NULL,
    notes TEXT NULL,
    UNIQUE (record_id, lang_code),
    FOREIGN KEY (record_id) REFERENCES maintenance_records(id) ON DELETE CASCADE
);
CREATE INDEX idx_record_translations_record_lang ON maintenance_record_translations(record_id, lang_code);

--------------------------------------------------

-- 8. Tabel Issue_Reports (Tanpa field yang bisa diterjemahkan)
CREATE TABLE issue_reports (
    id VARCHAR(26) PRIMARY KEY,
    asset_id VARCHAR(26) NOT NULL,
    reported_by VARCHAR(26) NOT NULL,
    reported_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    issue_type VARCHAR(50) NOT NULL,
    priority issue_priority DEFAULT 'Medium',
    status issue_status DEFAULT 'Open',
    resolved_date TIMESTAMP WITH TIME ZONE NULL,
    resolved_by VARCHAR(26) NULL,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    FOREIGN KEY (reported_by) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (resolved_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_issue_reports_asset_id ON issue_reports(asset_id);
CREATE INDEX idx_issue_reports_status ON issue_reports(status);
CREATE INDEX idx_issue_reports_priority ON issue_reports(priority);

-- 8a. Tabel Terjemahan untuk Issue_Reports
CREATE TABLE issue_report_translations (
    id VARCHAR(26) PRIMARY KEY,
    report_id VARCHAR(26) NOT NULL,
    lang_code VARCHAR(5) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NULL,
    resolution_notes TEXT NULL,
    UNIQUE (report_id, lang_code),
    FOREIGN KEY (report_id) REFERENCES issue_reports(id) ON DELETE CASCADE
);
CREATE INDEX idx_report_translations_report_lang ON issue_report_translations(report_id, lang_code);

--------------------------------------------------

-- 9. Tabel Scan_Logs (Tidak perlu diterjemahkan, DENGAN KOORDINAT GPS UNTUK PELACAKAN)
CREATE TABLE scan_logs (
    id VARCHAR(26) PRIMARY KEY,
    asset_id VARCHAR(26) NULL,
    scanned_value VARCHAR(255) NOT NULL,
    scan_method scan_method_type NOT NULL,
    scanned_by VARCHAR(26) NOT NULL,
    scan_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    scan_location_lat DECIMAL(11,8) NULL,
    scan_location_lng DECIMAL(11,8) NULL,
    scan_result scan_result_type NOT NULL,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE SET NULL,
    FOREIGN KEY (scanned_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_scan_logs_scan_timestamp ON scan_logs(scan_timestamp);
CREATE INDEX idx_scan_logs_scanned_by ON scan_logs(scanned_by);
CREATE INDEX idx_scan_logs_result ON scan_logs(scan_result);
CREATE INDEX idx_scan_logs_location ON scan_logs(scan_location_lat, scan_location_lng);

--------------------------------------------------

-- 10. Tabel Notifications (Tanpa field yang bisa diterjemahkan)
CREATE TABLE notifications (
    id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
    related_asset_id VARCHAR(26) NULL,
    type notification_type NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (related_asset_id) REFERENCES assets(id) ON DELETE SET NULL
);

CREATE INDEX idx_notifications_user_id_is_read ON notifications(user_id, is_read);
CREATE INDEX idx_notifications_type ON notifications(type);

-- 10a. Tabel Terjemahan untuk Notifications
CREATE TABLE notification_translations (
    id VARCHAR(26) PRIMARY KEY,
    notification_id VARCHAR(26) NOT NULL,
    lang_code VARCHAR(5) NOT NULL,
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    UNIQUE (notification_id, lang_code),
    FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE
);
CREATE INDEX idx_notification_translations_notification_lang ON notification_translations(notification_id, lang_code);
```

## 7. CONTOH QUERY UNTUK MULTIBAHASA DAN GOOGLE MAPS

### Query untuk mendapatkan lokasi dengan koordinat untuk Google Maps
```sql
-- Mendapatkan semua lokasi dengan koordinat untuk Google Maps
SELECT
    l.id,
    l.location_code,
    l.building,
    l.floor,
    l.latitude,
    l.longitude,
    COALESCE(lt_user.location_name, lt_default.location_name) as location_name,
    COUNT(a.id) as asset_count
FROM locations l
LEFT JOIN location_translations lt_user ON l.id = lt_user.location_id AND lt_user.lang_code = 'en-US'
LEFT JOIN location_translations lt_default ON l.id = lt_default.location_id AND lt_default.lang_code = 'id-ID'
LEFT JOIN assets a ON l.id = a.location_id AND a.status = 'Active'
WHERE l.latitude IS NOT NULL AND l.longitude IS NOT NULL
GROUP BY l.id, l.location_code, l.building, l.floor, l.latitude, l.longitude, lt_user.location_name, lt_default.location_name;
```

### Query untuk tracking pergerakan aset berdasarkan scan history
```sql
-- Mendapatkan riwayat pergerakan aset dengan koordinat GPS dari scan logs
SELECT
    sl.asset_id,
    a.asset_tag,
    a.asset_name,
    sl.scan_timestamp,
    sl.scan_location_lat,
    sl.scan_location_lng,
    COALESCE(lt_user.location_name, lt_default.location_name) as current_location,
    u.full_name as scanned_by
FROM scan_logs sl
JOIN assets a ON sl.asset_id = a.id
JOIN users u ON sl.scanned_by = u.id
LEFT JOIN locations l ON a.location_id = l.id
LEFT JOIN location_translations lt_user ON l.id = lt_user.location_id AND lt_user.lang_code = 'en-US'
LEFT JOIN location_translations lt_default ON l.id = lt_default.location_id AND lt_default.lang_code = 'id-ID'
WHERE sl.asset_id = 'ASSET_ID_HERE'
  AND sl.scan_location_lat IS NOT NULL
  AND sl.scan_location_lng IS NOT NULL
  AND sl.scan_result = 'Success'
ORDER BY sl.scan_timestamp DESC;
```

### Query untuk mendapatkan kategori dengan terjemahan
```sql
-- Mendapatkan kategori dalam bahasa yang dipilih user (fallback ke bahasa default)
SELECT
    c.id,
    c.category_code,
    COALESCE(ct_user.category_name, ct_default.category_name) as category_name,
    COALESCE(ct_user.description, ct_default.description) as description
FROM categories c
LEFT JOIN category_translations ct_user ON c.id = ct_user.category_id AND ct_user.lang_code = 'en-US'
LEFT JOIN category_translations ct_default ON c.id = ct_default.category_id AND ct_default.lang_code = 'id-ID'
WHERE c.parent_id IS NULL;
```

### Query untuk pencarian aset multibahasa
```sql
-- Pencarian aset dengan nama kategori dan lokasi dalam bahasa user
SELECT DISTINCT
    a.id,
    a.asset_tag,
    a.asset_name,
    a.data_matrix_value,
    COALESCE(ct_user.category_name, ct_default.category_name) as category_name,
    COALESCE(lt_user.location_name, lt_default.location_name) as location_name,
    l.latitude,
    l.longitude
FROM assets a
JOIN categories c ON a.category_id = c.id
LEFT JOIN category_translations ct_user ON c.id = ct_user.category_id AND ct_user.lang_code = 'en-US'
LEFT JOIN category_translations ct_default ON c.id = ct_default.category_id AND ct_default.lang_code = 'id-ID'
LEFT JOIN locations l ON a.location_id = l.id
LEFT JOIN location_translations lt_user ON l.id = lt_user.location_id AND lt_user.lang_code = 'en-US'
LEFT JOIN location_translations lt_default ON l.id = lt_default.location_id AND lt_default.lang_code = 'id-ID'
WHERE
    a.asset_name ILIKE '%laptop%'
    OR a.brand ILIKE '%laptop%'
    OR ct_user.category_name ILIKE '%laptop%'
    OR ct_default.category_name ILIKE '%laptop%'
    OR lt_user.location_name ILIKE '%laptop%'
    OR lt_default.location_name ILIKE '%laptop%';
```

### Query untuk dashboard analytics dengan geolocation
```sql
-- Dashboard dengan statistik aset per lokasi dan koordinat untuk heatmap
SELECT
    COALESCE(lt_user.location_name, lt_default.location_name) as location_name,
    l.latitude,
    l.longitude,
    COUNT(a.id) as total_assets,
    COUNT(CASE WHEN a.status = 'Active' THEN 1 END) as active_assets,
    COUNT(CASE WHEN a.status = 'Maintenance' THEN 1 END) as maintenance_assets,
    COUNT(CASE WHEN a.warranty_end <= CURRENT_DATE + INTERVAL '30 days' THEN 1 END) as warranty_expiring_soon
FROM locations l
LEFT JOIN location_translations lt_user ON l.id = lt_user.location_id AND lt_user.lang_code = 'en-US'
LEFT JOIN location_translations lt_default ON l.id = lt_default.location_id AND lt_default.lang_code = 'id-ID'
LEFT JOIN assets a ON l.id = a.location_id
WHERE l.latitude IS NOT NULL AND l.longitude IS NOT NULL
GROUP BY l.id, l.latitude, l.longitude, lt_user.location_name, lt_default.location_name
ORDER BY total_assets DESC;
```
