import socket
import time
import json
import logging
import sys
from pathlib import Path

# Константы
CONFIG_FILE = "configs/client_config.json"

# Настройка логгера
logger = logging.getLogger("client")
stdout_handler = logging.StreamHandler(stream=sys.stdout)
fmt = logging.Formatter(
    "%(name)s: %(asctime)s | %(levelname)s | %(filename)s:%(lineno)s | %(message)s"
)
stdout_handler.setFormatter(fmt)
logger.addHandler(stdout_handler)
logger.setLevel(logging.INFO)

# Загрузка конфигурации
def load_config():
    default_config = {"host": "127.0.0.1", "port": 8080, "interval": 5}
    if not Path(CONFIG_FILE).exists():
        logger.warning(f"Config file not found. Creating default config at {CONFIG_FILE}")
        Path(CONFIG_FILE).parent.mkdir(parents=True, exist_ok=True)
        with open(CONFIG_FILE, "w") as f:
            json.dump(default_config, f, indent=4)

    with open(CONFIG_FILE) as f:
        return json.load(f)

# Основная логика работы клиента
def main():
    config = load_config()
    host = config["host"]
    port = config["port"]
    interval = config["interval"]

    logger.info(f"Starting client with config: {config}")

    while True:
        try:
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as client_socket:
                logger.info(f"Connecting to server at {host}:{port}")
                client_socket.connect((host, port))
                logger.info(f"Connected to server at {host}:{port}")

                while True:
                    message = "Сообщение от клиента"
                    try:
                        # Отправка сообщения
                        send_time = time.time()
                        client_socket.sendall(message.encode())
                        logger.info(f"Message sent: {message}")

                        # Чтение ответа
                        response = client_socket.recv(1024).decode()
                        receive_time = time.time()
                        logger.info(f"Response received: {response}")
                        logger.info(f"Round-trip time: {receive_time - send_time:.2f} seconds")

                        time.sleep(interval)
                    except (socket.error, BrokenPipeError) as e:
                        logger.error(f"Socket error during communication: {e}")
                        break
        except (ConnectionRefusedError, socket.error) as e:
            logger.error(f"Failed to connect to server at {host}:{port}. Error: {e}")
            logger.info(f"Retrying in {interval} seconds...")
            time.sleep(interval)
        except Exception as e:
            logger.exception(f"An unexpected error occurred: {e}")
            break

if __name__ == "__main__":
    main()