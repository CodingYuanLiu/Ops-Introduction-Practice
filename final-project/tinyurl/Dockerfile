FROM golang:latest

ADD ./ /home/app/tinyurl/

RUN export GO111MODULE=on && export GOPROXY=https://goproxy.io && cd /home/app/tinyurl && go build main.go

EXPOSE 4396

CMD /home/app/tinyurl/main
