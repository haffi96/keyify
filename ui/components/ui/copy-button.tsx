'use client';

import React from "react";
import { Button } from "@/components/ui/button";
import { CheckIcon, CopyIcon } from "@radix-ui/react-icons";
import { cn } from "@/lib/utils";
import { useToast } from "@/components/ui/use-toast";

interface CopyButtonProps extends React.HTMLAttributes<HTMLButtonElement> {
  value: string;
  src?: string;
}

export async function copyToClipboardWithMeta(value: string) {
  navigator.clipboard.writeText(value);
}

export function CopyToClipboardButton({
  value,
  className,
  src,
  ...props
}: CopyButtonProps) {
  const [hasCopied, setHasCopied] = React.useState(false);
  const { toast } = useToast();

  React.useEffect(() => {
    setTimeout(() => {
      setHasCopied(false);
    }, 2000);
  }, [hasCopied]);

  return (
    <Button
      size="icon"
      variant="ghost"
      className={cn(
        "relative z-10 h-6 w-6 text-zinc-50 hover:bg-zinc-700 hover:text-zinc-50",
        className,
      )}
      onClick={() => {
        copyToClipboardWithMeta(
          value,
        );
        setHasCopied(true);
        toast({
          title: "Copied to clipboard",
          duration: 1000,
        });
      }}
      {...props}
    >
      <span className="sr-only">Copy</span>
      {hasCopied ? (
        <CheckIcon className="h-3 w-3" />
      ) : (
        <CopyIcon className="h-3 w-3" />
      )}
    </Button>
  );
}


export function ButtonWithCopy({ text, className }: { text: string, className?: string }) {
  return (
    <div className={cn("flex w-max flex-row items-center space-x-2 rounded bg-zinc-900 p-2 ", className)}>
      <p>{text}</p>
      <CopyToClipboardButton value={text} />
    </div>
  );
}
