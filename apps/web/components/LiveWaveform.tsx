"use client";

import { useEffect, useRef } from "react";

type Props = {
  audioBlob: Blob | null;
};

export default function Waveform({ audioBlob }: Props) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  useEffect(() => {
    if (!audioBlob || !canvasRef.current) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    const audioCtx = new AudioContext();

    audioBlob.arrayBuffer().then((buffer) => {
      audioCtx.decodeAudioData(buffer).then((audioBuffer) => {
        const data = audioBuffer.getChannelData(0); // mono
        drawWaveform(data, ctx, canvas.width, canvas.height);
        console.log("Decoded samples:", audioBuffer.length);
        console.log("Sample rate:", audioBuffer.sampleRate);

      });
    });

    


    return () => {
      audioCtx.close();
    };
  }, [audioBlob]);

  return (
    <canvas
      ref={canvasRef}
      width={600}
      height={150}
      style={{ border: "1px solid #333", marginTop: "1rem" }}
    />
  );
}

function drawWaveform(
  data: Float32Array,
  ctx: CanvasRenderingContext2D,
  width: number,
  height: number
) {
  ctx.clearRect(0, 0, width, height);

  // Midline
  ctx.strokeStyle = "#444";
  ctx.beginPath();
  ctx.moveTo(0, height / 2);
  ctx.lineTo(width, height / 2);
  ctx.stroke();

  ctx.strokeStyle = "#00ffcc";
  ctx.lineWidth = 1;

  const step = Math.ceil(data.length / width);
  const amp = height / 2;

  ctx.beginPath();

  for (let i = 0; i < width; i++) {
    let min = 1.0;
    let max = -1.0;

    for (let j = 0; j < step; j++) {
      const datum = data[i * step + j];
      if (datum < min) min = datum;
      if (datum > max) max = datum;
    }

    ctx.moveTo(i, amp + min * amp);
    ctx.lineTo(i, amp + max * amp);
  }

  ctx.stroke();
}
