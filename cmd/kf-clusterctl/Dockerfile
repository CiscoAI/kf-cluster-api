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

FROM google/cloud-sdk:278.0.0-alpine
COPY --from=build /go/bin/kf-clusterctl /usr/bin/kf-clusterctl
COPY --from=build /go/src/kf-clusterctl/cmd/kf-clusterctl/gcp_entrypoint.sh /gcp_entrypoint.sh
RUN chmod +x /gcp_entrypoint.sh
RUN chmod +x /usr/bin/kf-clusterctl

# Download kubectl linux binary
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.17.0/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/
# Download kops linux binary
RUN curl -Lo kops https://github.com/kubernetes/kops/releases/download/$(curl -s https://api.github.com/repos/kubernetes/kops/releases/latest | grep tag_name | cut -d '"' -f 4)/kops-linux-amd64
RUN chmod +x ./kops
RUN mv ./kops /usr/local/bin/
# Download kfctl linux binary
RUN wget https://storage.googleapis.com/kubernetes-jenkins/pr-logs/pull/kubeflow_kfctl/173/kubeflow-kfctl-presubmit/1215810676291801090/artifacts/build_bin/kfctl
RUN chmod +x ./kfctl
RUN mv ./kfctl /usr/local/bin/
