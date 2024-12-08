import socket
import time
import json
import structlog
from pathlib import Path

structlog.configure(
    processors=[
        # structlog.processors.TimeStamper(),
        # structlog.processors.add_log_level,
        structlog.processors.TimeStamper(fmt="%Y-%m-%dT%H:%M:%S", utc=False),
        structlog.processors.JSONRenderer()
    ],
    cache_logger_on_first_use=True
)

logger = structlog.get_logger()

# Загрузка конфигурации
CONFIG_FILE = "configs/client_config.json"

if not Path(CONFIG_FILE).exists():
    default_config = {"host": "127.0.0.1", "port": 8080, "interval": 5}
    with open(CONFIG_FILE, "w") as f:
        json.dump(default_config, f, indent=4)

with open(CONFIG_FILE) as f:
    config = json.load(f)

HOST = config["host"]
PORT = config["port"]
INTERVAL = config["interval"]

def main():
    try:
        logger.info("Starting client")
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as client_socket:
            logger.info("Connecting to server", addr=f"{HOST}:{PORT}")
            client_socket.connect((HOST, PORT))
            logger.info("Connected to server", addr=f"{HOST}:{PORT}")

            while True:
                message = "Сообщение от клиента" 
                send_time = time.time()
                client_socket.sendall(message.encode())
                logger.info("Message sent")
                response = client_socket.recv(1024).decode()
                receive_time = time.time()
                logger.info("Message received", response)

                time.sleep(INTERVAL)
    except Exception as e:
        logger.error("Socket closed", error=str(e))

if __name__ == "__main__":
    main()