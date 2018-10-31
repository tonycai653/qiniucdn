package parse

import (
	"cdn/log/parse"
	"testing"
	"time"
)

func TestLogSize(t *testing.T) {
	tables := map[uint64]string{
		1:                            "1Bytes",
		20:                           "20Bytes",
		1024:                         "1KB",
		1025:                         "1KB 1Bytes",
		2048:                         "2KB",
		2049:                         "2KB 1Bytes",
		1024 * 1024:                  "1MB",
		1024*1024 + 1:                "1MB 1Bytes",
		1024*1024 + 1024 + 2:         "1MB 1KB 2Bytes",
		1024*1024*1024*2 + 1024*1024: "2GB 1MB",
	}

	for key, value := range tables {
		var size parse.LogSize = parse.NewLogSize()
		size.Add(key)
		if size.String() != value {
			t.Errorf("expected: %s, real: %s", value, size.String())
		}
	}
}

func TestTokenExpired(t *testing.T) {
	tm, _ := time.ParseInLocation("02/Jan/2006:15:04:05", "03/Aug/2018:14:33:42", time.Local)

	logLine := parse.LogLine{
		ClientIp:  "222.172.134.143",
		Hit:       "Hit",
		RespTime:  3878,
		ReqTime:   tm,
		Method:    "GET",
		Url:       "https://pro-app-qn.fir.im/de6c6da598066b501aa2543bb7407e0f5adbe749.apk?attname=huayu_an.weima3d.     com_201805152219.apk_1.0.apk&e=1532414273&token=LOvmia8oXF4xnLh0IdH05XMYpH6ENHNpARlmPc-T:oNtZGTqGu4Wg5lRImrDysMt_bDk=",
		RespCode:  200,
		RespSize:  1365604,
		Referer:   "https://fir.im     /aucr",
		UserAgent: "Mozilla/5.0 (Linux; U; Android 6.0.1; zh-cn; OPPO R9s Build/MMB29M) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/53.0.2785.134 Mobile      Safari/537.36 OppoBrowser/4.6.5.3",
	}

	expired, _ := logLine.TokenExpired()
	if !expired {
		t.Errorf("expected: true, real: %v\n", expired)
	}
}

func TestParseLine(t *testing.T) {
	var log parse.LogLine
	url := "112.20.202.37 HIT 0 [01/Jul/2018:00:05:50 +0800] \"GET https://pro-app-qn.fir.im/a46a39fc3aae3a34447d9ab10533c2133906a0df.apk?attname=com.fanwe.hybrid.app.App.apk_2.4.2.apk&e=1529742356&token=LOvmia8oXF4xnLh0IdH05XMYpH6ENHNpARlmPc-T:yspjYQgAm3o-z4OR3gkwh5f1JGI= HTTP/1.1\" 206 66716 \"https://fir.im/qiansho\" \"AndroidDownloadManager/7.1.2+(Linux;+U;+Android+7.1.2;+MI+5X+Build/N2G47H)\""
	parse.ParseLine(url, &log)
	t.Log(log)
}
