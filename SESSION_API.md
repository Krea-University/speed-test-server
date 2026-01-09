# Session-Based Speed Test API

## Overview

The speed test server now uses a **session-based approach** for collecting and storing test results. This provides a secure, organized way to track speed test data without exposing API keys in the frontend.

## How It Works

1. **Start Session**: Frontend calls `/test/start` to create a new test session
2. **Get Token**: Backend returns a unique session token (valid for 15 minutes)
3. **Run Tests**: Frontend uses this token for all test requests
4. **Complete Session**: Frontend calls `/test/complete?token=xxx` to save all results to database

## API Endpoints

### 1. Start Test Session

**Endpoint**: `POST /test/start`

**Description**: Creates a new test session and returns a unique token

**Request**: No body required

**Response**:
```json
{
  "token": "e4c27968-82d6-4d9d-818e-6b6e8e1b45d2",
  "expires_at": "2026-01-09T08:44:59Z",
  "message": "Test session created successfully"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/test/start
```

### 2. Complete Test Session

**Endpoint**: `POST /test/complete?token=<session_token>`

**Description**: Completes the test session and saves all collected data to the database

**Query Parameters**:
- `token` (required): The session token received from `/test/start`

**Response**:
```json
{
  "id": "f00122eb-3a1d-400c-88ba-4f3bd622d1f1",
  "message": "Speed test completed and saved successfully"
}
```

**Example**:
```bash
TOKEN="e4c27968-82d6-4d9d-818e-6b6e8e1b45d2"
curl -X POST "http://localhost:8080/test/complete?token=$TOKEN"
```

## Frontend Integration

### Basic Flow

```javascript
// 1. Start test session
const startResponse = await fetch('/test/start', {
    method: 'POST'
});
const { token } = await startResponse.json();

// 2. Run your speed tests
// (ping, download, upload, etc.)
// Backend automatically tracks data using the session

// 3. Complete the test
const completeResponse = await fetch(`/test/complete?token=${token}`, {
    method: 'POST'
});
const { id } = await completeResponse.json();
console.log('Test saved with ID:', id);
```

## Benefits

### 🔒 **Security**
- No API keys exposed in frontend code
- Session tokens are temporary (15-minute expiry)
- Tokens automatically cleaned up after use

### 📊 **Automatic Tracking**
- Backend collects all test data in the session
- No need to manually aggregate results
- Clean, organized data storage

### 🚀 **Simple Integration**
- Just 2 API calls: start and complete
- Token-based authentication
- No complex state management needed

## Session Data

The backend automatically tracks the following data in each session:

- `token` - Unique session identifier
- `client_ip` - Client IP address
- `user_agent` - Browser/client user agent
- `ping_latency_ms` - Ping latency
- `jitter_ms` - Jitter measurement
- `download_speed_mbps` - Download speed
- `upload_speed_mbps` - Upload speed
- `download_bytes` - Total bytes downloaded
- `upload_bytes` - Total bytes uploaded
- `isp` - Internet Service Provider
- `country` - Country code
- `region` - Region/state
- `city` - City name
- `created_at` - Session start time
- `expires_at` - Session expiry time

## Error Handling

### Invalid Token
```json
{
  "error": "Invalid or expired token"
}
```
**HTTP Status**: 404

### Missing Token
```json
{
  "error": "Token required"
}
```
**HTTP Status**: 400

### Database Unavailable
```json
{
  "error": "Database not available"
}
```
**HTTP Status**: 503

## Session Lifecycle

1. **Created**: Session is created when `/test/start` is called
2. **Active**: Session remains active for 15 minutes
3. **Completed**: Session is deleted after `/test/complete` is called
4. **Expired**: Sessions older than 15 minutes are automatically cleaned up

## Notes

- Sessions are stored in memory (not persisted to database until completion)
- Each session can only be completed once
- Expired sessions are automatically removed every 5 minutes
- Session tokens are UUID v4 format

## Testing

```bash
# Create a session
TOKEN=$(curl -s -X POST http://localhost:8080/test/start | jq -r '.token')
echo "Session Token: $TOKEN"

# Complete the session
curl -s -X POST "http://localhost:8080/test/complete?token=$TOKEN" | jq .

# Verify in database
docker-compose exec mysql mysql -uspeedtest -pspeedtest speedtest -e "SELECT * FROM speed_tests ORDER BY created_at DESC LIMIT 1;"
```

## Migration from Old API

If you're migrating from the old API key-based approach:

**Old Way**:
```javascript
// Had to include API key in every request
fetch('/result', {
    method: 'POST',
    headers: {
        'X-API-Key': 'demo-api-key-2025'
    },
    body: JSON.stringify(testData)
});
```

**New Way**:
```javascript
// Get token once
const { token } = await fetch('/test/start', { method: 'POST' }).then(r => r.json());

// Use token for completion
await fetch(`/test/complete?token=${token}`, { method: 'POST' });
```

## Future Enhancements

Potential improvements for the session system:

- [ ] Update session data during tests (real-time tracking)
- [ ] Add session status endpoint to check progress
- [ ] Support partial test completion
- [ ] Add session metadata (browser, device, etc.)
- [ ] Implement session resume capability
