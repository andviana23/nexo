/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Login - Redesign Completo
 *
 * Formulário de autenticação com validação Zod + React Hook Form.
 * Design System NEXO v1.0 - Clean & Professional
 */

'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowRight, Eye, EyeOff, Loader2, Lock, Mail } from 'lucide-react';
import Link from 'next/link';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useAuth } from '@/hooks/use-auth';
import { cn } from '@/lib/utils';

// =============================================================================
// VALIDAÇÃO
// =============================================================================

const loginSchema = z.object({
  email: z.string().min(1, 'Email é obrigatório').email('Email inválido'),
  password: z
    .string()
    .min(1, 'Senha é obrigatória')
    .min(6, 'Senha deve ter pelo menos 6 caracteres'),
});

type LoginFormData = z.infer<typeof loginSchema>;

// =============================================================================
// COMPONENTE
// =============================================================================

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  const { login, isLoggingIn, loginError } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const onSubmit = async (data: LoginFormData) => {
    try {
      await login(data);
    } catch {
      // Erro já é tratado pelo hook e exibido via loginError
    }
  };

  return (
    <div className="space-y-8 animate-in fade-in-0 slide-in-from-bottom-4 duration-500">
      {/* Header */}
      <div className="space-y-2 text-center lg:text-left">
        <h1 className="text-3xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
          Bem-vindo de volta
        </h1>
        <p className="text-slate-500 dark:text-slate-400">
          Entre com suas credenciais para acessar o sistema
        </p>
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Erro geral */}
        {loginError && (
          <div className="flex items-center gap-3 p-4 text-sm text-red-600 bg-red-50 dark:bg-red-950/50 dark:text-red-400 rounded-xl border border-red-200 dark:border-red-900 animate-in fade-in-0 slide-in-from-top-2 duration-300">
            <div className="h-2 w-2 rounded-full bg-red-500 animate-pulse flex-shrink-0" />
            <span>{loginError}</span>
          </div>
        )}

        {/* Email */}
        <div className="space-y-2">
          <Label htmlFor="email" className="text-sm font-medium text-slate-700 dark:text-slate-300">
            Email
          </Label>
          <div className="relative">
            <Mail className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-400 pointer-events-none" />
            <Input
              id="email"
              type="email"
              placeholder="seu@email.com"
              autoComplete="email"
              disabled={isLoggingIn}
              {...register('email')}
              className={cn(
                'h-12 pl-12 text-base bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl transition-all duration-200',
                'focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500',
                'placeholder:text-slate-400',
                errors.email && 'border-red-500 focus:ring-red-500/20 focus:border-red-500'
              )}
            />
          </div>
          {errors.email && (
            <p className="text-sm text-red-500 animate-in fade-in-0 slide-in-from-top-1 duration-200">
              {errors.email.message}
            </p>
          )}
        </div>

        {/* Senha */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <Label htmlFor="password" className="text-sm font-medium text-slate-700 dark:text-slate-300">
              Senha
            </Label>
            <Link
              href="/forgot-password"
              className="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium transition-colors"
              tabIndex={-1}
            >
              Esqueceu a senha?
            </Link>
          </div>
          <div className="relative">
            <Lock className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-400 pointer-events-none" />
            <Input
              id="password"
              type={showPassword ? 'text' : 'password'}
              placeholder="••••••••"
              autoComplete="current-password"
              disabled={isLoggingIn}
              {...register('password')}
              className={cn(
                'h-12 pl-12 pr-12 text-base bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl transition-all duration-200',
                'focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500',
                'placeholder:text-slate-400',
                errors.password && 'border-red-500 focus:ring-red-500/20 focus:border-red-500'
              )}
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 transition-colors p-1"
              tabIndex={-1}
              aria-label={showPassword ? 'Ocultar senha' : 'Mostrar senha'}
            >
              {showPassword ? (
                <EyeOff className="h-5 w-5" />
              ) : (
                <Eye className="h-5 w-5" />
              )}
            </button>
          </div>
          {errors.password && (
            <p className="text-sm text-red-500 animate-in fade-in-0 slide-in-from-top-1 duration-200">
              {errors.password.message}
            </p>
          )}
        </div>

        {/* Submit Button */}
        <Button
          type="submit"
          className="w-full h-12 text-base font-semibold rounded-xl bg-slate-900 hover:bg-slate-800 dark:bg-slate-100 dark:hover:bg-slate-200 dark:text-slate-900 transition-all duration-200 shadow-lg shadow-slate-900/10 hover:shadow-xl hover:shadow-slate-900/20"
          disabled={isLoggingIn}
        >
          {isLoggingIn ? (
            <>
              <Loader2 className="mr-2 h-5 w-5 animate-spin" />
              Entrando...
            </>
          ) : (
            <>
              Entrar
              <ArrowRight className="ml-2 h-5 w-5" />
            </>
          )}
        </Button>
      </form>

      {/* Divider */}
      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <span className="w-full border-t border-slate-200 dark:border-slate-800" />
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-slate-50 dark:bg-slate-950 px-4 text-slate-400">
            Precisa de ajuda?
          </span>
        </div>
      </div>

      {/* Support Link */}
      <p className="text-sm text-slate-500 dark:text-slate-400 text-center">
        Entre em contato com{' '}
        <a
          href="mailto:suporte@nexo.app"
          className="text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium transition-colors"
        >
          suporte@nexo.app
        </a>
      </p>
    </div>
  );
}
