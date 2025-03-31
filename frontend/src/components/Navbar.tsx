/*
 * Copyright 2025 Praveen Kumar
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

"use client";

import { useAuth } from "@/contexts/AuthContext";
import Link from "next/link";
import { Button } from "./ui/button";
import { LogOut, Link as LinkIcon } from "lucide-react";

export default function Navbar({
  isAuthenticated,
}: {
  isAuthenticated: boolean;
}) {
  const { logout, isLoading } = useAuth();

  return (
    <header className="bg-background border-b">
      <div className="container flex h-16 items-center justify-between">
        <div className="flex items-center gap-2">
          <LinkIcon className="h-6 w-6" />
          <Link
            href={isAuthenticated ? "/dashboard" : "/"}
            className="text-xl font-bold"
          >
            Shortly
          </Link>
        </div>
        <nav className="flex items-center gap-4">
          {isLoading ? (
            <></>
          ) : isAuthenticated ? (
            <Button
              variant="ghost"
              onClick={() => logout()}
              className="flex items-center gap-2"
            >
              <LogOut className="h-4 w-4" />
              Logout
            </Button>
          ) : (
            <div className="flex items-center gap-2">
              <Link href="/login">
                <Button variant="ghost">Login</Button>
              </Link>
              <Link href="/register">
                <Button>Register</Button>
              </Link>
            </div>
          )}
        </nav>
      </div>
    </header>
  );
}
