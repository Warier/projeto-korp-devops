FROM golang:1.27.1-alpine3.24 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -mod=readonly -trimpath -tags=nomsgpack \
    -ldflags="-s -w" -o /out/http-server-projeto-korp ./cmd/http-server-projeto-korp

FROM scratch AS runtime

COPY --from=build /out/http-server-projeto-korp /http-server-projeto-korp

USER 65532:65532
ENV GIN_MODE=release
EXPOSE 8080

ENTRYPOINT ["/http-server-projeto-korp"]
