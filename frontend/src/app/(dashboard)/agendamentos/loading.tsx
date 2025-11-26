/**
 * Loading state para a página de Agendamentos
 */

import { Skeleton } from '@/components/ui/skeleton';

export default function AgendamentosLoading() {
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="space-y-2">
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-4 w-64" />
        </div>
        <div className="flex gap-2">
          <Skeleton className="h-9 w-24" />
          <Skeleton className="h-9 w-36" />
        </div>
      </div>

      {/* Calendário */}
      <div className="rounded-lg border bg-card p-4 shadow-sm">
        {/* Header do calendário */}
        <div className="flex justify-between items-center mb-4">
          <Skeleton className="h-10 w-32" />
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-10 w-64" />
        </div>

        {/* Grid do calendário */}
        <div className="space-y-2">
          {/* Header das colunas */}
          <div className="flex gap-2">
            <Skeleton className="h-6 w-16" />
            <Skeleton className="h-6 flex-1" />
            <Skeleton className="h-6 flex-1" />
            <Skeleton className="h-6 flex-1" />
          </div>

          {/* Linhas do calendário */}
          {Array.from({ length: 12 }).map((_, i) => (
            <div key={i} className="flex gap-2">
              <Skeleton className="h-12 w-16" />
              <Skeleton className="h-12 flex-1" />
              <Skeleton className="h-12 flex-1" />
              <Skeleton className="h-12 flex-1" />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
