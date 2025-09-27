# PRD Backend - Waw Store Topup Game Online API

## Status: 📋 SPECIFICATION READY

Dokumen spesifikasi backend API untuk aplikasi web topup game online **Waw Store** berdasarkan analisis frontend yang sudah ada.

### Tech Stack Requirements
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL (wajib)
- **ORM**: GORM (Go ORM)
- **Authentication**: JWT (Access + Refresh Token)
- **File Upload**: AWS S3 (dengan fallback ke local storage /cdn)
- **Caching**: Redis (optional, dijelaskan di section caching)
- **Validation**: Go Validator v10
- **Documentation**: Swagger/OpenAPI dengan gin-swagger
- **Testing**: Go testing package + Testify

## 1. 🗄️ Database Schema (PostgreSQL + GORM)

### GORM Model Structures

#### Users Model
```go
type User struct {
    ID                uint      `gorm:"primaryKey" json:"id"`
    Name              string    `gorm:"size:100;not null" json:"name" validate:"required,min=2,max=100"`
    Phone             string    `gorm:"size:15;unique;not null" json:"phone" validate:"required,phone"`
    PasswordHash      string    `gorm:"size:255;not null" json:"-"`
    Email             *string   `gorm:"size:100" json:"email" validate:"omitempty,email"`
    Status            string    `gorm:"size:20;default:active" json:"status" validate:"oneof=active inactive suspended"`
    EmailVerifiedAt   *time.Time `json:"email_verified_at"`
    PhoneVerifiedAt   *time.Time `json:"phone_verified_at"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
    
    // Relations
    GameAccounts      []GameAccount `gorm:"foreignKey:UserID" json:"game_accounts,omitempty"`
    Transactions      []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}
```

#### Admin Users Model
```go
type AdminUser struct {
    ID            uint       `gorm:"primaryKey" json:"id"`
    Username      string     `gorm:"size:50;unique;not null" json:"username" validate:"required,min=3,max=50"`
    Name          string     `gorm:"size:100;not null" json:"name" validate:"required,min=2,max=100"`
    Email         *string    `gorm:"size:100" json:"email" validate:"omitempty,email"`
    PasswordHash  string     `gorm:"size:255;not null" json:"-"`
    Role          string     `gorm:"size:20;default:operator" json:"role" validate:"oneof=super_admin admin operator moderator"`
    Status        string     `gorm:"size:20;default:active" json:"status" validate:"oneof=active inactive"`
    LastLoginAt   *time.Time `json:"last_login_at"`
    CreatedBy     *uint      `gorm:"index" json:"created_by"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
    
    // Relations
    Creator       *AdminUser `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}
```

#### Categories Model
```go
type Category struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Name         string    `gorm:"size:100;not null" json:"name" validate:"required,min=2,max=100"`
    Slug         string    `gorm:"size:100;unique;not null" json:"slug"`
    Description  *string   `gorm:"type:text" json:"description"`
    IconURL      *string   `gorm:"size:255" json:"icon_url"`
    DisplayOrder int       `gorm:"default:0" json:"display_order"`
    Status       string    `gorm:"size:20;default:active" json:"status" validate:"oneof=active inactive"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    
    // Relations
    Products     []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
```

#### Products Model
```go
type FormField struct {
    Field       string   `json:"field" validate:"required"`
    Label       string   `json:"label" validate:"required"`
    Type        string   `json:"type" validate:"required,oneof=text select number"`
    Required    bool     `json:"required"`
    Placeholder string   `json:"placeholder,omitempty"`
    Options     []string `json:"options,omitempty"`
}

type Product struct {
    ID                 uint        `gorm:"primaryKey" json:"id"`
    CategoryID         uint        `gorm:"not null;index:idx_category_status" json:"category_id" validate:"required"`
    Name               string      `gorm:"size:200;not null" json:"name" validate:"required,min=2,max=200"`
    Slug               string      `gorm:"size:200;unique;not null" json:"slug"`
    Description        *string     `gorm:"type:text" json:"description"`
    Price              float64     `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,min=0"`
    OriginalPrice      *float64    `gorm:"type:decimal(10,2)" json:"original_price" validate:"omitempty,min=0"`
    DiscountPercentage int         `gorm:"default:0" json:"discount_percentage" validate:"min=0,max=100"`
    ImageURL           *string     `gorm:"size:255" json:"image_url"`
    FormConfig         []FormField `gorm:"type:jsonb;not null" json:"form_config" validate:"required,dive"`
    Status             string      `gorm:"size:20;default:active;index:idx_category_status,idx_status_order" json:"status" validate:"oneof=active inactive out_of_stock"`
    DisplayOrder       int         `gorm:"default:0;index:idx_status_order" json:"display_order"`
    CreatedAt          time.Time   `json:"created_at"`
    UpdatedAt          time.Time   `json:"updated_at"`
    
    // Relations
    Category           Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    Transactions       []Transaction `gorm:"foreignKey:ProductID" json:"transactions,omitempty"`
}
```

#### Game Accounts Model
```go
type GameAccount struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null;index:idx_user_game" json:"user_id" validate:"required"`
    GameName  string    `gorm:"size:100;not null;index:idx_user_game" json:"game_name" validate:"required,max=100"`
    GameID    string    `gorm:"size:100;not null" json:"game_id" validate:"required,max=100"`
    Server    *string   `gorm:"size:50" json:"server" validate:"omitempty,max=50"`
    ZoneID    *string   `gorm:"size:50" json:"zone_id" validate:"omitempty,max=50"`
    Nickname  *string   `gorm:"size:100" json:"nickname" validate:"omitempty,max=100"`
    IsPrimary bool      `gorm:"default:false" json:"is_primary"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // Relations
    User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

#### Vouchers Model
```go
type Voucher struct {
    ID                    uint                   `gorm:"primaryKey" json:"id"`
    Code                  string                 `gorm:"size:50;unique;not null;index:idx_code_status" json:"code" validate:"required,uppercase,min=3,max=50"`
    Type                  string                 `gorm:"size:20;not null" json:"type" validate:"required,oneof=percentage fixed"`
    Value                 float64                `gorm:"type:decimal(10,2);not null" json:"value" validate:"required,min=0"`
    Description           *string                `gorm:"type:text" json:"description"`
    ApplicationType       string                 `gorm:"size:20;default:all" json:"application_type" validate:"oneof=all category product"`
    MinTransactionAmount  float64                `gorm:"type:decimal(10,2);default:0" json:"min_transaction_amount" validate:"min=0"`
    MaxDiscountAmount     *float64               `gorm:"type:decimal(10,2)" json:"max_discount_amount" validate:"omitempty,min=0"`
    Quota                 int                    `gorm:"not null" json:"quota" validate:"required,min=1"`
    UsedCount             int                    `gorm:"default:0" json:"used_count"`
    MaxUsesPerUser        int                    `gorm:"default:1" json:"max_uses_per_user" validate:"min=1"`
    StartDate             time.Time              `gorm:"type:date;not null;index:idx_dates" json:"start_date" validate:"required"`
    EndDate               time.Time              `gorm:"type:date;not null;index:idx_dates" json:"end_date" validate:"required,gtfield=StartDate"`
    Status                string                 `gorm:"size:20;default:active;index:idx_code_status" json:"status" validate:"oneof=active inactive expired"`
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
    
    // Relations
    Applications          []VoucherApplication   `gorm:"foreignKey:VoucherID" json:"applications,omitempty"`
    Usages                []VoucherUsage         `gorm:"foreignKey:VoucherID" json:"usages,omitempty"`
}
```

#### Voucher Applications Model
```go
type VoucherApplication struct {
    ID             uint      `gorm:"primaryKey" json:"id"`
    VoucherID      uint      `gorm:"not null" json:"voucher_id" validate:"required"`
    ApplicableID   uint      `gorm:"not null" json:"applicable_id" validate:"required"`
    ApplicableType string    `gorm:"size:20;not null" json:"applicable_type" validate:"required,oneof=category product"`
    CreatedAt      time.Time `json:"created_at"`
    
    // Relations
    Voucher        Voucher   `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"`
}

#### Voucher Usage Model
```go
type VoucherUsage struct {
    ID             uint      `gorm:"primaryKey" json:"id"`
    VoucherID      uint      `gorm:"not null;uniqueIndex:idx_unique_usage" json:"voucher_id" validate:"required"`
    UserID         uint      `gorm:"not null;uniqueIndex:idx_unique_usage" json:"user_id" validate:"required"`
    TransactionID  uint      `gorm:"not null;uniqueIndex:idx_unique_usage" json:"transaction_id" validate:"required"`
    DiscountAmount float64   `gorm:"type:decimal(10,2);not null" json:"discount_amount" validate:"required,min=0"`
    UsedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"used_at"`
    
    // Relations
    Voucher        Voucher     `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"`
    User           User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Transaction    Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
}
```

#### Transactions Model
```go
type GameAccountData struct {
    GameAccount string `json:"game_account" validate:"required"`
    GameZone    string `json:"game_zone,omitempty"`
    GameServer  string `json:"game_server,omitempty"`
    Nickname    string `json:"nickname,omitempty"`
}

type Transaction struct {
    ID               uint            `gorm:"primaryKey" json:"id"`
    TransactionCode  string          `gorm:"size:20;unique;not null;index:idx_transaction_code" json:"transaction_code"`
    UserID           uint            `gorm:"not null;index:idx_user_status" json:"user_id" validate:"required"`
    ProductID        uint            `gorm:"not null" json:"product_id" validate:"required"`
    
    // Game Account Info (denormalized for transaction history)
    GameAccountData  GameAccountData `gorm:"type:jsonb;not null" json:"game_account_data" validate:"required"`
    
    // Pricing
    ProductPrice     float64         `gorm:"type:decimal(10,2);not null" json:"product_price" validate:"required,min=0"`
    PaymentFee       float64         `gorm:"type:decimal(10,2);default:0" json:"payment_fee" validate:"min=0"`
    VoucherDiscount  float64         `gorm:"type:decimal(10,2);default:0" json:"voucher_discount" validate:"min=0"`
    TotalAmount      float64         `gorm:"type:decimal(10,2);not null" json:"total_amount" validate:"required,min=0"`
    
    // Payment
    PaymentMethod    string          `gorm:"size:50;not null" json:"payment_method" validate:"required,oneof=gopay ovo dana bca mandiri bni"`
    PaymentStatus    string          `gorm:"size:20;default:pending" json:"payment_status" validate:"oneof=pending paid failed expired refunded"`
    PaymentReference *string         `gorm:"size:100" json:"payment_reference"`
    PaymentURL       *string         `gorm:"type:text" json:"payment_url"`
    
    // Contact
    WhatsApp         string          `gorm:"size:15;not null" json:"whatsapp" validate:"required,phone"`
    
    // Status & Processing
    Status           string          `gorm:"size:20;default:pending;index:idx_user_status,idx_status_created" json:"status" validate:"oneof=pending processing completed failed cancelled"`
    ProcessedAt      *time.Time      `json:"processed_at"`
    CompletedAt      *time.Time      `json:"completed_at"`
    ExpiredAt        *time.Time      `json:"expired_at"`
    
    // Metadata
    UserAgent        *string         `gorm:"type:text" json:"user_agent"`
    IPAddress        *string         `gorm:"size:45" json:"ip_address"`
    
    CreatedAt        time.Time       `gorm:"index:idx_status_created" json:"created_at"`
    UpdatedAt        time.Time       `json:"updated_at"`
    
    // Relations
    User             User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Product          Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
    Logs             []TransactionLog `gorm:"foreignKey:TransactionID" json:"logs,omitempty"`
    VoucherUsages    []VoucherUsage  `gorm:"foreignKey:TransactionID" json:"voucher_usages,omitempty"`
}
```

#### Transaction Logs Model
```go
type TransactionLog struct {
    ID              uint       `gorm:"primaryKey" json:"id"`
    TransactionID   uint       `gorm:"not null" json:"transaction_id" validate:"required"`
    StatusFrom      *string    `gorm:"size:20" json:"status_from"`
    StatusTo        string     `gorm:"size:20;not null" json:"status_to" validate:"required"`
    Message         *string    `gorm:"type:text" json:"message"`
    Metadata        *string    `gorm:"type:jsonb" json:"metadata"` // Store as JSON string for flexibility
    CreatedByAdmin  *uint      `json:"created_by_admin"`
    CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    
    // Relations
    Transaction     Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
    Admin           *AdminUser  `gorm:"foreignKey:CreatedByAdmin" json:"admin,omitempty"`
}
```

## 2. 🔌 API Endpoints Specification

### Authentication Endpoints

#### Customer Authentication
```
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/logout
POST   /api/auth/refresh-token
GET    /api/auth/profile
PUT    /api/auth/profile
POST   /api/auth/change-password
```

#### Admin Authentication
```
POST   /api/admin/auth/login
POST   /api/admin/auth/logout
POST   /api/admin/auth/refresh-token
GET    /api/admin/auth/profile
```

### Customer Endpoints

#### Products & Categories
```
GET    /api/categories
GET    /api/products
GET    /api/products/:id
GET    /api/products/search?q=mobile+legends&category=1
```

#### Game Accounts Management
```
GET    /api/game-accounts
POST   /api/game-accounts
PUT    /api/game-accounts/:id
DELETE /api/game-accounts/:id
```

#### Vouchers
```
POST   /api/vouchers/validate
  Body: { code: "NEWUSER10", productId: 1, amount: 20000 }
  Response: { valid: true, discount: 2000, finalAmount: 18000 }
```

#### Transactions
```
POST   /api/transactions
GET    /api/transactions
GET    /api/transactions/:id
POST   /api/transactions/:id/cancel
```

### Admin Endpoints

#### Dashboard Analytics
```
GET    /api/admin/analytics/dashboard
  Response: {
    totalSales: 2450000000,
    totalTransactions: 125000,
    totalUsers: 15000,
    totalProducts: 89,
    dailySales: [...],
    topProducts: [...],
    recentTransactions: [...]
  }
```

#### Admin Users CRUD
```
GET    /api/admin/users
POST   /api/admin/users
GET    /api/admin/users/:id
PUT    /api/admin/users/:id
DELETE /api/admin/users/:id
```

#### Categories CRUD
```
GET    /api/admin/categories
POST   /api/admin/categories
PUT    /api/admin/categories/:id
DELETE /api/admin/categories/:id
POST   /api/admin/categories/:id/upload-icon
```

#### Products CRUD
```
GET    /api/admin/products
POST   /api/admin/products
GET    /api/admin/products/:id
PUT    /api/admin/products/:id
DELETE /api/admin/products/:id
POST   /api/admin/products/:id/upload-image
```

#### Vouchers CRUD
```
GET    /api/admin/vouchers
POST   /api/admin/vouchers
PUT    /api/admin/vouchers/:id
DELETE /api/admin/vouchers/:id
GET    /api/admin/vouchers/:id/usage-stats
```

#### Transactions Management
```
GET    /api/admin/transactions
GET    /api/admin/transactions/:id
PUT    /api/admin/transactions/:id/status
GET    /api/admin/transactions/export?format=csv&start_date=2024-01-01
```

## 3. 💼 Business Logic Requirements

### Dynamic Game Account Validation
```javascript
// Product form_config examples
{
  "Mobile Legends": [
    { "field": "gameAccount", "label": "User ID", "type": "text", "required": true },
    { "field": "gameZone", "label": "Zone ID", "type": "text", "required": true }
  ],
  "Free Fire": [
    { "field": "gameAccount", "label": "Player ID", "type": "text", "required": true }
  ],
  "Genshin Impact": [
    { "field": "gameAccount", "label": "UID", "type": "text", "required": true },
    { "field": "gameServer", "label": "Server", "type": "select", "required": true, 
      "options": ["Asia", "America", "Europe", "TW/HK/MO"] }
  ]
}
```

### Voucher System Logic
```javascript
// Voucher validation algorithm
validateVoucher(code, userId, productId, amount) {
  1. Check voucher exists and active
  2. Check date validity (start_date <= now <= end_date)
  3. Check usage quota (used_count < quota)
  4. Check user usage limit (per user max uses)
  5. Check minimum transaction amount
  6. Check applicability (all/category/product)
  7. Calculate discount amount
  8. Apply max discount cap if percentage type
  9. Return discount details
}
```

### Transaction Processing Flow
```javascript
// Transaction creation workflow
createTransaction(userId, productId, gameAccountData, paymentMethod, voucherCode) {
  1. Validate product availability
  2. Validate game account data against product form_config
  3. Apply voucher discount if provided
  4. Calculate total amount (product_price - voucher + payment_fee)
  5. Create transaction record
  6. Generate payment URL/reference
  7. Set expiry time (24 hours)
  8. Send WhatsApp notification
  9. Return transaction details
}
```

### Payment Status Updates
```javascript
// Webhook handler for payment gateway
handlePaymentWebhook(transactionCode, status, reference) {
  1. Verify webhook signature
  2. Find transaction by code
  3. Update payment status
  4. Log status change
  5. Process order if paid
  6. Send notification to user
  7. Update analytics cache
}
```

## 4. 🔐 Security & Validation

### Authentication Flow
```javascript
// JWT Token Structure
AccessToken: {
  sub: userId,
  type: "user", // or "admin"
  role: "customer", // or admin role
  exp: 15 minutes
}

RefreshToken: {
  sub: userId,
  type: "refresh",
  exp: 7 days
}
```

### Input Validation Schemas
```javascript
// Mirror frontend Yup schemas
const userRegistrationSchema = {
  name: required().min(2).max(100),
  phone: required().regex(/^\d{10,15}$/),
  password: required().min(6).max(100)
}

const transactionSchema = {
  productId: required().integer(),
  gameAccountData: required().object(),
  whatsapp: required().regex(/^\d{10,15}$/),
  paymentMethod: required().oneOf(['gopay', 'ovo', 'dana', 'bca', 'mandiri', 'bni'])
}
```

### Rate Limiting
```javascript
// API Rate limits
const rateLimits = {
  "/api/auth/login": "5 req/min per IP",
  "/api/transactions": "10 req/min per user",
  "/api/vouchers/validate": "20 req/min per user",
  "/api/admin/*": "100 req/min per admin"
}
```

## 5. 📊 Analytics & Reporting

### Dashboard Metrics Calculation
```sql
-- Real-time analytics queries
SELECT 
  COUNT(*) as total_transactions,
  SUM(total_amount) as total_sales,
  COUNT(DISTINCT user_id) as unique_customers
FROM transactions 
WHERE status = 'completed' 
  AND DATE(created_at) = CURRENT_DATE;

-- Daily sales trend (last 7 days)
SELECT 
  DATE(created_at) as date,
  SUM(total_amount) as daily_sales
FROM transactions 
WHERE status = 'completed' 
  AND created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 7 DAY)
GROUP BY DATE(created_at)
ORDER BY date;

-- Top products by revenue
SELECT 
  p.name,
  COUNT(t.id) as total_orders,
  SUM(t.total_amount) as total_revenue
FROM transactions t
JOIN products p ON t.product_id = p.id
WHERE t.status = 'completed'
  AND t.created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
GROUP BY p.id, p.name
ORDER BY total_revenue DESC
LIMIT 10;
```

## 6. 🚀 Performance Optimization

### Caching Strategy (Redis - Optional)

**Redis digunakan untuk meningkatkan performa dengan menyimpan data yang sering diakses:**

```go
type CacheService interface {
    Set(key string, value interface{}, ttl time.Duration) error
    Get(key string, dest interface{}) error
    Delete(key string) error
    Clear(pattern string) error
}

// Cache keys dan TTL (Time To Live)
var CacheKeys = map[string]time.Duration{
    "products:active":       1 * time.Hour,    // Daftar produk aktif
    "categories:active":     6 * time.Hour,    // Daftar kategori aktif  
    "voucher:%s":           15 * time.Minute,  // Validasi voucher (%s = voucher code)
    "analytics:dashboard":   5 * time.Minute,  // Data dashboard admin
    "user:profile:%d":      30 * time.Minute,  // Profile user (%d = user_id)
    "product:%d":           2 * time.Hour,     // Detail produk (%d = product_id)
}

// Kegunaan Redis untuk masing-masing data:

// 1. Products & Categories - Mengurangi query database untuk data yang jarang berubah
// 2. Voucher validation - Mencegah spam validation dan mengurangi beban database
// 3. Dashboard analytics - Data statistik mahal untuk dihitung real-time
// 4. User profiles - Data user yang sering diakses saat authenticated requests
// 5. Product details - Detail produk yang sering dilihat customer

// Jika Redis tidak tersedia, aplikasi tetap berjalan normal (direct database query)
```

**Manfaat Caching:**
- **Performance**: Mengurangi response time dari database queries
- **Scalability**: Mengurangi beban database saat traffic tinggi  
- **Cost**: Mengurangi computational cost untuk complex queries (analytics)
- **User Experience**: Response time lebih cepat untuk end-users

**Redis Setup (Optional):**
- Jika `REDIS_URL` di-set: gunakan Redis untuk caching
- Jika tidak: skip caching, langsung ke database (tetap functional)

### Database Indexing (PostgreSQL)
```sql
-- GORM akan auto-create indexes berdasarkan struct tags, tapi ini adalah manual indexes untuk performance

-- Transactions indexes
CREATE INDEX CONCURRENTLY idx_transactions_user_status ON transactions(user_id, status);
CREATE INDEX CONCURRENTLY idx_transactions_status_created ON transactions(status, created_at DESC);
CREATE INDEX CONCURRENTLY idx_transactions_code ON transactions(transaction_code);

-- Products indexes  
CREATE INDEX CONCURRENTLY idx_products_category_status ON products(category_id, status);
CREATE INDEX CONCURRENTLY idx_products_status_order ON products(status, display_order);

-- Vouchers indexes
CREATE INDEX CONCURRENTLY idx_vouchers_code_status ON vouchers(code, status);
CREATE INDEX CONCURRENTLY idx_vouchers_dates ON vouchers(start_date, end_date);

-- Game accounts indexes
CREATE INDEX CONCURRENTLY idx_game_accounts_user_game ON game_accounts(user_id, game_name);

-- Voucher usage unique constraint
CREATE UNIQUE INDEX CONCURRENTLY idx_unique_voucher_usage 
ON voucher_usage(voucher_id, user_id, transaction_id);

-- Full-text search untuk products (optional)
CREATE INDEX CONCURRENTLY idx_products_search 
ON products USING gin(to_tsvector('indonesian', name || ' ' || COALESCE(description, '')));
```

## 7. 🔄 Integration Points

### Payment Gateway Integration
```javascript
// Payment provider interfaces
interface PaymentProvider {
  createPayment(amount, method, reference): PaymentResponse
  checkStatus(reference): PaymentStatus
  handleWebhook(payload, signature): WebhookResult
}

// Support multiple payment gateways
const paymentProviders = {
  midtrans: new MidtransProvider(),
  xendit: new XenditProvider(),
  duitku: new DuitkuProvider()
}
```

### WhatsApp Integration
```javascript
// WhatsApp Business API integration
const whatsappService = {
  sendTransactionConfirmation(whatsapp, transactionData),
  sendPaymentInstructions(whatsapp, paymentData),
  sendCompletionNotification(whatsapp, orderData)
}
```

### File Upload Service (S3 + Local Fallback)
```go
type FileUploadService interface {
    UploadProductImage(file multipart.File, filename string) (*FileUploadResult, error)
    UploadCategoryIcon(file multipart.File, filename string) (*FileUploadResult, error)
    DeleteFile(fileURL string) error
}

type FileUploadResult struct {
    URL      string `json:"url"`
    FileName string `json:"filename"`
    Size     int64  `json:"size"`
}

// Implementation priority:
// 1. AWS S3 (jika AWS_S3_BUCKET di-set)
// 2. Local storage /cdn (fallback)

type S3UploadService struct {
    Client   *s3.Client
    Bucket   string
    Region   string
    BaseURL  string
}

type LocalUploadService struct {
    UploadDir string // "./cdn"
    BaseURL   string // "http://localhost:8080/cdn"
}

// Contoh struktur folder local:
// ./cdn/
// ├── products/
// │   ├── product-1-image.jpg
// │   └── product-2-image.png
// └── categories/
//     ├── category-1-icon.svg
//     └── category-2-icon.png

// File dapat diakses melalui:
// GET /cdn/products/product-1-image.jpg
// GET /cdn/categories/category-1-icon.svg
```

## 8. 📱 Mobile API Considerations

### Response Format Standardization
```javascript
// Standard API response format
{
  success: boolean,
  message: string,
  data: any,
  pagination?: {
    page: number,
    limit: number,
    total: number,
    totalPages: number
  },
  errors?: ValidationError[]
}
```

### Mobile-Specific Endpoints
```javascript
// Optimized for mobile app
GET /api/mobile/dashboard // Compressed dashboard data
GET /api/mobile/products/featured // Featured products only
GET /api/mobile/categories/minimal // Category names and IDs only
```

## 9. 🧪 Testing Strategy

### Test Categories
- **Unit Tests**: Business logic, validation, utilities
- **Integration Tests**: Database operations, external APIs
- **E2E Tests**: Complete transaction flow
- **Load Tests**: High concurrency scenarios
- **Security Tests**: Authentication, authorization, input validation

### Test Coverage Requirements
- **Controllers**: 90%+ coverage
- **Services**: 95%+ coverage
- **Utilities**: 100% coverage
- **Critical Flows**: Transaction processing, payment handling

## 10. 🚚 Deployment & DevOps

### Environment Configuration
```javascript
// Environment variables
const config = {
  NODE_ENV: process.env.NODE_ENV,
  PORT: process.env.PORT || 3000,
  DATABASE_URL: process.env.DATABASE_URL,
  JWT_SECRET: process.env.JWT_SECRET,
  REDIS_URL: process.env.REDIS_URL,
  
  // Payment gateways
  MIDTRANS_SERVER_KEY: process.env.MIDTRANS_SERVER_KEY,
  MIDTRANS_CLIENT_KEY: process.env.MIDTRANS_CLIENT_KEY,
  
  // File upload
  CLOUDINARY_URL: process.env.CLOUDINARY_URL,
  
  // WhatsApp
  WHATSAPP_API_KEY: process.env.WHATSAPP_API_KEY,
  WHATSAPP_PHONE_ID: process.env.WHATSAPP_PHONE_ID
}
```

### Docker Configuration (Go)
```dockerfile
# Multi-stage build for Go backend
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for go modules)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates untuk HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Create cdn directory untuk local file storage
RUN mkdir -p ./cdn/products ./cdn/categories

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]
```

### Docker Compose untuk Development
```yaml
version: '3.8'

services:
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_USER=topup_user
      - DB_PASSWORD=topup_password
      - DB_NAME=topup_db
      - DB_PORT=5432
      - JWT_SECRET=your-super-secret-jwt-key
      - REDIS_URL=redis://redis:6379
      # AWS S3 (optional)
      # - AWS_S3_BUCKET=your-bucket-name
      # - AWS_ACCESS_KEY_ID=your-access-key
      # - AWS_SECRET_ACCESS_KEY=your-secret-key
      # - AWS_REGION=ap-southeast-1
    depends_on:
      - postgres
      - redis
    volumes:
      - ./cdn:/root/cdn  # Mount local cdn directory

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=topup_db
      - POSTGRES_USER=topup_user
      - POSTGRES_PASSWORD=topup_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

---

## 📋 Implementation Checklist

### Phase 1: Core Backend (Week 1-2)
- [ ] Database schema & migrations
- [ ] Authentication system (JWT)
- [ ] Basic CRUD operations
- [ ] Input validation & error handling
- [ ] API documentation (Swagger)

### Phase 2: Business Logic (Week 3-4)
- [ ] Dynamic game account validation
- [ ] Voucher system implementation
- [ ] Transaction processing flow
- [ ] Payment gateway integration
- [ ] WhatsApp notifications

### Phase 3: Admin Features (Week 5-6)
- [ ] Admin authentication & RBAC
- [ ] Analytics dashboard APIs
- [ ] Admin CRUD operations
- [ ] File upload handling
- [ ] Reporting & export features

### Phase 4: Optimization & Testing (Week 7-8)
- [ ] Performance optimization
- [ ] Caching implementation
- [ ] Comprehensive testing
- [ ] Security audit
- [ ] Load testing
- [ ] Production deployment

### Go Dependencies (go.mod)
```go
module topup-backend

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/go-playground/validator/v10 v10.15.5
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.3
    github.com/go-redis/redis/v8 v8.11.5
    github.com/aws/aws-sdk-go-v2 v1.21.0
    github.com/aws/aws-sdk-go-v2/service/s3 v1.40.0
    github.com/swaggo/gin-swagger v1.6.0
    github.com/swaggo/swag v1.8.12
    golang.org/x/crypto v0.14.0
    github.com/google/uuid v1.3.1
    github.com/joho/godotenv v1.4.0
    github.com/stretchr/testify v1.8.4
)
```

---

**Status**: 📋 **SPECIFICATION COMPLETE**  
**Tech Stack**: Go 1.21+ + Gin + GORM + PostgreSQL + Redis (optional) + AWS S3/Local Storage + JWT  
**Estimated Timeline**: 8 weeks  
**Team Requirement**: 2-3 Go Backend Developers