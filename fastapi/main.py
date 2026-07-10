import time
from datetime import datetime

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware

from core import logging as app_logging
from core.config import settings
from routers import process

app_logging.setup()

app = FastAPI(
    title="FastLearn AI",
    version="1.0.0",
    docs_url="/docs",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:8080"],
    allow_methods=["*"],
    allow_headers=["*"],
)

_METHOD_COLOR = {
    "GET":    "\033[34m",
    "POST":   "\033[32m",
    "PUT":    "\033[33m",
    "PATCH":  "\033[36m",
    "DELETE": "\033[31m",
}

def _status_color(status: int) -> str:
    if status < 300: return "\033[32m"
    if status < 400: return "\033[36m"
    if status < 500: return "\033[33m"
    return "\033[31m"


@app.middleware("http")
async def log_requests(request: Request, call_next):
    if request.url.path == "/health":
        return await call_next(request)

    start = time.perf_counter()
    response = await call_next(request)
    ms = (time.perf_counter() - start) * 1000

    ts  = datetime.now().strftime("%H:%M:%S")
    mc  = _METHOD_COLOR.get(request.method, "\033[2m")
    sc  = _status_color(response.status_code)

    print(
        f"\033[90m{ts}\033[0m  "
        f"{mc}\033[1m{request.method:<7}\033[0m  "
        f"{str(request.url.path):<44}  "
        f"{sc}\033[1m{response.status_code}\033[0m  "
        f"\033[90m{ms:.0f}ms\033[0m",
        flush=True,
    )
    return response


app.include_router(process.router)


@app.get("/health")
async def health():
    return {"status": "ok"}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=settings.port, reload=True)
