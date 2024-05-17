FROM golang:latest
LABEL authors="huxin"

WORKDIR /bitkv

COPY . .

RUN cd bitkv-server && go build -o ../bin/bitkv-server
RUN cd bitkv-ctl && go build -o ../bin/bitkv-ctl

ENV PATH=/bitkv/bin/

CMD [ "bitkv-server" ]