FROM alpine:latest


RUN apk --update --allow-untrusted --repository http://dl-4.alpinelinux.org/alpine/edge/community/ add \
      tor \
&& rm -rf /var/cache/apk/* /tmp/* /var/tmp/*



WORKDIR /opt

ADD main .
ADD templates ./templates


VOLUME /configs

ENTRYPOINT ["/opt/main"]
