import { Button } from "@/components/ui/button";
import { getSignInUrl } from "@workos-inc/authkit-nextjs";

export async function SignInButton({ large }: { large?: boolean }) {
  const authorizationUrl = await getSignInUrl();

  return (
    <Button className="bg-zinc-400 dark:bg-zinc-900 text-black dark:text-white">
      <a href={authorizationUrl}>Sign In {large && "with AuthKit"}</a>
    </Button>
  );
}