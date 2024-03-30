"use client";

import { MessageSquare, Settings, Server, Home } from "lucide-react";
import UserItem from "@/components/userItem";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { type User } from '@workos-inc/node';
import { useRouter } from "next/navigation";

export function Sidebar() {
  const router = useRouter();
  const menuList = [
    {
      group: "General",
      items: [
        {
          link: "/",
          icon: <Home />,
          text: "Home",
        },
        {
          link: "/apis",
          icon: <Server />,
          text: "APIs",
        },
        {
          link: "/logs",
          icon: <MessageSquare />,
          text: "Logs",
        },
      ],
    },
    {
      group: "Settings",
      items: [
        {
          link: "/",
          icon: <Settings />,
          text: "General Settings",
        },
      ],
    },
  ];

  return (
    <div className={`bg fixed flex min-h-screen w-[300px] min-w-[300px] flex-col gap-4 p-4
                  text-black dark:text-white bg-zinc-200 border-r-2 border-r-black/25
                  dark:bg-zinc-800 dark:border-r-white/25`}
    >
      <div className="grow">
        {menuList.map((menu, index) => (
          <div className="mx-1 flex flex-col py-2" key={index}>
            <p>{menu.group}</p>
            {menu.items.map((item, i) => (
              <Link className="py-2" key={i} href={item.link}>
                <div className={`flex flex-row hover:text-zinc-400 rounded bg-zinc-300 p-3
                                dark:bg-zinc-900
                                hover:ring-1 hover:ring-black/25`}>
                  {item.icon}
                  <p className="px-2">{item.text}</p>
                </div>
              </Link>
            ))}
          </div>
        ))}
      </div>
      <div>
        <Button
          className="bg-zinc-300 text-black w-full hover:text-white dark:bg-zinc-900 dark:text-white dark:hover:bg-black/25"
          onClick={async () => {
            const resp = await fetch("/logout", {
              headers: {
                "Content-Type": "application/json",
                "Access-Control-Allow-Origin": "*",
              },
            })
            const respJson = await resp.json()

            console.log("resp", respJson);
            console.log("status", resp.status);

            if (resp.ok) {
              router.replace("/");
            }

          }}
        >
          Sign Out
        </Button>
      </div>
    </div>
  );
}
