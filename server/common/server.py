import socket
import logging
import signal

import common.utils as utils
import common.protocol as protocol


class Server:
    def __init__(self, port, listen_backlog):
        self._running = True
        signal.signal(signal.SIGTERM, self.__shutdown)
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(("", port))
        self._server_socket.listen(listen_backlog)
        self._agencies_done = [False, False, False, False, False]

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
            msg = protocol.read_message(client_sock)
            self.__handle_message(client_sock, msg)
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __handle_message(self, client_sock, msg):
        """
        Handles a message from the client
        """
        if msg.action == "BET":
            self.__handle_bet_message(client_sock, msg)
        elif msg.action == "FINISH":
            self.__handle_finish_message(client_sock, msg)
        elif msg.action == "WINNER":
            self.__handle_winner_message(client_sock, msg)

    def __handle_winner_message(self, client_sock, msg):
        logging.info(f"action: ganadores | agencies_done: {self._agencies_done}")
        if not all(self._agencies_done):
            logging.info(
                f"action: ganadores | result: fail | error: No todas las agencias han finalizado"
            )
            protocol.send_message(
                client_sock,
                "WINNERWAIT",
                5,  # TODO: mover a constante o configuracion (5 segundos)
            )
        else:
            logging.info(f"action: sorteo | result: success")
            bets = utils.load_bets()
            winners = [
                bet
                for bet in bets
                if utils.has_won(bet) and bet.agency == int(msg.payload)
            ]
            logging.info(
                f"action: ganadores | result: success | winners: {str(winners)}"
            )
            protocol.send_message(
                client_sock, "WINNER", protocol.bets_to_string(winners)
            )

    def __handle_finish_message(self, client_sock, msg):
        self._agencies_done[int(msg.payload) - 1] = True
        logging.info(
            f"action: finalizar_apuestas | result: success | client_id: {msg.payload}"
        )
        protocol.send_ok(client_sock)

    def __handle_bet_message(self, client_sock, msg):
        bets = protocol.bets_from_string(msg.payload)
        utils.store_bets(bets)
        logging.info(
            f"action: apuestas_almacenadas | result: success | client_id: {bets[0].agency}"
        )
        protocol.send_ok(client_sock)

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
