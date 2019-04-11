From golang:1.11-alpine
WORKDIR /go/src/apm_demo_gin

COPY . .

# RUN go get -u github.com/golang/dep/cmd/dep && dep ensure -v -vendor-only 
ENV 1GOOS=linux
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" .

From scratch
COPY --from=0 /go/src/apm_demo_gin/apm_demo_gin /apm_demo_gin
CMD ["/apm_demo_gin"]