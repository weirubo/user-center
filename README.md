# User Center

Go 语言实现的用户中心服务，基于 Kratos 框架。

## 功能

- 用户注册（密码/验证码）
- 用户登录（密码/验证码）
- JWT 认证
- 邮箱验证码
- OAuth 登录（WeChat, GitHub, Google）
- 密码修改
- 账号删除
- 密码锁定（5次失败后锁定5分钟）

## 技术栈

- Go 1.26
- Kratos
- MySQL
- Redis

## 快速开始

```bash
# 下载依赖
go mod download

# 运行服务
go run cmd/server/main.go
```

服务默认监听：
- HTTP: http://localhost:8000
- gRPC: localhost:9000

## 配置

配置文件位于 `configs/config.yaml`

```yaml
server:
  http:
    network: ":8000"
  grpc:
    network: ":9000"

database:
  source: "root:123456@tcp(127.0.0.1:3306)/user_center?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

auth:
  jwt_secret: "your-secret-key"
```

## API 文档

### 1. 用户注册

**请求**
```
POST /api/v1/register
```

**Request Body (密码模式)**
```json
{
  "email": "user@example.com",
  "password": "123456",
  "nickname": "username"
}
```

**Request Body (验证码模式)**
```json
{
  "email": "user@example.com",
  "password": "123456",
  "code": "123456",
  "nickname": "username"
}
```

**Response (成功)**
```json
{
  "id": 1,
  "message": "register success"
}
```

---

### 2. 用户登录

**请求**
```
POST /api/v1/login
```

**Request Body (密码模式)**
```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

**Request Body (验证码模式)**
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

**Response (成功)**
```json
{
  "id": 1,
  "nickname": "username",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (失败)**
```json
{
  "message": "invalid credentials"
}
```

---

### 3. 发送验证码

**请求**
```
POST /api/v1/verifycode/send
```

**Request Body**
```json
{
  "email": "user@example.com"
}
```

**Response (成功)**
```json
{
  "message": "verification code sent"
}
```

---

### 4. 获取用户信息

**请求**
```
GET /api/v1/userinfo
```

**Request Headers**
```
Authorization: Bearer <token>
```

**Response (成功)**
```json
{
  "id": 1,
  "email": "user@example.com",
  "phone": "",
  "nickname": "username"
}
```

**Response (失败)**
```
unauthorized (401)
```

---

### 5. 修改密码

**请求**
```
POST /api/v1/password/change
```

**Request Headers**
```
Authorization: Bearer <token>
```

**Request Body**
```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

**Response (成功)**
```json
{
  "message": "password changed successfully"
}
```

---

### 6. 删除账号

**请求**
```
POST /api/v1/account/delete
```

**Request Headers**
```
Authorization: Bearer <token>
```

**Response (成功)**
```json
{
  "message": "account deleted successfully"
}
```

---

### 7. OAuth 登录

**请求**
```
POST /api/v1/oauth/wechat
POST /api/v1/oauth/github
POST /api/v1/oauth/google
```

**Request Body**
```json
{
  "code": "oauth_code"
}
```

**Response (成功)**
```json
{
  "id": 1,
  "nickname": "username",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## 错误码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（Token 无效） |
| 500 | 服务器内部错误 |

## License

MIT
