/**
 * NEXO - Sistema de Gestão para Barbearias
 * Layout de Autenticação
 *
 * Layout para páginas públicas (login, forgot-password, etc.)
 * Design responsivo com painel lateral ilustrativo.
 */

import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Login | NEXO',
  description: 'Acesse sua conta NEXO - Sistema de Gestão para Barbearias',
};

interface AuthLayoutProps {
  children: React.ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="min-h-screen flex">
      {/* Painel Esquerdo - Ilustração (oculto em mobile) */}
      <div className="hidden lg:flex lg:w-1/2 xl:w-[55%] bg-linear-to-br from-primary via-primary/90 to-primary/80 relative overflow-hidden">
        {/* Padrão decorativo */}
        <div className="absolute inset-0 opacity-10">
          <div
            className="absolute inset-0"
            style={{
              backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.4'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
            }}
          />
        </div>

        {/* Conteúdo */}
        <div className="relative z-10 flex flex-col justify-between p-8 lg:p-12 xl:p-16 w-full">
          {/* Logo */}
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-xl bg-white/10 backdrop-blur-sm flex items-center justify-center">
              <span className="text-2xl font-bold text-white">N</span>
            </div>
            <span className="text-2xl font-bold text-white tracking-tight">
              NEXO
            </span>
          </div>

          {/* Mensagem central */}
          <div className="space-y-6">
            <h1 className="text-4xl lg:text-5xl xl:text-6xl font-bold text-white leading-tight">
              Gerencie sua
              <br />
              <span className="text-white/80">barbearia</span>
              <br />
              com inteligência
            </h1>
            <p className="text-lg text-white/70 max-w-md">
              Agendamentos, finanças, equipe e clientes em uma única plataforma
              projetada para o seu negócio crescer.
            </p>
          </div>

          {/* Footer */}
          <div className="flex items-center gap-8 text-sm text-white/50">
            <span>© 2025 NEXO</span>
            <a
              href="https://nexo.app/termos"
              className="hover:text-white/80 transition-colors"
            >
              Termos de Uso
            </a>
            <a
              href="https://nexo.app/privacidade"
              className="hover:text-white/80 transition-colors"
            >
              Privacidade
            </a>
          </div>
        </div>
      </div>

      {/* Painel Direito - Formulário */}
      <div className="flex-1 flex flex-col min-h-screen bg-background">
        {/* Header mobile */}
        <header className="lg:hidden p-4 flex items-center justify-center border-b">
          <div className="flex items-center gap-2">
            <div className="h-8 w-8 rounded-lg bg-primary flex items-center justify-center">
              <span className="text-lg font-bold text-primary-foreground">
                N
              </span>
            </div>
            <span className="text-xl font-bold tracking-tight">NEXO</span>
          </div>
        </header>

        {/* Conteúdo do formulário */}
        <main className="flex-1 flex items-center justify-center p-4 sm:p-6 lg:p-8">
          <div className="w-full max-w-[400px] space-y-8">{children}</div>
        </main>

        {/* Footer mobile */}
        <footer className="lg:hidden p-4 text-center text-sm text-muted-foreground border-t">
          © 2025 NEXO. Todos os direitos reservados.
        </footer>
      </div>
    </div>
  );
}
