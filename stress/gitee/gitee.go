package gitee

import (
	"net/http"
	"strconv"
	"time"

	"github.com/opensourceways/robot-gitee-test/stress/utils"
)

const salt = "123"

type Platform struct {
}

func (g Platform) SetHeader(req *http.Request) {
	timestamp := time.Now().UnixMilli()
	token := utils.PayloadSignature(strconv.FormatInt(timestamp, 10), salt)

	req.Header.Add("X-Gitee-Timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Add("X-Gitee-Event", "Note Hook")
	req.Header.Add("User-Agent", "Robot-Gitee-Access")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Gitee-Token", token)
}

func (g Platform) GetPayloadFile() string {
	return "./gitee/note_event_payload"
}
