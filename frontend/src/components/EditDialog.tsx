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

import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";

import { Dialog } from "@/components/ui/dialog";
import { Field } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useUpdateLink } from "@/hooks/useLinks";
import { updateUrlSchema, type UpdateUrlInput } from "@/lib/validators";
import { apiErrorMessage } from "@/lib/api";
import type { UrlResponse } from "@/types";

interface EditDialogProps {
  link: UrlResponse | null;
  onClose: () => void;
}

export function EditDialog({ link, onClose }: EditDialogProps) {
  const updateLink = useUpdateLink();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<UpdateUrlInput>({
    resolver: zodResolver(updateUrlSchema),
    defaultValues: { new_url: "" },
  });

  useEffect(() => {
    if (link) reset({ new_url: link.long_url });
  }, [link, reset]);

  const onSubmit = handleSubmit(async (values) => {
    if (!link) return;
    const next = values.new_url;
    if (next === link.long_url) {
      onClose();
      return;
    }
    try {
      await updateLink.mutateAsync({ id: link.id, new_url: next });
      toast.success("Destination updated");
      onClose();
    } catch (err) {
      toast.error(apiErrorMessage(err, "Could not update link"));
    }
  });

  return (
    <Dialog
      open={link !== null}
      onOpenChange={(open) => !open && onClose()}
      title="Edit destination"
      description={link ? `/${link.short_code}` : undefined}
    >
      <form onSubmit={onSubmit} className="space-y-4">
        <Field label="Destination URL" htmlFor="new_url" error={errors.new_url?.message}>
          <Input
            id="new_url"
            {...register("new_url")}
            placeholder="https://example.com/…"
            autoComplete="off"
            spellCheck={false}
            invalid={!!errors.new_url}
          />
        </Field>
        <div className="flex justify-end gap-2">
          <Button type="button" variant="secondary" size="sm" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" size="sm" loading={updateLink.isPending}>
            Save changes
          </Button>
        </div>
      </form>
    </Dialog>
  );
}
