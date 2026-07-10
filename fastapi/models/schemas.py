from pydantic import BaseModel
from typing import Optional


class ProcessRequest(BaseModel):
    video_id: str
    file_path: Optional[str] = None    # локальный путь /uploads/video/...
    youtube_url: Optional[str] = None
    language: str = "auto"


class KeyPoint(BaseModel):
    topic: str
    explanation: str


class Example(BaseModel):
    topic: str
    example: str


class QAItem(BaseModel):
    question: str
    answer: str


class SummaryContent(BaseModel):
    overview: str
    key_points: list[KeyPoint]
    key_takeaways: list[str]
    examples: list[Example]


class ProcessResult(BaseModel):
    title: str
    duration_seconds: int
    language: str
    content: SummaryContent
    qa: list[QAItem] = []
    transcript: str
