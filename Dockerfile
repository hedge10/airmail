FROM golang:1.20-alpine3.17 AS builder

LABEL org.opencontainers.image.source="https://github.com/hedge10/airmail"

ARG VERSION=SNAPSHOT
ENV CGO_ENABLED=0

WORKDIR /workspace

COPY go.mod go.sum /workspace/

RUN go mod download

COPY . /workspace/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'github.com/hedge10/airmail.Version=$VERSION'" -o airmail ./cmd/airmail

FROM gcr.io/distroless/static

COPY --from=builder /workspace/airmail .

EXPOSE 9900

ENTRYPOINT [ "/airmail" ]