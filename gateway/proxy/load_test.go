package proxy

import (
	"eff-gateway/gateway/proxy/types"
	"fmt"
	"testing"
)

func TestLocalConfig(t *testing.T) {

	var respond types.GlobalHttp
	err := loadYaml(&respond, "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(respond.HttpList) == 0 {
		loadJson(&respond, "")
	}
	storeLocalJson(respond)
	checkConfig()

}
