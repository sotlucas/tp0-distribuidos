import logging

LEN_BYTES = 4


def read_message(client_sock) -> str:
    """
    Reads a bet message from a specific client socket
    """
    addr = client_sock.getpeername()
    length = client_sock.recv(LEN_BYTES).decode("utf-8")
    msg = client_sock.recv(int(length)).decode("utf-8")

    logging.info(
        f"action: receive_message | result: success | ip: {addr[0]} | msg: {msg}"
    )

    return msg


def send_ok(client_sock):
    """
    Send OK message to client
    """
    client_sock.send("OK\n".encode("utf-8"))
