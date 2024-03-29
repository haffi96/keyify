"use client";

import { MessageSquare, Settings, Server } from "lucide-react";
import UserItem from "@/components/userItem";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { type User } from '@workos-inc/node';

export function Sidebar({ user }: { user: User }) {
  const menuList = [
    {
      group: "General",
      items: [
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
    <div className="bg fixed flex min-h-screen w-[300px] min-w-[300px] flex-col gap-4 border-r p-4 text-white">
      <div>
        <UserItem user={user} />
      </div>
      <div className="grow rounded bg-zinc-600">
        {menuList.map((menu, index) => (
          <div className="mx-4 flex flex-col py-2" key={index}>
            <p className="text-zinc-400">{menu.group}</p>
            {menu.items.map((item, i) => (
              <div className="flex flex-col p-2" key={i}>
                <Link key={i} href={item.link}>
                  <div className="flex flex-row rounded bg-zinc-500 p-2 hover:text-zinc-800">
                    {item.icon}
                    <p className="px-2">{item.text}</p>
                  </div>
                </Link>
              </div>
            ))}
          </div>
        ))}
      </div>
      <div>
        <Button
          onClick={async () => {
            const resp = await fetch("/logout", {
              headers: {
                "Content-Type": "application/json",
              },
            })
            if (resp.ok) {
              window.location.href = "/";
            }

          }}
        >
          Sign Out
        </Button>
      </div>
    </div>
  );
}
