"use client";

import { ThemeToggle } from "@/components/theme-toggle";

export default function Header() {
  return (
    <div className="flex w-full gap-4 border-b-black/25 dark:border-b-white/25 dark:border-b-1 border-b p-4">
      <div className="ml-auto">
        <ThemeToggle />
      </div>
    </div>
  );
}
