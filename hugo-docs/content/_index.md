---
title: "FITS Backend Documentation"
description: "Comprehensive documentation for the FITS (Flexible IT Training System) Backend API"
date: 2025-10-23
---

# FITS Backend Documentation

Welcome to the FITS (Flexible IT Training System) Backend API documentation. This documentation complements the [Swagger/OpenAPI specification](/docs) with architectural insights, development guides, and deployment instructions.

## What is FITS?

FITS is a flexible IT training management system that provides a RESTful API for managing students, teachers, and signing requests with comprehensive authentication and authorization.

## Key Features

- **Authentication & Authorization**: JWT-based authentication with role-based access control (RBAC)
- **User Management**: Separate domains for students and teachers with invitation-based registration
- **Signing Requests**: Digital signing workflow management
- **Rate Limiting**: Multi-tier rate limiting (global IP-based + per-user role-based)
- **Session Management**: Secure session handling with token refresh
- **API Documentation**: Auto-generated Swagger/OpenAPI 3.0 documentation
- **Security**: Comprehensive security headers, password validation, and secure defaults

## Technology Stack

- **Language**: Go 1.25.1
- **Web Framework**: Fiber v2 (Express-inspired, built on Fasthttp)
- **Database**: PostgreSQL with GORM ORM (SQLite for testing)
- **Authentication**: JWT tokens with golang-jwt/jwt
- **API Documentation**: Swag (Swagger/OpenAPI 3.0)
- **Logging**: Uber Zap (structured logging)
- **Monitoring**: Prometheus metrics endpoint

## Quick Links

- [Architecture Overview](/project-overview/architecture/)
- [Installation Guide](/getting-started/installation/)
- [Quick Start](/getting-started/quick-start/)
- [API Reference](/api/endpoints/)
- [Swagger UI](http://localhost:8080/docs)

## Documentation Structure

### Project Overview
Learn about the system architecture, design principles, and technical decisions.

### Getting Started
Step-by-step guides to install, configure, and run the FITS backend.

### Development
Development environment setup, testing strategies, and contribution guidelines.

### API Reference
Detailed API documentation complementing the Swagger specification.

### Infrastructure
Deployment guides, security considerations, and operational best practices.

## Version

Current Version: **1.0.1**

## Support

For issues, questions, or contributions, please visit our [GitHub repository](https://github.com/JustDoItBetter/FITS-backend).
