# 🎟️ TicketBooking API

**A modern, production-ready Event & Ticket Booking Backend** built with Go.

Multi-language events · Real-time messaging · Stripe payments · Cloudinary media · Role-based access · Blog & comments

---

## ✨ Features

| Module | Description |
|--------|-------------|
| **🔐 Authentication** | JWT-based auth, OTP email verification, role-based access (User / Manager / Admin) |
| **👤 Users** | Registration, profile management, status control |
| **🎭 Events** | Create & manage events with multi-language (i18n) title, description & location |
| **📂 Event Categories** | Organize events by category |
| **🎫 Bookings** | Ticket booking with unique booking codes & availability tracking |
| **💳 Payments** | Stripe checkout + webhook support, payment history for users/managers/admins |
| **💬 Messaging** | Real-time WebSocket chat + conversation history |
| **📝 Blog** | Full blog system with likes, comments & admin moderation |
| **🖼️ Media** | Image upload via Cloudinary (thumbnails + galleries) |
| **🌐 i18n** | Multi-language support (Language + Translation modules) |
| **💱 Currency** | Multi-currency support |
| **⚙️ Settings** | App & page settings management |
| **📧 Email** | OTP & notification emails (SMTP + Resend) |

---

## 🛠️ Tech Stack

- **Language:** Go
- **Framework:** [Echo v5](https://echo.labstack.com/)
- **ORM:** GORM
- **Database:** PostgreSQL (Neon recommended)
- **Auth:** JWT
- **Payments:** Stripe
- **Media:** Cloudinary
- **Email:** SMTP / Resend
- **Realtime:** WebSocket (custom Hub)
- **Validation:** go-playground/validator

---

## 📁 Project Structure
