import { SignInButton } from "@/components/signInButton";

export default async function Home() {
  return (
    <main className="flex">
      <div className="m-auto text-center align-middle items-center p-20">
        <SignInButton />
      </div>
    </main>
  );
}
