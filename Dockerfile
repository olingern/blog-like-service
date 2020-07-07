FROM golang:1.14 as builder
WORKDIR /app
COPY .* ./
ADD . /
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server ./cmd/server

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /server
COPY .env ./
COPY *.json ./
CMD ["/server"]