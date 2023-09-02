import socket
import logging
import signal

from common.utils import Bet, store_bets


class Server:
    def __init__(self, port, listen_backlog):
        self._running = True
        signal.signal(signal.SIGTERM, self.__shutdown)
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(("", port))
        self._server_socket.listen(listen_backlog)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError:
                break

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            addr = client_sock.getpeername()
            length = client_sock.recv(4).decode("utf-8")
            msg = client_sock.recv(int(length)).decode("utf-8")
            logging.info(
                f"action: receive_message | result: success | ip: {addr[0]} | msg: {msg}"
            )

            bet = Bet(*msg.split(":"))
            store_bets([bet])
            logging.info(
                f"action: apuesta_almacenada | result: success | client_id: {bet.agency} | dni: {bet.document} | numero: {bet.number}"
            )

            client_sock.send("OK\n".encode("utf-8"))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info("action: accept_connections | result: in_progress")
        c, addr = self._server_socket.accept()
        logging.info(f"action: accept_connections | result: success | ip: {addr[0]}")
        return c

    def __shutdown(self, *args):
        """
        Shutdown server socket
        """
        logging.info("action: shutdown_server | result: in_progress")
        self._running = False
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        logging.info("action: shutdown_server | result: success")
