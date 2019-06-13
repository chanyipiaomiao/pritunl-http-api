FROM busybox

WORKDIR /pritunl-http-api

ADD pritunl-http-api .
ADD conf ./conf

EXPOSE 30080

CMD ["./pritunl-http-api"]
