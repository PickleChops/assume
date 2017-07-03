FROM alpine:latest

RUN \
	mkdir -p /aws && \
	apk -Uuv add groff less python py-pip ca-certificates wget && \
	update-ca-certificates && \
	pip install awscli && \
	apk --purge -v del py-pip && \
	rm /var/cache/apk/*

ADD assume /usr/local/bin

WORKDIR /aws
ENTRYPOINT [ "sh", "-c" ]


