package config

type Config struct {
	AppId             string
	AppSecret         string
	MerchantId        string
	TPDomain          string // 请求支付域名 加http或者https前缀，比如：https://tp-pay.snssdk.com
	TPClientTimeoutMs int
}

type TTPay struct {
	AppID       string `json:"app_id"`
	Body        string `json:"body"`
	Currency    string `json:"currency"`
	LimitPay    string `json:"limit_pay"`
	MerchantID  string `json:"merchant_id"`
	NotifyURL   string `json:"notify_url"`
	OutOrderNo  string `json:"out_order_no"`
	PaymentType string `json:"payment_type"`
	ProductCode string `json:"product_code"`
	RiskInfo    string `json:"risk_info"`
	Sign        string `json:"sign"`
	SignType    string `json:"sign_type"`
	Subject     string `json:"subject"`
	Timestamp   string `json:"timestamp"`
	TotalAmount string `json:"total_amount"`
	TradeTime   string `json:"trade_time"`
	TradeType   string `json:"trade_type"`
	UID         string `json:"uid"`
	ValidTime   string `json:"valid_time"`
	Version     string `json:"version"`
	WxType      string `json:"wx_type"`
	WxURL       string `json:"wx_url"`
}
