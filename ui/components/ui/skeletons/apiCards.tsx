import { Skeleton } from "@/components/ui/skeleton";

export function ApiCardsSkeleton() {
  return (
    <div className="flex flex-row space-x-2">
      {
        [1, 2].map((_, index) => (
          <Skeleton key={index} className="rounded-lg w-1/6 bg-zinc-200 shadow-sm shadow-black/50 dark:text-white dark:bg-zinc-700 p-10" />
        ))
      }
    </div>
  )
}