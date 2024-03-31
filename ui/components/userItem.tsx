import { User } from "@workos-inc/node";
import Image from "next/image";

export default function UserItem({ user }: { user: User }) {
  return (
    <div className="flex items-center justify-between gap-2 rounded-[8px] p-1 text-black/50 dark:text-white">
      <Image src={user.profilePictureUrl!} alt="" width="50" height="50" className="rounded-full" />
      <div className="grow">
        <p className="text-[16px] font-bold">{user?.firstName}</p>
        <p className="text-[12px] text-neutral-500">
          {user?.email}
        </p>
      </div>
    </div>
  );
}
