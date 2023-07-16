FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o /service ./${SERVICE_NAME}/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /service /app/service

ARG ENVIRONMENT
ENV ENVIRONMENT=${ENVIRONMENT}
ENV SERVICE_NAME=${SERVICE_NAME}

EXPOSE 8080

ENTRYPOINT ["/app/service"]