package gitee

import (
	"net/http"
	"strconv"
	"time"
)

type Platform struct {
}

func (g Platform) SetHeader(req *http.Request) {
	timestamp := time.Now().UnixMilli()
	req.Header.Add("X-Gitee-Timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Add("X-Gitee-Event", "Note Hook")
	req.Header.Add("User-Agent", "Robot-Gitee-Access")
}

func (g Platform) GetPayloadFile() string {
	return "./gitee/note_event_payload"
}
