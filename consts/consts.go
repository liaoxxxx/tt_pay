package consts

const (
	MethodTradeQuery     = "tp.trade.query"
	MethodRefundCreate   = "tp.refund.create"
	MethodRefundQuery    = "tp.refund.query"
	MethodTradeCreate    = "tp.trade.create"
	MethodTradeConfirm   = "tp.trade.confirm"
	MethodWithdrawCreate = "tp.withdraw.create"
	MethodWithdrawQuery  = "tp.withdraw.query"

	TPDomain = "https://tp-pay.snssdk.com"
	TPPath   = "gateway"
	TPPathU  = "gateway-u"

	TtPayPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDOZZ7iAkS3oN970+yDONe5TPhP
rLHoNOZOjJjackEtgbptdy4PYGBGdeAUAz75TO7YUGESCM+JbyOz1YzkMfKl2HwY
doePEe8qzfk5CPq6VAhYJjDFA/M+BAZ6gppWTjKnwMcHVK4l2qiepKmsw6bwf/kk
LTV9l13r6Iq5U+vrmwIDAQAB
-----END PUBLIC KEY-----`
)

//// 以下为参数列表
//var (
//	TradeNotifyRespParams = map[string]bool{
//		"notify_id":    true,
//		"sign_type":    true,
//		"sign":         true,
//		"app_id":       true,
//		"event_code":   true,
//		"out_order_no": true,
//		"trade_no":     true,
//		"total_amount": true,
//		"pay_channel":  true,
//		"merchant_id":  true,
//		"pay_time":     true,
//		"pay_type":     true,
//		"trade_status": true,
//		"trade_msg":    true,
//		"extension":    true,
//	}
//)
