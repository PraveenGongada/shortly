# Shortly

<div align="center">
  <img src="https://raw.githubusercontent.com/PraveenGongada/shortly/refs/heads/main/docs/images/logo.svg" alt="Shortly Logo" width="200" />
  <p></p>

[![Next.js](https://img.shields.io/badge/Next.js-14.2.0-black?style=flat-square&logo=next.js)](https://nextjs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-^5-blue?style=flat-square&logo=typescript)](https://www.typescriptlang.org)
[![Tailwind](https://img.shields.io/badge/Tailwind-^3.3.0-38B2AC?style=flat-square&logo=tailwind-css)](https://tailwindcss.com)
[![License](https://img.shields.io/github/license/PraveenGongada/shortly?style=flat-square)](LICENSE)

  <p></p>
  <p>A modern, fast, and easy-to-use URL shortening website built with Next.js</p>
</div>

## ✨ Features

- **Instant URL Shortening**: Transform long URLs into compact, shareable links with a single click
- **User Dashboard**: Manage all your shortened URLs in one convenient location
- **Click Analytics**: Track the performance of your links with detailed redirect statistics
- **Secure Authentication**: Complete user authentication system with registration and login
- **Responsive Design**: Beautiful UI that works seamlessly across desktop and mobile devices
- **Modern Tech Stack**: Built with Next.js, TypeScript, and Tailwind CSS for the frontend

## 🚀 Getting Started

### Prerequisites

- Node.js 18.x or later
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

Open [http://localhost:3001](http://localhost:3001) in your browser.

## 🏗️ Project Structure

```
/
├── public/              # Static assets
└── src/
    ├── app/             # Next.js app directory
    │   ├── dashboard/   # Dashboard page
    │   ├── login/       # Login page
    │   ├── register/    # Registration page
    │   └── [shortUrl]/  # URL redirection page
    ├── components/      # Reusable UI components
    │   ├── ui/          # UI components (buttons, cards, etc.)
    │   └── Navbar.tsx   # Navigation component
    ├── contexts/        # React contexts
    │   └── AuthContext.tsx  # Authentication context
    ├── lib/             # Utility functions
    │   ├── api.ts       # API client
    │   └── utils.ts     # Helper functions
    └── middleware/      # Next.js middleware
```

## 🔧 Configuration

The application can be configured using environment variables:

```env
# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080/api  # URL of the backend API
```

## 📚 Tech Stack

- [Next.js](https://nextjs.org/) - React framework
- [TypeScript](https://www.typescriptlang.org/) - Type safety
- [Tailwind CSS](https://tailwindcss.com/) - Styling
- [Radix UI](https://www.radix-ui.com/) - Accessible UI components
- [Lucide React](https://lucide.dev/) - Icon library

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check [issues page](https://github.com/praveengongada/shortly/issues).

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [Shadcn UI](https://ui.shadcn.com/) - UI components
- [Vercel](https://vercel.com/) - Deployment platform
- [Next.js](https://nextjs.org/) - The React Framework

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/PraveenGongada">Praveen Kumar</a></p>
  <p>
    <a href="https://linkedin.com/in/praveengongada">LinkedIn</a> •
    <a href="https://praveengongada.com">Website</a>
  </p>
</div>
