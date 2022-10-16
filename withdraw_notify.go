package tt_pay

import (
	"context"
	"errors"
	"net/url"

	"github.com/liaoxxxx/tt_pay/consts"
	"github.com/liaoxxxx/tt_pay/util"
)

// 提现回调接口
func WithdrawNotify(ctx context.Context, req *WithdrawNotifyRequest) (*WithdrawNotifyResponse, error) {
	params, err := url.ParseQuery(req.Param)

	if err != nil {
		util.Debug("Parse params failed: err[%s]", err)
		return nil, err
	}

	resp := new(WithdrawNotifyResponse)
	resp.Param = make(map[string]string)
	signMap := make(map[string]interface{})

	for key, val := range params {
		resp.Param[key] = val[0]
		signMap[key] = interface{}(val[0])
	}

	resp.Decode()

	sign := resp.Get("sign")

	if valid := util.VerifyMd5WithRsa(signMap, sign, consts.TtPayPublicKey); !valid {
		return nil, errors.New("Invalid sign")
	}

	return resp, nil
}

// 提现回调请求
type WithdrawNotifyRequest struct {
	Param string
}

// SetParam 将回调的param参数赋值给该实例成员变量
func (req *WithdrawNotifyRequest) SetParam(s string) {
	req.Param = s
}

type WithdrawNotifyResponse struct {
	Param           map[string]string
	NotifyId        string
	SignType        string
	Sign            string
	EventCode       string
	MerchantId      string
	OutTradeNo      string
	WithdrawTradeNo string
	Amount          string
	WithdrawTime    string
	WithdrawStatus  string
	TradeMsg        string
	Extension       string `json:"extension"`
}

func (resp *WithdrawNotifyResponse) Decode() {
	resp.NotifyId = resp.Get("notify_id")
	resp.SignType = resp.Get("sign_type")
	resp.Sign = resp.Get("sign")
	resp.EventCode = resp.Get("event_code")
	resp.MerchantId = resp.Get("merchant_id")
	resp.OutTradeNo = resp.Get("out_trade_no")
	resp.WithdrawTradeNo = resp.Get("withdraw_trade_no")
	resp.Amount = resp.Get("amount")
	resp.WithdrawTime = resp.Get("withdraw_time")
	resp.WithdrawStatus = resp.Get("withdraw_status")
	resp.TradeMsg = resp.Get("trade_msg")
	resp.Extension = resp.Get("extension")
}

// 设置原始响应
func (resp *WithdrawNotifyResponse) Get(key string) string {
	return resp.Param[key]
}
