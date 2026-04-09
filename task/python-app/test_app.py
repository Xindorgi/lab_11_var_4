import pytest
from fastapi.testclient import TestClient

from app import app

client = TestClient(app)


def test_ping_status():
    response = client.get("/ping")
    assert response.status_code == 200


def test_ping_body():
    response = client.get("/ping")
    assert response.json() == {"message": "pong"}


def test_health_status():
    response = client.get("/health")
    assert response.status_code == 200


def test_health_body():
    response = client.get("/health")
    assert response.json() == {"status": "ok"}


def test_ping_content_type():
    response = client.get("/ping")
    assert "application/json" in response.headers["content-type"]


def test_health_content_type():
    response = client.get("/health")
    assert "application/json" in response.headers["content-type"]