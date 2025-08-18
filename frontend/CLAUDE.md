# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

- **Development server**: `yarn dev` (runs on port 3001)
- **Build**: `yarn build`
- **Production start**: `yarn start`
- **Linting**: `yarn lint`
- **Install dependencies**: `yarn install`

## Application Architecture

This is a Next.js 14 frontend application for a URL shortening service called "Shortly". The application uses the App Router pattern with the following key architectural components:

### API Integration
- Backend API runs on `http://localhost:8080/api` by default
- API client is centralized in `src/lib/api.ts` with organized modules:
  - `authApi`: Authentication endpoints (login, register, logout)
  - `urlApi`: URL management endpoints (create, read, update, delete, analytics, redirect)
- Uses cookie-based authentication with automatic token management
- API proxy configured via Next.js rewrites for development

### Authentication System
- Cookie-based authentication using browser cookies
- `AuthContext` (`src/contexts/AuthContext.tsx`) provides global authentication state
- Authentication middleware (`src/middleware/auth.ts`) handles route protection:
  - Redirects authenticated users away from login/register pages
  - Protects `/dashboard` route for authenticated users only
  - Runs on routes: `/dashboard`, `/login`, `/register`, `/`

### UI Architecture
- Uses Radix UI components as the foundation with custom styling
- UI components located in `src/components/ui/` (buttons, cards, dialogs, forms, etc.)
- Tailwind CSS with custom design system (CSS variables for theming)
- Toast notifications system integrated throughout the app
- Responsive design with mobile-first approach

### Page Structure
- **Home** (`/`): Landing page for unauthenticated users
- **Login** (`/login`): User authentication
- **Register** (`/register`): User registration  
- **Dashboard** (`/dashboard`): Main authenticated user interface for URL management
- **URL Redirect** (`/[shortUrl]`): Dynamic route for handling short URL redirects

### State Management
- React Context for authentication state
- Local component state for UI interactions
- API responses managed through custom hooks and context

### TypeScript Integration
- Comprehensive type definitions in `src/types/index.ts`
- Covers all API request/response interfaces
- Pagination and generic response types included

## Configuration Details

### Environment Variables
- `NEXT_PUBLIC_API_URL`: Backend API URL (defaults to `/api` for proxy or `http://localhost:8080/api`)

### Key Dependencies
- **Next.js 14.2.0**: React framework with App Router
- **TypeScript**: Full type safety
- **Tailwind CSS**: Utility-first styling
- **Radix UI**: Accessible component primitives
- **React Hook Form**: Form handling
- **Lucide React**: Icon library

### Build Configuration
- Standalone output mode for containerization
- API proxy rewrites for development
- React strict mode enabled
- Custom port 3001 for both dev and production