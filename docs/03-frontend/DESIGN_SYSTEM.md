# NEXO — Design System v1.0

> Identidade visual: **Clean & Professional** – foco na usabilidade, clareza e performance.
> Objetivo: Padronizar a interface do NEXO, garantindo consistência e velocidade de desenvolvimento.

---

## 1. Visão Geral

O **NEXO Design System** é o conjunto de decisões de design codificadas em tokens, componentes e padrões.

Ele garante que:

- O sistema tenha uma aparência coesa.
- O desenvolvimento seja rápido (copiar e colar componentes prontos).
- A manutenção seja simples (alterar tokens globais reflete em tudo).

---

## 2. Tecnologias

O Design System é construído sobre:

1.  **Tailwind CSS 4:** Motor de estilização utilitária.
2.  **shadcn/ui:** Biblioteca de componentes headless (Radix UI) com estilos Tailwind.
3.  **CSS Variables:** Tokens de design semânticos (ex: `--primary`, `--destructive`).
4.  **Lucide React:** Ícones consistentes.
5.  **Inter & JetBrains Mono:** Tipografia padrão.

---

## 3. Estrutura da Documentação

1.  **Foundations (`01-FOUNDATIONS.md`):**

    - Cores, Tipografia, Espaçamento, Sombras, Breakpoints.
    - A "física" do design system.

2.  **Architecture (`02-ARCHITECTURE.md`):**

    - Stack técnico, estrutura de pastas, gerenciamento de estado.
    - Como o código é organizado.

3.  **Components (`03-COMPONENTS.md`):**

    - Catálogo de componentes (Botões, Inputs, Cards, Tabelas).
    - Como usar os componentes do shadcn/ui.

4.  **Patterns (`04-PATTERNS.md`):**
    - Receitas de código para problemas comuns.
    - Formulários, Data Fetching, Tratamento de Erros.

---

## 4. Princípios de Design

### 4.1 Simplicidade

Menos é mais. Evite decorações desnecessárias. O foco é o conteúdo e a ação do usuário.

### 4.2 Consistência

Um botão "Salvar" deve parecer e funcionar da mesma forma em todas as telas. Use os componentes existentes, não invente novos sem necessidade.

### 4.3 Feedback

Toda ação do usuário deve ter uma resposta.

- Clique -> Estado de loading ou active.
- Sucesso -> Toast verde.
- Erro -> Toast vermelho ou mensagem inline.

### 4.4 Acessibilidade

Não é opcional.

- Use HTML semântico.
- Garanta contraste de cores.
- Suporte navegação por teclado.
- Os componentes do shadcn/ui (Radix) já resolvem 90% disso, não quebre o que já funciona.

---

## 5. Workflow de Desenvolvimento

1.  **Nova Tela:** Comece pelo `layout.tsx` ou `page.tsx`.
2.  **Estrutura:** Use os componentes de layout (`Card`, `Grid`).
3.  **Interatividade:** Adicione os componentes de UI (`Button`, `Input`).
4.  **Estado:** Conecte com `React Hook Form` e `TanStack Query`.
5.  **Refino:** Ajuste espaçamentos e responsividade com classes Tailwind.

---

## 6. Customização do Tema

O tema é controlado pelo arquivo `src/app/globals.css`.
Para alterar a cor primária do sistema, basta alterar a variável `--primary` (e `--primary-foreground`) neste arquivo. O Tailwind irá propagar a mudança para todos os componentes que usam `bg-primary`, `text-primary`, etc.
