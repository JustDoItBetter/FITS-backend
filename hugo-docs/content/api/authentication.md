---
title: "Authentication"
weight: 1
description: "JWT authentication guide for FITS Backend. Dual-token approach, authentication flow, endpoints, security best practices, and implementation examples."
---

# Authentication

FITS Backend uses JWT (JSON Web Tokens) for authentication with a dual-token approach.

## Authentication Flow

```
┌─────────┐                      ┌─────────┐                 ┌──────────┐
│ Client  │                      │   API   │                 │ Database │
└────┬────┘                      └────┬────┘                 └─────┬────┘
     │                                │                            │
     │  1. POST /auth/login           │                            │
     ├───────────────────────────────>│                            │
     │  {email, password}             │                            │
     │                                │  2. Validate credentials   │
     │                                ├───────────────────────────>│
     │                                │                            │
     │                                │<───────────────────────────┤
     │                                │  User data                 │
     │                                │                            │
     │                                │  3. Verify password        │
     │                                │                            │
     │                                │  4. Generate tokens        │
     │                                │     - Access token         │
     │                                │     - Refresh token        │
     │                                │                            │
     │                                │  5. Store session          │
     │                                ├───────────────────────────>│
     │                                │                            │
     │  6. Return tokens              │                            │
     │<───────────────────────────────┤                            │
     │  {token, refresh_token}        │                            │
     │                                │                            │
     │  7. Authenticated request      │                            │
     │  Authorization: Bearer {token} │                            │
     ├───────────────────────────────>│                            │
     │                                │  8. Validate JWT           │
     │                                │                            │
     │                                │  9. Check permissions      │
     │                                │                            │
     │                                │  10. Process request       │
     │                                ├───────────────────────────>│
     │                                │                            │
     │  11. Response                  │                            │
     │<───────────────────────────────┤                            │
```

## Token Types

### Access Token

Short-lived token for API authentication.

**Properties:**
- **Lifetime**: 15 minutes (default, configurable)
- **Purpose**: Authenticate API requests
- **Storage**: Memory or sessionStorage (never localStorage for security)
- **Transmission**: `Authorization: Bearer {token}` header

**Claims:**
```json
{
  "sub": "123",              // User ID
  "email": "user@example.com",
  "role": "student",
  "exp": 1698012345,         // Expiration timestamp
  "iat": 1698012045          // Issued at timestamp
}
```

### Refresh Token

Long-lived token for obtaining new access tokens.

**Properties:**
- **Lifetime**: 7 days (default, configurable)
- **Purpose**: Obtain new access tokens
- **Storage**: HttpOnly cookie or secure storage
- **Database**: Hashed and stored for validation

**Security Features:**
- Stored hashed in database
- Can be revoked
- Single-use (optional)
- Rotation on use (optional)

## Endpoints

### Bootstrap Admin

Create the initial admin user. Can only be called once when no users exist.

**Endpoint:** `POST /api/v1/bootstrap/init`

**Request:**
```json
{
  "email": "admin@example.com",
  "password": "SecurePassword123!",
  "first_name": "Admin",
  "last_name": "User"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Bootstrap admin created successfully",
  "data": {
    "user": {
      "id": 1,
      "email": "admin@example.com",
      "role": "admin",
      "is_active": true,
      "created_at": "2025-10-23T18:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-10-23T18:15:00Z"
  }
}
```

**Errors:**
- `409 Conflict`: Admin already exists

### Login

Authenticate with email and password.

**Endpoint:** `POST /api/v1/auth/login`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "UserPassword123!"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 2,
      "email": "user@example.com",
      "role": "student",
      "is_active": true,
      "created_at": "2025-10-23T18:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-10-23T18:15:00Z"
  }
}
```

**Errors:**
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid credentials
- `403 Forbidden`: Account inactive

### Register

Register a new user with an invitation token.

**Endpoint:** `POST /api/v1/auth/register`

**Request:**
```json
{
  "invitation_token": "inv_abc123def456",
  "password": "NewUserPassword123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 3,
      "email": "john@example.com",
      "role": "teacher",
      "is_active": true,
      "created_at": "2025-10-23T18:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-10-23T18:15:00Z"
  }
}
```

**Errors:**
- `400 Bad Request`: Invalid input or weak password
- `404 Not Found`: Invalid invitation token
- `409 Conflict`: Email already registered
- `410 Gone`: Invitation expired

### Refresh Token

Obtain a new access token using a refresh token.

**Endpoint:** `POST /api/v1/auth/refresh`

**Headers:**
```
Authorization: Bearer {current_access_token}
```

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-10-23T18:30:00Z"
  }
}
```

**Errors:**
- `401 Unauthorized`: Invalid or expired refresh token

### Logout

Invalidate the current session.

**Endpoint:** `POST /api/v1/auth/logout`

**Headers:**
```
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

### Logout All

Invalidate all sessions for the current user.

**Endpoint:** `POST /api/v1/auth/logout-all`

**Headers:**
```
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "message": "All sessions logged out successfully"
}
```

## Using Authentication

### JavaScript/TypeScript Example

```typescript
class APIClient {
  private baseURL = 'http://localhost:8080';
  private accessToken: string | null = null;
  private refreshToken: string | null = null;
  
  async login(email: string, password: string) {
    const response = await fetch(`${this.baseURL}/api/v1/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });
    
    const data = await response.json();
    
    if (data.success) {
      this.accessToken = data.data.token;
      this.refreshToken = data.data.refresh_token;
      
      // Store in memory or sessionStorage (NOT localStorage)
      sessionStorage.setItem('token', this.accessToken);
      sessionStorage.setItem('refreshToken', this.refreshToken);
    }
    
    return data;
  }
  
  async request(endpoint: string, options: RequestInit = {}) {
    // Add authorization header
    const headers = {
      ...options.headers,
      'Authorization': `Bearer ${this.accessToken}`,
    };
    
    let response = await fetch(`${this.baseURL}${endpoint}`, {
      ...options,
      headers,
    });
    
    // If token expired, refresh and retry
    if (response.status === 401) {
      await this.refreshAccessToken();
      
      headers.Authorization = `Bearer ${this.accessToken}`;
      response = await fetch(`${this.baseURL}${endpoint}`, {
        ...options,
        headers,
      });
    }
    
    return response.json();
  }
  
  async refreshAccessToken() {
    const response = await fetch(`${this.baseURL}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.accessToken}`,
      },
      body: JSON.stringify({
        refresh_token: this.refreshToken,
      }),
    });
    
    const data = await response.json();
    
    if (data.success) {
      this.accessToken = data.data.token;
      this.refreshToken = data.data.refresh_token;
      
      sessionStorage.setItem('token', this.accessToken);
      sessionStorage.setItem('refreshToken', this.refreshToken);
    } else {
      // Refresh failed, redirect to login
      this.logout();
    }
  }
  
  async logout() {
    await fetch(`${this.baseURL}/api/v1/auth/logout`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    });
    
    this.accessToken = null;
    this.refreshToken = null;
    sessionStorage.clear();
  }
}
```

### curl Example

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password"
  }' | jq -r '.data.token' > token.txt

# Use token
TOKEN=$(cat token.txt)

curl -X GET http://localhost:8080/api/v1/student \
  -H "Authorization: Bearer $TOKEN"
```

## Security Best Practices

### Client-Side

1. **Never store tokens in localStorage**
   - Use sessionStorage or memory
   - HttpOnly cookies are preferred

2. **Always use HTTPS in production**
   - Tokens transmitted over TLS
   - Prevents token interception

3. **Implement token refresh**
   - Refresh before expiration
   - Handle refresh failures gracefully

4. **Clear tokens on logout**
   - Remove from storage
   - Call logout endpoint

5. **Validate token expiration client-side**
   - Check `exp` claim
   - Proactively refresh

### Server-Side (Implemented)

1. **Short-lived access tokens** (15 minutes)
2. **Refresh token rotation** (optional)
3. **Session storage** in database
4. **Token revocation** support
5. **Rate limiting** on auth endpoints
6. **Password hashing** with bcrypt
7. **Secure password requirements**

## Password Requirements

- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

## Rate Limiting

Authentication endpoints have strict rate limits:

- **Bootstrap**: 5 attempts per hour per IP
- **Login**: 5 attempts per 15 minutes per IP
- **Register**: 10 attempts per hour per IP
- **Refresh**: 30 attempts per hour per user

## Error Responses

### Invalid Credentials

```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid email or password"
  }
}
```

### Token Expired

```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Token expired"
  }
}
```

### Weak Password

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Password does not meet requirements",
    "details": {
      "field": "password",
      "requirements": [
        "Minimum 8 characters",
        "At least one uppercase letter",
        "At least one lowercase letter",
        "At least one number",
        "At least one special character"
      ]
    }
  }
}
```

## Next Steps

- [API Endpoints](/api/endpoints/) - Full endpoint reference
- [Error Handling](/api/error-handling/) - Error codes and responses
- [Quick Start](/getting-started/quick-start/) - Try the authentication flow
