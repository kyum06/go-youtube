package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

const (
	jsvarStr   = "[a-zA-Z_\\$][a-zA-Z_0-9]*"
	reverseStr = ":function\\(a\\)\\{" +
		"(?:return )?a\\.reverse\\(\\)" +
		"\\}"
	spliceStr = ":function\\(a,b\\)\\{" +
		"a\\.splice\\(0,b\\)" +
		"\\}"
	swapStr = ":function\\(a,b\\)\\{" +
		"var c=a\\[0\\];a\\[0\\]=a\\[b(?:%a\\.length)?\\];a\\[b(?:%a\\.length)?\\]=c(?:;return a)?" +
		"\\}"
)

var (
	actionsObjRegexp = regexp.MustCompile(fmt.Sprintf(
		"var (%s)=\\{((?:(?:%s%s|%s%s|%s%s),?\\n?)+)\\};", jsvarStr, jsvarStr, swapStr, jsvarStr, spliceStr, jsvarStr, reverseStr))

	actionsFuncRegexp = regexp.MustCompile(fmt.Sprintf(
		"function(?: %s)?\\(a\\)\\{"+
			"a=a\\.split\\(\"\"\\);\\s*"+
			"((?:(?:a=)?%s\\.%s\\(a,\\d+\\);)+)"+
			"return a\\.join\\(\"\"\\)"+
			"\\}", jsvarStr, jsvarStr, jsvarStr))

	reverseReg = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, reverseStr))
	spliceReg  = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, spliceStr))
	swapReg    = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, swapStr))
)

func reverse(s string) (r string) {
	for _, v := range s {
		r = string(v) + r
	}
	return
}

func splice(s string, i int) string {
	return s[i:]
}

func swap(s string, i int) string {
	p := i % len(s)
	bs := []rune(s)
	bs[0], bs[p] = bs[p], bs[0]
	return string(bs)
}

func (y YOUTUBE) Decipher(format *Format, path []byte) {

	if (*format).Cipher == "" {
		return
	}

	URL := "http://youtube.com" + string(path)

	res, _ := http.Get(URL)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	obj := actionsObjRegexp.FindSubmatch(body)
	fun := actionsFuncRegexp.FindSubmatch(body)

	if len(obj) < 3 || len(fun) < 2 {
		return
	}

	var reverseKey, spliceKey, swapKey string

	if r := reverseReg.FindSubmatch(obj[2]); len(r) > 1 {
		reverseKey = string(r[1])
	}
	if r := spliceReg.FindSubmatch(obj[2]); len(r) > 1 {
		spliceKey = string(r[1])
	}
	if r := swapReg.FindSubmatch(obj[2]); len(r) > 1 {
		swapKey = string(r[1])
	}

	regex, _ := regexp.Compile(fmt.Sprintf("(?:a=)?%s\\.(%s|%s|%s)\\(a,(\\d+)\\)", obj[1], reverseKey, spliceKey, swapKey))

	params, _ := url.ParseQuery((*format).Cipher)
	fmt.Println((*format).Cipher, "b")
	bs := params.Get("s")
	for _, s := range regex.FindAllSubmatch(fun[1], -1) {
		fmt.Println(bs)
		switch string(s[1]) {
		case reverseKey:
			bs = reverse(bs)
		case swapKey:
			arg, _ := strconv.Atoi(string(s[2]))
			bs = swap(bs, arg)
		case spliceKey:
			arg, _ := strconv.Atoi(string(s[2]))
			bs = splice(bs, arg)
		}
	}

	(*format).URL = fmt.Sprintf("%s&%s=%s", params.Get("url"), params.Get("sp"), bs)
}
