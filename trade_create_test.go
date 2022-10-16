package tt_pay

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/liaoxxxx/tt_pay/config"
)

// Benchmark for version 1.0, 测试时在TPDomain填写Mock server的URL
func BenchmarkTradeCreate1_0(b *testing.B) {
	conf := config.Config{
		AppId:             "______________", // 支付方分配给业务方的ID，用于获取 签名/验签 的密钥信息
		AppSecret:         "______________", // 支付方密钥
		MerchantId:        "______________", // 支付方分配给业务方的商户编号
		TPDomain:          "______________",
		TPClientTimeoutMs: 6000,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := NewTradeCreateRequest(conf)
		req.Version = "1.0"                                       // 下单接口默认为2.0， 可更改为1.0
		req.OutOrderNo = fmt.Sprintf("%d", time.Now().Unix())     // 此处是随机生成的，使用时请填写您的商户订单号
		req.Uid = "123"                                           // 填写用户在头条的id
		req.TotalAmount = 1                                       // 填写订单金额
		req.Currency = "CNY"                                      // 填写币种，一般均为CNY
		req.Subject = "测试订单"                                      // 填写您的订单名称
		req.Body = "测试订单内容"                                       // 填写您的订单内容
		req.TradeTime = fmt.Sprintf("%d", time.Now().Unix())      // 交易时间，此处自动生成，您也可以根据需求赋值，但必须为Unix时间戳
		req.ValidTime = "36000"                                   // 填写您的订单有效时间（单位：秒）
		req.NotifyUrl = "https://google.com"                      // 填写您的异步通知地址
		req.RiskInfo = `{"ip":"127.0.0.1", "device_id":"122333"}` // 严格json字符串格式
		req.ProductCode = "pay"                                   // 固定值，不要改动
		req.PaymentType = "combine"                               // 固定值，不要改动
		// 支付方式（必填）：可选值：SDK|H5。
		// SDK：业务方App必须是头条主端App，或者具备头条主端支付SDK及ToutiaoJSBridge的能力
		// H5：业务方App不具备SDK支付的能力如果非法，默认为H5支付方式
		req.CashdeskTradeType = "H5"

		// 以下为附加功能，用来设置收银台展示风格和扩展字段
		// 不需要定制展示风格和扩展字段时，请忽略以下部分
		req.SetButtonColor("#F85959")                                                       // 设置按键颜色
		req.SetFontColor("#FFFFFF")                                                         // 设置字体颜色
		req.SetShowLeftTime(true)                                                           // 是否展示剩余时间,默认false
		req.SetCashdeskShowStyle(1)                                                         // 设置收银台展示风格
		req.SetResultPageStyle(1)                                                           // 设置结果页风格
		req.SetFirstDefaultPayType("wx")                                                    // 设置第一默认支付方式
		req.SetExtendParams(`{"enable_pay_channels":"pcredit,moneyFund,debitCardExpress"}`) // 设置支付宝扩展字段和 enable_pay_chanenels, disable_pay_channels
		// 请注意，调用以上Set函数设置部分字段后，再在CashdeskExts中设置的相应字段不会生效
		req.CashdeskExts = `{"123":"123"}` // 设置你需要传的其他扩展参数

		ctx := context.Background()
		TradeCreate(ctx, req)
	}
}

// Benchmark for version 2.0
func BenchmarkTradeCreate2_0(b *testing.B) {
	conf := config.Config{
		AppId:             "______________", // 支付方分配给业务方的ID，用于获取 签名/验签 的密钥信息
		AppSecret:         "______________", // 支付方密钥
		MerchantId:        "______________", // 支付方分配给业务方的商户编号
		TPDomain:          "______________",
		TPClientTimeoutMs: 6000,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := NewTradeCreateRequest(conf)
		req.Version = "2.0"                                       // 下单接口默认为2.0， 可更改为1.0
		req.OutOrderNo = fmt.Sprintf("%d", time.Now().Unix())     // 此处是随机生成的，使用时请填写您的商户订单号
		req.Uid = "123"                                           // 填写用户在头条的id
		req.TotalAmount = 1                                       // 填写订单金额
		req.Currency = "CNY"                                      // 填写币种，一般均为CNY
		req.Subject = "测试订单"                                      // 填写您的订单名称
		req.Body = "测试订单内容"                                       // 填写您的订单内容
		req.TradeTime = fmt.Sprintf("%d", time.Now().Unix())      // 交易时间，此处自动生成，您也可以根据需求赋值，但必须为Unix时间戳
		req.ValidTime = "36000"                                   // 填写您的订单有效时间（单位：秒）
		req.NotifyUrl = "https://google.com"                      // 填写您的异步通知地址
		req.RiskInfo = `{"ip":"127.0.0.1", "device_id":"122333"}` // 严格json字符串格式
		req.ProductCode = "pay"                                   // 固定值，不要改动
		req.PaymentType = "combine"                               // 固定值，不要改动
		// 支付方式（必填）：可选值：SDK|H5。
		// SDK：业务方App必须是头条主端App，或者具备头条主端支付SDK及ToutiaoJSBridge的能力
		// H5：业务方App不具备SDK支付的能力如果非法，默认为H5支付方式
		req.CashdeskTradeType = "SDK"

		// 以下为附加功能，用来设置收银台展示风格和扩展字段
		// 不需要定制展示风格和扩展字段时，请忽略以下部分
		req.SetButtonColor("#F85959")                                                       // 设置按键颜色
		req.SetFontColor("#FFFFFF")                                                         // 设置字体颜色
		req.SetShowLeftTime(true)                                                           // 是否展示剩余时间,默认false
		req.SetCashdeskShowStyle(1)                                                         // 设置收银台展示风格
		req.SetResultPageStyle(1)                                                           // 设置结果页风格
		req.SetFirstDefaultPayType("wx")                                                    // 设置第一默认支付方式
		req.SetExtendParams(`{"enable_pay_channels":"pcredit,moneyFund,debitCardExpress"}`) // 设置支付宝扩展字段和 enable_pay_chanenels, disable_pay_channels
		// 请注意，调用以上Set函数设置部分字段后，再在CashdeskExts中设置的相应字段不会生效
		req.CashdeskExts = `{"123":"123"}` // 设置你需要传的其他扩展参数

		ctx := context.Background()
		TradeCreate(ctx, req)
	}
}
