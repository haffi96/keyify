import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Keyify",
  description: "Api Keys as a Service",
};

export default async function ApisLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {


  return (
    <div className="flex items-start justify-between">
      <main className="grid h-full w-full pl-[300px]">
        <div className="p-8">{children}</div>
      </main>
    </div>
  );
}
