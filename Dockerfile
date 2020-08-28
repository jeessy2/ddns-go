# build stage(Only for server)
FROM golang AS builder
WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go get -d -v . \
    && go install -v . \
    && go build -v .

# final stage, build server
FROM golang
WORKDIR /app
COPY --from=builder /app/ddns-go /app/ddns-go
EXPOSE 9876
ENTRYPOINT /app/ddns-go
LABEL Name=ddns-go Version=0.0.1
