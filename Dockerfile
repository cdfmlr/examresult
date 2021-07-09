# golang:latest as build-env
FROM golang:alpine AS build-env

# Redundant, current golang images already include ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /

COPY . ./examresult

ENV CGO_ENABLED=0
RUN cd examresult \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
    && go build -o examresult . \
    && cd ..

FROM scratch

WORKDIR /

COPY --from=build-env /examresult/examresult /
COPY --from=build-env /examresult/config.json /

COPY --from=build-env /examresult/asset /
COPY --from=build-env /examresult/asset/register.html /asset/

# copy the ca-certificate.crt from the build stage
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV GIN_MODE=release
ENTRYPOINT ["/examresult", "-c", "/config.json"]

