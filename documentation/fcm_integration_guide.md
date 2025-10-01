# FCM Integration with Notification System

## Overview

Integrasi Firebase Cloud Messaging (FCM) dengan sistem notifikasi inventory API untuk mengirimkan push notification ke user device.

## Features

1. **Database Notification**: Semua notifikasi disimpan di database dengan support multi-language
2. **Push Notification**: Notifikasi dikirim via FCM ke device user secara real-time
3. **Asset Assignment Notification**: Otomatis mengirim notifikasi ketika asset di-assign ke user
4. **Non-blocking**: FCM notification dikirim secara asynchronous (tidak memblokir operasi utama)
5. **Graceful Degradation**: Jika FCM tidak tersedia, sistem tetap berfungsi (hanya menyimpan di database)

## Architecture

```
┌─────────────┐
│ Asset       │
│ Service     │
└──────┬──────┘
       │
       │ (when asset assigned)
       ▼
┌─────────────────┐     ┌──────────────┐
│ Notification    │────▶│  Database    │
│ Service         │     │  (notifs)    │
└────────┬────────┘     └──────────────┘
         │
         │ (async)
         ▼
    ┌────────┐
    │  FCM   │
    │ Client │
    └───┬────┘
        │
        ▼
   ┌─────────┐
   │  User   │
   │ Device  │
   └─────────┘
```

## Database Schema

### Users Table
```sql
ALTER TABLE users ADD COLUMN fcm_token TEXT;
```

FCM token disimpan di user table dan di-update via endpoint `/api/v1/users/profile/fcm-token`

## API Endpoints

### 1. Update FCM Token
Update FCM token untuk current user (authenticated)

**Endpoint**: `PATCH /api/v1/users/profile/fcm-token`

**Headers**:
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "fcmToken": "dXyZ1234567890abcdef..."
}
```

**Response**:
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": "01HTXYZ...",
    "name": "john_doe",
    "email": "john@example.com",
    "fullName": "John Doe",
    ...
  }
}
```

### 2. Create Asset with Assignment
Ketika asset dibuat dengan `assignedTo`, notifikasi otomatis terkirim

**Endpoint**: `POST /api/v1/assets`

**Request Body**:
```json
{
  "assetTag": "LAPTOP-001",
  "assetName": "MacBook Pro 2023",
  "categoryId": "01HTXYZ...",
  "assignedTo": "01HTUSER123...",
  ...
}
```

**Result**:
- Asset created di database
- Notification created di database
- Push notification dikirim ke user device (jika FCM token tersedia)

### 3. Update Asset Assignment
Ketika asset di-assign ke user berbeda

**Endpoint**: `PATCH /api/v1/assets/:id`

**Request Body**:
```json
{
  "assignedTo": "01HTUSER456..."
}
```

**Result**:
- Asset updated di database
- Notification created untuk user baru
- Push notification dikirim ke user device

## Notification Format

### Database (Multi-language)
```json
{
  "id": "01HTNOTIF...",
  "userId": "01HTUSER...",
  "relatedAssetId": "01HTASSET...",
  "type": "STATUS_CHANGE",
  "isRead": false,
  "createdAt": "2025-01-10T12:00:00Z",
  "translations": [
    {
      "langCode": "en",
      "title": "Asset Assigned",
      "message": "Asset 'MacBook Pro 2023' (Tag: LAPTOP-001) has been assigned to you."
    },
    {
      "langCode": "id-ID",
      "title": "Aset Ditugaskan",
      "message": "Aset 'MacBook Pro 2023' (Tag: LAPTOP-001) telah ditugaskan kepada Anda."
    }
  ]
}
```

### FCM Push Notification
```json
{
  "notification": {
    "title": "Asset Assigned",
    "body": "Asset 'MacBook Pro 2023' (Tag: LAPTOP-001) has been assigned to you."
  },
  "data": {
    "notification_id": "01HTNOTIF...",
    "user_id": "01HTUSER...",
    "related_asset_id": "01HTASSET...",
    "type": "STATUS_CHANGE",
    "is_read": "false",
    "click_action": "FLUTTER_NOTIFICATION_CLICK"
  }
}
```

## Environment Variables

```env
# Firebase Configuration
ENABLE_FCM=true
FIREBASE_TYPE=service_account
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_PRIVATE_KEY_ID=your-private-key-id
FIREBASE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n"
FIREBASE_CLIENT_EMAIL=firebase-adminsdk-xxxxx@your-project.iam.gserviceaccount.com
FIREBASE_CLIENT_ID=123456789
FIREBASE_AUTH_URI=https://accounts.google.com/o/oauth2/auth
FIREBASE_TOKEN_URI=https://oauth2.googleapis.com/token
FIREBASE_AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs
FIREBASE_CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/...
FIREBASE_UNIVERSE_DOMAIN=googleapis.com
```

## Client Implementation Example

### Flutter/Dart
```dart
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:http/http.dart' as http;

class NotificationService {
  final FirebaseMessaging _messaging = FirebaseMessaging.instance;

  // Request permission and get token
  Future<void> initialize() async {
    // Request permission
    NotificationSettings settings = await _messaging.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );

    if (settings.authorizationStatus == AuthorizationStatus.authorized) {
      // Get FCM token
      String? token = await _messaging.getToken();

      if (token != null) {
        // Send token to backend
        await updateFCMToken(token);
      }

      // Listen for token refresh
      _messaging.onTokenRefresh.listen((newToken) {
        updateFCMToken(newToken);
      });
    }
  }

  // Update FCM token on backend
  Future<void> updateFCMToken(String token) async {
    final response = await http.patch(
      Uri.parse('https://api.example.com/api/v1/users/profile/fcm-token'),
      headers: {
        'Authorization': 'Bearer $jwtToken',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({'fcmToken': token}),
    );

    if (response.statusCode == 200) {
      print('FCM token updated successfully');
    }
  }

  // Handle foreground messages
  void setupForegroundMessageHandler() {
    FirebaseMessaging.onMessage.listen((RemoteMessage message) {
      print('Got a message in foreground!');
      print('Message data: ${message.data}');

      if (message.notification != null) {
        print('Title: ${message.notification!.title}');
        print('Body: ${message.notification!.body}');

        // Show local notification
        showLocalNotification(message);
      }
    });
  }

  // Handle background messages
  static Future<void> firebaseMessagingBackgroundHandler(
    RemoteMessage message
  ) async {
    print('Handling background message: ${message.messageId}');
  }

  // Handle notification tap
  void setupNotificationTapHandler() {
    FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) {
      print('Notification tapped!');

      // Navigate to specific screen based on data
      if (message.data['related_asset_id'] != null) {
        // Navigate to asset detail screen
        navigateToAssetDetail(message.data['related_asset_id']);
      }
    });
  }
}
```

### React Native / JavaScript
```javascript
import messaging from '@react-native-firebase/messaging';
import axios from 'axios';

class NotificationService {
  async initialize() {
    // Request permission
    const authStatus = await messaging().requestPermission();
    const enabled =
      authStatus === messaging.AuthorizationStatus.AUTHORIZED ||
      authStatus === messaging.AuthorizationStatus.PROVISIONAL;

    if (enabled) {
      // Get FCM token
      const token = await messaging().getToken();

      if (token) {
        // Send token to backend
        await this.updateFCMToken(token);
      }

      // Listen for token refresh
      messaging().onTokenRefresh(async (newToken) => {
        await this.updateFCMToken(newToken);
      });
    }
  }

  async updateFCMToken(token) {
    try {
      await axios.patch(
        'https://api.example.com/api/v1/users/profile/fcm-token',
        { fcmToken: token },
        {
          headers: {
            'Authorization': `Bearer ${jwtToken}`,
            'Content-Type': 'application/json',
          },
        }
      );
      console.log('FCM token updated successfully');
    } catch (error) {
      console.error('Error updating FCM token:', error);
    }
  }

  setupForegroundMessageHandler() {
    // Handle foreground messages
    messaging().onMessage(async (remoteMessage) => {
      console.log('Foreground message received:', remoteMessage);

      // Show local notification
      // Use @react-native-community/push-notification-ios or similar
    });
  }

  setupBackgroundMessageHandler() {
    // Handle background messages
    messaging().setBackgroundMessageHandler(async (remoteMessage) => {
      console.log('Background message received:', remoteMessage);
    });
  }

  setupNotificationTapHandler() {
    // Handle notification opened app
    messaging().onNotificationOpenedApp((remoteMessage) => {
      console.log('Notification tapped:', remoteMessage);

      // Navigate based on data
      if (remoteMessage.data?.related_asset_id) {
        navigation.navigate('AssetDetail', {
          assetId: remoteMessage.data.related_asset_id,
        });
      }
    });

    // Check if app was opened from notification
    messaging()
      .getInitialNotification()
      .then((remoteMessage) => {
        if (remoteMessage) {
          console.log('App opened from notification:', remoteMessage);
          // Handle initial notification
        }
      });
  }
}
```

## Testing

### 1. Test FCM Integration
```bash
# 1. Start the server
go run app/main.go

# 2. Login to get JWT token
curl -X POST http://localhost:5000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'

# 3. Update FCM token
curl -X PATCH http://localhost:5000/api/v1/users/profile/fcm-token \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "fcmToken": "your-device-fcm-token"
  }'

# 4. Create asset with assignment
curl -X POST http://localhost:5000/api/v1/assets \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "assetTag": "TEST-001",
    "assetName": "Test Asset",
    "categoryId": "01HTCAT...",
    "assignedTo": "<user_id>",
    "status": "Active",
    "condition": "Good"
  }'
```

### 2. Check Notification in Database
```sql
-- Get latest notifications for user
SELECT n.id, n.user_id, n.type, n.is_read, n.created_at,
       nt.lang_code, nt.title, nt.message
FROM notifications n
JOIN notification_translations nt ON n.id = nt.notification_id
WHERE n.user_id = '01HTUSER...'
ORDER BY n.created_at DESC
LIMIT 10;
```

## Use Cases

### Use Case 1: Asset Assignment (Current Implementation)
✅ Implemented
- User A creates/updates asset
- Asset assigned to User B
- Notification created in database
- Push notification sent to User B's device

### Use Case 2: Maintenance Reminder (Future)
- System checks maintenance schedules daily
- Sends notification to assigned technician
- Reminder 24 hours before maintenance due

### Use Case 3: Warranty Expiration (Future)
- System checks warranty dates weekly
- Sends notification to admin 30 days before expiry
- Helps plan equipment replacement

### Use Case 4: Issue Report (Future)
- User reports issue on asset
- Notification sent to admin/manager
- Quick response to problems

### Use Case 5: Asset Movement (Future)
- Asset moved to new location
- Notification sent to asset owner
- Location tracking updates

## Troubleshooting

### FCM Token Not Working
1. Check Firebase credentials in `.env`
2. Verify ENABLE_FCM=true
3. Check server logs for FCM errors
4. Verify device token is valid

### Notification Not Received on Device
1. Check if FCM token was updated successfully
2. Verify notification was created in database
3. Check FCM logs in Firebase Console
4. Verify device has internet connection
5. Check app notification permissions

### Notification Created but FCM Not Sent
- This is expected behavior if FCM is disabled or not configured
- System will continue to work (notification stored in database)
- User can still see notifications via API

## Future Enhancements

1. **Notification Topics**: Subscribe users to topics (e.g., "maintenance_alerts")
2. **Scheduled Notifications**: Send notifications at specific times
3. **Rich Notifications**: Add images, buttons, and custom actions
4. **Notification Preferences**: Let users choose which notifications to receive
5. **Notification History**: Track notification delivery status
6. **Batch Notifications**: Send to multiple users efficiently
7. **Notification Templates**: Reusable notification formats
8. **Analytics**: Track notification open rates and engagement

## License

MIT
