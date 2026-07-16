FROM golang:1.26.5@sha256:ae5a2316d12f3e78fd99177dad452e6ad4f240af2d71d57b480c3477f250fec6
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
