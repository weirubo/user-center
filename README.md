# User Center

Go 语言实现的用户中心服务，基于 Kratos 框架。

## 功能

- 用户注册
- 用户登录
- JWT 认证
- 邮箱验证
- OAuth 登录（WeChat, GitHub, Google）
- 密码管理
- 账号删除

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

## 配置

配置文件位于 `configs/config.yaml`

## API 文档

### 注册

```
POST /api/v1/register
```

### 登录

```
POST /api/v1/login
```

### 验证码

```
POST /api/v1/verifycode/send
```

## License

MIT
