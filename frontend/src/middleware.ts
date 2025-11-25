/**
 * NEXO - Sistema de Gestão para Barbearias
 * Middleware Next.js
 *
 * Proteção de rotas no Edge - verifica cookie de autenticação.
 */

import type { NextRequest } from 'next/server';
import { NextResponse } from 'next/server';

// =============================================================================
// CONFIGURAÇÃO
// =============================================================================

const TOKEN_COOKIE_NAME = 'nexo-token';

/**
 * Rotas públicas que não requerem autenticação
 */
const PUBLIC_ROUTES = [
  '/login',
  '/forgot-password',
  '/reset-password',
  '/terms',
  '/privacy',
];

/**
 * Rotas de API e assets que devem ser ignoradas
 */
const IGNORED_ROUTES = [
  '/api',
  '/_next',
  '/favicon.ico',
  '/robots.txt',
  '/sitemap.xml',
];

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Verifica se a rota deve ser ignorada pelo middleware
 */
function shouldIgnoreRoute(pathname: string): boolean {
  return IGNORED_ROUTES.some((route) => pathname.startsWith(route));
}

/**
 * Verifica se a rota é pública
 */
function isPublicRoute(pathname: string): boolean {
  return PUBLIC_ROUTES.some((route) => pathname.startsWith(route));
}

/**
 * Obtém o token do cookie
 */
function getTokenFromCookie(request: NextRequest): string | null {
  return request.cookies.get(TOKEN_COOKIE_NAME)?.value ?? null;
}

// =============================================================================
// MIDDLEWARE
// =============================================================================

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Ignora rotas de API e assets
  if (shouldIgnoreRoute(pathname)) {
    return NextResponse.next();
  }

  const token = getTokenFromCookie(request);
  const isAuthenticated = !!token;
  const isPublic = isPublicRoute(pathname);

  // Rota raiz - redireciona baseado na autenticação
  if (pathname === '/') {
    if (isAuthenticated) {
      return NextResponse.redirect(new URL('/dashboard', request.url));
    }
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // Usuário autenticado tentando acessar rota pública (login, etc.)
  // Redireciona para dashboard
  if (isAuthenticated && isPublic) {
    return NextResponse.redirect(new URL('/dashboard', request.url));
  }

  // Usuário não autenticado tentando acessar rota protegida
  // Redireciona para login com return URL
  if (!isAuthenticated && !isPublic) {
    const loginUrl = new URL('/login', request.url);

    // Salva a URL original para redirecionar após login
    if (pathname !== '/dashboard') {
      loginUrl.searchParams.set('returnUrl', pathname);
    }

    return NextResponse.redirect(loginUrl);
  }

  // Permite acesso
  return NextResponse.next();
}

// =============================================================================
// CONFIG
// =============================================================================

export const config = {
  matcher: [
    /*
     * Match all request paths except:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     */
    '/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
  ],
};
