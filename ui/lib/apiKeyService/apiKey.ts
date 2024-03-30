import { getRootKey } from "@/app/auth";

interface CreateApiKeyProps {
  apiId: string;
  keyName: string;
  prefix?: string;
  permissions?: string;
}


export async function createApiKey(createApiKeyProps: CreateApiKeyProps) {
  const rootKey = await getRootKey();

  const baseUrl = process.env.NEXT_PUBLIC_API_URL;
  const response = await fetch(`${baseUrl}/key`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${rootKey}`,
    },
    body: JSON.stringify({
      "apiId": createApiKeyProps.apiId,
      "name": createApiKeyProps.keyName,
      "prefix": createApiKeyProps.prefix,
      "roles": [createApiKeyProps.permissions],
    }),
  })

  const data = await response.json();
  console.log(data);

  if (!response.ok) {
    throw new Error("Failed to create a new api key");
  }
}