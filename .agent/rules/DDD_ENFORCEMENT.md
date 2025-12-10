O NEXO segue Domain-Driven Design.

Regras:
- Usar linguagem ubíqua igual ao negócio.
- Cada módulo é um Bounded Context independente (Financeiro, Assinaturas, Estoque, Lista da Vez, Agendamento, CRM).
- Entidades representam conceitos do domínio e possuem identidade.
- Value Objects são imutáveis, sem identidade, e responsáveis por validações.
- Aggregates controlam invariantes e devem ser a única forma de mutação de entidades internas.
- Regras de negócio residem exclusivamente no domínio ou em Domain Services.
- Não permitir que infraestrutura ou handlers modifiquem entidades diretamente.
