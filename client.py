import grpc

def create_client_channel(addr: str, crt: bytes) -> grpc.Channel:
    creds = grpc.ssl_channel_credentials(crt)
    channel = grpc.secure_channel(addr, creds)
    return channel


def send_rpc(channel: grpc.Channel) -> None:
    messages, servcies = grpc.protos_and_services("credit.proto")
    request = messages.CreditRequest()
    stub = servcies.CreditServiceStub(channel)

    try:
        response = stub.Credit(request)
    except grpc.RpcError as rpc_error:
        print('Received error: %s', rpc_error)
    else:
        print('Received message: %s', response)


def main() -> None:
    crt: bytes = None
    with open("root.crt", "rb") as f:
        crt = f.read()

    channel = create_client_channel("localhost:50051", crt)
    send_rpc(channel)
    channel.close()


if __name__ == '__main__':
    main()
