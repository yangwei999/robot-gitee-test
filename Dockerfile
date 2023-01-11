FROM golang:1.18 as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/robot-gitee-test
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o robot-gitee-test .

# copy binary config and utils
FROM alpine:3.14
COPY  --from=BUILDER /go/src/github.com/opensourceways/robot-gitee-test/robot-gitee-test /opt/app/robot-gitee-test

ENTRYPOINT ["/opt/app/robot-gitee-test"]
