import { NextRequest, NextResponse } from "next/server";
import { getAuthUrl, getUser, verifyJwtToken } from "@/app/auth";

async function getRootKey(user_id: string): Promise<string> {
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
  return respJson.rootKey;
}

export async function middleware(request: NextRequest) {
  const { cookies } = request;
  const { value: token } = cookies.get("token") ?? { value: null };

  const hasVerifiedToken = token && (await verifyJwtToken(token));

  // Redirect unauthenticated users to the AuthKit flow
  if (!hasVerifiedToken) {
    const authorizationUrl = await getAuthUrl();
    const response = NextResponse.redirect(authorizationUrl);

    response.cookies.delete("token");
    response.cookies.delete("rootKey");

    return response;
  }

  // Get user
  const { user } = await getUser();

  const nextResponse = NextResponse.next();
  // Get Root key
  if (!request.cookies.get("rootKey")) {
    const rootKey = await getRootKey(user?.id!);

    nextResponse.cookies.set("rootKey", rootKey);
  }

  return nextResponse;
}

// Match against the account page
export const config = { matcher: ['/apis/:path*'] };