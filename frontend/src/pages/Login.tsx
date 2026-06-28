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

import { Link, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";

import { AuthCard } from "@/components/AuthCard";
import { Field } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/auth/AuthContext";
import { loginSchema, type LoginInput } from "@/lib/validators";
import { apiErrorMessage } from "@/lib/api";

export function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    formState: { errors, touchedFields, isValid, isSubmitting },
  } = useForm<LoginInput>({
    resolver: zodResolver(loginSchema),
    mode: "onChange",
    defaultValues: { email: "", password: "" },
  });

  const onSubmit = handleSubmit(async (values) => {
    try {
      await login(values);
      navigate("/dashboard", { replace: true });
    } catch (err) {
      toast.error(
        apiErrorMessage(err, "Could not sign in", {
          401: "Invalid email or password",
        }),
      );
    }
  });

  return (
    <AuthCard
      title="Welcome back"
      footer={
        <>
          No account?{" "}
          <Link to="/register" className="text-fg underline-offset-4 hover:underline">
            Create one
          </Link>
        </>
      }
    >
      <form onSubmit={onSubmit} className="space-y-4" noValidate>
        <Field
          label="Email"
          htmlFor="email"
          error={touchedFields.email ? errors.email?.message : undefined}
        >
          <Input
            id="email"
            type="email"
            autoComplete="email"
            placeholder="you@example.com"
            invalid={touchedFields.email && !!errors.email}
            {...register("email")}
          />
        </Field>
        <Field
          label="Password"
          htmlFor="password"
          error={touchedFields.password ? errors.password?.message : undefined}
        >
          <Input
            id="password"
            type="password"
            autoComplete="current-password"
            placeholder="••••••••"
            invalid={touchedFields.password && !!errors.password}
            {...register("password")}
          />
        </Field>
        <Button
          type="submit"
          loading={isSubmitting}
          disabled={!isValid}
          className="w-full justify-center"
        >
          Sign in
        </Button>
      </form>
    </AuthCard>
  );
}
