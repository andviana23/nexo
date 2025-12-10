import { z } from 'zod';

export const customerSchema = z.object({
    nome: z.string().min(3, 'Nome deve ter pelo menos 3 caracteres'),
    telefone: z.string().min(10, 'Telefone inválido'), // Assume phone mask
    email: z.string().email('Email inválido').optional().or(z.literal('')),
    cpf: z.string().optional().or(z.literal('')),
    data_nascimento: z.string().optional().or(z.literal('')),
    genero: z.enum(['M', 'F', 'NB', 'PNI']).optional(),
    endereco_cep: z.string().optional().or(z.literal('')),
    endereco_logradouro: z.string().optional().or(z.literal('')),
    endereco_numero: z.string().optional().or(z.literal('')),
    endereco_complemento: z.string().optional().or(z.literal('')),
    endereco_bairro: z.string().optional().or(z.literal('')),
    endereco_cidade: z.string().optional().or(z.literal('')),
    endereco_estado: z.string().optional().or(z.literal('')),
    observacoes: z.string().max(1000, 'Máximo de 1000 caracteres').optional(),
    tags: z.array(z.string()).optional(),
    ativo: z.boolean().default(true).optional(),
});

export type CustomerFormValues = z.infer<typeof customerSchema>;

// Masks and Utilities can be imported or re-implemented if needed, 
// usually they are in component or shared utils.
