package main

import (
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type T66y2 struct {
	proxy *url.URL

	regTd *regexp.Regexp
	regA  *regexp.Regexp
}

func NewT66y2() (*T66y2, error) {
	regTd, e := regexp.Compile(`<td class="tal" style="padding-left:8px" id="">[\d\D]*?</td>`)
	if e != nil {
		return nil, e
	}

	regA, e := regexp.Compile(`<a href="htm_data/[\d\D]*?</a>`)
	if e != nil {
		return nil, e
	}

	return &T66y2{
		regTd: regTd,
		regA:  regA,
	}, nil
}

//返回唯一的 類別標識
func (t *T66y2) GetId() string {
	//返回唯一的 類別標識
	return "草榴社區-亞洲無碼原創區"
}

//設置代理
func (t *T66y2) SetProxy(proxy *url.URL) {
	t.proxy = proxy
}

//發送get請求 並解析數據
func (t *T66y2) Get(i int) error {
	addr := fmt.Sprintf("http://t66y.com/thread0806.php?fid=2&search=&page=%v", i+1)
	id := t.GetId()
	return t.get(id, addr)
}

func (t *T66y2) get(id, addr string) error {
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(t.proxy),
		},
	}
	r, e := c.Get(addr)
	if e != nil {
		return e
	}

	dec := mahonia.NewDecoder("gbk")
	reader := dec.NewReader(r.Body)
	b, e := ioutil.ReadAll(reader)
	if e != nil {
		return e
	}
	return t.analyze(id, b)
}
func (t *T66y2) analyze(id string, b []byte) error {
	reg := t.regTd
	arrs := reg.FindAll(b, -1)
	if arrs != nil {
		for i := 0; i < len(arrs); i++ {
			if e := t.analyzeTd(id, arrs[i]); e != nil {
				log.Println(e)
			}
		}

	}
	return nil
}
func (t *T66y2) analyzeTd(id string, b []byte) error {
	reg := t.regA
	arrs := reg.FindAll(b, 1)
	if arrs == nil {
		return nil
	}
	b = arrs[0]
	size := len(b)
	if size < 10 || b[size-5] == '>' {
		return nil
	}

	pos := bytes.Index(b, []byte("href="))
	if pos == -1 {
		return nil
	}
	b = b[pos+6:]
	pos = bytes.Index(b, []byte(`"`))
	if pos == -1 {
		return nil
	}
	node := Node{
		Url: string(b[:pos]),
	}
	b = b[pos:]
	pos = bytes.Index(b, []byte(`>`))
	if pos == -1 {
		return nil
	}
	b = b[pos+1:]
	b = b[:len(b)-4]
	node.Name = string(b)
	if _, e := GetEngine().InsertOne(node); e != nil {
		return e
	}
	log.Println("ok", node.Name)
	return nil
}
