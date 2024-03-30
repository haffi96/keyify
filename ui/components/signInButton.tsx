import { Button } from "@/components/ui/button";
import WorkOS from "@workos-inc/node";
import { headers } from "next/headers";

const workos = new WorkOS(process.env.WORKOS_API_KEY);

async function getAuthorizationUrl(returnPathname?: string) {
  const pageHeaders = headers();

  const protocol = pageHeaders.get("x-forwarded-proto") || "http";
  const host = pageHeaders.get("host") || "localhost:3000";
  const origin = `${protocol}://${host}`;
  const callbackRedirectUrl = `${origin}/callback`

  return workos.userManagement.getAuthorizationUrl({
    provider: 'authkit',
    clientId: process.env.WORKOS_CLIENT_ID as string,
    redirectUri: callbackRedirectUrl,
    state: returnPathname ? btoa(JSON.stringify({ returnPathname })) : undefined,
  });
}

export async function SignInButton({ large }: { large?: boolean }) {
  const authorizationUrl = await getAuthorizationUrl();

  return (
    <Button className="bg-zinc-400 dark:bg-zinc-900 text-black dark:text-white">
      <a href={authorizationUrl}>Sign In {large && "with AuthKit"}</a>
    </Button>
  );
}