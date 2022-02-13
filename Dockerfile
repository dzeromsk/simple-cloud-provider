# syntax=docker/dockerfile:experimental

FROM golang:1.16-alpine as dev
RUN apk add --no-cache git ca-certificates
RUN adduser -D appuser
COPY . /src/
WORKDIR /src

ENV GO111MODULE=on
RUN --mount=type=cache,sharing=locked,id=gomod,target=/go/pkg/mod/cache \
    --mount=type=cache,sharing=locked,id=goroot,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags '-s -w -extldflags -static' -o simple-cloud-provider

FROM scratch
COPY --from=dev /src/simple-cloud-provider /
CMD ["/simple-cloud-provider"]