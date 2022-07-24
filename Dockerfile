FROM golangci/golangci-lint:v1.47-alpine AS go-builder

RUN apk update
RUN apk upgrade
RUN apk add ca-certificates git gcc g++ libc-dev

WORKDIR /go/uiassignment

ENV GO111MODULE=on

RUN mkdir -p /artifact/uiassignment
COPY ./ /go/uiassignment

RUN \
    go mod init uiassignment || true && \
    go mod tidy && \
    go mod download && \
    (cd cmd/uiassignment && go build -o uiassignment-binary) && \
    mv ./cmd/uiassignment/uiassignment-binary /artifact

FROM alpine:3.13.2

RUN apk add tzdata

RUN mkdir -p /swagger-docs && \
    mkdir -p /logs

#copy artifact
COPY --from=go-builder artifact/ /app/uiassignment

ENTRYPOINT [ "/app/uiassignment/uiassignment-binary" ]