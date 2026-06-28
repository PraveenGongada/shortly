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
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { urlApi } from "@/lib/api";
import type { CreateUrlRequest, UrlUpdateRequest } from "@/types";

export const PAGE_SIZE = 9;

export function useLinks(offset: number) {
  return useQuery({
    queryKey: ["links", offset],
    queryFn: async () => {
      const rows = (await urlApi.list(PAGE_SIZE + 1, offset)) ?? [];
      return {
        items: rows.slice(0, PAGE_SIZE),
        hasNext: rows.length > PAGE_SIZE,
      };
    },
    placeholderData: keepPreviousData,
  });
}

export function useCreateLink() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateUrlRequest) => urlApi.create(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["links"] }),
  });
}

export function useUpdateLink() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: UrlUpdateRequest) => urlApi.update(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["links"] }),
  });
}

export function useDeleteLink() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => urlApi.remove(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["links"] }),
  });
}
