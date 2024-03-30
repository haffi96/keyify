export async function getNewRootKey(user_id: string): Promise<string> {
  // Check if the root key is already stored in a cookie
  const baseUrl = process.env.NEXT_PUBLIC_API_URL;
  const response = await fetch(`${baseUrl}/rootKey`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      "workspaceId": user_id,
    }),
  })

  if (!response.ok) {
    throw new Error("Failed to fetch root key");
  }

  const respJson = await response.json();
  const newRootKey = respJson.rootKey as string;
  return newRootKey;
}