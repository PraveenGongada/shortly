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

import { Link as LinkIcon } from "lucide-react";
import { cn } from "@/lib/utils";

export function LogoMark({ className }: { className?: string }) {
  return (
    <div className={cn("relative flex items-center justify-center", className)}>
      <div
        aria-hidden
        className="pointer-events-none absolute h-48 w-48 rounded-full blur-2xl sm:h-56 sm:w-56"
        style={{
          background:
            "radial-gradient(circle, rgba(255,255,255,0.10) 0%, transparent 70%)",
        }}
      />
      <LinkIcon
        className="relative h-24 w-24 text-white sm:h-28 sm:w-28"
        strokeWidth={1.5}
        style={{
          filter:
            "drop-shadow(0 0 6px rgba(255,255,255,0.4)) drop-shadow(0 0 16px rgba(255,255,255,0.2))",
        }}
      />
    </div>
  );
}
