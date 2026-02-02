import Image from "next/image";
import MicRecorder from "@/components/MicRecorder";

export default function Home() {
  return (
     <main style={{padding: "2rem"}}>
      <h1>Vocal Genome Engine</h1>
      <MicRecorder/>
     </main>
  );
}
