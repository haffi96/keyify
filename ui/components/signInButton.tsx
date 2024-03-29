import { Button } from "@/components/ui/button";
import { clearCookie, getAuthUrl, getUser } from "@/app/auth";

export async function SignInButton({ large }: { large?: boolean }) {
  const { isAuthenticated } = await getUser();


  const authorizationUrl = await getAuthUrl();

  if (isAuthenticated) {
    return (
      <div>
        <form
          action={async () => {
            "use server";
            await clearCookie();
          }}
        >
          <Button type="submit">
            Sign Out
          </Button>
        </form>
      </div>
    );
  }

  return (
    <Button>
      <a href={authorizationUrl}>Sign In {large && "with AuthKit"}</a>
    </Button>
  );
}