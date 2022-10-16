package tt_pay

import (
	"context"
	"errors"
	"net/url"

	"github.com/liaoxxxx/tt_pay/consts"
	"github.com/liaoxxxx/tt_pay/util"
)

// 下单回调请求
type TradeNotifyRequest struct {
	Param string
}

// 下单回调接口
func TradeNotify(ctx context.Context, req *TradeNotifyRequest) (*TradeNotifyResponse, error) {
	// 解析回调参数
	params, err := url.ParseQuery(req.Param)
	if err != nil {
		util.Debug("TradeNotify failed: params: [%s] err: [%s]", req.Param, err)
		return nil, util.Wrap(err, "TradeNotify failed when [ParseQuery()]")
	}

	resp := new(TradeNotifyResponse)
	resp.Param = make(map[string]string)
	signMap := make(map[string]interface{})
	for key, val := range params {
		resp.Param[key] = val[0]
		signMap[key] = interface{}(val[0])
	}

	sign := resp.Get("sign")

	if valid := util.VerifyMd5WithRsa(signMap, sign, consts.TtPayPublicKey); !valid {
		return nil, errors.New("Invalid sign")
	}

	resp.Decode()

	return resp, nil
}

// SetParam 将回调的param参数赋值给该实例成员变量
func (req *TradeNotifyRequest) SetParam(s string) {
	req.Param = s
}

// 下单回调响应
type TradeNotifyResponse struct {
	Param       map[string]string
	NotifyId    string
	SignType    string
	Sign        string
	AppId       string
	EventCode   string
	MerchantId  string
	OutOrderNo  string
	TradeNo     string
	TotalAmount string
	PayChannel  string
	PayTime     string
	PayType     string
	TradeStatus string
	TradeMsg    string
	Extension   string `json:"extension"`
}

// 解析响应中的参数
func (resp *TradeNotifyResponse) Decode() {
	resp.NotifyId = resp.Get("notify_id")
	resp.SignType = resp.Get("sign_type")
	resp.Sign = resp.Get("sign")
	resp.AppId = resp.Get("app_id")
	resp.EventCode = resp.Get("event_code")
	resp.OutOrderNo = resp.Get("out_order_no")
	resp.TradeNo = resp.Get("trade_no")
	resp.TotalAmount = resp.Get("total_amount")
	resp.PayChannel = resp.Get("pay_channel")
	resp.MerchantId = resp.Get("merchant_id")
	resp.PayTime = resp.Get("pay_time")
	resp.PayType = resp.Get("pay_type")
	resp.TradeStatus = resp.Get("trade_status")
	resp.TradeMsg = resp.Get("trade_msg")
	resp.Extension = resp.Get("extension")
}

// 提取Param内的值
func (resp *TradeNotifyResponse) Get(key string) string {
	return resp.Param[key]
}
