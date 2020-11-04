FROM golang:1.15-alpine as build

# install dependencies
# ENV GO111MODULE=on
# WORKDIR $GOPATH/src/github.com/minhdanh/vault-template
# COPY go.mod go.sum ./
# RUN go get .

# build binary
WORKDIR /src
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /vault-template

FROM scratch

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /vault-template /vault-template

ENTRYPOINT ["/vault-template"]
