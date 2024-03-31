import { NextRequest, NextResponse } from 'next/server';
import { createApiKey } from '@/lib/apiKeyService/apiKey';

export async function POST(request: NextRequest) {
  const req = await request.json();
  await createApiKey({
    apiId: req.apiId,
    keyName: req.keyName,
    prefix: req.prefix,
    permissions: req.permissions,
  });
  return NextResponse.json({ message: "success" });
}