"use client";

import { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Copy, Check, ExternalLink } from "lucide-react";

export default function OAuthAuthorizePage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [code, setCode] = useState("");
  const [copied, setCopied] = useState(false);
  const [clientId, setClientId] = useState("");
  const [redirectUri, setRedirectUri] = useState("");
  const [state, setState] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Parse OAuth parameters from URL
    const clientIdParam = searchParams.get("client_id");
    const redirectUriParam = searchParams.get("redirect_uri");
    const stateParam = searchParams.get("state");
    const responseType = searchParams.get("response_type");
    const scope = searchParams.get("scope");

    if (!clientIdParam || !redirectUriParam || responseType !== "code") {
      setError("Invalid OAuth request");
      setLoading(false);
      return;
    }

    setClientId(clientIdParam);
    setRedirectUri(redirectUriParam);
    setState(stateParam || "");

    // Generate authorization code
    const authCode = generateAuthCode();
    setCode(authCode);
    setLoading(false);
  }, [searchParams]);

  const generateAuthCode = () => {
    const chars =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let result = "";
    for (let i = 0; i < 32; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy code:", err);
    }
  };

  const handleComplete = () => {
    // Redirect back to CLI with authorization code
    const callbackUrl = `${redirectUri}?code=${code}&state=${state}`;
    window.location.href = callbackUrl;
  };

  const [error, setError] = useState("");

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-2 text-gray-600">Loading authorization...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle className="text-red-600">Authorization Error</CardTitle>
            <CardDescription>
              There was a problem with the authorization request
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center mb-4">
            <svg
              className="w-6 h-6 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              />
            </svg>
          </div>
          <CardTitle>Aether Vault Authorization</CardTitle>
          <CardDescription>
            Enter this code in your CLI to complete authentication
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="auth-code" className="text-sm font-medium">
              Authorization Code
            </Label>
            <div className="flex space-x-2">
              <Input
                id="auth-code"
                value={code}
                readOnly
                className="font-mono text-center text-lg tracking-wider bg-gray-50"
              />
              <Button
                variant="outline"
                size="icon"
                onClick={copyToClipboard}
                className="shrink-0"
              >
                {copied ? (
                  <Check className="h-4 w-4" />
                ) : (
                  <Copy className="h-4 w-4" />
                )}
              </Button>
            </div>
            <p className="text-xs text-gray-500">
              This code will expire in 10 minutes
            </p>
          </div>

          <div className="bg-blue-50 p-3 rounded-lg">
            <div className="flex items-start space-x-2">
              <ExternalLink className="h-4 w-4 text-blue-600 mt-0.5" />
              <div className="text-sm">
                <p className="text-blue-800 font-medium">CLI Instructions</p>
                <p className="text-blue-600">
                  Run this command in your terminal and enter the code when
                  prompted:
                </p>
                <code className="block mt-1 text-xs bg-blue-100 p-1 rounded">
                  vault login --method oauth
                </code>
              </div>
            </div>
          </div>

          <div className="space-y-2 text-xs text-gray-500">
            <div className="flex justify-between">
              <span>Client ID:</span>
              <span className="font-mono">{clientId.substring(0, 8)}...</span>
            </div>
            <div className="flex justify-between">
              <span>Redirect URI:</span>
              <span className="font-mono truncate max-w-[200px]">
                {redirectUri}
              </span>
            </div>
          </div>

          <Button onClick={handleComplete} className="w-full" disabled={!code}>
            Complete Authorization
          </Button>

          <div className="text-center text-xs text-gray-500">
            <p>
              By authorizing, you allow this CLI to access your vault secrets
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
