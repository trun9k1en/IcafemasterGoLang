# API Documentation

**Base URL:** `http://localhost:8080/api/v1`

---

## 1. Authentication APIs

### 1.1 Đăng ký tài khoản

**Endpoint:** `POST /auth/register`
**Access:** Public

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "password123",
  "phone": "0901234567",
  "full_name": "John Doe"
}
```

**Response Success (201):**
```json
{
  "statusCode": 201,
  "message": "Registration successful",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "username": "john_doe",
    "phone": "0901234567",
    "full_name": "John Doe",
    "role": "customer",
    "permissions": ["registration:read", "file:read"],
    "is_active": true,
    "created_on": "2024-01-15T10:30:00Z",
    "modified_on": "2024-01-15T10:30:00Z"
  }
}
```

**Response Error (409):**
```json
{
  "statusCode": 409,
  "message": "Username already exists",
  "data": "resource already exists"
}
```

---

### 1.2 Đăng nhập

**Endpoint:** `POST /auth/login`
**Access:** Public

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "password123"
}
```

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "username": "john_doe",
      "email": "",
      "full_name": "John Doe",
      "role": "customer",
      "permissions": ["registration:read", "file:read"]
    }
  }
}
```

**Response Error (401):**
```json
{
  "statusCode": 401,
  "message": "invalid username or password",
  "data": "invalid username or password"
}
```

---

### 1.3 Refresh Token

**Endpoint:** `POST /auth/refresh`
**Access:** Public

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": { ... }
  }
}
```

---

### 1.4 Đăng xuất

**Endpoint:** `POST /auth/logout`
**Access:** Public

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "Logged out successfully",
  "data": null
}
```

---

## 2. User Management APIs (Admin Only)

> **Yêu cầu:** Header `Authorization: Bearer <access_token>`
> **Quyền:** Chỉ user có role `admin` mới truy cập được

---

### 2.1 Tạo user mới

**Endpoint:** `POST /users`
**Access:** Admin only

**Request Body:**
```json
{
  "username": "staff_user",
  "password": "password123",
  "phone": "0912345678",
  "email": "staff@example.com",
  "full_name": "Staff User",
  "role": "staff"
}
```

**Các role hợp lệ:** `admin`, `manager`, `sale`, `staff`, `customer`

**Response Success (201):**
```json
{
  "statusCode": 201,
  "message": "User created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439012",
    "username": "staff_user",
    "email": "staff@example.com",
    "phone": "0912345678",
    "full_name": "Staff User",
    "role": "staff",
    "permissions": ["registration:read", "registration:write", "file:read", "file:write"],
    "is_active": true,
    "created_on": "2024-01-15T10:30:00Z",
    "modified_on": "2024-01-15T10:30:00Z"
  }
}
```

---

### 2.2 Lấy danh sách users

**Endpoint:** `GET /users`
**Access:** Admin only

**Query Parameters:**
| Param | Type | Default | Description |
|-------|------|---------|-------------|
| limit | int | 10 | Số lượng kết quả |
| offset | int | 0 | Vị trí bắt đầu |

**Example:** `GET /users?limit=20&offset=0`

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "username": "john_doe",
      "phone": "0901234567",
      "full_name": "John Doe",
      "role": "customer",
      "permissions": ["registration:read", "file:read"],
      "is_active": true,
      "created_on": "2024-01-15T10:30:00Z",
      "modified_on": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "total": 50,
    "limit": 10,
    "offset": 0
  }
}
```

---

### 2.3 Lấy chi tiết user

**Endpoint:** `GET /users/:id`
**Access:** Admin only

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "username": "john_doe",
    "phone": "0901234567",
    "full_name": "John Doe",
    "role": "customer",
    "permissions": ["registration:read", "file:read"],
    "custom_permissions": [],
    "is_active": true,
    "created_on": "2024-01-15T10:30:00Z",
    "modified_on": "2024-01-15T10:30:00Z",
    "last_login": "2024-01-15T12:00:00Z"
  }
}
```

---

### 2.4 Cập nhật thông tin user

**Endpoint:** `PUT /users/:id`
**Access:** Admin only

**Request Body:** (tất cả fields là optional)
```json
{
  "email": "newemail@example.com",
  "phone": "0987654321",
  "full_name": "New Name",
  "role": "manager",
  "is_active": true,
  "custom_permissions": ["user:manage"]
}
```

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "User updated successfully",
  "data": { ... }
}
```

---

### 2.5 Cập nhật quyền user

**Endpoint:** `PUT /users/:id/role`
**Access:** Admin only

**Request Body:**
```json
{
  "role": "manager",
  "custom_permissions": ["user:manage", "file:delete"]
}
```

**Giải thích:**
- `role`: Thay đổi role của user → permissions sẽ tự động cập nhật theo role
- `custom_permissions`: Thêm quyền bổ sung ngoài quyền của role

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "User role updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "username": "john_doe",
    "role": "manager",
    "permissions": [
      "registration:read",
      "registration:write",
      "registration:delete",
      "file:read",
      "file:write",
      "file:delete"
    ],
    "custom_permissions": ["user:manage", "file:delete"]
  }
}
```

---

### 2.6 Đổi mật khẩu user

**Endpoint:** `PUT /users/:id/password`
**Access:** Admin only

**Request Body:**
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "Password changed successfully",
  "data": null
}
```

---

### 2.7 Xóa user

**Endpoint:** `DELETE /users/:id`
**Access:** Admin only

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "User deleted successfully",
  "data": null
}
```

---

## 3. File Management APIs

### 3.1 Upload file

**Endpoint:** `POST /files/upload`
**Access:** Public
**Content-Type:** `multipart/form-data`

**Request:**
```
file: <binary file>
```

**Response Success (201):**
```json
{
  "statusCode": 201,
  "message": "File uploaded successfully",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "file_name": "1705312200_document.pdf",
    "original_name": "document.pdf",
    "file_path": "./uploads/files/1705312200_document.pdf",
    "file_type": "document",
    "mime_type": "application/pdf",
    "size": 1024000,
    "url": "http://localhost:8080/api/v1/files/serve/1705312200_document.pdf",
    "created_on": "2024-01-15T10:30:00Z"
  }
}
```

---

### 3.2 Upload video

**Endpoint:** `POST /videos/upload`
**Access:** Public
**Content-Type:** `multipart/form-data`

**Request:**
```
file: <binary video file>
```

**Response Success (201):**
```json
{
  "statusCode": 201,
  "message": "Video uploaded successfully",
  "data": {
    "id": "507f1f77bcf86cd799439014",
    "file_name": "1705312200_video.mp4",
    "original_name": "video.mp4",
    "file_path": "./uploads/videos/1705312200_video.mp4",
    "file_type": "video",
    "mime_type": "video/mp4",
    "size": 50240000,
    "url": "http://localhost:8080/api/v1/videos/serve/1705312200_video.mp4",
    "created_on": "2024-01-15T10:30:00Z"
  }
}
```

---

### 3.3 Lấy danh sách files

**Endpoint:** `GET /files`
**Access:** Public

**Query Parameters:**
| Param | Type | Default | Description |
|-------|------|---------|-------------|
| limit | int | 10 | Số lượng kết quả |
| offset | int | 0 | Vị trí bắt đầu |

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "Files retrieved successfully",
  "data": [ ... ],
  "meta": {
    "total": 100,
    "limit": 10,
    "offset": 0
  }
}
```

---

### 3.4 Lấy danh sách videos

**Endpoint:** `GET /videos`
**Access:** Public

**Response:** Tương tự GET /files

---

### 3.5 Lấy chi tiết file

**Endpoint:** `GET /files/:id`
**Access:** Public

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "File retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "file_name": "1705312200_document.pdf",
    "original_name": "document.pdf",
    "file_type": "document",
    "mime_type": "application/pdf",
    "size": 1024000,
    "url": "http://localhost:8080/api/v1/files/serve/1705312200_document.pdf",
    "created_on": "2024-01-15T10:30:00Z"
  }
}
```

---

### 3.6 Download file

**Endpoint:** `GET /files/serve/:filename`
**Access:** Public

**Response:** Binary file với headers:
```
Content-Description: File Transfer
Content-Disposition: attachment; filename=document.pdf
Content-Type: application/octet-stream
```

---

### 3.7 Stream video

**Endpoint:** `GET /videos/serve/:filename`
**Access:** Public

**Response:** Video stream với headers:
```
Accept-Ranges: bytes
Content-Type: video/mp4
```

---

### 3.8 Xóa file

**Endpoint:** `DELETE /files/:id`
**Access:** Public (nên chuyển sang Protected)

**Response Success (200):**
```json
{
  "statusCode": 200,
  "message": "File deleted successfully",
  "data": null
}
```

---

## 4. Hệ thống phân quyền

### 4.1 Roles

| Role | Mô tả |
|------|-------|
| `admin` | Quản trị viên - toàn quyền |
| `manager` | Quản lý - quyền cao |
| `sale` | Nhân viên kinh doanh |
| `staff` | Nhân viên |
| `customer` | Khách hàng (mặc định khi đăng ký) |

### 4.2 Permissions

| Permission | Mô tả |
|------------|-------|
| `registration:read` | Xem đăng ký |
| `registration:write` | Tạo/sửa đăng ký |
| `registration:delete` | Xóa đăng ký |
| `file:read` | Xem file |
| `file:write` | Upload file |
| `file:delete` | Xóa file |
| `user:manage` | Quản lý users |

### 4.3 Role-Permission Mapping

| Role | Permissions |
|------|-------------|
| `admin` | Tất cả permissions |
| `manager` | registration:*, file:* |
| `sale` | registration:read/write, file:read/write |
| `staff` | registration:read/write, file:read/write |
| `customer` | registration:read, file:read |

### 4.4 Custom Permissions

Admin có thể gán thêm `custom_permissions` cho user ngoài permissions từ role.

**Ví dụ:** User có role `staff` nhưng được gán thêm `user:manage`:
```json
{
  "role": "staff",
  "permissions": ["registration:read", "registration:write", "file:read", "file:write"],
  "custom_permissions": ["user:manage"]
}
```

→ User này có tất cả permissions của `staff` + `user:manage`

---

## 5. Error Codes

| HTTP Code | Mô tả |
|-----------|-------|
| 200 | Thành công |
| 201 | Tạo mới thành công |
| 400 | Request không hợp lệ |
| 401 | Chưa xác thực / Token không hợp lệ |
| 403 | Không có quyền truy cập |
| 404 | Không tìm thấy |
| 409 | Conflict (duplicate) |
| 500 | Lỗi server |

---

## 6. Health Check

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "ok",
  "message": "Server is running"
}
```
