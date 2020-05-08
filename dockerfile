FROM golang:alpine AS builder
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/gitery
COPY . .
RUN go mod download
RUN mkdir -p bin && cp ./configs/configs.yaml ./bin/
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/gitery ./cmd/gitery

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=builder /go/src/gitery/bin /go/bin
EXPOSE 80
ENTRYPOINT /go/bin/gitery -env=production --port 80