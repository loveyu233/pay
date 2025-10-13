# Go Pay Server

这是一个基于 Go 语言的支付服务器，为微信和支付宝支付提供统一的接口。它使用 Gin 框架构建，并集成了流行的支付 SDK。

## 功能

- **微信支付:**
  - 小程序支付
  - 公众号支付
  - 支付和退款通知
  - 订单查询
- **支付宝:**
  - 小程序支付
  - 支付和退款通知
  - 订单查询
- **微信小程序:**
  - 用户登录和会话管理
  - 手机号解密
  - 二维码生成
- **微信公众号:**
  - 事件处理（例如，关注）
  - 消息处理
  - 模板消息推送

## 依赖

- [gin-gonic/gin](https://github.com/gin-gonic/gin): HTTP Web 框架。
- [ArtisanCloud/PowerWeChat/v3](https://github.com/ArtisanCloud/PowerWeChat): 适用于 Go 的微信 SDK。
- [go-pay/gopay](https://github.com/go-pay/gopay): 适用于 Go 的支付宝和微信支付 SDK。
- [loveyu233/gb](https://github.com/loveyu233/gb): Go 的实用程序库。

## 安装

1.  克隆存储库：
    ```bash
    git clone <repository-url>
    ```
2.  安装依赖：
    ```bash
    go mod tidy
    ```

## 配置

您需要通过提供微信和支付宝的必要凭据来配置应用程序。

### 微信

微信的配置由 `InitWXMiniProgramService`、`InitWXOfficialAccountAppService` 和 `InitWXWXPaymentApp` 函数处理。您需要为您的微信小程序、公众号和支付帐户提供 AppID、Secret 和其他相关详细信息。

### 支付宝

支付宝的配置由 `InitAliClient` 函数处理。您需要提供您的 AppID、私钥、AES 密钥以及您的公钥证书的路径。

## 用法

1.  使用您的配置初始化服务。
2.  为 Gin 路由器注册处理程序。
3.  运行 Gin 服务器。

```go
package main

import (
	"github.com/gin-gonic/gin"
	"your-project/pay"
)

func main() {
	// ... 您的配置设置 ...

	// 初始化微信和支付宝服务
	// pay.InitWXMiniProgramService(...)
	// pay.InitWXOfficialAccountAppService(...)
	// pay.InitWXWXPaymentApp(...)
	// pay.InitAliClient(...)

	r := gin.Default()
	apiGroup := r.Group("/api")

	// 注册处理程序
	pay.InsWX.WXMini.RegisterHandlers(apiGroup)
	pay.InsWX.WXOfficial.RegisterHandlers(apiGroup)
	pay.InsWX.WXPay.RegisterHandlers(apiGroup)
	pay.InsZFB.RegisterHandlers(apiGroup)

	r.Run(":8080")
}
```

## API 端点

### 微信小程序

- `POST /wx/login`: 小程序登录。

### 微信公众号

- `GET /wx/callback`: 回调验证。
- `POST /wx/callback`: 接收消息和事件。
- `POST /wx/push`: 推送消息。

### 微信支付

- `POST /wx/notify/payment`: 支付通知回调。
- `POST /wx/notify/refund`: 退款通知回调。
- `POST /wx/pay`: 创建支付。
- `POST /wx/refund`: 创建退款。

### 支付宝

- `POST /zfb/login`: 支付宝小程序登录。
- `POST /zfb/notify`: 支付和退款通知回调。
- `POST /zfb/pay`: 创建支付。
- `POST /zfb/refund`: 创建退款。