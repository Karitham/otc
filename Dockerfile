FROM golang:1.16-alpine as builder
WORKDIR /otc
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /otc/otc ./cmd/otc/.

FROM alpine:latest
WORKDIR /otc
COPY --from=builder /otc/otc ./otc
ENTRYPOINT ["/otc/otc"]