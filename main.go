package main

import (
	"flag"
	"fmt"
	"king-go/go-xorm/params"
	"log"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	const (
		defaultHelp   = false
		usageHelp     = "show help"
		defaultAction = "filter"
		usageAction   = "get/g(do http request) or filter/f(do sql select)"
		defaultLimit  = "10"
		usageLimit    = "start,n work width get or select"
		defaultProxy  = "socks5://127.0.0.1:1080"
		usageProxy    = "http request use proxy"

		defaultName = ""
		usageName   = "sql where name like"
		defaultGid  = ""
		usageGid    = "sql where gid like"

		defaultCount = false
		usageCount   = "sql count"
	)
	var help bool
	flag.BoolVar(&help, "help", defaultHelp, usageHelp)
	flag.BoolVar(&help, "h", defaultHelp, usageHelp+"	(shorthand)")

	var action string
	flag.StringVar(&action, "action", defaultAction, usageAction)
	flag.StringVar(&action, "a", defaultAction, usageAction+"	(shorthand)")

	var limit string
	flag.StringVar(&limit, "limit", defaultLimit, usageLimit)

	var proxy string
	flag.StringVar(&proxy, "proxy", defaultProxy, usageProxy)

	var name, gid string
	flag.StringVar(&name, "name", defaultName, usageName)
	flag.StringVar(&gid, "gid", defaultGid, usageGid)

	var count bool
	flag.BoolVar(&count, "count", defaultCount, usageCount)

	flag.Parse()
	if help {
		flag.PrintDefaults()
	} else if action == "get" || action == "g" {
		doGet(proxy, limit)
	} else if action == "filter" || action == "f" {
		doFilter(name, gid, limit, count)
	} else {
		flag.PrintDefaults()
	}
}
func getLimit(str string) (int, int) {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0, 10
	}
	strs := strings.Split(str, ",")
	start, _ := strconv.ParseInt(strs[0], 10, 64)
	var n int64
	if len(strs) > 1 {
		n, _ = strconv.ParseInt(strs[1], 10, 64)
		if n < 1 {
			n = 10
		}
	} else {
		n = start
		start = 0
	}

	return int(start), int(n)
}
func doGet(proxyVal, limit string) {
	spiders := make([]Spider, 0)
	type newfunc func() (*T66y, error)
	funcs := []newfunc{
		NewT66y2,  //"草榴社區-亞洲無碼原創區"
		NewT66y15, //"草榴社區-亞洲有碼原創區"
		NewT66y4,  //"草榴社區-歐美創區"
	}
	for i := 0; i < len(funcs); i++ {
		spider, e := funcs[i]()
		if e != nil {
			log.Fatal(e)
		}
		spiders = append(spiders, spider)
	}

	if proxyVal != "" {
		for _, spider := range spiders {
			spider.SetProxy(proxyVal)
		}
	}

	start, n := getLimit(limit)
	//4
	for i := 0; i < n; i++ {
		fmt.Printf("request page %v\n", i)
		for _, spider := range spiders {
			spider.Get(start + i)
		}
	}
}
func doFilter(name, gid, limit string, isCount bool) {
	name = strings.TrimSpace(name)
	gid = strings.TrimSpace(gid)
	wh := params.NewParams(2)
	okName := false
	if name != "" {
		wh.WriteWhere("name like ?")
		wh.WriteParam(name)
		okName = true
	}
	if gid != "" {
		if okName {
			wh.WriteWhere("and gid like ?")
		} else {
			wh.WriteWhere("gid like ?")
		}
		wh.WriteParam(gid)
	}

	start, n := getLimit(limit)
	var beans []Node
	//GetEngine().ShowSQL(true)

	session := NewSession()
	defer session.Close()
	wh.Where(session)
	if isCount {
		if count, e := session.Count(Node{}); e != nil {
			log.Fatal(e)
		} else {
			fmt.Printf("%v rows\n\n\n", count)
		}
		return
	}
	e := session.
		Limit(n, start).
		Decr("Id").
		Find(&beans)
	if e != nil {
		log.Fatal(e)
	}

	size := len(beans)
	for i := 0; i < size; i++ {
		fmt.Printf("row %v\n", i+1)
		fmt.Printf("	%v\n", beans[i].Gid)
		fmt.Printf("	%v\n", beans[i].Name)
		fmt.Printf("	%v\n\n", beans[i].Url)
	}
	fmt.Printf("%v rows\n\n\n", size)
}
