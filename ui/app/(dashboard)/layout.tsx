import type { Metadata } from "next";
import { Inter } from "next/font/google";
import Header from "@/components/ui/header";
import { Sidebar } from "@/components/sidebar";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Keyify",
  description: "Api Keys as a Service",
};

export default async function DashboardLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {

  return (
    <main
      className={`${inter.className} bg-zinc-200 dark:bg-zinc-800`}
    >
      <Sidebar />
      <Header />
      {children}
    </main>
  );
}
