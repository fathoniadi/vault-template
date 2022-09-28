FROM golang:1.16 as build

# build binary
WORKDIR /src
COPY . ./

#RUN CGO_ENABLED=0 GOOS=linux go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /vault-template

FROM alpine

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /vault-template /bin/vault-template

ENTRYPOINT ["/bin/vault-template"]
