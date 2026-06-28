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

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { useNavigate } from "react-router-dom";

import { authApi } from "@/lib/api";
import { queryClient, setUnauthorizedHandler } from "@/lib/queryClient";
import type { LoginRequest, RegisterRequest } from "@/types";

const AUTH_FLAG = "shortly_authed";

interface AuthContextValue {
  isAuthed: boolean;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const navigate = useNavigate();
  const [isAuthed, setIsAuthed] = useState(
    () => localStorage.getItem(AUTH_FLAG) === "1",
  );

  const markAuthed = useCallback(() => {
    localStorage.setItem(AUTH_FLAG, "1");
    setIsAuthed(true);
  }, []);

  const clearSession = useCallback(() => {
    localStorage.removeItem(AUTH_FLAG);
    setIsAuthed(false);
    queryClient.clear();
  }, []);

  const login = useCallback(
    async (data: LoginRequest) => {
      await authApi.login(data);
      markAuthed();
    },
    [markAuthed],
  );

  const register = useCallback(
    async (data: RegisterRequest) => {
      await authApi.register(data);
      markAuthed();
    },
    [markAuthed],
  );

  const logout = useCallback(async () => {
    try {
      await authApi.logout();
    } catch {}
    clearSession();
    navigate("/login", { replace: true });
  }, [clearSession, navigate]);

  useEffect(() => {
    setUnauthorizedHandler(() => {
      clearSession();
      navigate("/login", { replace: true });
    });
    return () => setUnauthorizedHandler(null);
  }, [clearSession, navigate]);

  const value = useMemo(
    () => ({ isAuthed, login, register, logout }),
    [isAuthed, login, register, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within an AuthProvider");
  return ctx;
}
