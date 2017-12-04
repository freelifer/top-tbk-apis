package tbk

import (
	"github.com/freelifer/top-tbk-apis/api"
	"testing"
)

func Test_TbkMap(t *testing.T) {
	tbkMap := NewTbkMap()
	tbkMap.Put("key", "value")
	tbkMap.Put("aaa", "")

	t.Log(tbkMap.Get("key"))
	t.Log(tbkMap.Get("aaa"))
}

// go test -run="Test_taobao_tbk_dg_item_coupon_get" *go -v
func Test_taobao_tbk_dg_item_coupon_get(t *testing.T) {
	var request api.TbkDgItemCouponGetRequest
	request.InitDefaultValue("http://gw.api.taobao.com/router/rest", "24659164", "cbe2b136be37cd2b66fd4490b8fbfb94")
	request.SetAdzoneId("148758292")
	request.SetPlatform("2")
	resp, _, err := request.Response()

	// t.Log(string(data))
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.Results)
}
