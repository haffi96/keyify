import { getRootKey } from "@/app/auth";
import { DataTable } from "@/components/ui/dataTable";

interface GetApiKeysProps {
  apiId: string;
  limit?: number;
  offset?: number;
}

async function getApiKeys(getApiKeysProps: GetApiKeysProps) {
  const baseUrl = process.env.NEXT_PUBLIC_API_URL;
  const rootKey = await getRootKey();

  const url = `${baseUrl}/keys?apiId=${getApiKeysProps.apiId}`;

  const response = await fetch(url, {
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

export default async function apiIdPage({
  params,
}: {
  params: { apiId: string };
}) {
  const { apiId } = params;

  const keys = await getApiKeys({ apiId });
  return (
    <main className="flex">
      <DataTable data={keys} />
    </main>
  );
}
