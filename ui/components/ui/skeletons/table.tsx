import { Skeleton } from "@/components/ui/skeleton"

export function TableSkeleton() {
  return (
    <div className="flex flex-col space-y-5">
      <Skeleton className="h-10 w-[250px] bg-zinc-700" />
      <Skeleton className="h-[200px] rounded-xl w-full bg-zinc-200 dark:bg-zinc-700" />
    </div>
  )
}
