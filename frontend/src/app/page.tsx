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
import { Button } from "@/components/ui/button";
import Link from "next/link";
import {
  ArrowRight,
  Link as LinkIcon,
  BarChart3,
  Clock,
  Shield,
} from "lucide-react";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import Navbar from "@/components/Navbar";

export default function HomePage() {
  const { isLoading } = useAuth();

  if (isLoading) {
    return (
      <>
        <Navbar isAuthenticated={false} />
        <div className="flex items-center justify-center min-h-[calc(100vh-4rem)]">
          <LoadingSpinner size={32} />
        </div>
      </>
    );
  }

  return (
    <>
      <Navbar isAuthenticated={false} />
      <div className="container py-12 px-4 md:px-6 lg:px-8">
        <div className="flex flex-col items-center text-center space-y-4 py-8 md:py-12">
          <LinkIcon className="h-12 w-12 text-primary" />
          <h1 className="text-3xl md:text-5xl font-bold tracking-tighter">
            Simplify Your Links
          </h1>
          <p className="text-muted-foreground md:text-xl max-w-[600px] mx-auto">
            Create short, memorable links that redirect anywhere on the web,
            with powerful analytics to track performance.
          </p>
          <div className="flex flex-col sm:flex-row items-center gap-4 mt-6">
            <Link href="/register">
              <Button size="lg" className="gap-2">
                Get Started <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
            <Link href="/login">
              <Button variant="outline" size="lg">
                Login
              </Button>
            </Link>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 py-12">
          <div className="flex flex-col items-center text-center p-6 border rounded-lg">
            <Clock className="h-10 w-10 text-primary mb-4" />
            <h3 className="text-xl font-medium mb-2">Save Time</h3>
            <p className="text-muted-foreground">
              Transform long, complex URLs into short, easy-to-share links with
              just a single click.
            </p>
          </div>
          <div className="flex flex-col items-center text-center p-6 border rounded-lg">
            <BarChart3 className="h-10 w-10 text-primary mb-4" />
            <h3 className="text-xl font-medium mb-2">Track Performance</h3>
            <p className="text-muted-foreground">
              Monitor link clicks and engagement with detailed analytics to
              optimize your content.
            </p>
          </div>
          <div className="flex flex-col items-center text-center p-6 border rounded-lg">
            <Shield className="h-10 w-10 text-primary mb-4" />
            <h3 className="text-xl font-medium mb-2">Stay Secure</h3>
            <p className="text-muted-foreground">
              All links are protected and can be updated or deleted anytime for
              complete control.
            </p>
          </div>
        </div>
      </div>
    </>
  );
}
