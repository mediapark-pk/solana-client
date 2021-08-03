FROM golang:1.16-alpine AS builder

RUN go version
########
# Prep
########
ENV -w GO111MODULE=auto
ENV GO111MODULE=on


# add the source
COPY . /go/src/solana-sdk
WORKDIR /go/src/solana-sdk/

########
# Build Go Wrapper
########

# Install go dependencies
RUN go get github.com/gorilla/mux 
RUN GO111MODULE=on go get -v github.com/portto/solana-go-sdk
RUN go get -v github.com/portto/solana-go-sdk/client
RUN go get -v github.com/portto/solana-go-sdk/types

  


#build the go app
RUN GOOS=linux GOARCH=amd64 go build -o ./sdk ./createAndSendTransaction.go createKeyPair.go getBalance.go utils.go tokensTransfer.go  server.go

########
# Package into runtime image
########
FROM alpine

# copy the executable from the builder image
COPY --from=builder /go/src/solana-sdk/sdk .

ENTRYPOINT ["/sdk"]

EXPOSE 12345 3004