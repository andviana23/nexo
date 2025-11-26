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

  // Debug: Log detalhado para diagnosticar cookie
  const cookies = request.cookies.getAll();
  const logData = {
    pathname,
    hasToken: !!token,
    tokenLength: token?.length || 0,
    tokenPreview: token ? token.substring(0, 20) + '...' : '❌ NO TOKEN',
    isAuthenticated,
    isPublic,
    cookieCount: cookies.length,
    allCookies: cookies.map(c => ({
      name: c.name,
      valueLength: c.value.length,
      valuePreview: c.value.substring(0, 20) + '...'
    })),
  };
  
  // Log no console do navegador
  console.log('🔐 [Middleware]', logData);
  
  // Log no terminal do Next.js (stderr para garantir visibilidade)
  console.error('🔐 [MW]', JSON.stringify(logData, null, 2));

  // Rota raiz - se autenticado, continua; se não, vai para login
  if (pathname === '/') {
    if (isAuthenticated) {
      return NextResponse.next();
    }
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // Usuário autenticado tentando acessar rota pública (login, etc.)
  // Redireciona para a página inicial
  if (isAuthenticated && isPublic) {
    return NextResponse.redirect(new URL('/', request.url));
  }

  // Usuário não autenticado tentando acessar rota protegida
  // Redireciona para login com return URL
  if (!isAuthenticated && !isPublic) {
    const loginUrl = new URL('/login', request.url);

    // Salva a URL original para redirecionar após login
    if (pathname !== '/') {
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
