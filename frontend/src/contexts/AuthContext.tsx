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

import { createContext, useContext, useState, ReactNode } from "react";
import { authApi } from "@/lib/api";
import { useRouter } from "next/navigation";
import { UserLoginRequest, UserRegisterRequest } from "@/types";
import { useToast } from "@/components/ui/use-toast";

interface AuthContextType {
  isLoading: boolean;
  login: (data: UserLoginRequest) => Promise<boolean>;
  register: (data: UserRegisterRequest) => Promise<boolean>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const router = useRouter();
  const { toast } = useToast();

  const login = async (data: UserLoginRequest): Promise<boolean> => {
    try {
      setIsLoading(true);

      await authApi.login(data);

      toast({
        title: "Success",
        description: "You have successfully logged in",
      });

      router.push("/dashboard");
      return true;
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to login";
      toast({
        title: "Login Failed",
        description: message,
        variant: "destructive",
      });
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (data: UserRegisterRequest): Promise<boolean> => {
    try {
      setIsLoading(true);

      await authApi.register(data);

      toast({
        title: "Welcome!",
        description: "Your account has been created successfully!",
      });

      router.push("/dashboard");
      return true;
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to register";
      toast({
        title: "Registration Failed",
        description: message,
        variant: "destructive",
      });
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async (): Promise<void> => {
    try {
      setIsLoading(true);
      await authApi.logout();

      toast({
        title: "Logged Out",
        description: "You have been logged out successfully",
      });

      router.push("/login");
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to logout";
      toast({
        title: "Logout Failed",
        description: message,
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        isLoading,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
