package main

import (
	"net/url"
)

//各web蜘蛛定義
type Spider interface {
	//返回唯一的 類別標識
	GetId() string

	//設置代理
	SetProxy(proxy *url.URL)

	//發送get請求 並解析數據 記錄到數據庫
	Get(i int) error
}
