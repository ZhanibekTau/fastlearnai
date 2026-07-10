import json
import logging

from openai import AsyncOpenAI
import openai
from core.config import settings

logger = logging.getLogger(__name__)

SUMMARY_PROMPT = """You are an expert educational content summarizer.

Given the transcript below, generate a structured summary in the SAME language as the transcript (EN/RU/KK).

Return ONLY valid JSON with this exact structure:
{{
  "title": "Video title (infer from content)",
  "overview": "2-3 sentence overview of the video",
  "key_points": [
    {{"topic": "Topic name", "explanation": "Detailed explanation (2-3 paragraphs)"}}
  ],
  "key_takeaways": ["Takeaway 1", "Takeaway 2", "Takeaway 3"],
  "examples": [
    {{"topic": "Topic", "example": "Concrete example from the video"}}
  ],
  "qa": [
    {{"question": "Self-check question", "answer": "Answer from the video"}}
  ]
}}

TRANSCRIPT:
{transcript}
"""


class OpenAIService:
    def __init__(self):
        self.client = AsyncOpenAI(
            api_key=settings.deep_seek_api,
            base_url=settings.deep_seek_host,
        )

    async def generate_summary(self, transcript: str) -> dict:
        tokens_approx = len(transcript) // 4
        logger.info(f"deepseek-chat  tokens≈{tokens_approx}")

        try:
            response = await self.client.chat.completions.create(
                model="deepseek-chat",
                messages=[
                    {
                        "role": "user",
                        "content": SUMMARY_PROMPT.format(transcript=transcript[:100_000]),
                    }
                ],
                response_format={"type": "json_object"},
                max_tokens=4096,
            )
        except openai.RateLimitError as e:
            raise RuntimeError(f"DeepSeek rate limit: {e}")

        usage = response.usage
        logger.info(f"deepseek-chat done  in={usage.prompt_tokens}  out={usage.completion_tokens}")

        return json.loads(response.choices[0].message.content)
