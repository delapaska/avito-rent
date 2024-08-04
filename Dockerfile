FROM golang:alpine as builder


WORKDIR /app

 
COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o /main ./cmd/main.go

RUN go build -o /migrate ./cmd/migrate/main.go


FROM alpine:latest

WORKDIR /root/


RUN apk --no-cache add postgresql-client


COPY --from=builder /main .
COPY --from=builder /migrate .


COPY --from=builder /app/cmd/migrate/migrations ./migrations


COPY --from=builder /app/.env .

CMD ["./main"]