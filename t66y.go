package main

import (
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type T66y struct {
	proxy string

	regTd *regexp.Regexp
	regA  *regexp.Regexp

	id  string
	gid int
}

func NewT66y2() (*T66y, error) {
	t66y, e := initBaseT66y()
	if e != nil {
		return nil, e
	}
	t66y.id = "草榴社區-亞洲無碼原創區"
	t66y.gid = 2
	return t66y, nil
}
func NewT66y15() (*T66y, error) {
	t66y, e := initBaseT66y()
	if e != nil {
		return nil, e
	}
	t66y.id = "草榴社區-亞洲有碼原創區"
	t66y.gid = 15
	return t66y, nil
}
func NewT66y4() (*T66y, error) {
	t66y, e := initBaseT66y()
	if e != nil {
		return nil, e
	}
	t66y.id = "草榴社區-歐美創區"
	t66y.gid = 4
	return t66y, nil
}

func initBaseT66y() (*T66y, error) {
	regTd, e := regexp.Compile(`<td class="tal" style="padding-left:8px" id="">[\d\D]*?</td>`)
	if e != nil {
		return nil, e
	}

	regA, e := regexp.Compile(`<a href="htm_data/[\d\D]*?</a>`)
	if e != nil {
		return nil, e
	}

	return &T66y{
		regTd: regTd,
		regA:  regA,
	}, nil
}

//返回唯一的 類別標識
func (t *T66y) GetId() string {
	//返回唯一的 類別標識
	return t.id
}

//設置代理
func (t *T66y) SetProxy(proxy string) {
	t.proxy = proxy
}

//發送get請求 並解析數據
func (t *T66y) Get(i int) error {
	addr := fmt.Sprintf("http://t66y.com/thread0806.php?fid=%v&search=&page=%v", t.gid, i+1)
	id := t.GetId()
	return t.get(id, addr)
}
func (t *T66y) get(id, addr string) error {
	var c *http.Client

	if strings.HasPrefix(t.proxy, "http") {
		proxyUrl, e := url.Parse(t.proxy)
		if e != nil {
			return e
		}
		c = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else if strings.HasPrefix(t.proxy, "socks5://") {
		//connect proxy
		proxyAddr := t.proxy[9:]
		dialer, e := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
		if e != nil {
			log.Fatalln(e)
		}
		//create dial
		httpTransport := &http.Transport{}
		httpTransport.Dial = dialer.Dial

		//http client
		c = &http.Client{Transport: httpTransport}
	} else {
		c = &http.Client{
			Transport: &http.Transport{},
		}
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
func (t *T66y) analyze(id string, b []byte) error {
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
func (t *T66y) analyzeTd(id string, b []byte) error {
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
		Url: "http://t66y.com/" + string(b[:pos]),
		Gid: t.GetId(),
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
		//return e
		return nil
	}
	fmt.Println(t.GetId(), node.Name)
	return nil
}
