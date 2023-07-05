FROM golang:1.19 AS builder
COPY . /src/cmd/main.go
RUN go build  -o test  src/cmd/main.go

FROM golang:1.19
COPY --from=builder test test
ENTRYPOINT ./test