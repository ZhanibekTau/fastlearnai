import logging
import sys
from datetime import datetime

RESET   = "\033[0m"
BOLD    = "\033[1m"
RED     = "\033[31m"
YELLOW  = "\033[33m"
CYAN    = "\033[36m"
GRAY    = "\033[90m"

LEVEL_BADGE = {
    "DEBUG":    f"{GRAY}DBUG{RESET}",
    "INFO":     f"{CYAN}{BOLD}INFO{RESET}",
    "WARNING":  f"{YELLOW}{BOLD}WARN{RESET}",
    "ERROR":    f"{RED}{BOLD}ERRO{RESET}",
    "CRITICAL": f"{RED}{BOLD}CRIT{RESET}",
}


class ColorFormatter(logging.Formatter):
    def format(self, record: logging.LogRecord) -> str:
        ts = datetime.fromtimestamp(record.created).strftime("%H:%M:%S")
        badge = LEVEL_BADGE.get(record.levelname, record.levelname)
        msg = record.getMessage()
        return f"{GRAY}{ts}{RESET}  {badge}  {BOLD}{msg}{RESET}"


def setup() -> None:
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(ColorFormatter())

    root = logging.getLogger()
    root.setLevel(logging.INFO)
    root.handlers = [handler]

    # Silence uvicorn's own access log — we print it ourselves via middleware
    logging.getLogger("uvicorn.access").handlers = []
    logging.getLogger("uvicorn.access").propagate = False

    # Route uvicorn/fastapi logs through our formatter
    for name in ("uvicorn", "uvicorn.error", "fastapi"):
        lg = logging.getLogger(name)
        lg.handlers = [handler]
        lg.propagate = False
