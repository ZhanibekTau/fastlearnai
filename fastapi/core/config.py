from pathlib import Path
from pydantic_settings import BaseSettings

_ENV_FILE = Path(__file__).resolve().parents[2] / ".env"


class Settings(BaseSettings):
    openai_api_key: str = ""
    groq_api_key: str = ""
    deep_seek_api: str = ""
    deep_seek_host: str = "https://api.deepseek.com"
    gcs_bucket: str = "fastlearn-videos"
    google_application_credentials: str = ""
    max_file_size_mb: int = 500
    port: int = 8001

    model_config = {"env_file": str(_ENV_FILE), "extra": "ignore"}


settings = Settings()
