import logging

from common.utils import Bet

LEN_BYTES = 4
BET_DELIMITER = ";"
BET_PARAMS_DELIMITER = ":"


def read_message(client_sock) -> str:
    """
    Reads a bet message from a specific client socket
    """
    addr = client_sock.getpeername()
    length_bytes = client_sock.recv(LEN_BYTES)
    length = int.from_bytes(length_bytes, "big")
    msg = client_sock.recv(int(length)).decode("utf-8")

    logging.debug(
        f"action: receive_message | result: success | ip: {addr[0]} | msg: {msg} | length: {length}"
    )

    return msg


def send_ok(client_sock):
    """
    Send OK message to client
    """
    client_sock.send("OK\n".encode("utf-8"))


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
