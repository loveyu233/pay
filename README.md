# Go 支付服务

这是一个使用 Go 语言编写的支付服务库，集成了微信支付和支付宝支付功能。它被设计为可轻松集成到现有的 Go Web 项目中（尤其是使用 Gin 框架的项目）。

## ✨ 功能特性

- **微信支付 (WeChat Pay)**
  - JSAPI 交易下单
  - 交易退款
  - 支付成功异步回调通知处理
  - 退款成功异步回调通知处理
  - 订单查询（支付与退款）
  - 基于 `ArtisanCloud/PowerWeChat` 库

- **支付宝 (Alipay)**
  - 小程序用户授权与登录
  - 交易下单 (`TradeCreate`)
  - 交易退款
  - 支付与退款异步回调通知处理
  - 订单查询（支付与退款）
  - 基于 `go-pay/gopay` 库

- **框架集成**
  - 提供 Gin 框架的 `RouterGroup` 注册，方便快速集成 API 路由。

## 📂 项目结构

```
.
├── go.mod
├── go.sum
├── payment.go          # 微信支付的核心结构与初始化
├── patment_server.go   # 微信支付的 Gin API 路由与处理逻辑 (文件名疑似拼写错误, 应为 payment_server.go)
├── zfb_server.go       # 支付宝支付的 Gin API 路由与处理逻辑
└── README.md           # 本文档
```

## 🚀 快速开始

### 1. 安装依赖

```bash
go get github.com/ArtisanCloud/PowerWeChat/v3
go get github.com/go-pay/gopay
go get github.com/gin-gonic/gin
go get github.com/loveyu233/gb
```

### 2. 初始化与配置

#### 微信支付

你需要先实现 `WXPayImp` 接口，然后调用 `InitWXWXPaymentApp` 来初始化一个微信支付客户端。

```go
import "your/path/to/pay"

// 1. 实现 WXPayImp 接口
type MyWXPayHandler struct{}

func (h *MyWXPayHandler) PayNotify(orderId string, attach string) error {
    // 处理支付成功逻辑，例如更新订单状态
    return nil
}
func (h *MyWXPayHandler) RefundNotify(orderId string) error {
    // 处理退款成功逻辑
    return nil
}
func (h *MyWXPayHandler) Pay(c *gin.Context) (*pay.PayRequest, error) {
    // 从请求中解析并返回支付参数
    var req pay.PayRequest
    // ... 解析逻辑 ...
    return &req, nil
}
func (h *MyWXPayHandler) Refund(c *gin.Context) (*pay.RefundRequest, error) {
    // 从请求中解析并返回退款参数
    var req pay.RefundRequest
    // ... 解析逻辑 ...
    return &req, nil
}


// 2. 配置并初始化
wxPayConfig := pay.WXPaymentAppConfig{
    Payment: pay.Payment{
        AppID:       "你的微信 AppID",
        MchID:       "你的商户号",
        MchApiV3Key: "你的 APIv3 密钥",
        // ... 其他证书和配置
    },
    WXPayImp: &MyWXPayHandler{},
}

wxPayClient, err := pay.InitWXWXPaymentApp(wxPayConfig)
if err != nil {
    // 处理错误
}

// 3. 注册到 Gin 路由
router := gin.Default()
wxGroup := router.Group("/api")
wxPayClient.RegisterHandlers(wxGroup)
```

#### 支付宝

你需要先实现 `ZfbMiniImp` 接口，然后调用 `InitAliClient` 来初始化支付宝客户端。

```go
import "your/path/to/pay"

// 1. 实现 ZfbMiniImp 接口
type MyZFBPayHandler struct{}

// ... 实现接口的所有方法，例如 Pay, Refund, PayNotify 等 ...

func (h *MyZFBPayHandler) Pay(c *gin.Context) (*pay.ZFBPayParam, error) {
    // 从请求中解析并返回支付参数
    var param pay.ZFBPayParam
    // ... 解析逻辑 ...
    return &param, nil
}


// 2. 初始化
err := pay.InitAliClient(
    "你的支付宝 AppID",
    "你的应用私钥",
    "你的 AES 密钥",
    "应用公钥证书路径",
    "支付宝公钥证书路径",
    "支付宝根证书路径",
    "异步通知回调 URL",
    true, // 是否保存日志
    &MyZFBPayHandler{},
)
if err != nil {
    // 处理错误
}

// 3. 注册到 Gin 路由
router := gin.Default()
zfbGroup := router.Group("/api")
pay.InsZFB.RegisterHandlers(zfbGroup)

```

### 3. 运行

你可以将此项目作为库导入到你的主应用中，或者直接运行（如果包含 `main` 函数）。

```bash
# 编译
go build

# 运行 (假设你的主文件是 main.go)
go run main.go
```

## 📦 API 端点

### 微信支付 (前缀: `/wx`)

- `POST /pay`: 创建支付订单
- `POST /refund`: 发起退款
- `POST /notify/payment`: 支付异步回调
- `POST /notify/refund`: 退款异步回调

### 支付宝 (前缀: `/zfb`)

- `POST /login`: 小程序登录/授权
- `POST /pay`: 创建支付订单
- `POST /refund`: 发起退款
- `POST /notify`: 支付/退款异步回调

---
*该 README 文件由 AI 根据项目代码自动生成。*
