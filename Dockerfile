FROM golang:1.26.0@sha256:c83e68f3ebb6943a2904fa66348867d108119890a2c6a2e6f07b38d0eb6c25c5
WORKDIR /build
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=./cmd,target=cmd \
    --mount=type=bind,source=./server,target=server \
    cd cmd/spaghettini && \
    CGO_ENABLED=0 go build -o /spaghettini .

FROM scratch
COPY --from=0 /spaghettini /spaghettini
EXPOSE 8080
ENTRYPOINT ["/spaghettini"]
