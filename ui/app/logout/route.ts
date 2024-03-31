import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

export async function GET(_: Request) {
  cookies().delete("wos-session");
  cookies().delete('keyify-auth');

  return new NextResponse(JSON.stringify({
    "message": "Logged out successfully."
  }), {
    status: 200,
  });
}