import  {NextResponse} from "next/server"

export async function POST(req: Request){
    const formData  = await req.formData();
    const audio = formData.get("audio") as Blob;

    const goResp = await fetch("http://localhost:8080/analyze", {
        method: "POST",
        body: audio,
    });

    const data = await goResp.json();
    return NextResponse.json(data);
}