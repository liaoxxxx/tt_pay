package tt_pay

import (
	"context"
	"errors"
	"net/url"

	"github.com/zoooozz/tt_pay/consts"
	"github.com/zoooozz/tt_pay/util"
)

// 退款回调请求
type RefundNotifyRequest struct {
	Param string
}

// 退款回调接口
func RefundNotify(ctx context.Context, req *RefundNotifyRequest) (*RefundNotifyResponse, error) {
	params, err := url.ParseQuery(req.Param)

	if err != nil {
		util.Debug("RefundNotify failed: params: [%s] err: [%s]", req.Param, err)
		return nil, util.Wrap(err, "RefundNotify failed when [ParseQuery()]")
	}

	resp := new(RefundNotifyResponse)
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

// 将回调的param参数赋值给该实例成员变量
func (req *RefundNotifyRequest) SetParam(s string) {
	req.Param = s
}

// 退款回调响应
type RefundNotifyResponse struct {
	Param        map[string]string
	NotifyId     string
	SignType     string
	Sign         string
	AppId        string
	EventCode    string
	OutRefundNo  string
	RefundNo     string
	RefundAmount string
	RefundTime   string
	MerchantId   string
	RefundStatus string
}

// 解析响应中的参数
func (resp *RefundNotifyResponse) Decode() {
	resp.NotifyId = resp.Get("notify_id")
	resp.SignType = resp.Get("sign_type")
	resp.Sign = resp.Get("sign")
	resp.AppId = resp.Get("app_id")
	resp.EventCode = resp.Get("event_code")
	resp.OutRefundNo = resp.Get("out_refund_no")
	resp.RefundNo = resp.Get("refund_no")
	resp.RefundAmount = resp.Get("refund_amount")
	resp.RefundTime = resp.Get("refund_time")
	resp.MerchantId = resp.Get("merchant_id")
	resp.RefundStatus = resp.Get("refund_status")
}

// 提取Param内的值
func (resp *RefundNotifyResponse) Get(key string) string {
	return resp.Param[key]
}
