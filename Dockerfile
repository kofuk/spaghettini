FROM golang:1.26.5@sha256:3aff6657219a4d9c14e27fb1d8976c49c29fddb70ba835014f477e1c70636647
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
