FROM golang:1.26.3@sha256:313faae491b410a35402c05d35e7518ae99103d957308e940e1ae2cfa0aac29b
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
