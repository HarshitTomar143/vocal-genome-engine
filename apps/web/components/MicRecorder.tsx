"use client";

import { useRef, useState } from "react";
import Waveform from "./LiveWaveform";

export default function MicRecorder() {
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const chunksRef = useRef<Blob[]>([]);

  const [isRecording, setIsRecording] = useState(false);
  const [audioURL, setAudioURL] = useState<string | null>(null);
  const [audioBlob, setAudioBlob] = useState<Blob | null>(null);

  async function blobToWav(blob: Blob): Promise<any>{
    const arrayBuffer = await blob.arrayBuffer();
    const audioCtx= new AudioContext();
    const audioBuffer = await audioCtx.decodeAudioData(arrayBuffer);
    return encodeWAV(audioBuffer)
  }

  function encodeWAV(audioBuffer: AudioBuffer): ArrayBuffer {
    const samples = audioBuffer.getChannelData(0);
    const sampleRate = audioBuffer.sampleRate;

    const buffer = new ArrayBuffer(44 + samples.length * 2 );
    const view = new DataView(buffer);

    let offset = 0;
    const write = ()
  }

  async function startRecording() {
    setAudioURL(null);
    chunksRef.current = [];

    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });

    const mediaRecorder = new MediaRecorder(stream);

    mediaRecorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        chunksRef.current.push(event.data);
      }
    };

    mediaRecorder.onstop = () => {
      const blob = new Blob(chunksRef.current);
      console.log("Recorded blog type: ",blob.type)
      const url = URL.createObjectURL(blob);
      setAudioBlob(blob);
      setAudioURL(url);
    };

    mediaRecorder.start();
    mediaRecorderRef.current = mediaRecorder;
    setIsRecording(true);
  }

  function stopRecording() {
    mediaRecorderRef.current?.stop();
    setIsRecording(false);
  }

  async function sendToBackend(){
    if(!audioBlob){
        return;
    }

    const formData = new FormData();
    formData.append("audio", audioBlob)

    const res = await fetch("/api/analyze", {
        method: "POST",
        body: formData,
    });

    const data  = await res.json();
    console.log("DSP result: ", data);
  }

  return (
    <div style={{ padding: "1rem", border: "1px solid #333", borderRadius: 8 }}>
      <h3>Mic Recorder</h3>

      {!isRecording ? (
        <button onClick={startRecording}>🎙️ Start Recording</button>
      ) : (
        <button onClick={stopRecording}>⏹️ Stop Recording</button>
      )}

      {audioURL && (
        <div style={{ marginTop: "1rem" }}>
          <audio controls src={audioURL} />
        </div>
      )}

      {audioBlob && <Waveform audioBlob={audioBlob}/>}

      <button onClick={sendToBackend}>Analyze</button>
    </div>
  );
}
