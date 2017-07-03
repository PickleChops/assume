FROM golang:1.8.3 as builder
WORKDIR /go/src/assume
COPY assume.go .
RUN go get -d -v .
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:latest
RUN \
	mkdir -p /aws && \
	apk -Uuv add groff less python py-pip ca-certificates wget && \
	update-ca-certificates && \
	pip install awscli && \
	apk --purge -v del py-pip && \
	rm /var/cache/apk/*

COPY --from=builder /go/src/assume/assume /usr/local/bin/assume

WORKDIR /aws
ENTRYPOINT [ "sh", "-c" ]