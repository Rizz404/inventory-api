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

3. **Kelola Maintenance Schedules (Dengan Terjemahan & Interval Fleksibel)**
   - Tap menu "Maintenance" → "Schedules".
   - Tap tombol "+" untuk menambah jadwal baru.
   - Form input diisi:
     - **Data Utama**:
       - `asset_id`: Pilih asset
       - `maintenance_type`: "Preventive" (pilihan: Preventive, Corrective, Inspection, Calibration)
       - `is_recurring`: TRUE (jika jadwal berulang)
       - **Interval Fleksibel** (jika recurring):
         - `interval_value`: 15 (contoh: setiap 15 unit)
         - `interval_unit`: "Days" (pilihan: Minutes, Hours, Days, Weeks, Months, Years)
       - `scheduled_time`: "09:00:00" (waktu spesifik, null jika kapan saja)
       - `next_scheduled_date`: "2025-09-15 09:00:00"
       - `state`: "Active" (pilihan: Active, Paused, Stopped, Completed)
       - `auto_complete`: FALSE (TRUE jika otomatis completed setelah 1x maintenance)
       - `estimated_cost`: 500000
     - **Terjemahan**:
       - `title` (ID): "Perawatan Rutin Laptop"
       - `description` (ID): "Pembersihan dan update sistem operasi setiap 15 hari."
       - `title` (EN): "Routine Laptop Maintenance"
       - `description` (EN): "Cleaning and operating system update every 15 days."
       - `title` (JA): "ラップトップの定期メンテナンス"
       - `description` (JA): "15日ごとのクリーニングとオペレーティングシステムの更新。"
   - Simpan. Data masuk ke `maintenance_schedules` dan `maintenance_schedule_translations`.

4. **Kelola User**
   - Tap menu "Master Data" → "Users".
   - Tap tombol "+" untuk menambah user baru.
   - Form input diisi:
     - `name`: "budi.staff"
     - `email`: "budi@fujiyama.co.id"
     - `password`: (akan di-hash)
     - `full_name`: "Budi Santoso"
     - `role`: "Staff"
     - `employee_id`: "EMP-2025-001"
     - `preferred_lang`: "id-ID"
     - `is_active`: TRUE
     - `fcm_token`: (untuk push notification)
   - Simpan. Data masuk ke tabel `users`.

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
     - `status`: "Active" (pilihan: Active, Maintenance, Disposed, Lost)
     - `condition_status`: "Good" (pilihan: Good, Fair, Poor, Damaged)
     - `location_id`: Pilih "Gudang IT" (nama lokasi ditampilkan sesuai bahasa).
     - `assigned_to`: NULL (belum ditugaskan)
   - Simpan.

2. **Menyiapkan Identifier Fisik (Data Matrix)**
   - Sistem generate Data Matrix code dengan value "LPT-2025-001".
   - Data Matrix image disimpan dan URL-nya diupdate ke field `data_matrix_image_url`.
   - Data Matrix dapat menyimpan lebih banyak data dibanding QR Code dan lebih tahan terhadap kerusakan fisik.
   - Cetak stiker Data Matrix dan tempel di aset.

### Alur 3: Monitoring dan Reporting
**Langkah-langkah:**
1. **Dashboard Analytics (Contoh Bahasa Jepang)**
   - Chart menunjukkan: "カテゴリ ラップトップ: 50 ユニット", "カテゴリ モニター: 45 ユニット".
   - Menampilkan aset yang warranty-nya akan expired dalam 30 hari.
   - Menampilkan maintenance schedules yang akan datang.

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
   - Timeline menampilkan kapan dan dimana aset pernah di-scan.

### Alur 5: Mengelola Notifikasi
**Langkah-langkah:**
1. **Membuat Notifikasi Manual**
   - Tap menu "Notifications" → "Create".
   - Form input diisi:
     - `user_id`: Pilih user penerima
     - `type`: "MAINTENANCE" (pilihan: MAINTENANCE, WARRANTY, ISSUE, MOVEMENT, STATUS_CHANGE, LOCATION_CHANGE, CATEGORY_CHANGE)
     - `priority`: "NORMAL" (pilihan: LOW, NORMAL, HIGH, URGENT)
     - `related_entity_type`: "asset"
     - `related_entity_id`: ID aset terkait
     - `expires_at`: "2025-12-31" (jika ada expiration)
     - **Terjemahan**:
       - `title` (ID): "Perawatan Terjadwal"
       - `message` (ID): "Laptop Anda akan menjalani perawatan rutin besok."
       - `title` (EN): "Scheduled Maintenance"
       - `message` (EN): "Your laptop will undergo routine maintenance tomorrow."
       - `title` (JA): "定期メンテナンス"
       - `message` (JA): "お使いのラップトップは明日定期メンテナンスを受けます。"
   - Simpan. Data masuk ke `notifications` dan `notification_translations`.

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
   - Arahkan kamera ke stiker Data Matrix di laptop.
   - Aplikasi membaca `scanned_value`: "LPT-2025-001".
   - Aplikasi query ke `assets` menggunakan `asset_tag`.
   - **Data scan masuk ke `scan_logs`**:
     - `asset_id`: ID aset yang berhasil ditemukan
     - `scanned_value`: "LPT-2025-001"
     - `scan_method`: "DATA_MATRIX"
     - `scanned_by`: ID Bu Siti
     - `scan_timestamp`: timestamp saat scan
     - `scan_location_lat`: -6.3734 (dari GPS device)
     - `scan_location_lng`: 106.8284 (dari GPS device)
     - `scan_result`: "Success" (pilihan: Success, Invalid ID, Asset Not Found)
   - Halaman detail aset ditampilkan dengan nama kategori dan lokasi sesuai bahasa Bu Siti.

2. **Scan dengan Input Manual**
   - Jika kamera tidak tersedia atau Data Matrix rusak.
   - Bu Siti input manual: "LPT-2025-001".
   - Data masuk ke `scan_logs` dengan `scan_method`: "MANUAL_INPUT".

### Alur 2: Mengelola Perpindahan dan Penugasan Aset
**Langkah-langkah:**
1. **Pindahkan dan Tugaskan Aset**
   - Bu Siti scan laptop LPT-2025-001.
   - Dari halaman detail, tap "Pindahkan / Tugaskan".
   - Form perubahan diisi:
     - `from_location_id`: Lokasi saat ini (auto-filled)
     - `to_location_id`: Pilih "Marketing Room 3rd Floor"
     - `from_user_id`: User yang saat ini memegang (auto-filled jika ada)
     - `to_user_id`: Pilih "Budi Santoso"
     - `movement_date`: "2025-08-08 10:00:00"
     - `moved_by`: ID Bu Siti (auto-filled)
     - **Terjemahan Notes**:
       - `notes` (ID): "Penyerahan untuk karyawan baru Budi Santoso"
       - `notes` (EN): "Handover to new employee Budi Santoso"
       - `notes` (JA): "新入社員ブディ・サントソへの引き継ぎ"
   - Simpan. Data masuk ke `asset_movements` dan `asset_movement_translations`.
   - Sistem otomatis update `location_id` dan `assigned_to` di tabel `assets`.

### Alur 3: Mengelola Perawatan Aset
**Langkah-langkah:**
1. **Membuat Record Perawatan**
   - Bu Siti scan aset yang akan dirawat.
   - Tap "Add Maintenance Record".
   - Form input diisi:
     - **Data Utama**:
       - `schedule_id`: Pilih jadwal terkait (jika ada)
       - `asset_id`: Auto-filled dari aset yang di-scan
       - `maintenance_date`: "2025-08-08 09:00:00"
       - `completion_date`: "2025-08-08 11:30:00"
       - `duration_minutes`: 150
       - `performed_by_user`: "Siti Aminah"
       - `performed_by_vendor`: NULL (jika dilakukan internal)
       - `result`: "Success" (pilihan: Success, Partial, Failed, Rescheduled)
       - `actual_cost`: 500000
     - **Terjemahan**:
       - `title` (ID): "Pembersihan dan Pengecekan Hardware"
       - `notes` (ID): "Membersihkan keyboard, layar, dan mengecek kondisi baterai."
       - `title` (EN): "Hardware Cleaning and Check"
       - `notes` (EN): "Cleaned keyboard, screen, and checked battery condition."
       - `title` (JA): "ハードウェアクリーニングとチェック"
       - `notes` (JA): "キーボード、スクリーンをクリーニングし、バッテリー状態をチェックしました。"
   - Simpan. Data masuk ke `maintenance_records` dan `maintenance_record_translations`.
   - Jika ada `schedule_id` dan schedule memiliki `auto_complete` = TRUE, sistem update schedule state menjadi 'Completed'.
   - Jika schedule `is_recurring` = TRUE, sistem hitung `next_scheduled_date` berikutnya berdasarkan `interval_value` dan `interval_unit`.

---

## 4. ALUR EMPLOYEE (Role: 'Employee')

### Dashboard Employee
- Header: "Halo, Budi Santoso".
- Daftar "Aset Saya".
- Notifikasi yang belum dibaca.

### Alur 1: Melihat Aset yang Digunakan
**Langkah-langkah:**
1. Budi login (misal dengan bahasa Indonesia).
2. Dashboard menampilkan: "Laptop Dell Latitude 5430" yang berada di lokasi "Ruang Marketing Lantai 3".
3. Budi tap dan melihat detail aset, dengan nama kategori "Laptop".
4. **Fitur Google Maps**: Budi dapat melihat lokasi aset di peta dengan koordinat yang tersimpan di tabel `locations`.
5. Budi dapat melihat riwayat scan aset dari `scan_logs`.

### Alur 2: Melaporkan Masalah Aset
**Langkah-langkah:**
1. **Buat Laporan Masalah**
   - Budi tap aset bermasalah di dashboard.
   - Tap "Laporkan Masalah".
   - Form input diisi:
     - **Data Utama**:
       - `asset_id`: Auto-filled
       - `reported_by`: ID Budi (auto-filled)
       - `reported_date`: Auto-filled (CURRENT_TIMESTAMP)
       - `issue_type`: "Hardware"
       - `priority`: "Medium" (pilihan: Low, Medium, High, Critical)
       - `status`: "Open" (default, pilihan: Open, In Progress, Resolved, Closed)
     - **Terjemahan**:
       - `title` (ID): "Keyboard Tidak Responsif"
       - `description` (ID): "Beberapa tombol keyboard tidak berfungsi dengan baik."
       - `title` (EN): "Keyboard Not Responsive"
       - `description` (EN): "Several keyboard keys are not working properly."
       - `title` (JA): "キーボードが反応しない"
       - `description` (JA): "いくつかのキーボードキーが正常に動作しません。"
   - Simpan. Data masuk ke `issue_reports` dan `issue_report_translations`.
   - Sistem otomatis membuat notifikasi untuk Admin dan Staff terkait.

2. **Tracking Status Laporan**
   - Budi dapat melihat status issue report: Open → In Progress → Resolved → Closed.
   - Ketika Staff/Admin menyelesaikan masalah:
     - Update `status` = "Resolved"
     - Isi `resolved_date` dan `resolved_by`
     - Tambahkan `resolution_notes` di `issue_report_translations`
   - Budi menerima notifikasi bahwa issue telah resolved.

### Alur 3: Melihat dan Menandai Notifikasi
**Langkah-langkah:**
1. Budi tap ikon notifikasi di dashboard.
2. Menampilkan daftar notifikasi sesuai `preferred_lang` Budi.
3. Notifikasi yang belum dibaca (`is_read` = FALSE) ditampilkan dengan highlight.
4. Budi tap salah satu notifikasi:
   - Sistem update `is_read` = TRUE dan isi `read_at` dengan timestamp.
   - Menampilkan detail notifikasi dengan title dan message dalam bahasa Budi.
5. Notifikasi yang sudah expired (`expires_at` < NOW) dapat disembunyikan atau ditampilkan dengan style berbeda.

---

## 5. FITUR BERSAMA (ADMIN, STAFF, EMPLOYEE)

### Pencarian Aset Multibahasa
- User bisa mengetik "Latitude" (EN), "Dell" (brand), "ラップトップ" (JA), atau "Ruang Marketing" di search bar.
- Aplikasi akan mencari di:
  - `assets.asset_name`, `assets.brand`, `assets.model`, `assets.serial_number`
  - `category_translations.category_name`, `category_translations.description`
  - `location_translations.location_name`
- Hasil pencarian menampilkan informasi dalam bahasa sesuai `preferred_lang` user.

### Notifikasi Multibahasa & Push Notification
- Sistem akan mengirim notifikasi dalam bahasa sesuai `preferred_lang` user di tabel `users`.
- Title dan message notifikasi disimpan dalam `notification_translations` untuk semua bahasa.
- Push notification menggunakan `fcm_token` dari tabel `users`.
- Notifikasi dapat di-filter berdasarkan `type` dan `priority`.
- Notifikasi expired otomatis tidak ditampilkan (berdasarkan `expires_at`).

### Google Maps Integration
- **Peta Real-time**: Menampilkan semua lokasi aset dengan marker di Google Maps menggunakan koordinat dari `locations`.
- **Riwayat Pergerakan**: Tracking pergerakan aset berdasarkan:
  - Data di `asset_movements` (perpindahan formal)
  - Data di `scan_logs` dengan koordinat GPS (scan history)
- **Heatmap**: Visualisasi konsentrasi aset di berbagai lokasi.
- **Geofencing** (opsional): Notifikasi jika aset keluar dari area yang ditentukan.

---

## 6. DATABASE SCHEMA (SESUAI DENGAN MIGRATION FILES)

### 1. Tabel Users
```sql
CREATE TYPE user_role AS ENUM ('Admin', 'Staff', 'Employee');

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
  avatar_url VARCHAR(255) NULL,
  fcm_token TEXT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_role_active ON users(role, is_active);
```

### 2. Tabel Categories & Translations
```sql
CREATE TABLE categories (
  id VARCHAR(26) PRIMARY KEY,
  parent_id VARCHAR(26) NULL,
  category_code VARCHAR(20) UNIQUE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_parent_category FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);

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
```

### 3. Tabel Locations & Translations
```sql
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

CREATE TABLE location_translations (
  id VARCHAR(26) PRIMARY KEY,
  location_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  location_name VARCHAR(100) NOT NULL,
  UNIQUE (location_id, lang_code),
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);

CREATE INDEX idx_location_translations_location_lang ON location_translations(location_id, lang_code);
```

### 4. Tabel Assets
```sql
CREATE TYPE asset_status AS ENUM ('Active', 'Maintenance', 'Disposed', 'Lost');
CREATE TYPE asset_condition AS ENUM ('Good', 'Fair', 'Poor', 'Damaged');

CREATE TABLE assets (
  id VARCHAR(26) PRIMARY KEY,
  asset_tag VARCHAR(50) UNIQUE NOT NULL,
  data_matrix_image_url VARCHAR(255) NULL,
  asset_name VARCHAR(200) NOT NULL,
  category_id VARCHAR(26) NOT NULL,
  brand VARCHAR(100) NULL,
  model VARCHAR(100) NULL,
  serial_number VARCHAR(100) UNIQUE NULL,
  purchase_date DATE NULL,
  purchase_price DECIMAL(15, 2) NULL,
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
```

### 5. Tabel Asset Movements & Translations
```sql
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
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
  CONSTRAINT check_recurring_interval CHECK (
    (is_recurring = FALSE) OR
    (is_recurring = TRUE AND interval_value IS NOT NULL AND interval_unit IS NOT NULL)
  )
);

CREATE INDEX idx_maintenance_schedules_asset_id ON maintenance_schedules(asset_id);
CREATE INDEX idx_maintenance_schedules_state ON maintenance_schedules(state);
CREATE INDEX idx_maintenance_schedules_next_date ON maintenance_schedules(next_scheduled_date);

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
```

### 7. Tabel Maintenance Records & Translations
```sql
CREATE TYPE maintenance_result AS ENUM ('Success', 'Partial', 'Failed', 'Rescheduled');

CREATE TABLE maintenance_records (
  id VARCHAR(26) PRIMARY KEY,
  schedule_id VARCHAR(26) NULL,
  asset_id VARCHAR(26) NOT NULL,
  maintenance_date TIMESTAMP WITH TIME ZONE NOT NULL,
  completion_date TIMESTAMP WITH TIME ZONE NULL,
  duration_minutes INT NULL,
  performed_by_user VARCHAR(26) NULL,
  performed_by_vendor VARCHAR(150) NULL,
  result maintenance_result DEFAULT 'Success',
  actual_cost DECIMAL(12, 2) NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (schedule_id) REFERENCES maintenance_schedules(id) ON DELETE SET NULL,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
  FOREIGN KEY (performed_by_user) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_maintenance_records_schedule_id ON maintenance_records(schedule_id);
CREATE INDEX idx_maintenance_records_asset_id ON maintenance_records(asset_id);
CREATE INDEX idx_maintenance_records_date ON maintenance_records(maintenance_date);

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
```

### 8. Tabel Issue Reports & Translations
```sql
CREATE TYPE issue_priority AS ENUM ('Low', 'Medium', 'High', 'Critical');
CREATE TYPE issue_status AS ENUM ('Open', 'In Progress', 'Resolved', 'Closed');

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
```

### 9. Tabel Scan Logs
```sql
CREATE TYPE scan_method_type AS ENUM ('DATA_MATRIX', 'MANUAL_INPUT');
CREATE TYPE scan_result_type AS ENUM ('Success', 'Invalid ID', 'Asset Not Found');

CREATE TABLE scan_logs (
  id VARCHAR(26) PRIMARY KEY,
  asset_id VARCHAR(26) NULL,
  scanned_value VARCHAR(255) NOT NULL,
  scan_method scan_method_type NOT NULL,
  scanned_by VARCHAR(26) NOT NULL,
  scan_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  scan_location_lat DECIMAL(11, 8) NULL,
  scan_location_lng DECIMAL(11, 8) NULL,
  scan_result scan_result_type NOT NULL,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE SET NULL,
  FOREIGN KEY (scanned_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_scan_logs_scan_timestamp ON scan_logs(scan_timestamp);
CREATE INDEX idx_scan_logs_scanned_by ON scan_logs(scanned_by);
CREATE INDEX idx_scan_logs_result ON scan_logs(scan_result);
CREATE INDEX idx_scan_logs_location ON scan_logs(scan_location_lat, scan_location_lng);
```

### 10. Tabel Notifications & Translations
```sql
CREATE TYPE notification_type AS ENUM (
  'MAINTENANCE', 'WARRANTY', 'ISSUE', 'MOVEMENT',
  'STATUS_CHANGE', 'LOCATION_CHANGE', 'CATEGORY_CHANGE'
);
CREATE TYPE notification_priority AS ENUM ('LOW', 'NORMAL', 'HIGH', 'URGENT');

CREATE TABLE notifications (
  id VARCHAR(26) PRIMARY KEY,
  user_id VARCHAR(26) NOT NULL,
  related_entity_type VARCHAR(50) NULL,
  related_entity_id VARCHAR(26) NULL,
  related_asset_id VARCHAR(26) NULL,
  type notification_type NOT NULL,
  priority notification_priority DEFAULT 'NORMAL',
  is_read BOOLEAN DEFAULT FALSE,
  read_at TIMESTAMP WITH TIME ZONE NULL,
  expires_at TIMESTAMP WITH TIME ZONE NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (related_asset_id) REFERENCES assets(id) ON DELETE SET NULL
);

CREATE INDEX idx_notifications_user_id_is_read ON notifications(user_id, is_read);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_priority ON notifications(priority);
CREATE INDEX idx_notifications_related_entity ON notifications(related_entity_type, related_entity_id);
CREATE INDEX idx_notifications_expires_at ON notifications(expires_at);

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
