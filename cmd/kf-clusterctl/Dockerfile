FROM golang:1.13-alpine as build

WORKDIR /go/src/kf-clusterctl
ADD . /go/src/kf-clusterctl

RUN GO111MODULE=on go mod download
RUN GO111MODULE=on go mod vendor

RUN GO111MODULE=on \
 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
 go build -a -ldflags '-w -extldflags "-static"' \
 -o /go/bin/kf-clusterctl  \
 /go/src/kf-clusterctl/cmd/kf-clusterctl/kf-clusterctl.go

FROM gcr.io/distroless/static:nonroot

COPY --from=build /go/bin/kf-clusterctl /
USER nonroot:nonroot

ENTRYPOINT ["/kf-clusterctl"]