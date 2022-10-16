package tt_pay

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/liaoxxxx/tt_pay/util"
)

// 本SDK中的默认Client配置
var httpClient = http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:        2000,
		MaxIdleConnsPerHost: 2000,
		IdleConnTimeout:     65 * time.Second,
		TLSHandshakeTimeout: 3 * time.Second,
	},
}

// 可用此函数自定义HttpClient
func SetHttpClient(c http.Client) {
	httpClient = c
	log.Printf("Set HttpClient to: %v", c)
}

// TPRequest接口
// Encode: 将Request编码成POST请求的Body
type TPRequest interface {
	Encode() (string, error)
	GetUrl() string
	GetLogId() string
}

// TPResponse接口
// Decode：从收到的JSON格式响应中解析参数
type TPResponse interface {
	SetData(data *simplejson.Json)
	Decode() error
}

// 执行请求
func Execute(ctx context.Context, timeout int, req TPRequest, resp TPResponse) error {
	if timeout <= 0 {
		return errors.New("ClientTimeout must be a positive number")
	}

	body, err := req.Encode()
	if err != nil {
		return util.Wrap(err, "Execute failed when [TPRequest.Encode()]")
	}

	logId := req.GetLogId()
	if _logId, ok := ctx.Value("K_LOGID").(string); ok {
		logId = _logId
	}

	statusCode, respBytes, err := HttpPost(req.GetUrl(), "application/x-www-form-urlencoded", body, logId, timeout)
	if err != nil {
		return util.Wrap(err, "Execute failed when [HttpPost()]")
	}

	util.Debug("statusCode[%v] resp[%s]", statusCode, string(respBytes))

	respJson, err := simplejson.NewJson(respBytes)
	if err != nil {
		return util.Wrap(err, "Execute failed when [simplejson.NewJson()]")
	}

	// 判定此次请求是否成功
	// 当一次请求进行到这里时，说明已经与财经后端建立了网络连接并进行了一次成功交互，但该次请求可能成功也可能失败
	// 这里将网络连接成功但请求失败的情况也当做error处理
	if err := success(respJson, req); err != nil {
		return err
	}

	resp.SetData(respJson)
	if err := resp.Decode(); err != nil {
		return util.Wrap(err, "Execute failed when [HttpPost()]")
	}

	return nil
}

func HttpPost(url, contentType, body string, logId string, timeoutMs int) (cnt int, respBytes []byte, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		util.Debug("HttpTypePost NewRequest url[%s] body[%s] err[%s]\n", url, body, err)
		return 0, nil, util.Wrap(err, "HttpPost failed when [http.NewRequest()]")
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)
	req.Header.Set("Content-type", contentType)
	req.Header.Set("X-Tt-Logid", logId)

	resp, err := httpClient.Do(req)
	if err != nil {
		util.Debug("HttpPost client.Do err: %v, url: %s\n", err, url)
		return 0, nil, util.Wrap(err, "HttpPost failed when [client.Do()]")
	}
	// 如果关闭Body失败，将错误信息打印到log中
	// 这里考虑下出现error要不要返回以及如何handle
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			util.Debug("HttpPost resp.Body.Close failed err: %v", cerr)
		}
	}()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Debug("HttpPost ioutil.ReadAll err: %v, url: %s\n", err, url)
		return resp.StatusCode, nil, util.Wrap(err, "HttpPost failed when [ioutil.ReadAll()]")
	}

	util.Debug("HttpTypePost url[%s] contentType[%s] body[%s] code[%d] resp body[%s]\n",
		url, contentType, body, resp.StatusCode, string(respBody))
	return resp.StatusCode, respBody, nil
}

// 提取响应中的错误信息，如果本次请求成功，则返回nil
// 走网关的接口和二维码下单接口返回参数格式不一样，采用switch'区分
// 走网关的接口返回参数格式为：
//		{
//			"response": {
//				"code": "20000",
//				"msg": "Service Currently Unavailable",
//				"sub_code": "TP.SYSTEM_ERROR",
//				"sub_msg": "接口返回错误"
//			},
//			"sign": "ERITJKEIJKJHKKKKKKKHJEREEEEEEEEEEE"
//		}
// 二维码下单接口返回参数格式为：
//		{
//			"data": {
//				"url": "https://tp-pay.snssdk.com/cashdesk/openapi/qrcode?trade_no=xxxxxxxxxxx"
//				"trade_no": "xxxxxxxxxxx",
////		},
//			"code": 0,
//			"msg": "",
//		}
func success(respJson *simplejson.Json, req TPRequest) error {
	switch respJson.Get("response").Interface() {
	case nil:
		// 默认值设为-1，不与其他返回码冲突
		if respJson.Get("code").MustInt(-1) != 0 {
			ret := new(util.Error)
			ret.Code = respJson.Get("code").MustString("")
			ret.Msg = respJson.Get("msg").MustString("")
			ret.Detail = "log_id:" + req.GetLogId()
			return ret
		}
	default:
		if respJson.Get("response").Get("code").MustString("") != "10000" {
			ret := new(util.Error)
			ret.Code = respJson.Get("response").Get("code").MustString("")
			ret.Msg = respJson.Get("response").Get("msg").MustString("")
			ret.SubCode = respJson.Get("response").Get("sub_code").MustString("")
			ret.SubMsg = respJson.Get("response").Get("sub_msg").MustString("")
			ret.Detail = "log_id:" + req.GetLogId()
			return ret
		}
	}
	return nil
}
