import signal
import uvicorn

from fastapi import FastAPI
from fastapi.responses import JSONResponse

app = FastAPI()


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

    def _handle_sigterm(self, sig, frame):
        self.server.should_exit = True

    def run(self):
        self.register_signals()
        self.server.run()


if __name__ == "__main__":
    Server().run()