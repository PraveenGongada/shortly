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

import { useState } from "react";
import { Check, Copy, Globe, Pencil, QrCode, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import { useCopy } from "@/hooks/useCopy";
import { displayUrl, shortLink } from "@/lib/utils";
import { toast } from "sonner";
import type { UrlResponse } from "@/types";

function faviconUrl(longUrl: string): string | null {
  try {
    return `${new URL(longUrl).origin}/favicon.ico`;
  } catch {
    return null;
  }
}

interface LinkRowProps {
  link: UrlResponse;
  maxRedirects: number;
  onQr: (link: UrlResponse) => void;
  onEdit: (link: UrlResponse) => void;
  onDelete: (link: UrlResponse) => void;
}

export function LinkRow({ link, maxRedirects, onQr, onEdit, onDelete }: LinkRowProps) {
  const [iconError, setIconError] = useState(false);
  const { copied, copy } = useCopy();

  const full = shortLink(link.short_code);
  const favicon = faviconUrl(link.long_url);
  const pct = maxRedirects > 0 ? (link.redirects / maxRedirects) * 100 : 0;
  const barWidth = link.redirects > 0 ? Math.max(pct, 8) : 0;

  const onCopy = () => copy(full, () => toast.success("Short link copied"));

  return (
    <div className="group flex items-center gap-3 px-4 py-3 transition-colors hover:bg-white/[0.02] sm:gap-4 sm:px-5">
      {favicon && !iconError ? (
        <img
          src={favicon}
          alt=""
          loading="lazy"
          onError={() => setIconError(true)}
          className="h-5 w-5 shrink-0 rounded-sm"
        />
      ) : (
        <Globe className="h-5 w-5 shrink-0 text-faint" />
      )}

      <div className="flex min-w-0 flex-1 flex-col gap-0.5 sm:flex-row sm:items-center sm:gap-3">
        <div className="flex min-w-0 items-center gap-1.5">
          <a
            href={full}
            target="_blank"
            rel="noreferrer noopener"
            className="min-w-0 truncate font-mono text-sm font-medium text-fg hover:underline"
          >
            <span className="sm:hidden">{link.short_code}</span>
            <span className="hidden sm:inline">{displayUrl(full)}</span>
          </a>
          <button
            onClick={onCopy}
            aria-label="Copy short link"
            className="shrink-0 rounded-md p-1 text-faint opacity-100 transition hover:bg-white/[0.06] hover:text-fg focus-visible:opacity-100 sm:opacity-0 sm:group-hover:opacity-100"
          >
            {copied ? (
              <Check className="h-3.5 w-3.5 text-green-400" />
            ) : (
              <Copy className="h-3.5 w-3.5" />
            )}
          </button>
        </div>

        <p
          className="min-w-0 truncate text-sm text-faint sm:flex-1"
          title={link.long_url}
        >
          {displayUrl(link.long_url)}
        </p>
      </div>

      <div
        className="hidden w-32 shrink-0 items-center gap-2 sm:flex"
        title={`${link.redirects} clicks`}
      >
        <div className="h-1.5 flex-1 overflow-hidden rounded-full bg-elevated">
          <div
            className="h-full rounded-full bg-white/70"
            style={{ width: `${barWidth}%` }}
          />
        </div>
        <span className="w-7 text-right text-xs tabular-nums text-muted">
          {link.redirects}
        </span>
      </div>

      <div className="flex shrink-0 items-center">
        <Button variant="ghost" size="icon" onClick={() => onQr(link)} aria-label="Show QR code">
          <QrCode className="h-4 w-4" />
        </Button>
        <Button variant="ghost" size="icon" onClick={() => onEdit(link)} aria-label="Edit destination">
          <Pencil className="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          onClick={() => onDelete(link)}
          aria-label="Delete link"
          className="hover:text-red-400"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}
