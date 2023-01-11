package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/opensourceways/community-robot-lib/interrupts"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/sirupsen/logrus"
)

type options struct {
	service liboptions.ServiceOptions
}

func (o *options) Validate() error {

	return nil
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.Parse(args)
	return o
}

func main() {
	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	http.HandleFunc("/gitee-hook", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("request num %s", r.Header.Get("request_num"))
	})

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(o.service.Port)}

	defer interrupts.WaitForGracefulShutdown()
	interrupts.ListenAndServe(httpServer, o.service.GracePeriod)
}
