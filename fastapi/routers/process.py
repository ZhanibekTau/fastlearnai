import logging
import os
import time

from fastapi import APIRouter, HTTPException
from models.schemas import ProcessRequest, ProcessResult, SummaryContent, KeyPoint, Example, QAItem
from services.whisper_service import WhisperService
from services.openai_service import OpenAIService
from services.youtube_service import download_youtube

router = APIRouter()
logger = logging.getLogger(__name__)

whisper = WhisperService()
openai_svc = OpenAIService()


@router.post("/process", response_model=ProcessResult)
async def process_video(req: ProcessRequest):
    tmp_audio = None
    audio_path = None
    duration_seconds = 0
    start = time.perf_counter()

    src = f"youtube:{req.youtube_url[:50]}" if req.youtube_url else f"file:{os.path.basename(req.file_path or '')}"
    logger.info(f"job started  video={req.video_id}  src={src}  lang={req.language}")

    try:
        if req.youtube_url:
            audio_path, duration_seconds = await download_youtube(req.youtube_url)
            tmp_audio = audio_path
        elif req.file_path:
            if not os.path.exists(req.file_path):
                raise HTTPException(status_code=404, detail=f"File not found: {req.file_path}")
            audio_path = req.file_path
        else:
            raise HTTPException(status_code=400, detail="Provide file_path or youtube_url")

        transcription = await whisper.transcribe(audio_path, req.language)
        transcript_text = transcription["text"]
        language = transcription["language"]
        if duration_seconds == 0:
            duration_seconds = int(transcription.get("duration", 0))

        summary_data = await openai_svc.generate_summary(transcript_text)

        content = SummaryContent(
            overview=summary_data.get("overview", ""),
            key_points=[KeyPoint(**kp) for kp in summary_data.get("key_points", [])],
            key_takeaways=summary_data.get("key_takeaways", []),
            examples=[Example(**ex) for ex in summary_data.get("examples", [])],
        )
        qa = [QAItem(**q) for q in summary_data.get("qa", [])]

        elapsed = time.perf_counter() - start
        title = summary_data.get("title", "Untitled")
        logger.info(f"job done  video={req.video_id}  title={title!r}  elapsed={elapsed:.1f}s")

        return ProcessResult(
            title=title,
            duration_seconds=duration_seconds,
            language=language,
            content=content,
            qa=qa,
            transcript=transcript_text,
        )

    except HTTPException:
        raise
    except Exception as exc:
        elapsed = time.perf_counter() - start
        logger.error(f"job failed  video={req.video_id}  elapsed={elapsed:.1f}s  error={exc}")
        raise HTTPException(status_code=500, detail=str(exc))
    finally:
        if tmp_audio and os.path.exists(tmp_audio):
            os.unlink(tmp_audio)
