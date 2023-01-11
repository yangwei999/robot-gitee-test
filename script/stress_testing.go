package main

import (
	"bytes"
	"errors"
	"flag"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/opensourceways/community-robot-lib/utils"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-test/script/gitee"
)

var (
	client utils.HttpClient
	wg     sync.WaitGroup
)

type options struct {
	endpoint   string
	concurrent int
}

func (o *options) Validate() error {
	if o.endpoint == "" {
		return errors.New("missing endpoint")
	}

	return nil
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	fs.StringVar(&o.endpoint, "endpoint", "", "endpoint of delivery robot")
	fs.IntVar(&o.concurrent, "concurrent", 10, "concurrent per second")

	fs.Parse(args)

	return o
}

func init() {
	client = utils.NewHttpClient(3)
}

func main() {
	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("invalid options")
	}

	plat := gitee.Platform{}
	requestNum := 1
	for {
		req := buildRequest(plat, o.endpoint, requestNum)

		wg.Add(1)
		go send(req)

		if requestNum%o.concurrent == 0 {
			wg.Wait()
			time.Sleep(time.Second)
			logrus.Infof("send request count:%d", requestNum)
		}

		requestNum++
	}
}

type platformer interface {
	SetHeader(r *http.Request)
	GetPayloadFile() string
}

func buildRequest(plat platformer, endpoint string, requestNum int) *http.Request {
	body, err := os.ReadFile(plat.GetPayloadFile())
	if err != nil {
		logrus.WithError(err).Fatal("read note event file failed")
	}

	req, _ := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)

	plat.SetHeader(req)
	req.Header.Add("request_num", strconv.Itoa(requestNum))

	return req
}

func send(req *http.Request) {
	defer wg.Done()
	if _, err := client.ForwardTo(req, nil); err != nil {
		logrus.WithError(err).Error()
	}
}
