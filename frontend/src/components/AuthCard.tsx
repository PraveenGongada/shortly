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

import { type ReactNode } from "react";
import { Logo } from "@/components/Logo";

interface AuthCardProps {
  title: string;
  children: ReactNode;
  footer: ReactNode;
}

export function AuthCard({ title, children, footer }: AuthCardProps) {
  return (
    <div className="relative flex min-h-full flex-col items-center justify-center px-4 py-10">
      <Logo to="/" className="absolute left-6 top-5" />

      <div className="w-full max-w-sm animate-fade-in">
        <div className="rounded-xl border border-border bg-surface p-6 sm:p-7">
          <div className="mb-6 space-y-1 text-center">
            <h1 className="text-xl font-semibold tracking-tight">{title}</h1>
          </div>
          {children}
        </div>
        <p className="mt-6 text-center text-sm text-muted">{footer}</p>
      </div>
    </div>
  );
}
