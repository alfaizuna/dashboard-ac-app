# Dashboard AC Backend - Daftar Perintah CURL

<style>
.copy-btn {
  background: #007acc;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  margin-left: 10px;
  float: right;
}
.copy-btn:hover {
  background: #005a9e;
}
.code-container {
  position: relative;
  margin-bottom: 15px;
}
.copy-success {
  background: #28a745 !important;
}
</style>

<script>
function copyToClipboard(text, button) {
  navigator.clipboard.writeText(text).then(function() {
    button.textContent = 'Copied!';
    button.classList.add('copy-success');
    setTimeout(() => {
      button.textContent = 'Copy';
      button.classList.remove('copy-success');
    }, 2000);
  });
}
</script>

## Base URL
```
http://localhost:8080
```

## 1. Health Check

### Cek Status Server
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/health\" \\
  -H \"Content-Type: application/json\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/health" \
  -H "Content-Type: application/json"
```
</div>

## 2. Authentication (Public Endpoints)

### Register User Baru
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/auth/register\" \\
  -H \"Content-Type: application/json\" \\
  -d '{
    \"name\": \"John Doe\",
    \"email\": \"john@example.com\",
    \"password\": \"password123\",
    \"role\": \"admin\"
  }'`, this)">Copy</button>

```bash
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "role": "admin"
  }'
```
</div>

### Login
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/auth/login\" \\
  -H \"Content-Type: application/json\" \\
  -d '{
    \"email\": \"john@example.com\",
    \"password\": \"password123\"
  }'`, this)">Copy</button>

```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```
</div>

### Refresh Token
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/auth/refresh\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_REFRESH_TOKEN\"`, this)">Copy</button>

```bash

```
</div>

## 3. User Profile (Protected)

### Get Profile Saya
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/me\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/me" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

## 4. User Management (Admin Only)

### Buat User Baru
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/users\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"name\": \"Jane Doe\",
    \"email\": \"jane@example.com\",
    \"password\": \"password123\",
    \"role\": \"technician\"
  }'`, this)">Copy</button>

```bash
curl -X POST "http://localhost:8080/api/v1/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "password123",
    "role": "technician"
  }'
```
</div>

### List Semua Users
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/users?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get User by ID
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/users/USER_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/users/USER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Update User
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X PUT \"http://localhost:8080/api/v1/users/USER_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"name\": \"Jane Smith\",
    \"email\": \"jane.smith@example.com\",
    \"role\": \"admin\"
  }'`, this)">Copy</button>

```bash
curl -X PUT "http://localhost:8080/api/v1/users/USER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Jane Smith",
    "email": "jane.smith@example.com",
    "role": "admin"
  }'
```
</div>

### Delete User
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X DELETE \"http://localhost:8080/api/v1/users/USER_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X DELETE "http://localhost:8080/api/v1/users/USER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Users by Role
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/users/role/admin\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/users/role/admin" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

## 5. Customer Management (Admin & Technician)

### Buat Customer Baru
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/customers\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"name\": \"Customer Name\",
    \"phone\": \"081234567890\",
    \"email\": \"customer@example.com\",
    \"address\": \"Jl. Customer Address No. 123\"
  }'`, this)">Copy</button>

```bash
curl -X POST "http://localhost:8080/api/v1/customers" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Customer Name",
    "phone": "081234567890",
    "email": "customer@example.com",
    "address": "Jl. Customer Address No. 123"
  }'
```
</div>

### List Semua Customers
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/customers?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/customers?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Customer by ID
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/customers/CUSTOMER_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/customers/CUSTOMER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Update Customer
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X PUT \"http://localhost:8080/api/v1/customers/CUSTOMER_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"name\": \"Updated Customer Name\",
    \"phone\": \"081234567891\",
    \"email\": \"updated@example.com\",
    \"address\": \"Updated Address\"
  }'`, this)">Copy</button>

```bash
curl -X PUT "http://localhost:8080/api/v1/customers/CUSTOMER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Updated Customer Name",
    "phone": "081234567891",
    "email": "updated@example.com",
    "address": "Updated Address"
  }'
```
</div>

### Delete Customer
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X DELETE \"http://localhost:8080/api/v1/customers/CUSTOMER_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X DELETE "http://localhost:8080/api/v1/customers/CUSTOMER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Search Customers
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/customers/search?name=John&phone=081&email=example&page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/customers/search?name=John&phone=081&email=example&page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

## 6. Technician Management (Admin Only)

### Buat Technician Baru
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/technicians\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"name\": \"Technician Name\",
    \"phone\": \"081234567890\",
    \"email\": \"tech@example.com\",
    \"specialization\": \"AC Repair\"
  }'`, this)">Copy</button>

```bash
curl -X POST "http://localhost:8080/api/v1/technicians" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Technician Name",
    "phone": "081234567890",
    "email": "tech@example.com",
    "specialization": "AC Repair"
  }'
```
</div>

### List Semua Technicians
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/technicians?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/technicians?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Technician by ID
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/technicians/TECHNICIAN_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/technicians/TECHNICIAN_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Update Technician
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X PUT \"http://localhost:8080/api/v1/technicians/TECHNICIAN_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"name\": \"Updated Technician Name\",
    \"phone\": \"081234567891\",
    \"email\": \"updated.tech@example.com\",
    \"specialization\": \"AC Installation\"
  }'`, this)">Copy</button>

```bash
curl -X PUT "http://localhost:8080/api/v1/technicians/TECHNICIAN_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Updated Technician Name",
    "phone": "081234567891",
    "email": "updated.tech@example.com",
    "specialization": "AC Installation"
  }'
```
</div>

### Delete Technician
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X DELETE \"http://localhost:8080/api/v1/technicians/TECHNICIAN_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X DELETE "http://localhost:8080/api/v1/technicians/TECHNICIAN_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Search Technicians
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/technicians/search?name=John&specialization=AC&page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/technicians/search?name=John&specialization=AC&page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

## 7. Service Management (Admin Only)

### Buat Service Baru
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/services\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"name\": \"AC Cleaning\",
    \"description\": \"Complete AC cleaning service\",
    \"price\": 150000
  }'`, this)">Copy</button>

```bash
curl -X POST "http://localhost:8080/api/v1/services" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "AC Cleaning",
    "description": "Complete AC cleaning service",
    "price": 150000
  }'
```
</div>

### List Semua Services
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/services?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/services?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Service by ID
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/services/SERVICE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/services/SERVICE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Update Service
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X PUT \"http://localhost:8080/api/v1/services/SERVICE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"name\": \"AC Deep Cleaning\",
    \"description\": \"Complete AC deep cleaning service\",
    \"price\": 200000
  }'`, this)">Copy</button>

```bash
curl -X PUT "http://localhost:8080/api/v1/services/SERVICE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "AC Deep Cleaning",
    "description": "Complete AC deep cleaning service",
    "price": 200000
  }'
```
</div>

### Delete Service
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`ccurl -X POST "http://localhost:8080/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_REFRESH_TOKEN"url -X DELETE \"http://localhost:8080/api/v1/services/SERVICE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X DELETE "http://localhost:8080/api/v1/services/SERVICE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Search Services
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/services/search?name=cleaning&min_price=100000&max_price=300000&page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/services/search?name=cleaning&min_price=100000&max_price=300000&page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

## 8. Schedule Management (Admin & Technician)

### Buat Schedule Baru
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/schedules\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"customer_id\": \"CUSTOMER_ID\",
    \"technician_id\": \"TECHNICIAN_ID\",
    \"service_id\": \"SERVICE_ID\",
    \"date\": \"2024-01-15\",
    \"time\": \"10:00\",
    \"notes\": \"Regular maintenance\"
  }'`, this)">Copy</button>

```bash
curl -X POST "http://localhost:8080/api/v1/schedules" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "customer_id": "CUSTOMER_ID",
    "technician_id": "TECHNICIAN_ID",
    "service_id": "SERVICE_ID",
    "date": "2024-01-15",
    "time": "10:00",
    "notes": "Regular maintenance"
  }'
```
</div>

### List Semua Schedules
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/schedules?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/schedules?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Schedule by ID
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/schedules/SCHEDULE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/schedules/SCHEDULE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Update Schedule
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X PUT \"http://localhost:8080/api/v1/schedules/SCHEDULE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"date\": \"2024-01-16\",
    \"time\": \"14:00\",
    \"status\": \"completed\",
    \"notes\": \"Service completed successfully\"
  }'`, this)">Copy</button>

```bash
curl -X PUT "http://localhost:8080/api/v1/schedules/SCHEDULE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "date": "2024-01-16",
    "time": "14:00",
    "status": "completed",
    "notes": "Service completed successfully"
  }'
```
</div>

### Delete Schedule
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X DELETE \"http://localhost:8080/api/v1/schedules/SCHEDULE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X DELETE "http://localhost:8080/api/v1/schedules/SCHEDULE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Search Schedules
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/schedules/search?customer_id=CUSTOMER_ID&technician_id=TECHNICIAN_ID&service_id=SERVICE_ID&status=pending&date_from=2024-01-01&date_to=2024-01-31&page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/schedules/search?customer_id=CUSTOMER_ID&technician_id=TECHNICIAN_ID&service_id=SERVICE_ID&status=pending&date_from=2024-01-01&date_to=2024-01-31&page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Schedules by Customer
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/schedules/customer/CUSTOMER_ID?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/schedules/customer/CUSTOMER_ID?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Schedules by Technician
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/schedules/technician/TECHNICIAN_ID?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/schedules/technician/TECHNICIAN_ID?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Schedules by Status
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/schedules/status/pending?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/schedules/status/pending?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

## 9. Invoice Management (Admin & Technician)

### Buat Invoice Baru
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/invoices\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"customer_id\": \"CUSTOMER_ID\",
    \"schedule_id\": \"SCHEDULE_ID\",
    \"total_amount\": 150000,
    \"notes\": \"Payment for AC cleaning service\"
  }'`, this)">Copy</button>

```bash
curl -X POST "http://localhost:8080/api/v1/invoices" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "customer_id": "CUSTOMER_ID",
    "schedule_id": "SCHEDULE_ID",
    "total_amount": 150000,
    "notes": "Payment for AC cleaning service"
  }'
```
</div>

### List Semua Invoices
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/invoices?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/invoices?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Invoice by ID
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/invoices/INVOICE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/invoices/INVOICE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Update Invoice
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X PUT \"http://localhost:8080/api/v1/invoices/INVOICE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"total_amount\": 175000,
    \"status\": \"paid\",
    \"notes\": \"Payment completed\"
  }'`, this)">Copy</button>

```bash
curl -X PUT "http://localhost:8080/api/v1/invoices/INVOICE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "total_amount": 175000,
    "status": "paid",
    "notes": "Payment completed"
  }'
```
</div>

### Delete Invoice
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X DELETE \"http://localhost:8080/api/v1/invoices/INVOICE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X DELETE "http://localhost:8080/api/v1/invoices/INVOICE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Search Invoices
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/invoices/search?customer_id=CUSTOMER_ID&status=pending&page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/invoices/search?customer_id=CUSTOMER_ID&status=pending&page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Invoices by Customer
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/invoices/customer/CUSTOMER_ID?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/invoices/customer/CUSTOMER_ID?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Invoices by Schedule
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/invoices/schedule/SCHEDULE_ID?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/invoices/schedule/SCHEDULE_ID?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Invoices by Status
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/invoices/status/paid?page=1&limit=10\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/invoices/status/paid?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

## 10. Invoice Detail Management (Admin & Technician)

### Buat Invoice Detail Baru
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X POST \"http://localhost:8080/api/v1/invoice-details\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"invoice_id\": \"INVOICE_ID\",
    \"service_id\": \"SERVICE_ID\",
    \"quantity\": 1,
    \"unit_price\": 150000,
    \"total_price\": 150000
  }'`, this)">Copy</button>

```bash
curl -X POST "http://localhost:8080/api/v1/invoice-details" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "invoice_id": "INVOICE_ID",
    "service_id": "SERVICE_ID",
    "quantity": 1,
    "unit_price": 150000,
    "total_price": 150000
  }'
```
</div>

### Get Invoice Detail by ID
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/invoice-details/INVOICE_DETAIL_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/invoice-details/INVOICE_DETAIL_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Update Invoice Detail
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X PUT \"http://localhost:8080/api/v1/invoice-details/INVOICE_DETAIL_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" \\
  -d '{
    \"quantity\": 2,
    \"unit_price\": 150000,
    \"total_price\": 300000
  }'`, this)">Copy</button>

```bash
curl -X PUT "http://localhost:8080/api/v1/invoice-details/INVOICE_DETAIL_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "quantity": 2,
    "unit_price": 150000,
    "total_price": 300000
  }'
```
</div>

### Delete Invoice Detail
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X DELETE \"http://localhost:8080/api/v1/invoice-details/INVOICE_DETAIL_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X DELETE "http://localhost:8080/api/v1/invoice-details/INVOICE_DETAIL_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Get Invoice Details by Invoice
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X GET \"http://localhost:8080/api/v1/invoice-details/invoice/INVOICE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X GET "http://localhost:8080/api/v1/invoice-details/invoice/INVOICE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

### Delete All Invoice Details by Invoice
<div class="code-container">
<button class="copy-btn" onclick="copyToClipboard(`curl -X DELETE \"http://localhost:8080/api/v1/invoice-details/invoice/INVOICE_ID\" \\
  -H \"Content-Type: application/json\" \\
  -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\"`, this)">Copy</button>

```bash
curl -X DELETE "http://localhost:8080/api/v1/invoice-details/invoice/INVOICE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```
</div>

## Catatan Penting:

1. **Authorization**: Ganti `YOUR_ACCESS_TOKEN` dengan token yang didapat dari endpoint login
2. **IDs**: Ganti placeholder seperti `USER_ID`, `CUSTOMER_ID`, dll. dengan ID yang sebenarnya
3. **Roles**: 
   - Admin: Akses penuh ke semua endpoint
   - Technician: Akses terbatas (tidak bisa manage users, technicians, services)
4. **Pagination**: Semua endpoint list mendukung parameter `page` dan `limit`
5. **Status Values**: 
   - Schedule: `pending`, `in_progress`, `completed`, `cancelled`
   - Invoice: `pending`, `paid`, `cancelled`

## Contoh Workflow Lengkap:

1. Register/Login untuk mendapatkan token
2. Buat customer, technician, dan service (admin only)
3. Buat schedule untuk customer
4. Buat invoice berdasarkan schedule
5. Tambahkan invoice details
6. Update status sesuai kebutuhan