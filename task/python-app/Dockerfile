ARG BASE=python:3.12-alpine
FROM ${BASE}

WORKDIR /app

RUN if [ -f /etc/alpine-release ]; then \
        apk add --no-cache curl; \
    else \
        apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*; \
    fi

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 8000

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD curl --fail http://localhost:8000/health || exit 1

CMD ["python", "app.py"]