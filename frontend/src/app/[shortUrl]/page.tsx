/*
 * Copyright 2025 Praveen Kumar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { Button } from "@/components/ui/button";
import { LinkIcon } from "lucide-react";
import Link from "next/link";
import { urlApi } from "@/lib/api";

export default function RedirectPage() {
  const params = useParams();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const shortUrl = params.shortUrl as string;

  useEffect(() => {
    const fetchUrlAndRedirect = async () => {
      try {
        const response = await urlApi.redirect(shortUrl);
        window.location.href = response.data;
      } catch (error) {
        setError(
          "The URL you're trying to visit doesn't exist or has been removed.",
        );
      } finally {
        setIsLoading(false);
      }
    };

    fetchUrlAndRedirect();
  }, [shortUrl]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center p-4">
        <LinkIcon className="h-12 w-12 text-primary mb-4" />
        <h1 className="text-2xl font-bold mb-2">Redirecting you...</h1>
        <LoadingSpinner size={32} className="mb-4" />
        <p className="text-muted-foreground text-center">
          You're being redirected to your destination.
        </p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center p-4">
        <LinkIcon className="h-12 w-12 text-primary mb-4" />
        <h1 className="text-2xl font-bold mb-2">URL Not Found</h1>
        <p className="text-muted-foreground text-center mb-6">{error}</p>
        <Link href="/">
          <Button>Return to Homepage</Button>
        </Link>
      </div>
    );
  }
}
