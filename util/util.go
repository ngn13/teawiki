package util

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	if st, err := os.Stat(path); err != nil {
		return false
	} else {
		return st.IsDir()
	}
}

func UrlJoin(u *url.URL, p ...string) string {
	for _, path := range p {
		u = u.JoinPath(path)
	}
	return u.String()
}

func Html(content string) template.HTML {
	return template.HTML(content)
}

func First(s string, n int) string {
	return s[:n]
}

func GetDuration(d string) (time.Duration, error) {
	var (
		dur uint64
		err error
	)

	dur_chr := d[len(d)-1]
	dur_num := strings.TrimSuffix(d, string(dur_chr))

	if dur, err = strconv.ParseUint(dur_num, 10, 64); err != nil {
		return 0, err
	}

	switch dur_chr {
	case 's':
		return time.Duration(dur) * (time.Second), nil

	case 'm':
		return time.Duration(dur) * (time.Second * 60), nil

	case 'h':
		return time.Duration(dur) * ((time.Second * 60) * 60), nil
	}

	return 0, fmt.Errorf("invalid time duration format")
}

func Host(u string) string {
	if urlp, err := url.Parse(u); err != nil {
		return u
	} else {
		return urlp.Hostname()
	}
}

func Map(args ...any) map[string]any {
	var (
		key string         = "undefined"
		res map[string]any = map[string]any{}
	)

	for i, a := range args {
		if i%2 == 0 {
			key = a.(string)
		} else {
			res[key] = a
		}
	}

	return res
}

func Add(nums ...int) int {
	num := 0

	for _, n := range nums {
		num += n
	}

	return num
}
