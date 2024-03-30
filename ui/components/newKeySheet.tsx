import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/select"
import { useParams } from 'next/navigation'
import { useToast } from "@/components/ui/use-toast"

const formSchema = z.object({
  keyName: z.string().min(2, {
    message: "Key name must be at least 2 characters.",
  }),
  prefix: z.string().optional(),
  permissions: z.string().optional(),
})

export function CreateNewKeySheet() {
  const apiId = useParams().apiId;
  const { toast } = useToast();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      keyName: "",
      prefix: "",
      permissions: "Admin",
    },
  })
  async function onSubmit(values: z.infer<typeof formSchema>) {
    toast({
      title: "Creating key...",
      duration: 1000,
    });

    const resp = await fetch(`/api/createApiKey`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        apiId: apiId,
        keyName: values.keyName,
        prefix: values.prefix,
        permissions: values.permissions,
      }),
    })

    if (resp.ok) {
      window.location.reload()
      return
    } else {
      toast({
        title: "Failed to create key. Please try again.",
        duration: 2000,
        variant: "destructive",
      });
    }
  }

  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button className="bg-zinc-400 dark:bg-zinc-900 text-black dark:text-white hover:bg-zinc-700">Create key</Button>
      </SheetTrigger>
      <SheetContent className="w-3/4 bg-zinc-400 dark:bg-zinc-800 text-black dark:text-white">
        <SheetHeader>
          <SheetTitle className="text-black dark:text-white">Create a new API key</SheetTitle>
          <SheetDescription>
            Define the properties of your new API key.
          </SheetDescription>
        </SheetHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 text-white">
            <FormField
              control={form.control}
              name="keyName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Key name</FormLabel>
                  <FormControl>
                    <Input {...field} className="text-black dark:bg-zinc-900 dark:text-white" />
                  </FormControl>
                  <FormDescription>
                    This is your key display name.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="prefix"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Prefix</FormLabel>
                  <FormControl>
                    <Input {...field} className="text-black dark:bg-zinc-900 dark:text-white" />
                  </FormControl>
                  <FormDescription>
                    This is the prefix for the key
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="permissions"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Permissions</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl className="text-black dark:bg-zinc-900 dark:text-white">
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent ref={field.ref}>
                      <SelectItem key="1" value="Admin">
                        Admin
                      </SelectItem>
                      <SelectItem key="2" value="Member">
                        Member
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </FormItem>
              )}
            />
            <Button className="text-black dark:text-white bg-zinc-400 dark:bg-zinc-900" type="submit">Submit</Button>
          </form>
        </Form>
      </SheetContent>
    </Sheet>
  )
}
