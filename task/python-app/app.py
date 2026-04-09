import signal
import uvicorn
import os
import httpx

from fastapi import FastAPI
from fastapi.responses import JSONResponse

app = FastAPI()


@app.get("/check-go")
async def check_go():
    go_url = os.getenv("GO_APP_URL", "http://go-app:8080/time")
    async with httpx.AsyncClient() as client:
        try:
            response = await client.get(go_url, timeout=5.0)
            return {"status": "ok", "go_response": response.json()}
        except Exception as e:
            return {"status": "error", "message": str(e)}


@app.get("/ping")
def ping():
    return {"message": "pong"}


@app.get("/health")
def health():
    return JSONResponse(status_code=200, content={"status": "ok"})


class Server:
    def __init__(self):
        self.config = uvicorn.Config(app, host="0.0.0.0", port=8000)
        self.server = uvicorn.Server(self.config)

    def register_signals(self):
        signal.signal(signal.SIGTERM, self._handle_sigterm)
        signal.signal(signal.SIGINT, self._handle_sigterm)

    def _handle_sigterm(self, sig, frame):
        self.server.should_exit = True

    def run(self):
        self.register_signals()
        self.server.run()


if __name__ == "__main__":
    Server().run()