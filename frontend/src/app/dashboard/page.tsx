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

import { useEffect, useState } from "react";
import Navbar from "@/components/Navbar";
import { urlApi } from "@/lib/api";
import { UrlResponse, CreateUrlRequest } from "@/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { useToast } from "@/components/ui/use-toast";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Copy, ExternalLink, Pencil, Trash2, Plus } from "lucide-react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

export default function DashboardPage() {
  const { toast } = useToast();
  const [urls, setUrls] = useState<UrlResponse[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [limit] = useState(10);
  const [offset, setOffset] = useState(0);
  const [newUrl, setNewUrl] = useState("");
  const [editingUrl, setEditingUrl] = useState<UrlResponse | null>(null);
  const [editUrl, setEditUrl] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isDeletingId, setIsDeletingId] = useState<string | null>(null);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [urlToDelete, setUrlToDelete] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);

  useEffect(() => {
    fetchUrls();
  }, [limit, offset]);

  const fetchUrls = async () => {
    try {
      setIsLoading(true);
      const response = await urlApi.getUrls({ limit, offset });
      setUrls(response.data);

      // Determine if there are more items
      // If the returned data length equals the requested limit, assume there are more
      setHasMore(response.data.length === limit);
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to fetch URLs",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateUrl = async () => {
    try {
      setIsSubmitting(true);

      if (!newUrl.trim()) {
        toast({
          title: "Error",
          description: "Please enter a valid URL",
          variant: "destructive",
        });
        return;
      }

      const data: CreateUrlRequest = {
        long_url: newUrl,
      };

      await urlApi.createShortUrl(data);

      setNewUrl("");
      setIsCreateDialogOpen(false);
      toast({
        title: "Success",
        description: "URL shortened successfully",
      });

      fetchUrls();
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to create short URL",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEditUrl = async () => {
    try {
      setIsSubmitting(true);

      if (!editUrl.trim() || !editingUrl) {
        toast({
          title: "Error",
          description: "Please enter a valid URL",
          variant: "destructive",
        });
        return;
      }

      await urlApi.updateUrl({
        id: editingUrl.id,
        new_url: editUrl,
      });

      setEditingUrl(null);
      setEditUrl("");
      setIsEditDialogOpen(false);

      toast({
        title: "Success",
        description: "URL updated successfully",
      });

      fetchUrls();
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to update URL",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteUrl = async (id: string) => {
    try {
      setIsDeletingId(id);

      await urlApi.deleteUrl(id);

      toast({
        title: "Success",
        description: "URL deleted successfully",
      });

      fetchUrls();
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to delete URL",
        variant: "destructive",
      });
    } finally {
      setIsDeletingId(null);
      setIsDeleteDialogOpen(false);
    }
  };

  const confirmDelete = (id: string) => {
    setUrlToDelete(id);
    setIsDeleteDialogOpen(true);
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toast({
      title: "Copied!",
      description: "URL copied to clipboard",
    });
  };

  const handlePageChange = (page: number) => {
    // Only allow moving forward if we know there are more items
    if (page > currentPage && !hasMore) return;

    // Don't allow going below page 1
    if (page < 1) return;

    setCurrentPage(page);
    setOffset((page - 1) * limit);
  };

  const renderPagination = () => {
    const pages = [];

    // Always show current page
    pages.push(
      <PaginationItem key={currentPage}>
        <PaginationLink isActive={true}>{currentPage}</PaginationLink>
      </PaginationItem>,
    );

    // If we know there are more pages, show the next page number
    if (hasMore) {
      pages.push(
        <PaginationItem key={currentPage + 1}>
          <PaginationLink
            isActive={false}
            onClick={() => handlePageChange(currentPage + 1)}
          >
            {currentPage + 1}
          </PaginationLink>
        </PaginationItem>,
      );

      // Show an ellipsis to indicate there might be more
      pages.push(
        <PaginationItem key="ellipsis">
          <PaginationEllipsis />
        </PaginationItem>,
      );
    }

    return pages;
  };

  // If we're loading the initial page data, show a full-page loading spinner
  if (isLoading) {
    return (
      <>
        <Navbar isAuthenticated={true} />
        <div className="flex justify-center items-center min-h-[calc(100vh-4rem)]">
          <LoadingSpinner size={40} />
        </div>
      </>
    );
  }

  return (
    <>
      <Navbar isAuthenticated={true} />
      <div className="container py-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-semibold">Your Short URLs</h1>
          <Dialog
            open={isCreateDialogOpen}
            onOpenChange={setIsCreateDialogOpen}
          >
            <DialogTrigger asChild>
              <Button className="flex items-center gap-2">
                <Plus className="h-4 w-4" />
                Create New URL
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Create Short URL</DialogTitle>
                <DialogDescription>
                  Enter a long URL to generate a short one.
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4 py-4">
                <div className="space-y-2">
                  <Label htmlFor="url">Long URL</Label>
                  <Input
                    id="url"
                    placeholder="https://example.com/very/long/url"
                    value={newUrl}
                    onChange={(e) => setNewUrl(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === "Enter" && newUrl.trim()) {
                        handleCreateUrl();
                      }
                    }}
                  />
                </div>
              </div>
              <DialogFooter>
                <Button
                  onClick={handleCreateUrl}
                  disabled={isSubmitting || !newUrl.trim()}
                >
                  {isSubmitting ? (
                    <LoadingSpinner className="mr-2" size={16} />
                  ) : null}
                  {isSubmitting ? "Creating..." : "Create Short URL"}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>

        {isLoading ? (
          <div className="flex justify-center items-center py-12">
            <LoadingSpinner size={32} />
          </div>
        ) : urls.length === 0 ? (
          <Card>
            <CardContent className="flex flex-col items-center justify-center py-12">
              <p className="text-muted-foreground mb-4">
                You haven't created any short URLs yet.
              </p>
              <Button
                onClick={() => setIsCreateDialogOpen(true)}
                className="flex items-center gap-2"
              >
                <Plus className="h-4 w-4" />
                Create Your First URL
              </Button>
            </CardContent>
          </Card>
        ) : (
          <>
            <div className="rounded-md border">
              <div className="grid grid-cols-1 md:grid-cols-5 p-4 font-medium bg-muted">
                <div className="md:col-span-2">Long URL</div>
                <div>Short URL</div>
                <div>Redirects</div>
                <div>Actions</div>
              </div>
              <div className="divide-y">
                {urls.map((url) => (
                  <div
                    key={url.id}
                    className="grid grid-cols-1 md:grid-cols-5 p-4 items-center"
                  >
                    <div className="md:col-span-2 truncate">
                      <a
                        href={url.longUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="hover:underline flex items-center gap-1 text-primary"
                      >
                        {url.longUrl}
                        <ExternalLink className="h-3 w-3" />
                      </a>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="truncate">{url.shortUrl}</span>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => copyToClipboard(url.shortUrl)}
                      >
                        <Copy className="h-4 w-4" />
                      </Button>
                    </div>
                    <div>{url.redirects}</div>
                    <div className="flex items-center gap-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => {
                          setEditingUrl(url);
                          setEditUrl(url.longUrl);
                          setIsEditDialogOpen(true);
                        }}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => confirmDelete(url.id)}
                        disabled={isDeletingId === url.id}
                      >
                        {isDeletingId === url.id ? (
                          <LoadingSpinner size={16} />
                        ) : (
                          <Trash2 className="h-4 w-4" />
                        )}
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {(urls.length > 0 || currentPage > 1) && (
              <div className="mt-6">
                <Pagination>
                  <PaginationContent>
                    <PaginationItem>
                      <PaginationPrevious
                        onClick={() => handlePageChange(currentPage - 1)}
                        className={
                          currentPage === 1
                            ? "pointer-events-none opacity-50"
                            : ""
                        }
                      />
                    </PaginationItem>

                    {renderPagination()}

                    <PaginationItem>
                      <PaginationNext
                        onClick={() => handlePageChange(currentPage + 1)}
                        className={
                          !hasMore ? "pointer-events-none opacity-50" : ""
                        }
                      />
                    </PaginationItem>
                  </PaginationContent>
                </Pagination>
              </div>
            )}
          </>
        )}

        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Edit URL</DialogTitle>
              <DialogDescription>
                Update the destination URL for {editingUrl?.shortUrl}
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="edit-url">New URL</Label>
                <Input
                  id="edit-url"
                  placeholder="https://example.com/new/destination"
                  value={editUrl}
                  onChange={(e) => setEditUrl(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" && editUrl.trim()) {
                      handleEditUrl();
                    }
                  }}
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                onClick={handleEditUrl}
                disabled={isSubmitting || !editUrl.trim()}
              >
                {isSubmitting ? (
                  <LoadingSpinner className="mr-2" size={16} />
                ) : null}
                {isSubmitting ? "Updating..." : "Update URL"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
        <AlertDialog
          open={isDeleteDialogOpen}
          onOpenChange={setIsDeleteDialogOpen}
        >
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This will permanently delete this URL. This action cannot be
                undone.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                onClick={() => urlToDelete && handleDeleteUrl(urlToDelete)}
                className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              >
                {isDeletingId === urlToDelete ? (
                  <LoadingSpinner className="mr-2" size={16} />
                ) : null}
                Delete
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </>
  );
}
