package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/opensourceways/community-robot-lib/utils"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-test/stress/gitee"
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
	fs.IntVar(&o.concurrent, "concurrent", 100, "concurrent per second")

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

	replaceBody := strings.Replace(string(body), "{replace_content}", strconv.Itoa(requestNum), 1)
	req, _ := http.NewRequest(
		http.MethodPost, endpoint, strings.NewReader(replaceBody),
	)

	plat.SetHeader(req)

	return req
}

func send(req *http.Request) {
	defer wg.Done()
	if _, err := client.ForwardTo(req, nil); err != nil {
		logrus.WithError(err).Error()
	}
}
