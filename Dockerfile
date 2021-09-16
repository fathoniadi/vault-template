FROM golang:1.15-alpine as build

# build binary
WORKDIR /src
COPY . ./

#RUN CGO_ENABLED=0 GOOS=linux go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /vault-template

FROM alpine

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /vault-template /vault-template

ENTRYPOINT ["/vault-template"]
