'use client';

import { BellIcon, Cookie, CreditCard, Inbox, MessageSquare, Settings, User } from "lucide-react";
import UserItem from "@/components/userItem";
import { Command, CommandGroup, CommandItem, CommandList } from "./ui/command";

export function Sidebar() {
  const menuList = [
    {
      group: "General",
      items: [
        {
          link: "/",
          icon: <User />,
          text: "Profile"
        },
        {
          link: "/",
          icon: <Inbox />,
          text: "Inbox"
        },
        {
          link: "/",
          icon: <CreditCard />,
          text: "Billing"
        },
        {
          link: "/",
          icon: <BellIcon />,
          text: "Notifications"
        }
      ]
    },
    {
      group: "Settings",
      items: [
        {
          link: "/",
          icon: <Settings />,
          text: "General Settings"
        },
        {
          link: "/",
          icon: <Cookie />,
          text: "Privacy"
        },
        {
          link: "/",
          icon: <MessageSquare />,
          text: "Logs"
        }
      ]
    }
  ]

  return <div className="fixed bg flex flex-col gap-4 w-[300px] min-w-[300px] border-r min-h-screen p-4 text-white">
    <div>
      <UserItem />
    </div>
    <div className="grow bg-zinc-600 rounded">
      <Command className="bg-zinc-600" style={{ overflow: 'visible' }}>
        <CommandList className="text-white" style={{ overflow: 'visible' }}>
          {menuList.map((menu: any, key: number) => (
            <CommandGroup key={key} heading={menu.group}>
              {menu.items.map((option: any, optionKey: number) =>
                <CommandItem key={optionKey} className="text-white flex gap-2 cursor-pointer">
                  {option.icon}
                  {option.text}
                </CommandItem>
              )}
            </CommandGroup>
          ))}
        </CommandList>
      </Command>
    </div>
    <div>Settings / Notifications</div>
  </div>;
}