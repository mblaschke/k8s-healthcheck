FROM golang:1.10 as build
WORKDIR /go/src/k8s-healthcheck/src
COPY ./src /go/src/k8s-healthcheck/src
RUN curl https://glide.sh/get | sh \
    && glide install
RUN mkdir /app/ \
    && cp -a entrypoint.sh /app/ \
    && chmod 555 /app/entrypoint.sh \
    && go build -o /app/k8s-healthcheck

#############################################
# FINAL IMAGE
#############################################
FROM alpine
RUN apk add --no-cache \
        libc6-compat \
    	ca-certificates
COPY --from=build /app/ /app/
USER 1000
ENTRYPOINT ["/app/entrypoint.sh"]
