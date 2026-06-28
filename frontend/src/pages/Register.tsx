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
import { Check } from "lucide-react";
import { toast } from "sonner";

import { AuthCard } from "@/components/AuthCard";
import { Field } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/auth/AuthContext";
import { registerSchema, passwordRules, type RegisterInput } from "@/lib/validators";
import { apiErrorMessage } from "@/lib/api";
import { cn } from "@/lib/utils";

export function Register() {
  const { register: registerUser } = useAuth();
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, touchedFields, isValid, isSubmitting },
  } = useForm<RegisterInput>({
    resolver: zodResolver(registerSchema),
    mode: "onChange",
    defaultValues: { name: "", email: "", password: "" },
  });

  const password = watch("password");

  const onSubmit = handleSubmit(async (values) => {
    try {
      await registerUser({
        name: values.name,
        email: values.email,
        password: values.password,
      });
      navigate("/dashboard", { replace: true });
    } catch (err) {
      toast.error(
        apiErrorMessage(err, "Could not create account", {
          409: "That email is already registered",
        }),
      );
    }
  });

  return (
    <AuthCard
      title="Create your account"
      footer={
        <>
          Already have an account?{" "}
          <Link to="/login" className="text-fg underline-offset-4 hover:underline">
            Sign in
          </Link>
        </>
      }
    >
      <form onSubmit={onSubmit} className="space-y-4" noValidate>
        <Field
          label="Name"
          htmlFor="name"
          error={touchedFields.name ? errors.name?.message : undefined}
        >
          <Input
            id="name"
            autoComplete="name"
            placeholder="Jane Doe"
            invalid={touchedFields.name && !!errors.name}
            {...register("name")}
          />
        </Field>
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
            autoComplete="new-password"
            placeholder="••••••••"
            invalid={touchedFields.password && !!errors.password}
            {...register("password")}
          />
        </Field>

        <ul className="grid grid-cols-2 gap-x-3 gap-y-1.5 pt-1">
          {passwordRules.map((rule) => {
            const ok = rule.test(password || "");
            return (
              <li
                key={rule.label}
                className={cn(
                  "flex items-center gap-1.5 text-xs transition-colors",
                  ok ? "text-green-400" : "text-faint",
                )}
              >
                <Check className={cn("h-3.5 w-3.5", !ok && "opacity-40")} />
                {rule.label}
              </li>
            );
          })}
        </ul>

        <Button
          type="submit"
          loading={isSubmitting}
          disabled={!isValid}
          className="w-full justify-center"
        >
          Create account
        </Button>
      </form>
    </AuthCard>
  );
}
