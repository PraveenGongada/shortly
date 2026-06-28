# Shortly

<div align="center">
  <img src="https://raw.githubusercontent.com/PraveenGongada/shortly/refs/heads/main/docs/images/logo.svg" alt="Shortly Logo" width="200" />
  <p></p>

[![Vite](https://img.shields.io/badge/Vite-^6.4.3-646CFF?style=flat-square&logo=vite)](https://vitejs.dev)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat-square&logo=react)](https://react.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-^5-blue?style=flat-square&logo=typescript)](https://www.typescriptlang.org)
[![Tailwind](https://img.shields.io/badge/Tailwind-^3.4.0-38B2AC?style=flat-square&logo=tailwind-css)](https://tailwindcss.com)
[![License](https://img.shields.io/github/license/PraveenGongada/shortly?style=flat-square)](LICENSE)

  <p></p>
  <p>A modern, fast, and easy-to-use URL shortening website built with React and Vite</p>
</div>

## ✨ Features

- **Instant URL Shortening**: Transform long URLs into compact, shareable links with a single click
- **User Dashboard**: Manage all your shortened URLs in one convenient location
- **Click Analytics**: Track the performance of your links with detailed redirect statistics
- **QR Codes**: Generate a scannable QR code for any short link, client-side
- **Secure Authentication**: Complete user authentication system with registration and login
- **Responsive Design**: Beautiful dark UI that works seamlessly across desktop and mobile devices
- **Modern Tech Stack**: Built as a static React SPA with Vite, TypeScript, and Tailwind CSS, served by nginx

## 🚀 Getting Started

### Prerequisites

- Node.js 20.x or later
- yarn

### Installation and Setup

1. **Clone the repository**

```bash
git clone https://github.com/praveengongada/shortly.git
cd shortly/frontend
```

2. **Install frontend dependencies**

```bash
yarn install
```

3. **Start the development server**

```bash
yarn dev
```

4. **Access the application**

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Build

```bash
yarn build      # type-checks, then emits the static build to dist/
yarn preview    # serve the production build locally
```

## 🏗️ Project Structure

```
/
├── public/                  # Static assets (favicons)
├── nginx.conf               # Production web server config
├── Dockerfile               # Build (Vite) → serve (nginx) image
└── src/
    ├── main.tsx             # App entry (router + query/auth providers)
    ├── App.tsx              # Route definitions
    ├── pages/               # Route pages
    │   ├── Landing.tsx      # Public landing page
    │   ├── Login.tsx        # Login
    │   ├── Register.tsx     # Registration
    │   ├── Dashboard.tsx    # Authenticated link management
    │   ├── RedirectShort.tsx  # /{code} short-link resolver
    │   └── NotFound.tsx     # 404
    ├── components/          # Reusable UI components
    │   ├── ui/              # Primitives (button, dialog, input, ...)
    │   └── ...              # CreateBar, LinkRow, Navbar, dialogs, ...
    ├── auth/                # Auth context + route guards
    │   ├── AuthContext.tsx
    │   └── ProtectedRoute.tsx
    ├── hooks/               # Data + UI hooks (useLinks, useCopy)
    ├── lib/                 # api client, query client, utils, validators
    └── types/               # Shared TypeScript types
```

## 📚 Tech Stack

- [Vite](https://vitejs.dev/) - Build tool & dev server
- [React](https://react.dev/) - UI library
- [TypeScript](https://www.typescriptlang.org/) - Type safety
- [Tailwind CSS](https://tailwindcss.com/) - Styling
- [TanStack Query](https://tanstack.com/query) - Data fetching & caching
- [React Hook Form](https://react-hook-form.com/) + [Zod](https://zod.dev/) - Forms & validation
- [Radix UI](https://www.radix-ui.com/) - Accessible UI components
- [Lucide React](https://lucide.dev/) - Icon library
- [nginx](https://nginx.org/) - Production static serving

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check [issues page](https://github.com/praveengongada/shortly/issues).

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [Shadcn UI](https://ui.shadcn.com/) - UI component inspiration
- [Radix UI](https://www.radix-ui.com/) - Accessible primitives
- [Vite](https://vitejs.dev/) - Build tooling

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/PraveenGongada">Praveen Kumar</a></p>
  <p>
    <a href="https://linkedin.com/in/praveengongada">LinkedIn</a> •
    <a href="https://praveengongada.com">Website</a>
  </p>
</div>
