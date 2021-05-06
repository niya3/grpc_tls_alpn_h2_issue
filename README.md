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

## prerequisites

```bash
make deps
```

### with h2, via tls.Listener

```bash
$ make server-h2-listener
2021/05/06 15:49:11 TLS config includes h2
2021/05/06 15:49:11 Pass TLS config via tls.NewListener
2021/05/06 15:49:18 Request: 0

$ make client
Received message: %s confirmation: "Credited 0"
```

### without h2, via tls.Listener


```bash
$ make server-listener
2021/05/06 15:49:52 TLS config doesn't include h2
2021/05/06 15:49:52 Pass TLS config via tls.NewListener

$ make client
I0506 15:49:54.290646018   61520 subchannel.cc:1012]         Connect failed: {"created":"@1620305394.290596812","description":"Cannot check peer: missing selected ALPN property.","file":"src/core/lib/security/security_connector/ssl_utils.cc","file_line":161}
I0506 15:49:54.290681420   61520 subchannel.cc:957]          Subchannel 0x55f3d3d6d750: Retry in 973 milliseconds
Received error: %s <_InactiveRpcError of RPC that terminated with:
        status = StatusCode.UNAVAILABLE
        details = "failed to connect to all addresses"
        debug_error_string = "{"created":"@1620305394.290677436","description":"Failed to pick subchannel","file":"src/core/ext/filters/client_channel/client_channel.cc","file_line":5419,"referenced_errors":[{"created":"@1620305394.290674361","description":"failed to connect to all addresses","file":"src/core/ext/filters/client_channel/lb_policy/pick_first/pick_first.cc","file_line":397,"grpc_status":14}]}"
>
```

### with h2, via grpc.Creds

```bash
$ make server-h2-creds
2021/05/06 15:50:06 TLS config includes h2
2021/05/06 15:50:06 Pass TLS config via grpc.Creds
2021/05/06 15:50:08 Request: 0

$ make client
Received message: %s confirmation: "Credited 0"
```

### without h2, via grpc.Creds

```bash
$ make server-creds
2021/05/06 15:50:12 TLS config doesn't include h2
2021/05/06 15:50:12 Pass TLS config via grpc.Creds
2021/05/06 15:50:13 Request: 0

$ make client
Received message: %s confirmation: "Credited 0"
```


## grpcurl client run

Clients written in Golang succesfully connect to server with any server's options combination:

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
