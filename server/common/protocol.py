import logging

from common.utils import Bet

LEN_BYTES = 4
BET_DELIMITER = ";"
BET_PARAMS_DELIMITER = ":"


class Message:
    """
    Represents a message sent by a client
    """

    def __init__(self, action: str, payload: str):
        self.action = action
        self.payload = payload


def recv_wrapper(client_sock, length: int):
    """
    Wrapper around recv to prevent short reads
    """
    msg = b""
    while len(msg) < length:
        chunk = client_sock.recv(length - len(msg))
        if not chunk:
            raise ConnectionError
        msg += chunk
    return msg


def read_message(client_sock) -> Message:
    """
    Reads a bet message from a specific client socket
    """
    addr = client_sock.getpeername()
    length_bytes = recv_wrapper(client_sock, LEN_BYTES)
    length = int.from_bytes(length_bytes, "big")
    msg = recv_wrapper(client_sock, length).decode("utf-8")

    logging.debug(
        f"action: receive_message | result: success | ip: {addr[0]} | msg: {msg} | length: {length}"
    )

    if not msg:
        return None

    action, payload = msg.split("::")
    return Message(action, payload)


def send_ok(client_sock):
    """
    Send OK message to client
    """
    client_sock.sendall("OK\n".encode("utf-8"))


def send_message(client_sock, action: str, payload: str):
    """
    Send message to client
    """
    msg = f"{action}::{payload}"
    length = len(msg)
    client_sock.sendall(length.to_bytes(LEN_BYTES, "big"))
    client_sock.sendall(msg.encode("utf-8"))


def bet_from_string(bet_str: str) -> Bet:
    """
    Parses a bet string into a Bet object
    """
    return Bet(*bet_str.split(BET_PARAMS_DELIMITER))


def bets_from_string(bets_str: str) -> list[Bet]:
    """
    Parses a bet string into a list of Bet objects
    """
    return [
        bet_from_string(bet_str) for bet_str in bets_str.split(BET_DELIMITER) if bet_str
    ]


def bets_to_string(bets: list[Bet]) -> str:
    """
    Parses a list of Bet objects into a bet string
    """
    return BET_DELIMITER.join([str(bet) for bet in bets])
