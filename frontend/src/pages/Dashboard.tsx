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
import { AlertCircle, LinkIcon } from "lucide-react";

import { Navbar } from "@/components/Navbar";
import { CreateBar } from "@/components/CreateBar";
import { LinkRow } from "@/components/LinkRow";
import { Pagination } from "@/components/Pagination";
import { QrDialog } from "@/components/QrDialog";
import { EditDialog } from "@/components/EditDialog";
import { DeleteDialog } from "@/components/DeleteDialog";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { PAGE_SIZE, useLinks } from "@/hooks/useLinks";
import type { UrlResponse } from "@/types";

export function Dashboard() {
  const [page, setPage] = useState(0);
  const [qrCode, setQrCode] = useState<string | null>(null);
  const [editLink, setEditLink] = useState<UrlResponse | null>(null);
  const [deleteLink, setDeleteLink] = useState<UrlResponse | null>(null);

  const { data, isLoading, isError, isFetching, refetch } = useLinks(page * PAGE_SIZE);
  const items = data?.items ?? [];
  const maxRedirects = items.reduce((m, l) => Math.max(m, l.redirects), 0);

  useEffect(() => {
    if (!isLoading && !isError && page > 0 && items.length === 0) {
      setPage((p) => Math.max(0, p - 1));
    }
  }, [isLoading, isError, page, items.length]);

  return (
    <div className="flex min-h-full flex-col">
      <Navbar />

      <main className="w-full flex-1 px-4 py-8 sm:px-6 lg:px-8">
        <CreateBar />

        <h1 className="mb-4 mt-10 text-lg font-semibold tracking-tight">My links</h1>

        <section className="overflow-hidden rounded-xl border border-border bg-surface">
          {isLoading ? (
            <div className="flex items-center justify-center py-16">
              <Spinner />
            </div>
          ) : isError ? (
            <div className="flex flex-col items-center gap-3 py-16 text-center">
              <AlertCircle className="h-6 w-6 text-red-400" />
              <p className="text-sm text-muted">Couldn't load your links.</p>
              <Button variant="secondary" size="sm" onClick={() => void refetch()}>
                Try again
              </Button>
            </div>
          ) : items.length === 0 ? (
            <div className="flex flex-col items-center gap-2 py-16 text-center">
              <div className="flex h-10 w-10 items-center justify-center rounded-full border border-border bg-elevated">
                <LinkIcon className="h-4 w-4 text-faint" />
              </div>
              <p className="text-sm font-medium text-fg">No links yet</p>
              <p className="text-sm text-faint">Paste a URL above to create your first one.</p>
            </div>
          ) : (
            <>
              <div className="divide-y divide-border">
                {items.map((link) => (
                  <LinkRow
                    key={link.id}
                    link={link}
                    maxRedirects={maxRedirects}
                    onQr={(l) => setQrCode(l.short_code)}
                    onEdit={setEditLink}
                    onDelete={setDeleteLink}
                  />
                ))}
              </div>
              <div className="border-t border-border">
                <Pagination
                  page={page}
                  hasNext={data?.hasNext ?? false}
                  disabled={isFetching}
                  onPrev={() => setPage((p) => Math.max(0, p - 1))}
                  onNext={() => setPage((p) => p + 1)}
                />
              </div>
            </>
          )}
        </section>
      </main>

      <QrDialog code={qrCode} onClose={() => setQrCode(null)} />
      <EditDialog link={editLink} onClose={() => setEditLink(null)} />
      <DeleteDialog link={deleteLink} onClose={() => setDeleteLink(null)} />
    </div>
  );
}
