/*
 * Copyright 2026 Praveen Kumar
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

import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import { ApiError, urlApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { NotFound } from "@/pages/NotFound";

type Status = "loading" | "notFound" | "error";

export function RedirectShort() {
  const { code } = useParams<{ code: string }>();
  const [status, setStatus] = useState<Status>("loading");

  useEffect(() => {
    if (!code) {
      setStatus("notFound");
      return;
    }
    let active = true;
    setStatus("loading");
    urlApi
      .resolve(code)
      .then((longUrl) => {
        if (!active) return;
        if (longUrl && /^https?:\/\//i.test(longUrl)) {
          window.location.replace(longUrl);
        } else {
          setStatus("notFound");
        }
      })
      .catch((err) => {
        if (!active) return;
        // Only a real 404 means the link doesn't exist; network/5xx is transient.
        setStatus(err instanceof ApiError && err.status === 404 ? "notFound" : "error");
      });
    return () => {
      active = false;
    };
  }, [code]);

  if (status === "notFound") return <NotFound />;

  if (status === "error") {
    return (
      <div className="flex min-h-full flex-col items-center justify-center gap-3 text-center">
        <p className="text-sm text-muted">Couldn't reach the server.</p>
        <Button variant="secondary" size="sm" onClick={() => window.location.reload()}>
          Try again
        </Button>
      </div>
    );
  }

  return (
    <div className="flex min-h-full flex-col items-center justify-center gap-3">
      <Spinner />
      <p className="text-sm text-faint">Redirecting…</p>
    </div>
  );
}
