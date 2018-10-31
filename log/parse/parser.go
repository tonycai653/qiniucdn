package parse

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogLine struct {
	ClientIp  string
	Hit       string
	RespTime  uint32    // 请求耗时，毫秒
	ReqTime   time.Time // 请求时间
	Method    string
	Url       string
	RespCode  int
	RespSize  uint64
	Referer   string
	UserAgent string
}

func (l LogLine) String() string {
	return fmt.Sprintf("ClientIp: %s\nHit: %s\nRespTime: %d\nReqTime: %v\nMethod: %s\nUrl: %s\nRespCode: %d\nRespSize: %d\nReferer: %s\nUserAgent: %s\n", l.ClientIp, l.Hit, l.RespTime, l.ReqTime, l.Method, l.Url, l.RespCode, l.RespSize, l.Referer, l.UserAgent)
}

func (l LogLine) TokenExpired() (expired bool, err error) {

	ul, err := url.Parse(l.Url)
	if err != nil {
		err = fmt.Errorf("parse url: %s, failed: %v\n", l.Url, err)
		return
	}
	v, err := url.ParseQuery(ul.RawQuery)
	if err != nil {
		err = fmt.Errorf("parse query: %s, failed: %v\n, url: %s", ul.RawQuery, err, l.Url)
		return
	}
	seconds := v.Get("e")
	if seconds == "" {
		return false, nil
	}
	secInt, err := strconv.Atoi(seconds)
	if err != nil {
		err = fmt.Errorf("convert str: %s to int, failed: %v\n, url: %s", seconds, err, l.Url)
		return
	}
	if time.Unix(int64(secInt), 0).Before(l.ReqTime) {
		return true, nil
	} else {
		return false, nil
	}
}

//Bytes, KB, MB, GB, TB
type LogSize []uint64

func NewLogSize() LogSize {
	return make(LogSize, 10, 10)
}

func (lsize LogSize) Add(respSize uint64) {
	lsize.add(respSize, 0)
}

func (lsize LogSize) add(respSize uint64, ind int) {
	if ind > len(lsize)-1 {
		panic("index too large")
	}
	lsize[ind] += respSize
	if lsize[ind] >= 1024 {
		size := lsize[ind] / 1024
		lsize[ind] %= 1024
		lsize.add(size, ind+1)
	}
}

func (lsize LogSize) String() (ret string) {
	unit := ""
	for i := len(lsize) - 1; i >= 0; i-- {
		if lsize[i] > 0 {
			switch i {
			case 4:
				unit = "TB"
			case 3:
				unit = "GB"
			case 2:
				unit = "MB"
			case 1:
				unit = "KB"
			case 0:
				unit = "Bytes"
			}
			if ret != "" {
				ret = fmt.Sprintf("%s %d%s", ret, lsize[i], unit)
			} else {
				ret = fmt.Sprintf("%s%d%s", ret, lsize[i], unit)
			}
		}
	}
	return
}

var reg = regexp.MustCompile(`(?P<ClientIp>\d+\.\d+\.\d+\.\d+)\s+(?P<Hit>(HIT|MISS|UNKNOWN|-))\s+(?P<RespTime>\d+)\s+(?P<ReqTime>\[.+\])\s+"(?P<Method>GET|HEAD|POST|OPTIONS)\s+(?P<Url>.+)\s+HTTP(s)?/(1.0|1.1|2.0)"\s+(?P<RespCode>\d+)\s+(?P<RespSize>\d+)\s+(?P<Referer>".*")\s+(?P<UserAgent>".*")\s*$`)

var names = reg.SubexpNames()

func ParseLine(line string, logLine *LogLine) error {
	fields := reg.FindStringSubmatch(line)
	if len(fields) < len(names) {
		return fmt.Errorf("ParseLine: %s\n", line)
	}
	for ind, name := range names {
		if ind == 0 || name == "" {
			continue
		}
		switch name {
		case "ClientIp":
			logLine.ClientIp = fields[ind]
		case "Hit":
			logLine.Hit = fields[ind]
		case "RespTime":
			respTime, _ := strconv.Atoi(fields[ind])
			logLine.RespTime = uint32(respTime)
		case "ReqTime":
			t := strings.TrimLeft(strings.Fields(fields[ind])[0], "[")
			tm, err := time.ParseInLocation("02/Jan/2006:15:04:05", t, time.Local)
			if err != nil {
				return fmt.Errorf("parse time: %s, error: %v\n", t, err)
			}
			logLine.ReqTime = tm
		case "Method":
			logLine.Method = fields[ind]
		case "Url":
			logLine.Url = fields[ind]
		case "RespCode":
			respCode, _ := strconv.Atoi(fields[ind])
			logLine.RespCode = respCode
		case "RespSize":
			respSize, _ := strconv.Atoi(fields[ind])
			logLine.RespSize = uint64(respSize)
		case "Referer":
			logLine.Referer = fields[ind]
		case "UserAgent":
			logLine.UserAgent = fields[ind]
		}
	}
	return nil
}
