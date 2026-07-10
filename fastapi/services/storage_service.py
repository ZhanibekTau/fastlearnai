import os
import tempfile
from google.cloud import storage
from core.config import settings


class StorageService:
    def __init__(self):
        self.client = storage.Client()
        self.bucket = self.client.bucket(settings.gcs_bucket)

    async def download(self, storage_key: str) -> str:
        tmp = tempfile.NamedTemporaryFile(delete=False, suffix=".mp4")
        blob = self.bucket.blob(storage_key)
        blob.download_to_filename(tmp.name)
        return tmp.name

    async def upload(self, local_path: str, destination_key: str) -> str:
        blob = self.bucket.blob(destination_key)
        blob.upload_from_filename(local_path)
        return f"gs://{settings.gcs_bucket}/{destination_key}"

    def delete(self, storage_key: str) -> None:
        blob = self.bucket.blob(storage_key)
        blob.delete()
