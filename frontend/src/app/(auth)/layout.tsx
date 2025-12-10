/**
 * NEXO - Sistema de Gestão para Barbearias
 * Layout de Autenticação - Redesign Completo
 *
 * Layout para páginas públicas (login, forgot-password, etc.)
 * Design responsivo com painel lateral ilustrativo.
 */

import { Calendar, DollarSign, Scissors, Users } from 'lucide-react';
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
      {/* Painel Esquerdo - Hero Section (oculto em mobile) */}
      <div className="hidden lg:flex lg:w-1/2 xl:w-[55%] bg-slate-900 relative overflow-hidden">
        {/* Gradiente de fundo */}
        <div className="absolute inset-0 bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900" />
        
        {/* Padrão decorativo sutil */}
        <div className="absolute inset-0 opacity-[0.03]">
          <div
            className="absolute inset-0"
            style={{
              backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
            }}
          />
        </div>

        {/* Círculos decorativos com gradiente */}
        <div className="absolute -top-32 -left-32 w-[500px] h-[500px] bg-gradient-to-br from-blue-500/20 to-transparent rounded-full blur-3xl" />
        <div className="absolute -bottom-32 -right-32 w-[600px] h-[600px] bg-gradient-to-tl from-emerald-500/10 to-transparent rounded-full blur-3xl" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[400px] h-[400px] bg-gradient-to-r from-purple-500/5 to-blue-500/5 rounded-full blur-3xl" />

        {/* Conteúdo */}
        <div className="relative z-10 flex flex-col justify-between p-8 lg:p-12 xl:p-16 w-full">
          {/* Logo */}
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-xl bg-white flex items-center justify-center shadow-lg shadow-white/10">
              <Scissors className="h-5 w-5 text-slate-900" />
            </div>
            <span className="text-2xl font-bold text-white tracking-tight">
              NEXO
            </span>
          </div>

          {/* Hero Content */}
          <div className="space-y-8">
            <div className="space-y-4">
              <h1 className="text-4xl lg:text-5xl xl:text-6xl font-bold text-white leading-[1.1] tracking-tight">
                A gestão da sua
                <br />
                <span className="bg-gradient-to-r from-blue-400 to-emerald-400 bg-clip-text text-transparent">
                  barbearia
                </span>
                <br />
                simplificada.
              </h1>
              <p className="text-lg text-slate-400 max-w-md leading-relaxed">
                Agendamentos, finanças, equipe e clientes em uma única plataforma 
                inteligente projetada para o seu negócio crescer.
              </p>
            </div>

            {/* Features Grid */}
            <div className="grid grid-cols-2 gap-4 pt-4">
              <div className="flex items-center gap-3 p-4 rounded-xl bg-white/5 backdrop-blur-sm border border-white/10 hover:bg-white/10 transition-colors">
                <div className="h-10 w-10 rounded-lg bg-blue-500/20 flex items-center justify-center">
                  <Calendar className="h-5 w-5 text-blue-400" />
                </div>
                <div>
                  <p className="text-sm font-medium text-white">Agenda</p>
                  <p className="text-xs text-slate-500">Inteligente</p>
                </div>
              </div>
              
              <div className="flex items-center gap-3 p-4 rounded-xl bg-white/5 backdrop-blur-sm border border-white/10 hover:bg-white/10 transition-colors">
                <div className="h-10 w-10 rounded-lg bg-emerald-500/20 flex items-center justify-center">
                  <DollarSign className="h-5 w-5 text-emerald-400" />
                </div>
                <div>
                  <p className="text-sm font-medium text-white">Financeiro</p>
                  <p className="text-xs text-slate-500">Completo</p>
                </div>
              </div>
              
              <div className="flex items-center gap-3 p-4 rounded-xl bg-white/5 backdrop-blur-sm border border-white/10 hover:bg-white/10 transition-colors">
                <div className="h-10 w-10 rounded-lg bg-purple-500/20 flex items-center justify-center">
                  <Users className="h-5 w-5 text-purple-400" />
                </div>
                <div>
                  <p className="text-sm font-medium text-white">Equipe</p>
                  <p className="text-xs text-slate-500">Gestão total</p>
                </div>
              </div>
              
              <div className="flex items-center gap-3 p-4 rounded-xl bg-white/5 backdrop-blur-sm border border-white/10 hover:bg-white/10 transition-colors">
                <div className="h-10 w-10 rounded-lg bg-amber-500/20 flex items-center justify-center">
                  <Scissors className="h-5 w-5 text-amber-400" />
                </div>
                <div>
                  <p className="text-sm font-medium text-white">Serviços</p>
                  <p className="text-xs text-slate-500">Personalizados</p>
                </div>
              </div>
            </div>
          </div>

          {/* Footer */}
          <div className="flex items-center gap-6 text-sm text-slate-500">
            <span>© 2025 NEXO</span>
            <a
              href="https://nexo.app/termos"
              className="hover:text-slate-300 transition-colors"
            >
              Termos de Uso
            </a>
            <a
              href="https://nexo.app/privacidade"
              className="hover:text-slate-300 transition-colors"
            >
              Privacidade
            </a>
          </div>
        </div>
      </div>

      {/* Painel Direito - Formulário */}
      <div className="flex-1 flex flex-col min-h-screen bg-slate-50 dark:bg-slate-950">
        {/* Header mobile */}
        <header className="lg:hidden p-6 flex items-center justify-center">
          <div className="flex items-center gap-2">
            <div className="h-9 w-9 rounded-lg bg-slate-900 flex items-center justify-center">
              <Scissors className="h-4 w-4 text-white" />
            </div>
            <span className="text-xl font-bold tracking-tight">NEXO</span>
          </div>
        </header>

        {/* Conteúdo do formulário */}
        <main className="flex-1 flex items-center justify-center p-6 sm:p-8 lg:p-12">
          <div className="w-full max-w-[420px]">{children}</div>
        </main>

        {/* Footer mobile */}
        <footer className="lg:hidden p-6 text-center text-sm text-muted-foreground">
          © 2025 NEXO. Todos os direitos reservados.
        </footer>
      </div>
    </div>
  );
}
