# 🔒 CSRF Protection Implementation

## Overview

API telah diimplementasikan dengan CSRF (Cross-Site Request Forgery) protection untuk meningkatkan keamanan. Semua endpoint yang menggunakan method POST, PUT, dan PATCH memerlukan CSRF token untuk validasi.

## 🚀 Quick Start

### 1. Mendapatkan CSRF Token

```bash
GET /api/csrf
```

**Response:**
```json
{
    "success": true,
    "message": "CSRF token generated successfully",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "expires_at": "2025-09-28T14:22:24+07:00",
        "usage": {
            "header": "X-CSRF-Token",
            "form_field": "csrf_token",
            "query_param": "csrf_token"
        }
    }
}
```

### 2. Menggunakan CSRF Token

CSRF token dapat dikirimkan dalam 3 cara:

#### Header (Recommended)
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: YOUR_TOKEN_HERE" \
  -d '{"phone": "+628123456789", "password": "password123"}'
```

#### Form Field
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -d "phone=+628123456789" \
  -d "password=password123" \
  -d "csrf_token=YOUR_TOKEN_HERE"
```

#### Query Parameter
```bash
curl -X POST "http://localhost:8080/api/auth/register?csrf_token=YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"phone": "+628123456789", "password": "password123"}'
```

## 📋 Endpoint yang Memerlukan CSRF Token

### User Endpoints
- `POST /api/auth/logout`
- `PUT /api/auth/profile`
- `POST /api/auth/change-password`
- `POST /api/game-accounts`
- `PUT /api/game-accounts/:id`
- `DELETE /api/game-accounts/:id`
- `POST /api/transactions`
- `POST /api/transactions/:id/cancel`

### Admin Endpoints
- `POST /api/admin/auth/logout`
- `POST /api/admin/users`
- `PUT /api/admin/users/:id`
- `DELETE /api/admin/users/:id`
- `POST /api/admin/categories`
- `PUT /api/admin/categories/:id`
- `DELETE /api/admin/categories/:id`
- `POST /api/admin/categories/:id/upload-icon`
- `POST /api/admin/products`
- `PUT /api/admin/products/:id`
- `DELETE /api/admin/products/:id`
- `POST /api/admin/products/:id/upload-image`
- `POST /api/admin/vouchers`
- `PUT /api/admin/vouchers/:id`
- `DELETE /api/admin/vouchers/:id`
- `PUT /api/admin/transactions/:id/status`

## 🛡️ Endpoint yang Tidak Memerlukan CSRF Token

### Public Endpoints (GET Methods)
- `GET /health`
- `GET /api/csrf`
- `GET /api/categories`
- `GET /api/products`
- Semua endpoint dengan method GET

### Authentication Endpoints
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh-token`
- `POST /api/admin/auth/login`
- `POST /api/admin/auth/refresh-token`

### Validation Endpoints
- `POST /api/vouchers/validate`

## 🔧 Implementation Details

### Token Generation
- Token menggunakan HMAC-SHA256 untuk signing
- Token berisi timestamp untuk expiry (24 jam)
- Random bytes untuk uniqueness
- Base64 URL encoding untuk safety

### Validation Process
1. Token diambil dari header `X-CSRF-Token`, form field `csrf_token`, atau query parameter `csrf_token`
2. Token diverifikasi signature-nya dengan HMAC
3. Token expiry dicek (valid selama 24 jam)
4. Jika valid, request dilanjutkan, jika tidak, dikembalikan error 403

### Error Responses
```json
{
    "success": false,
    "message": "CSRF token is required for this operation",
    "error": "csrf_token_missing"
}
```

```json
{
    "success": false,
    "message": "CSRF token is invalid or expired",
    "error": "csrf_token_invalid"
}
```

## 🧪 Testing

### 1. Manual Testing
Buka file `test-csrf.html` di browser untuk testing interaktif:
```bash
# Serve the file locally or copy to server static folder
open test-csrf.html
```

### 2. API Testing dengan cURL

```bash
# 1. Get CSRF Token
TOKEN=$(curl -s http://localhost:8080/api/csrf | jq -r '.data.token')

# 2. Test with token
curl -X POST http://localhost:8080/api/csrf/validate \
  -H "X-CSRF-Token: $TOKEN"

# 3. Test without token (should fail)
curl -X POST http://localhost:8080/api/csrf/validate
```

## ⚙️ Configuration

### Environment Variables
```env
# CSRF Secret untuk signing token (wajib diubah di production)
CSRF_SECRET=your-super-secret-csrf-key-change-in-production
```

### Security Best Practices
1. **Selalu gunakan HTTPS di production** - CSRF token sensitive
2. **Ubah CSRF_SECRET di production** - gunakan secret yang strong dan unique
3. **Implement proper error handling** - jangan expose internal information
4. **Monitor CSRF failures** - untuk deteksi serangan

## 🔍 Troubleshooting

### Common Issues

#### 1. Token Missing Error
```json
{"success": false, "message": "CSRF token is required for this operation"}
```
**Solution:** Pastikan token disertakan dalam header `X-CSRF-Token`

#### 2. Token Invalid Error
```json
{"success": false, "message": "CSRF token is invalid or expired"}
```
**Solutions:**
- Token expired (> 24 jam) - ambil token baru
- Token corrupted - ambil token baru
- Wrong CSRF_SECRET - check environment configuration

#### 3. Token Generation Failed
```json
{"success": false, "message": "Failed to generate CSRF token"}
```
**Solution:** Check server logs dan pastikan CSRF middleware properly configured

## 📚 JavaScript/Frontend Integration

### Axios Interceptor
```javascript
// Set up axios interceptor untuk automatic CSRF token
let csrfToken = null;

// Function to get CSRF token
async function getCSRFToken() {
    const response = await axios.get('/api/csrf');
    csrfToken = response.data.data.token;
    return csrfToken;
}

// Request interceptor
axios.interceptors.request.use(async (config) => {
    if (['post', 'put', 'patch', 'delete'].includes(config.method)) {
        if (!csrfToken) {
            await getCSRFToken();
        }
        config.headers['X-CSRF-Token'] = csrfToken;
    }
    return config;
});

// Response interceptor untuk handle CSRF errors
axios.interceptors.response.use(
    response => response,
    async (error) => {
        if (error.response?.status === 403 && 
            error.response?.data?.error === 'csrf_token_invalid') {
            // Token expired, get new one and retry
            await getCSRFToken();
            const originalRequest = error.config;
            originalRequest.headers['X-CSRF-Token'] = csrfToken;
            return axios.request(originalRequest);
        }
        return Promise.reject(error);
    }
);
```

### React Hook
```javascript
import { useState, useCallback } from 'react';

export function useCSRF() {
    const [token, setToken] = useState(null);

    const getToken = useCallback(async () => {
        try {
            const response = await fetch('/api/csrf');
            const data = await response.json();
            if (data.success) {
                setToken(data.data.token);
                return data.data.token;
            }
            throw new Error('Failed to get CSRF token');
        } catch (error) {
            console.error('CSRF token error:', error);
            throw error;
        }
    }, []);

    const makeRequest = useCallback(async (url, options = {}) => {
        let currentToken = token;
        if (!currentToken) {
            currentToken = await getToken();
        }

        const response = await fetch(url, {
            ...options,
            headers: {
                ...options.headers,
                'X-CSRF-Token': currentToken,
            },
        });

        // If CSRF token is invalid, retry with new token
        if (response.status === 403) {
            const errorData = await response.json();
            if (errorData.error === 'csrf_token_invalid') {
                currentToken = await getToken();
                return fetch(url, {
                    ...options,
                    headers: {
                        ...options.headers,
                        'X-CSRF-Token': currentToken,
                    },
                });
            }
        }

        return response;
    }, [token, getToken]);

    return { token, getToken, makeRequest };
}
```

## 🔗 Related Documentation

- [Swagger API Documentation](http://localhost:8080/swagger/index.html)
- [Security Best Practices](./SECURITY.md)
- [API Authentication](./AUTH.md)