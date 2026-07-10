import logging
import os

from openai import AsyncOpenAI
import openai
from core.config import settings

logger = logging.getLogger(__name__)

GROQ_BASE_URL = "https://api.groq.com/openai/v1"


class WhisperService:
    def __init__(self):
        self.client = AsyncOpenAI(
            api_key=settings.groq_api_key,
            base_url=GROQ_BASE_URL,
        )

    async def transcribe(self, audio_path: str, language: str = "auto") -> dict:
        filename = os.path.basename(audio_path)
        size_mb = os.path.getsize(audio_path) / 1_000_000
        logger.info(f"whisper  file={filename}  size={size_mb:.1f}MB  lang={language}")

        if size_mb > 24.5:
            raise ValueError(
                f"Audio file is {size_mb:.1f}MB — Whisper API limit is 25MB. "
                f"Video is too long (approx >{int(size_mb / 0.24)} min at current quality)."
            )

        with open(audio_path, "rb") as f:
            kwargs = {
                "model": "whisper-large-v3",
                "file": f,
                "response_format": "verbose_json",
            }
            if language != "auto":
                kwargs["language"] = language

            try:
                response = await self.client.audio.transcriptions.create(**kwargs)
            except openai.RateLimitError as e:
                raise RuntimeError(f"Groq rate limit: {e}")
            except openai.APIStatusError as e:
                logger.error(f"groq status error  status={e.status_code}  body={e.body}")
                raise RuntimeError(f"Groq API error {e.status_code}: {e.message}")
            except openai.APIConnectionError as e:
                logger.error(f"groq connection error  type={type(e).__name__}  cause={e.__cause__}  detail={e}")
                raise RuntimeError(f"Groq connection failed: {e}")

        words = len(response.text.split())
        logger.info(f"whisper done  lang={response.language}  duration={response.duration:.0f}s  words≈{words}")

        return {
            "text": response.text,
            "language": response.language,
            "duration": response.duration,
            "segments": [
                {"start": s["start"], "end": s["end"], "text": s["text"]}
                for s in (response.segments or [])
            ],
        }
