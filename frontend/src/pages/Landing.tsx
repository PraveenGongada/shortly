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
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Logo } from "@/components/Logo";
import { LogoMark } from "@/components/LogoMark";

const LINES = ["Shorten links", "Track every click", "Grab a QR code"];

export function Landing() {
  const glowRef = useRef<HTMLDivElement>(null);

  const onMouseMove = (e: React.MouseEvent) => {
    const el = glowRef.current;
    if (!el) return;
    el.style.setProperty("--mx", `${e.clientX}px`);
    el.style.setProperty("--my", `${e.clientY}px`);
  };

  return (
    <div className="relative flex min-h-full flex-col" onMouseMove={onMouseMove}>
      <div
        ref={glowRef}
        aria-hidden
        className="pointer-events-none fixed inset-0 -z-10"
        style={{
          background:
            "radial-gradient(600px circle at var(--mx, 50%) var(--my, 50%), rgba(255,255,255,0.06) 0%, transparent 60%)",
        }}
      />

      <header className="absolute inset-x-0 top-0 z-20 flex h-16 items-center justify-between px-4 lg:px-6">
        <Logo to="/" />
        <div className="flex items-center gap-2">
          <Link to="/login" className="flex">
            <Button variant="secondary" size="sm" className="min-w-[84px] justify-center rounded-md">
              Log in
            </Button>
          </Link>
          <Link to="/register" className="flex">
            <Button size="sm" className="min-w-[84px] justify-center rounded-md">
              Sign up
            </Button>
          </Link>
        </div>
      </header>

      <main className="relative mx-auto flex w-full max-w-[1440px] flex-1 flex-col items-center justify-center gap-12 px-6 py-24 lg:flex-row lg:justify-between lg:gap-0 lg:px-8 lg:py-10">
        <div className="z-0 flex w-full items-center justify-center lg:absolute lg:inset-0 lg:w-auto">
          <LogoMark />
        </div>

        <div className="relative z-10 flex w-full max-w-[444px] animate-fade-in flex-col items-center gap-6 text-center lg:items-start lg:gap-8 lg:text-left">
          <h1 className="my-0 text-[48px] font-normal leading-[1.05] tracking-tight sm:text-[64px]">
            Your links,
            <br />
            shortened.
          </h1>
          <div className="flex items-center gap-3">
            <Link to="/register" className="flex">
              <Button className="h-10 min-w-[104px] justify-center rounded-full">Sign up</Button>
            </Link>
            <Link to="/login" className="flex">
              <Button variant="secondary" className="h-10 min-w-[104px] justify-center rounded-full">
                Log in
              </Button>
            </Link>
          </div>
        </div>

        <nav className="relative z-10 flex w-full max-w-[364px] animate-fade-in flex-col items-center gap-3 lg:items-start lg:text-left lg:-mr-8">
          {LINES.map((line) => (
            <span
              key={line}
              className="font-mono text-sm font-semibold uppercase leading-[1.4] tracking-[1px] text-white"
            >
              {line}
            </span>
          ))}
        </nav>
      </main>
    </div>
  );
}
