FROM golang:1.27.3@sha256:e519d3b0a5b285e5aa53ab9eef02bc6bbb9cb04608389b743a136a9ce75f0d1e
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
