import { z } from 'zod';

export const supplierSchema = z.object({
    razao_social: z
        .string()
        .min(2, 'Razão social deve ter pelo menos 2 caracteres')
        .max(200, 'Razão social deve ter no máximo 200 caracteres'),
    nome_fantasia: z.string().max(200).optional(),
    cnpj: z
        .string()
        .min(14, 'CNPJ deve ter 14 dígitos') // Assuming raw digit check or handle by regex
        .max(18, 'CNPJ muito longo') // Formatted
        .optional()
        .or(z.literal('')),
    email: z.string().email('Email inválido').optional().or(z.literal('')),
    telefone: z
        .string()
        .min(10, 'Telefone deve ter pelo menos 10 dígitos')
        .max(15, 'Telefone deve ter no máximo 15 dígitos'),
    endereco_logradouro: z.string().max(300).optional(),
    endereco_cidade: z.string().max(100).optional(),
    endereco_estado: z
        .string()
        .length(2, 'Estado deve ter 2 caracteres (UF)')
        .optional()
        .or(z.literal('')),
    endereco_cep: z
        .string()
        .min(8, 'CEP inválido')
        .optional()
        .or(z.literal('')),
});

export type SupplierFormValues = z.infer<typeof supplierSchema>;
