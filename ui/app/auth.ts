import { cookies } from 'next/headers';
import { jwtVerify } from 'jose';
import { WorkOS } from '@workos-inc/node';
import { redirect } from 'next/navigation';
import { sealData, unsealData } from 'iron-session';

export const workos = new WorkOS(process.env.WORKOS_API_KEY);

export function getJwtSecretKey() {
  const secret = process.env.JWT_SECRET_KEY;

  if (!secret) {
    throw new Error('JWT_SECRET_KEY is not set');
  }

  return new Uint8Array(Buffer.from(secret, 'base64'));
}

export async function verifyJwtToken(token: string) {
  try {
    const { payload } = await jwtVerify(token, getJwtSecretKey());
    return payload;
  } catch (error) {
    return null;
  }
}

export async function getAccessToken(): Promise<string | undefined> {
  return cookies().get('token')?.value;
}

export async function getRootKey(): Promise<string | undefined> {
  const key = cookies().get('keyify-auth')?.value;
  if (key) {
    return unsealData(key, {
      password: process.env.ROOTKEY_PASSWORD as string,
    });
  }

}

export async function clearCookie() {
  cookies().delete("token");
  redirect("/");
}