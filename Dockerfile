FROM golang:1.19-alpine as builder

WORKDIR /install

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 go build -o umf /install/cmd/umf/main.go


FROM gcr.io/distroless/static-debian11

WORKDIR /app

COPY --from=builder /install/umf /app
