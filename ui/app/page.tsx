import Link from "next/link";
import { Button } from "@/components/ui/button";
import { SignInButton } from "@/components/signInButton";

export default async function Home() {
  return (
    <main className="flex">
      <div>
        <SignInButton />
        <Link href="/apis">
          <Button>Go to APIs</Button>
        </Link>
      </div>
    </main>
  );
}
