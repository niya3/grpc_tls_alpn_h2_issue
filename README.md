# Essence

Demo of `missing selected ALPN property` issue for gRPC-clients using the C based library: https://github.com/grpc/grpc/issues/23172

Clients written in Python can't connect to gRPC server written in Golang if it uses manually configured TLS listener (e.g. minimal TLS protocol version is specified) without explicitly specified `h2` as next proto in [`Config.NextProtos`](https://github.com/golang/go/blob/5f1df260a91183c605c08af7b00741d2761b84e4/src/crypto/tls/common.go#L620-L622) field. Clients written in golang (e.g. https://github.com/fullstorydev/grpcurl) don't have such troubles.

## server.go

An example server writen in Golang based on http://www.inanzzz.com/index.php/post/gq4x/using-tls-ssl-certificates-for-grpc-client-and-server-communications-in-golang

Does have one flag `-with-h2=true` to highlight the issue.

## client.py

An example client writen in Python based on https://github.com/grpc/grpc/blob/1c49176a24a501676a13b114356ab0bccaede22b/examples/python/auth/customized_auth_client.py

Doesn't have any flags.

# Usage

The code in this repository has been checked in `python:3.9.4-debian` docker image. **NB** Run all `make` commands from the directory containing the `Makefile` file.

```bash
docker run -ti -v $PWD:/dir python:3.9.4-buster bash
cd /dir
```

## dependecies

```bash
make deps
```

## Python client OK run

```bash
make server
make client
# server: 2021/05/04 12:32:02 Run server with h2
#
# client: Received message: %s confirmation: "Credited 0"
```

## Python client failed run

```bash
make broken-server
make client

# server: 2021/05/04 12:33:36 Run server WITHOUT h2
#
# client: I0504 12:33:38.260658538   52867 subchannel.cc:1012]         Connect failed: {"created":"@1620120818.260611055","description":"Cannot check peer: missing selected ALPN property.","file":"src/core/lib/security/security_connector/ssl_utils.cc","file_line":161}
# client: I0504 12:33:38.260715165   52867 subchannel.cc:957]          Subchannel 0x56070f114270: Retry in 981 milliseconds
# client: Received error: %s <_InactiveRpcError of RPC that terminated with:
# client:         status = StatusCode.UNAVAILABLE
```

## grpcurl client run

Clients written in Golang succesfully connect to server regardless of server's `-with-h2` value:

```bash
make go-client

# server: 2021/05/04 12:35:13 Run server with h2
#
# client: {
# client:   "confirmation": "Credited 0"
# client: }

# server: 2021/05/04 12:35:45 Run server with h2
#
# client: {
# client:   "confirmation": "Credited 0"
# client: }
```
