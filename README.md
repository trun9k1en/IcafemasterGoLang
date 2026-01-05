# iCafe Registration API

Hệ thống API quản lý đăng ký và upload file cho iCafe, xây dựng bằng Go với kiến trúc Clean Architecture.

## Mục lục

- [Cấu trúc dự án](#cấu-trúc-dự-án)
- [Yêu cầu hệ thống](#yêu-cầu-hệ-thống)
- [Cài đặt Database](#cài-đặt-database)
- [Cấu hình môi trường](#cấu-hình-môi-trường)
- [Chạy dự án](#chạy-dự-án)
- [API Endpoints](#api-endpoints)
- [Demo](#demo)

---

## Cấu trúc dự án

```
icafe-registration/
├── cmd/api/                          # Entry point của ứng dụng
│   └── main.go                       # Khởi tạo server
├── internal/                         # Code nội bộ của ứng dụng
│   ├── config/                       # Quản lý cấu hình
│   │   ├── config.go                 # Load biến môi trường
│   │   └── database.go               # Kết nối MongoDB
│   ├── domain/                       # Domain layer (entities, interfaces)
│   │   ├── auth.go                   # Định nghĩa auth
│   │   ├── user.go                   # Entity và interface User
│   │   ├── registration.go           # Entity và interface Registration
│   │   ├── file.go                   # Entity và interface File
│   │   └── errors.go                 # Định nghĩa lỗi
│   ├── delivery/http/                # HTTP handler layer
│   │   ├── router.go                 # Định nghĩa routes
│   │   ├── registration_handler.go   # Handler cho registration
│   │   ├── file_handler.go           # Handler cho file upload/download
│   │   └── middleware.go             # CORS, logging, recovery
│   ├── repository/mongodb/           # Data access layer
│   │   ├── registration_repository.go
│   │   ├── file_repository.go
│   │   └── user_repository.go
│   └── usecase/                      # Business logic layer
│       ├── registration_usecase.go
│       ├── file_usecase.go
│       └── auth_usecase.go
├── pkg/                              # Shared utilities
│   ├── response/                     # API response helpers
│   └── validator/                    # Input validation
├── uploads/                          # Thư mục lưu file upload
│   ├── files/                        # Lưu documents
│   └── videos/                       # Lưu videos
├── go.mod                            # Dependencies
├── Dockerfile                        # Docker image
├── docker-compose.yml                # Docker compose
├── Makefile                          # Build commands
└── .env.example                      # Template biến môi trường
```

### Kiến trúc Clean Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Delivery (HTTP)                       │
│              Gin Framework + Middleware                  │
└─────────────────────────┬───────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────┐
│                      Usecase                             │
│                  Business Logic                          │
└─────────────────────────┬───────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────┐
│                    Repository                            │
│                  MongoDB Access                          │
└─────────────────────────┬───────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────┐
│                      Domain                              │
│            Entities, Interfaces, Errors                  │
└─────────────────────────────────────────────────────────┘
```

---

## Yêu cầu hệ thống

- **Go** 1.22+
- **MongoDB** 7.0+
- **Make** (tùy chọn, để sử dụng Makefile)
- **Docker & Docker Compose** (tùy chọn)

---

## Cài đặt Database

### Cách 1: Cài MongoDB local

**macOS (Homebrew):**
```bash
brew tap mongodb/brew
brew install mongodb-community@7.0
brew services start mongodb-community@7.0
```

**Ubuntu/Debian:**
```bash
# Import MongoDB public GPG key
curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg --dearmor

# Add repository
echo "deb [ signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list

# Install
sudo apt-get update
sudo apt-get install -y mongodb-org

# Start service
sudo systemctl start mongod
sudo systemctl enable mongod
```

**Windows:**
- Tải MongoDB từ [mongodb.com/try/download/community](https://www.mongodb.com/try/download/community)
- Chạy installer và làm theo hướng dẫn

### Cách 2: Sử dụng Docker (Khuyến nghị)

```bash
# Chạy MongoDB container
docker run -d \
  --name mongodb \
  -p 27017:27017 \
  -v mongodb_data:/data/db \
  mongo:7.0
```

### Kiểm tra kết nối

```bash
# Sử dụng mongosh
mongosh "mongodb://localhost:27017"

# Hoặc kiểm tra với docker
docker exec -it mongodb mongosh
```

### Cấu trúc Database

Database: `icafe_registration`

**Collections:**

1. **registrations** - Lưu thông tin đăng ký
   ```json
   {
     "_id": "ObjectId",
     "full_name": "string",
     "phone_number": "string",
     "email": "string (unique)",
     "address": "string",
     "workstation_num": "int",
     "created_on": "datetime",
     "modified_on": "datetime"
   }
   ```

2. **files** - Lưu metadata của file
   ```json
   {
     "_id": "ObjectId",
     "file_name": "string (unique)",
     "original_name": "string",
     "file_path": "string",
     "file_type": "document|video|image",
     "mime_type": "string",
     "size": "int64",
     "url": "string",
     "created_on": "datetime"
   }
   ```

3. **users** - Lưu thông tin người dùng
   ```json
   {
     "_id": "ObjectId",
     "username": "string",
     "email": "string",
     "password": "string (hashed)",
     "full_name": "string",
     "role": "admin|manager|staff|customer",
     "is_active": "boolean",
     "created_on": "datetime"
   }
   ```

---

## Cấu hình môi trường

Tạo file `.env` từ template:

```bash
cp .env.example .env
```

Nội dung file `.env`:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=icafe_registration

# File Upload Configuration
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=52428800  # 50MB

# Base URL (for file URLs)
BASE_URL=http://localhost:8080

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-key-change-in-production
JWT_ACCESS_TOKEN_DURATION=15    # minutes
JWT_REFRESH_TOKEN_DURATION=168  # hours (7 days)
```

---

## Chạy dự án

### Cách 1: Chạy trực tiếp với Go

```bash
# Clone project (nếu cần)
git clone <repository-url>
cd icafe-registration

# Cài đặt dependencies
go mod download
# hoặc
make deps

# Tạo file .env
cp .env.example .env

# Chạy server
go run cmd/api/main.go
# hoặc
make run
```

### Cách 2: Build và chạy binary

```bash
# Build
make build

# Chạy
./bin/api
```

### Cách 3: Sử dụng Docker Compose (Khuyến nghị)

```bash
# Khởi động tất cả services (MongoDB + API)
make docker-up
# hoặc
docker-compose up -d

# Xem logs
docker-compose logs -f

# Dừng services
make docker-down
# hoặc
docker-compose down
```

### Các lệnh Makefile hữu ích

```bash
make build        # Build binary
make run          # Chạy với go run
make test         # Chạy tests
make fmt          # Format code
make lint         # Chạy linter
make docker-build # Build Docker image
make docker-up    # Khởi động Docker Compose
make docker-down  # Dừng Docker Compose
make clean        # Xóa build artifacts
```

### Kiểm tra server đã chạy

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok"
}
```

---

## API Endpoints

### Health Check

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/health` | Kiểm tra trạng thái server |

### Registration

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | `/api/v1/registrations` | Tạo đăng ký mới |
| GET | `/api/v1/registrations` | Danh sách đăng ký (có phân trang) |
| GET | `/api/v1/registrations/:id` | Lấy thông tin đăng ký theo ID |
| PUT | `/api/v1/registrations/:id` | Cập nhật đăng ký |
| DELETE | `/api/v1/registrations/:id` | Xóa đăng ký |

### File Management

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | `/api/v1/files/upload` | Upload document |
| POST | `/api/v1/videos/upload` | Upload video |
| GET | `/api/v1/files` | Danh sách documents |
| GET | `/api/v1/videos` | Danh sách videos |
| GET | `/api/v1/files/:id` | Lấy thông tin file |
| DELETE | `/api/v1/files/:id` | Xóa file |
| GET | `/api/v1/files/serve/:filename` | Download file |
| GET | `/api/v1/videos/serve/:filename` | Stream video |

---

## Demo

### 1. Tạo đăng ký mới

```bash
curl -X POST http://localhost:8080/api/v1/registrations \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Nguyen Van A",
    "phone_number": "0901234567",
    "email": "nguyenvana@example.com",
    "address": "123 Nguyen Hue, Q1, HCM",
    "workstation_num": 5
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "id": "6789abc123def456...",
    "full_name": "Nguyen Van A",
    "phone_number": "0901234567",
    "email": "nguyenvana@example.com",
    "address": "123 Nguyen Hue, Q1, HCM",
    "workstation_num": 5,
    "created_on": "2024-01-15T10:30:00Z",
    "modified_on": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Lấy danh sách đăng ký (có phân trang)

```bash
curl "http://localhost:8080/api/v1/registrations?limit=10&offset=0"
```

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": "6789abc123def456...",
      "full_name": "Nguyen Van A",
      "email": "nguyenvana@example.com",
      ...
    }
  ],
  "meta": {
    "total": 1,
    "limit": 10,
    "offset": 0
  }
}
```

### 3. Lấy đăng ký theo ID

```bash
curl http://localhost:8080/api/v1/registrations/6789abc123def456
```

### 4. Cập nhật đăng ký

```bash
curl -X PUT http://localhost:8080/api/v1/registrations/6789abc123def456 \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Nguyen Van A",
    "phone_number": "0909876543",
    "email": "nguyenvana@example.com",
    "address": "456 Le Loi, Q1, HCM",
    "workstation_num": 10
  }'
```

### 5. Xóa đăng ký

```bash
curl -X DELETE http://localhost:8080/api/v1/registrations/6789abc123def456
```

### 6. Upload file

```bash
# Upload document
curl -X POST http://localhost:8080/api/v1/files/upload \
  -F "file=@/path/to/document.pdf"

# Upload video
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "file=@/path/to/video.mp4"
```

Response:
```json
{
  "success": true,
  "data": {
    "id": "abc123...",
    "file_name": "550e8400-e29b-41d4-a716-446655440000.pdf",
    "original_name": "document.pdf",
    "file_type": "document",
    "mime_type": "application/pdf",
    "size": 1024000,
    "url": "http://localhost:8080/api/v1/files/serve/550e8400-e29b-41d4-a716-446655440000.pdf",
    "created_on": "2024-01-15T10:30:00Z"
  }
}
```

### 7. Danh sách files

```bash
# Danh sách documents
curl "http://localhost:8080/api/v1/files?limit=10&offset=0"

# Danh sách videos
curl "http://localhost:8080/api/v1/videos?limit=10&offset=0"
```

### 8. Download/Stream file

```bash
# Download document
curl -O http://localhost:8080/api/v1/files/serve/550e8400-e29b-41d4-a716-446655440000.pdf

# Stream video (mở trong browser)
open http://localhost:8080/api/v1/videos/serve/550e8400-e29b-41d4-a716-446655440000.mp4
```

### 9. Xóa file

```bash
curl -X DELETE http://localhost:8080/api/v1/files/abc123
```

---

## File types được hỗ trợ

| Loại | MIME Types |
|------|------------|
| Images | image/jpeg, image/png, image/gif |
| Videos | video/mp4, video/mpeg, video/quicktime, video/webm |
| Documents | application/pdf, application/zip, application/msword, application/vnd.openxmlformats-officedocument.wordprocessingml.document |

**Giới hạn kích thước:** 50MB (có thể cấu hình qua `MAX_FILE_SIZE`)

---

## Troubleshooting

### Lỗi kết nối MongoDB

```bash
# Kiểm tra MongoDB đang chạy
mongosh --eval "db.adminCommand('ping')"

# Hoặc với Docker
docker ps | grep mongo
```

### Lỗi permission thư mục uploads

```bash
mkdir -p uploads/files uploads/videos
chmod 755 uploads/files uploads/videos
```

### Xem logs

```bash
# Nếu chạy với Docker
docker-compose logs -f api

# Nếu chạy local, logs sẽ hiển thị trực tiếp trong terminal
```

---

## License

MIT License
