import io
import logging
import os
from playsound import playsound # depends on PyObjC
import random
import simpleaudio
import tempfile
import wave

import auth
import diskcache
import google.cloud.texttospeech as tts

VOICES = [
    "cmn-CN-Wavenet-A",
    "cmn-CN-Wavenet-B",
    "cmn-CN-Wavenet-C",
    "cmn-CN-Wavenet-D",
    "cmn-CN-Standard-C",
    "cmn-CN-Standard-B",
    "cmn-CN-Standard-A",
    "cmn-CN-Standard-D",
]

_queue = VOICES[:]

class WavFile:
    def __init__(self, data):
        with tempfile.NamedTemporaryFile(delete=False) as f:
            f.file.write(data)
            f.file.flush()
            self._tempfile = f.name
            logging.debug(self._tempfile)

    def __enter__(self):
        return self

    def __exit__(self):
        os.unlink(self._tempfile)

    def play(self, nonblocking=False):
        return playsound(self._tempfile, not nonblocking)

class Connection:
    def __init__(self):
        creds = auth.creds()
        self._client = tts.TextToSpeechClient(credentials=creds)
        self._cache = diskcache.Cache(".cache")

    def fetchWord(self, word, voice):
        key = ('tts', word, voice)
        audio = self._cache.get(key)
        if not audio:
            logging.info(f'cache miss: {word}')
            response = self._client.synthesize_speech(
                input=tts.SynthesisInput(
                    text=word,
                ),
                voice=tts.VoiceSelectionParams(
                    name=voice,
                    language_code="cmn-CN",
                ),
                audio_config=tts.AudioConfig(
                    audio_encoding=tts.AudioEncoding.LINEAR16,
                ),
            )
            audio = self._cache[key] = response.audio_content
        return audio

    def getWordWav(self, word, voice=None):
        if voice is None:
            voice = _queue[0]
            i = random.randrange(len(_queue) // 2, len(_queue))
            del _queue[0]
            _queue.insert(i, voice)

        # Cache all other voices.
        logging.debug("fetching all voices")
        for v in set(VOICES) - {voice}:
            self.fetchWord(word, v)

        logging.debug(f"{word}: parsing WAV data")
        wav = self.fetchWord(word, voice)
        return WavFile(wav)
