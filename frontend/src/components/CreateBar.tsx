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

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Check, Copy, Link as LinkIcon } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useCreateLink } from "@/hooks/useLinks";
import { useCopy } from "@/hooks/useCopy";
import { createUrlSchema, type CreateUrlInput } from "@/lib/validators";
import { apiErrorMessage } from "@/lib/api";
import { shortLink } from "@/lib/utils";

export function CreateBar() {
  const createLink = useCreateLink();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateUrlInput>({
    resolver: zodResolver(createUrlSchema),
    defaultValues: { long_url: "" },
  });

  const onSubmit = handleSubmit(async (values) => {
    try {
      const res = await createLink.mutateAsync({ long_url: values.long_url });
      reset({ long_url: "" });

      const link = shortLink(res.short_code);
      const toastId = toast.success("Short link created", {
        action: (
          <CopyToastButton link={link} onDone={() => toast.dismiss(toastId)} />
        ),
      });
    } catch (err) {
      toast.error(apiErrorMessage(err, "Could not shorten URL"));
    }
  });

  return (
    <div>
      <form onSubmit={onSubmit} className="flex flex-col gap-3 sm:flex-row">
        <div className="relative flex-1">
          <LinkIcon className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-faint" />
          <Input
            {...register("long_url")}
            placeholder="Paste a long URL…"
            autoComplete="off"
            autoCapitalize="off"
            spellCheck={false}
            invalid={!!errors.long_url}
            className="pl-9"
            aria-label="Long URL"
          />
        </div>
        <Button
          type="submit"
          loading={createLink.isPending}
          className="justify-center rounded-lg"
        >
          Shrink
        </Button>
      </form>

      {errors.long_url && (
        <p className="mt-2 text-xs text-red-400">{errors.long_url.message}</p>
      )}
    </div>
  );
}

function CopyToastButton({
  link,
  onDone,
}: {
  link: string;
  onDone: () => void;
}) {
  const { copied, copy } = useCopy(1000);

  return (
    <button
      type="button"
      aria-label="Copy short link"
      onClick={() => copy(link, () => setTimeout(onDone, 1000))}
      className="ml-auto rounded-md p-1 text-faint transition hover:bg-white/[0.06] hover:text-fg"
    >
      {copied ? (
        <Check className="h-4 w-4 text-green-400" />
      ) : (
        <Copy className="h-4 w-4" />
      )}
    </button>
  );
}
