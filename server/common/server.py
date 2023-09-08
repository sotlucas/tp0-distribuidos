import errno
import socket
import logging
import signal
import threading

import common.utils as utils
import common.protocol as protocol


class Server:
    def __init__(self, config_params):
        self._running = True
        signal.signal(signal.SIGTERM, self.__shutdown)
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(("", config_params["port"]))
        self._server_socket.listen(config_params["listen_backlog"])
        self._agencies = config_params["agencies"]
        self._agencies_lock = threading.Lock()
        self._file_lock = threading.Lock()
        self._winner_wait_time_seconds = config_params["winner_wait_time_seconds"]

    def run(self):
        """
        Server main loop. Accepts new connections and handles each one in a separate thread.
        """
        threads = []
        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                t = threading.Thread(
                    target=self.__handle_client_connection, args=(client_sock,)
                )
                t.start()
                threads.append(t)
            except OSError:
                break
        for t in threads:
            t.join()

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        connected = True
        while connected:
            try:
                msg = protocol.read_message(client_sock)
                if msg is None:
                    connected = False
                else:
                    self.__handle_message(client_sock, msg)
            except ConnectionError:
                connected = False
            except OSError as e:
                if e.errno != errno.ECONNRESET:
                    logging.error(
                        f"action: receive_message | result: fail | error: {e}"
                    )
                connected = False
            except Exception as e:
                logging.error(f"action: receive_message | result: fail | error: {e}")
                connected = False
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
        self._agencies_lock.acquire()
        logging.info(
            f"action: consulta_ganadores | agencies_restantes: {self._agencies}"
        )
        if self._agencies > 0:
            logging.info(
                f"action: consulta_ganadores | result: in_progress | client_id: {msg.payload} | msg: No todas las agencias han finalizado"
            )
            protocol.send_message(
                client_sock,
                "WINNERWAIT",
                self._winner_wait_time_seconds,
            )
        else:
            self._file_lock.acquire()
            bets = utils.load_bets()
            self._file_lock.release()
            winners = [
                bet
                for bet in bets
                if utils.has_won(bet) and bet.agency == int(msg.payload)
            ]
            logging.info(
                f"action: consulta_ganadores | result: success | winners: {str(winners)}"
            )
            protocol.send_message(
                client_sock, "WINNER", protocol.bets_to_string(winners)
            )
        self._agencies_lock.release()

    def __handle_finish_message(self, client_sock, msg):
        self._agencies_lock.acquire()
        self._agencies -= 1
        logging.info(
            f"action: finalizar_apuestas | result: success | client_id: {msg.payload}"
        )
        protocol.send_ok(client_sock)
        if self._agencies <= 0:
            logging.info(f"action: sorteo | result: success")
        self._agencies_lock.release()

    def __handle_bet_message(self, client_sock, msg):
        bets = protocol.bets_from_string(msg.payload)
        self._file_lock.acquire()
        utils.store_bets(bets)
        self._file_lock.release()
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
