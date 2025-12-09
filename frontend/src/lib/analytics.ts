/**
 * NEXO - Sistema de Gestão para Barbearias
 * Telemetria / Analytics
 *
 * Sistema de rastreamento de eventos para análise de uso.
 * Suporta múltiplos providers (console, Sentry, PostHog, etc.)
 */

'use client';

// =============================================================================
// TIPOS
// =============================================================================

export type AnalyticsEventCategory =
  | 'auth'
  | 'navigation'
  | 'multi_unit'
  | 'appointment'
  | 'financial'
  | 'stock'
  | 'customer'
  | 'settings'
  | 'error';

export interface AnalyticsEvent {
  category: AnalyticsEventCategory;
  action: string;
  label?: string;
  value?: number;
  metadata?: Record<string, unknown>;
}

export interface AnalyticsUser {
  id: string;
  tenant_id: string;
  role?: string;
  plan?: string;
}

export interface AnalyticsProvider {
  name: string;
  track: (event: AnalyticsEvent) => void;
  identify: (user: AnalyticsUser) => void;
  reset: () => void;
}

// =============================================================================
// PROVIDERS
// =============================================================================

/**
 * Provider de console (para desenvolvimento)
 */
const consoleProvider: AnalyticsProvider = {
  name: 'console',
  track: (event) => {
    if (process.env.NODE_ENV === 'development') {
      console.log('📊 [Analytics]', event.category, event.action, event);
    }
  },
  identify: (user) => {
    if (process.env.NODE_ENV === 'development') {
      console.log('📊 [Analytics] Identify:', user);
    }
  },
  reset: () => {
    if (process.env.NODE_ENV === 'development') {
      console.log('📊 [Analytics] Reset');
    }
  },
};

/**
 * Provider do Sentry (para produção - se configurado)
 */
const sentryProvider: AnalyticsProvider = {
  name: 'sentry',
  track: (event) => {
    // Sentry breadcrumb
    if (typeof window !== 'undefined') {
      const win = window as unknown as Record<string, unknown>;
      if (win.Sentry) {
        const Sentry = win.Sentry as {
          addBreadcrumb: (breadcrumb: Record<string, unknown>) => void;
        };
        Sentry.addBreadcrumb({
          category: event.category,
          message: `${event.action}${event.label ? `: ${event.label}` : ''}`,
          level: 'info',
          data: event.metadata,
        });
      }
    }
  },
  identify: (user) => {
    if (typeof window !== 'undefined') {
      const win = window as unknown as Record<string, unknown>;
      if (win.Sentry) {
        const Sentry = win.Sentry as {
          setUser: (user: Record<string, unknown>) => void;
        };
        Sentry.setUser({
          id: user.id,
          tenant_id: user.tenant_id,
          role: user.role,
        });
      }
    }
  },
  reset: () => {
    if (typeof window !== 'undefined') {
      const win = window as unknown as Record<string, unknown>;
      if (win.Sentry) {
        const Sentry = win.Sentry as {
          setUser: (user: null) => void;
        };
        Sentry.setUser(null);
      }
    }
  },
};

// =============================================================================
// ANALYTICS SERVICE
// =============================================================================

class AnalyticsService {
  private providers: AnalyticsProvider[] = [];
  private user: AnalyticsUser | null = null;
  private enabled = true;

  constructor() {
    // Adiciona providers padrão
    this.providers.push(consoleProvider);

    // Adiciona Sentry se disponível
    if (process.env.NODE_ENV === 'production') {
      this.providers.push(sentryProvider);
    }
  }

  /**
   * Adiciona um provider de analytics
   */
  addProvider(provider: AnalyticsProvider): void {
    this.providers.push(provider);
  }

  /**
   * Habilita/desabilita tracking
   */
  setEnabled(enabled: boolean): void {
    this.enabled = enabled;
  }

  /**
   * Identifica o usuário atual
   */
  identify(user: AnalyticsUser): void {
    this.user = user;
    if (!this.enabled) return;

    this.providers.forEach((provider) => {
      try {
        provider.identify(user);
      } catch (error) {
        console.error(`[Analytics] Error in ${provider.name}.identify:`, error);
      }
    });
  }

  /**
   * Reseta o usuário (logout)
   */
  reset(): void {
    this.user = null;
    if (!this.enabled) return;

    this.providers.forEach((provider) => {
      try {
        provider.reset();
      } catch (error) {
        console.error(`[Analytics] Error in ${provider.name}.reset:`, error);
      }
    });
  }

  /**
   * Rastreia um evento
   */
  track(event: AnalyticsEvent): void {
    if (!this.enabled) return;

    // Adiciona contexto do usuário se disponível
    const enrichedEvent = {
      ...event,
      metadata: {
        ...event.metadata,
        user_id: this.user?.id,
        tenant_id: this.user?.tenant_id,
        timestamp: new Date().toISOString(),
      },
    };

    this.providers.forEach((provider) => {
      try {
        provider.track(enrichedEvent);
      } catch (error) {
        console.error(`[Analytics] Error in ${provider.name}.track:`, error);
      }
    });
  }

  // ===========================================================================
  // EVENTOS PRÉ-DEFINIDOS - MULTI-UNIT
  // ===========================================================================

  /**
   * Usuário trocou de unidade
   */
  trackUnitSwitch(fromUnitId: string, toUnitId: string, toUnitName: string): void {
    this.track({
      category: 'multi_unit',
      action: 'unit_switch',
      label: toUnitName,
      metadata: {
        from_unit_id: fromUnitId,
        to_unit_id: toUnitId,
        to_unit_name: toUnitName,
      },
    });
  }

  /**
   * Usuário definiu unidade padrão
   */
  trackSetDefaultUnit(unitId: string, unitName: string): void {
    this.track({
      category: 'multi_unit',
      action: 'set_default_unit',
      label: unitName,
      metadata: {
        unit_id: unitId,
        unit_name: unitName,
      },
    });
  }

  /**
   * Seletor de unidades foi aberto
   */
  trackUnitSelectorOpen(): void {
    this.track({
      category: 'multi_unit',
      action: 'unit_selector_open',
    });
  }

  // ===========================================================================
  // EVENTOS PRÉ-DEFINIDOS - AUTH
  // ===========================================================================

  trackLogin(method: 'email' | 'social' = 'email'): void {
    this.track({
      category: 'auth',
      action: 'login',
      label: method,
    });
  }

  trackLogout(): void {
    this.track({
      category: 'auth',
      action: 'logout',
    });
  }

  // ===========================================================================
  // EVENTOS PRÉ-DEFINIDOS - NAVEGAÇÃO
  // ===========================================================================

  trackPageView(pageName: string, path: string): void {
    this.track({
      category: 'navigation',
      action: 'page_view',
      label: pageName,
      metadata: { path },
    });
  }

  // ===========================================================================
  // EVENTOS PRÉ-DEFINIDOS - ERROS
  // ===========================================================================

  trackError(errorMessage: string, context?: Record<string, unknown>): void {
    this.track({
      category: 'error',
      action: 'error',
      label: errorMessage,
      metadata: context,
    });
  }
}

// Singleton
export const analytics = new AnalyticsService();

export default analytics;
