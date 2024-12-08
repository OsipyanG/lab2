# Используем официальный образ для Python
FROM python:3.12

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . ./
RUN pip install --no-cache-dir -r requirements.txt


# Запуск приложения
CMD ["python", "client.py"]
