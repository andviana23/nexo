/**
 * NEXO - Sistema de Gestão para Barbearias
 * Validações Zod: Professional (Profissional)
 * 
 * @module lib/validations/professional
 * @description Schemas de validação para formulários de profissionais
 */

import { z } from 'zod';

// ============================================================================
// SCHEMAS AUXILIARES
// ============================================================================

/** Validação de CPF ou CNPJ (11 ou 14 dígitos numéricos) */
const cpfCnpjSchema = z.string()
  .min(11, 'CPF/CNPJ deve ter 11 ou 14 dígitos')
  .max(14, 'CPF/CNPJ deve ter 11 ou 14 dígitos')
  .regex(/^\d+$/, 'CPF/CNPJ deve conter apenas números')
  .refine((val) => val.length === 11 || val.length === 14, {
    message: 'CPF deve ter 11 dígitos ou CNPJ deve ter 14 dígitos',
  });

/** Validação de telefone brasileiro (10-11 dígitos) */
const phoneSchema = z.string()
  .min(10, 'Telefone deve ter no mínimo 10 dígitos')
  .max(11, 'Telefone deve ter no máximo 11 dígitos')
  .regex(/^\d+$/, 'Telefone deve conter apenas números');

/** Validação de horário HH:MM */
const timeSchema = z.string()
  .regex(/^([01]\d|2[0-3]):([0-5]\d)$/, 'Formato inválido (HH:MM)');

/** Turno de trabalho */
const workShiftSchema = z.object({
  inicio: timeSchema,
  fim: timeSchema,
}).refine(data => data.inicio < data.fim, {
  message: 'Horário de início deve ser menor que fim',
  path: ['fim'],
});

/** Dia da semana */
const weekDaySchema = z.object({
  ativo: z.boolean(),
  turnos: z.array(workShiftSchema),
});

/** Horário de trabalho completo */
const workScheduleSchema = z.object({
  segunda: weekDaySchema,
  terca: weekDaySchema,
  quarta: weekDaySchema,
  quinta: weekDaySchema,
  sexta: weekDaySchema,
  sabado: weekDaySchema,
  domingo: weekDaySchema,
});

// ============================================================================
// SCHEMA PRINCIPAL
// ============================================================================

/** Schema do formulário de criação de profissional */
export const createProfessionalSchema = z.object({
  // Dados obrigatórios
  nome: z.string()
    .min(3, 'Nome deve ter no mínimo 3 caracteres')
    .max(255, 'Nome muito longo'),

  email: z.string()
    .email('Email inválido')
    .max(255, 'Email muito longo'),

  telefone: phoneSchema,

  cpf: cpfCnpjSchema.optional().or(z.literal('')),

  tipo: z.enum(['BARBEIRO', 'GERENTE', 'RECEPCIONISTA', 'OUTRO'], {
    required_error: 'Selecione o tipo de profissional',
  }),

  data_admissao: z.date({
    required_error: 'Data de admissão é obrigatória',
  }).max(new Date(), 'Data não pode ser futura'),

  // Campos opcionais
  foto: z.string().url('URL inválida').optional().nullable().or(z.literal('')),

  especialidades: z.array(z.string()).optional(),

  observacoes: z.string().max(500, 'Máximo 500 caracteres').optional(),

  // Campos condicionais (Gerente)
  tambem_barbeiro: z.boolean().default(false),

  // Comissão
  tipo_comissao: z.enum(['PERCENTUAL', 'FIXO']).optional(),

  comissao: z.number()
    .min(0, 'Comissão não pode ser negativa')
    .max(100, 'Comissão máxima é 100%')
    .optional()
    .nullable()
    .default(0),

  comissao_produtos: z.number()
    .min(0, 'Comissão não pode ser negativa')
    .max(100, 'Comissão máxima é 100%')
    .optional()
    .nullable(),

  // Comissões por categoria (Frontend only por enquanto)
  comissoes_por_categoria: z.array(z.object({
    categoria_id: z.string(),
    comissao: z.number().min(0).max(100),
  })).optional(),

  // Horário de trabalho
  horario_trabalho: workScheduleSchema.optional(),

});
// Nota: removido superRefine que exigia comissão obrigatória para barbeiros.
// A comissão agora tem valor padrão 0, permitindo cadastro sem preencher aba financeira.

/** Schema do formulário de edição de profissional */
export const updateProfessionalSchema = z.object({
  nome: z.string()
    .min(3, 'Nome deve ter no mínimo 3 caracteres')
    .max(255, 'Nome muito longo')
    .optional(),

  email: z.string()
    .email('Email inválido')
    .max(255, 'Email muito longo')
    .optional(),

  telefone: phoneSchema.optional(),

  foto: z.string().url('URL inválida').optional().nullable().or(z.literal('')),

  especialidades: z.array(z.string()).optional(),

  observacoes: z.string().max(500, 'Máximo 500 caracteres').optional(),

  tipo_comissao: z.enum(['PERCENTUAL', 'FIXO']).optional(),

  comissao: z.number()
    .min(0, 'Comissão não pode ser negativa')
    .max(100, 'Comissão máxima é 100%')
    .optional()
    .nullable(),

  comissao_produtos: z.number()
    .min(0, 'Comissão não pode ser negativa')
    .max(100, 'Comissão máxima é 100%')
    .optional()
    .nullable(),

  // Comissões por categoria (Frontend only por enquanto)
  comissoes_por_categoria: z.array(z.object({
    categoria_id: z.string(),
    comissao: z.number().min(0).max(100),
  })).optional(),

  horario_trabalho: workScheduleSchema.optional(),

  status: z.enum(['ATIVO', 'INATIVO', 'AFASTADO', 'DEMITIDO']).optional(),
});

// ============================================================================
// TYPES
// ============================================================================

export type CreateProfessionalFormData = z.infer<typeof createProfessionalSchema>;
export type UpdateProfessionalFormData = z.infer<typeof updateProfessionalSchema>;

// ============================================================================
// HELPERS
// ============================================================================

/** Remove máscara do CPF/CNPJ */
export function unmaskCPF(cpf: string): string {
  return cpf.replace(/\D/g, '');
}

/** Remove máscara do CPF/CNPJ (alias) */
export function unmaskCPFCNPJ(doc: string): string {
  return doc.replace(/\D/g, '');
}

/** Remove máscara do telefone */
export function unmaskPhone(phone: string): string {
  return phone.replace(/\D/g, '');
}

/** Aplica máscara ao CPF ou CNPJ */
export function maskCPF(cpf: string): string {
  const cleaned = cpf.replace(/\D/g, '');

  // CNPJ (14 dígitos): 00.000.000/0000-00
  if (cleaned.length === 14) {
    return cleaned
      .replace(/(\d{2})(\d)/, '$1.$2')
      .replace(/(\d{3})(\d)/, '$1.$2')
      .replace(/(\d{3})(\d)/, '$1/$2')
      .replace(/(\d{4})(\d{1,2})$/, '$1-$2');
  }

  // CPF (11 dígitos): 000.000.000-00
  return cleaned
    .replace(/(\d{3})(\d)/, '$1.$2')
    .replace(/(\d{3})(\d)/, '$1.$2')
    .replace(/(\d{3})(\d{1,2})$/, '$1-$2');
}

/** Aplica máscara ao CPF ou CNPJ (alias) */
export function maskCPFCNPJ(doc: string): string {
  return maskCPF(doc);
}

/** Aplica máscara ao telefone */
export function maskPhone(phone: string): string {
  const cleaned = phone.replace(/\D/g, '');
  if (cleaned.length <= 10) {
    return cleaned
      .replace(/(\d{2})(\d)/, '($1) $2')
      .replace(/(\d{4})(\d)/, '$1-$2');
  }
  return cleaned
    .replace(/(\d{2})(\d)/, '($1) $2')
    .replace(/(\d{5})(\d)/, '$1-$2');
}

/** Remove máscara do CEP */
export function unmaskCEP(cep: string): string {
  return cep.replace(/\D/g, '');
}

/** Aplica máscara ao CEP */
export function maskCEP(cep: string): string {
  const cleaned = cep.replace(/\D/g, '');
  return cleaned.replace(/^(\d{5})(\d)/, '$1-$2');
}

/** Remove máscara do CNPJ (alias) */
export function unmaskCNPJ(cnpj: string): string {
  return unmaskCPF(cnpj);
}

/** Aplica máscara ao CNPJ (alias) */
export function maskCNPJ(cnpj: string): string {
  return maskCPF(cnpj);
}
