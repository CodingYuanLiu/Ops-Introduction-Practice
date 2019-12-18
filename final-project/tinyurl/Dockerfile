FROM golang:latest

ADD ./ /home/app/tinyurl/

RUN cd /home/app/tinyurl && go build main.go

EXPOSE 4396

CMD /home/app/tinyurl/main
