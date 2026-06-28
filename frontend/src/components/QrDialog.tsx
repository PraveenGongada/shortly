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

import { useRef } from "react";
import { QRCodeCanvas } from "qrcode.react";
import { Download } from "lucide-react";
import { toast } from "sonner";

import { Dialog } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { shortLink } from "@/lib/utils";

interface QrDialogProps {
  code: string | null;
  onClose: () => void;
}

export function QrDialog({ code, onClose }: QrDialogProps) {
  const wrapRef = useRef<HTMLDivElement>(null);
  const link = code ? shortLink(code) : "";

  const download = () => {
    const canvas = wrapRef.current?.querySelector("canvas");
    if (!canvas) return;
    try {
      const url = canvas.toDataURL("image/png");
      const a = document.createElement("a");
      a.href = url;
      a.download = `shortly-${code}.png`;
      a.click();
    } catch {
      toast.error("Could not download QR code");
    }
  };

  return (
    <Dialog
      open={code !== null}
      onOpenChange={(open) => !open && onClose()}
      title="QR code"
      description="Scan or download to share this short link."
    >
      <div className="flex flex-col items-center gap-4">
        <div ref={wrapRef} className="rounded-xl bg-white p-4">
          {code && (
            <QRCodeCanvas value={link} size={200} level="M" marginSize={1} />
          )}
        </div>
        <span className="break-all text-center font-mono text-xs text-muted">{link}</span>
        <Button variant="secondary" size="sm" onClick={download} className="w-full justify-center">
          <Download className="h-4 w-4" />
          Download PNG
        </Button>
      </div>
    </Dialog>
  );
}
