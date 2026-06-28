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

import { toast } from "sonner";

import { Dialog } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { useDeleteLink } from "@/hooks/useLinks";
import { apiErrorMessage } from "@/lib/api";
import { displayUrl } from "@/lib/utils";
import type { UrlResponse } from "@/types";

interface DeleteDialogProps {
  link: UrlResponse | null;
  onClose: () => void;
}

export function DeleteDialog({ link, onClose }: DeleteDialogProps) {
  const deleteLink = useDeleteLink();

  const onConfirm = async () => {
    if (!link) return;
    try {
      await deleteLink.mutateAsync(link.id);
      toast.success("Link deleted");
      onClose();
    } catch (err) {
      toast.error(apiErrorMessage(err, "Could not delete link"));
    }
  };

  return (
    <Dialog
      open={link !== null}
      onOpenChange={(open) => !open && onClose()}
      title="Delete link?"
      description="This can't be undone. The short link will stop working immediately."
    >
      {link && (
        <p className="mb-5 truncate rounded-lg border border-border bg-elevated px-3 py-2 font-mono text-sm text-muted">
          /{link.short_code}
          <span className="text-faint"> → {displayUrl(link.long_url)}</span>
        </p>
      )}
      <div className="flex justify-end gap-2">
        <Button variant="secondary" size="sm" onClick={onClose}>
          Cancel
        </Button>
        <Button variant="danger" size="sm" loading={deleteLink.isPending} onClick={onConfirm}>
          Delete
        </Button>
      </div>
    </Dialog>
  );
}
