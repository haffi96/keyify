import Link from "next/link";
import { getRootKey } from "@/app/auth";

interface Api {
  apiId: string;
  apiName: string;
  createdAt: string;
}

async function getApis(): Promise<Api[]> {
  const baseUrl = process.env.NEXT_PUBLIC_API_URL;
  const rootKey = await getRootKey();

  const response = await fetch(`${baseUrl}/apis`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${rootKey}`,
    },
  });

  if (!response.ok) {
    throw new Error("Failed to fetch API keys");
  }

  const respJson = await response.json();

  return respJson;
}

export default async function apisPage() {
  const apis = await getApis();


  return (
    <main className="flex">
      <div className="flex flex-row space-x-2">
        {apis.map((api) => (
          <Link key={api.apiId} href={`/apis/${api.apiId}`}>
            <div
              key={api.apiId}
              className="cursor-pointer rounded-lg bg-gray-200 p-10 hover:bg-gray-300"
            >
              {api.apiName}
            </div>
          </Link>
        ))}
      </div>
    </main>
  );
}
