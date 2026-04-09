ARG BASE=python:3.12-alpine
FROM ${BASE}

RUN apk add --no-cache curl 2>/dev/null || apt-get update && apt-get install -y curl

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 8000

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD curl --fail http://localhost:8000/health || exit 1

CMD ["python", "app.py"]