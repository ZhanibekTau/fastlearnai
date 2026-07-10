import asyncio
import logging
import os
import tempfile

logger = logging.getLogger(__name__)


async def download_youtube(url: str) -> tuple[str, int]:
    logger.info(f"yt-dlp  url={url[:70]}")

    tmp_dir = tempfile.mkdtemp()
    output_path = os.path.join(tmp_dir, "audio.mp3")

    cmd = [
        "yt-dlp",
        "--extract-audio",
        "--audio-format", "mp3",
        "--audio-quality", "9",
        # Re-encode to 16kHz mono 16kbps CBR — Whisper's native rate, keeps file < 25MB for ~200 min videos
        "--postprocessor-args", "ffmpeg:-ar 16000 -ac 1 -b:a 16k -c:a libmp3lame",
        "--output", output_path,
        "--no-playlist",
        "--extractor-args", "youtube:player_client=web,default",
        "--retries", "3",
        "--socket-timeout", "30",
        url,
    ]

    proc = await asyncio.create_subprocess_exec(
        *cmd,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )
    _, stderr = await proc.communicate()

    if proc.returncode != 0:
        raw = stderr.decode().strip()
        # Extract the last meaningful ERROR line from yt-dlp output
        error_lines = [l for l in raw.splitlines() if "ERROR" in l]
        err = error_lines[-1] if error_lines else raw.splitlines()[-1] if raw else "unknown error"
        raise RuntimeError(f"yt-dlp: {err}")

    duration = await _get_duration(output_path)
    size_mb = os.path.getsize(output_path) / 1_000_000
    logger.info(f"yt-dlp done  duration={duration}s  size={size_mb:.1f}MB")

    return output_path, duration


async def _get_duration(path: str) -> int:
    proc = await asyncio.create_subprocess_exec(
        "ffprobe", "-v", "quiet", "-show_entries", "format=duration",
        "-of", "default=noprint_wrappers=1:nokey=1", path,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )
    stdout, _ = await proc.communicate()
    try:
        return int(float(stdout.decode().strip()))
    except (ValueError, AttributeError):
        return 0
