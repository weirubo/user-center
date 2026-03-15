# User Center

Go 语言实现的用户中心服务，基于 Kratos 框架。

## 功能

- 用户注册
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

data:
  database:
    source: "root:123456@tcp(127.0.0.1:3306)/user_center?charset=utf8mb4&parseTime=True&loc=Local"
  redis:
    addr: "127.0.0.1:6379"

auth:
  jwt:
    secret: "your-secret-key"
    expire: 604800  # 7天
```

## API 文档

### 用户注册

```
POST /api/v1/register
```

请求体：
```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

### 用户登录

```
POST /api/v1/login
```

请求体：
```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

### 验证码登录

```
POST /api/v1/login/code
```

请求体：
```json
{
  "email": "user@example.com"
}
```

### 发送验证码

```
POST /api/v1/verifycode/send
```

请求体：
```json
{
  "email": "user@example.com"
}
```

### 修改密码

```
POST /api/v1/password/change
```

请求体：
```json
{
  "email": "user@example.com",
  "old_password": "old123456",
  "new_password": "new123456"
}
```

### 删除账号

```
DELETE /api/v1/account
```

请求头：
```
Authorization: Bearer <token>
```

### OAuth 登录

```
POST /api/v1/oauth/wechat
POST /api/v1/oauth/github  
POST /api/v1/oauth/google
```

## License

MIT
