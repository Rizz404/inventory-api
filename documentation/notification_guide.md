### 2. **Issue Report Service** (issue_report_service.go)
   - **Issue Reported**: Notify assigned user (jika asset assigned) atau admin tentang laporan issue baru.
     - Title: "New Issue Reported"
     - Message: "A new issue has been reported for asset {assetName} ({assetTag})."
   - **Issue Updated**: Notify reporter atau assigned user tentang update issue (deskripsi, priority).
     - Title: "Issue Updated"
     - Message: "Issue report for asset {assetName} has been updated."
   - **Issue Resolved**: Notify reporter tentang penyelesaian issue.
     - Title: "Issue Resolved"
     - Message: "Issue report for asset {assetName} has been resolved. Resolution: {resolutionNotes}."
   - **Issue Reopened**: Notify assigned user atau admin tentang issue yang dibuka kembali.
     - Title: "Issue Reopened"
     - Message: "Issue report for asset {assetName} has been reopened."

### 3. **Maintenance Schedule Service** (maintenance_schedule_service.go)
   - **Maintenance Scheduled**: Notify assigned user tentang jadwal maintenance baru.
     - Title: "Maintenance Scheduled"
     - Message: "Maintenance for asset {assetName} is scheduled on {scheduledDate}."
   - **Maintenance Due Soon**: Reminder ke assigned user beberapa hari sebelum due date.
     - Title: "Maintenance Due Soon"
     - Message: "Maintenance for asset {assetName} is due on {scheduledDate}. Please prepare."
   - **Maintenance Overdue**: Notify assigned user jika maintenance terlewat.
     - Title: "Maintenance Overdue"
     - Message: "Maintenance for asset {assetName} is overdue. Scheduled date was {scheduledDate}."

### 4. **Maintenance Record Service** (maintenance_record_service.go)
   - **Maintenance Performed**: Notify assigned user tentang maintenance yang telah dilakukan.
     - Title: "Maintenance Completed"
     - Message: "Maintenance for asset {assetName} has been completed. Notes: {notes}."
   - **Maintenance Failed**: Notify admin atau assigned user jika maintenance gagal.
     - Title: "Maintenance Failed"
     - Message: "Maintenance for asset {assetName} could not be completed. Reason: {failureReason}."

### 5. **Asset Movement Service** (asset_movement_service.go)
   - **Asset Moved**: Notify assigned user tentang perpindahan asset
     - Title: "Asset Moved"
     - Message: "Asset {assetName} has been moved from {oldLocation} to {newLocation}."

### 7. **Location Service** (location_service.go)
   - **Location Created/Updated**: Notify semua admin tentang perubahan lokasi.
     - Title: "Location Updated"
     - Message: "Location {locationName} has been updated in the system."

### 8. **Category Service** (category_service.go)
   - **Category Created/Updated**: Notify semua admin tentang perubahan kategori.
     - Title: "Category Updated"
     - Message: "Category {categoryName} has been updated."

