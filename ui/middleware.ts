import { NextRequest, NextResponse } from "next/server";
import { updateSession } from "@workos-inc/authkit-nextjs/dist/cjs/session";
import { cookies } from "next/headers";
import { getNewRootKey } from "@/lib/apiKeyService/rootKey";
import { sealData, unsealData } from 'iron-session';
import { type Session } from "@workos-inc/authkit-nextjs/src/interfaces"

export async function middleware(request: NextRequest, next: NextResponse) {
  if (!cookies().get("wos-session")) {
    const protocol = request.headers.get("x-forwarded-proto") || "http";
    const host = request.headers.get("host") || "localhost";
    const url = `${protocol}://${host}/login`;
    return NextResponse.redirect(url);
  }

  const response = await updateSession(request, false);

  const sessionToken = cookies().get("wos-session")

  const sessionData = await unsealData<Session>(sessionToken?.value!, {
    password: process.env.WORKOS_COOKIE_PASSWORD as string,
  });

  const user = sessionData?.user;

  // Check if root key in cookies
  if (!cookies().get("keyify-auth")) {
    const rootKey = await getNewRootKey(user?.id as string);
    const encryptedRootKey = await sealData(rootKey, {
      password: process.env.ROOTKEY_PASSWORD as string,
    });
    response.cookies.set("keyify-auth", encryptedRootKey, {
      httpOnly: true,
      secure: true,
      sameSite: "lax",
    });
  }

  return response
}

// Match against the pages
export const config = { matcher: ['/((?!api|callback|login|home|_next/static|_next/image|favicon.ico).*)', "/apis/:path*"] };