FROM golang:alpine AS builder
WORKDIR /root
ADD . .
ARG VERSION
WORKDIR /root/cmd/atlas
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w -X main.version=${VERSION}"
FROM alpine
COPY --from=builder /root/cmd/atlas/atlas /atlas
ENTRYPOINT ["/atlas"]
