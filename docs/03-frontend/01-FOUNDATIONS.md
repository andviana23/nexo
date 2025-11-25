# NEXO — Design Foundations

> Base visual do NEXO: cores, tipografia, espaçamentos, efeitos visuais e breakpoints.
> Se você está mexendo com UI, começa por aqui.

---

## 1. Identidade Visual

A identidade visual do NEXO é **Clean & Professional**:

- **Precisão** – interfaces calculadas, alinhadas, sem excessos visuais.
- **Clareza** – hierarquia clara, espaçamentos generosos, foco no conteúdo.
- **Modernidade** – design contemporâneo com transições suaves e microinterações.
- **Acessibilidade** – contraste adequado, foco visível, navegação por teclado.

---

## 2. Paleta de Cores (Tokens)

A paleta de cores do NEXO é **imutável**. Você **não pode criar novos valores**, apenas usar os tokens existentes.

Todos os tokens são expostos via **CSS Variables `--nexo-*`** no arquivo `globals.css` (ou mapeados via Tailwind).

### 2.1 Light Mode — Tema Padrão

Tema padrão do sistema. Clean e profissional:

| Token            | Valor     | Uso                      |
| ---------------- | --------- | ------------------------ |
| `background`     | `#F1F5F9` | Fundo geral da aplicação |
| `surface`        | `#FFFFFF` | Cards, painéis, modais   |
| `surface-subtle` | `#F8FAFC` | Backgrounds secundários  |
| `text-primary`   | `#0F172A` | Texto principal          |
| `text-secondary` | `#64748B` | Texto auxiliar, labels   |
| `text-muted`     | `#94A3B8` | Placeholders, hints      |
| `primary`        | `#0F2A4A` | Cor principal da marca   |
| `accent`         | `#2563EB` | Links, ações, CTAs       |
| `border`         | `#E2E8F0` | Bordas, divisores        |
| `success`        | `#38D69B` | Estados de sucesso       |
| `warning`        | `#F4B23E` | Estados de alerta        |
| `danger`         | `#EF4444` | Estados de erro          |

> Light é **sempre** o tema inicial. Dark é opt-in.

### 2.2 Dark Mode — Premium

Tema de alto contraste para uso noturno:

| Token            | Valor     | Uso                     |
| ---------------- | --------- | ----------------------- |
| `background`     | `#020617` | Fundo geral             |
| `surface`        | `#1E293B` | Cards, painéis          |
| `surface-subtle` | `#0F172A` | Backgrounds secundários |
| `text-primary`   | `#F8FAFC` | Texto principal         |
| `text-secondary` | `#94A3B8` | Texto auxiliar          |
| `text-muted`     | `#64748B` | Placeholders            |
| `primary`        | `#3B82F6` | Cor principal           |
| `accent`         | `#60A5FA` | Links, ações            |
| `border`         | `#334155` | Bordas, divisores       |
| `success`        | `#38D69B` | Sucesso                 |
| `warning`        | `#F4B23E` | Alerta                  |
| `danger`         | `#EF4444` | Erro                    |

---

## 3. CSS Variables

As variáveis CSS são definidas em `:root` e sobrescritas com `.dark`:

```css
:root {
  /* Backgrounds */
  --background: 210 40% 96.1%; /* #F1F5F9 */
  --foreground: 222.2 84% 4.9%; /* #0F172A */

  --card: 0 0% 100%; /* #FFFFFF */
  --card-foreground: 222.2 84% 4.9%;

  --popover: 0 0% 100%;
  --popover-foreground: 222.2 84% 4.9%;

  --primary: 213 66% 17%; /* #0F2A4A */
  --primary-foreground: 210 40% 98%;

  --secondary: 210 40% 96.1%;
  --secondary-foreground: 222.2 47.4% 11.2%;

  --muted: 210 40% 96.1%;
  --muted-foreground: 215.4 16.3% 46.9%;

  --accent: 210 40% 96.1%;
  --accent-foreground: 222.2 47.4% 11.2%;

  --destructive: 0 84.2% 60.2%;
  --destructive-foreground: 210 40% 98%;

  --border: 214.3 31.8% 91.4%; /* #E2E8F0 */
  --input: 214.3 31.8% 91.4%;
  --ring: 222.2 84% 4.9%;

  --radius: 0.5rem;
}

.dark {
  --background: 222.2 84% 4.9%; /* #020617 */
  --foreground: 210 40% 98%; /* #F8FAFC */

  --card: 222.2 84% 4.9%;
  --card-foreground: 210 40% 98%;

  --popover: 222.2 84% 4.9%;
  --popover-foreground: 210 40% 98%;

  --primary: 217.2 91.2% 59.8%; /* #3B82F6 */
  --primary-foreground: 222.2 47.4% 11.2%;

  --secondary: 217.2 32.6% 17.5%;
  --secondary-foreground: 210 40% 98%;

  --muted: 217.2 32.6% 17.5%;
  --muted-foreground: 215 20.2% 65.1%;

  --accent: 217.2 32.6% 17.5%;
  --accent-foreground: 210 40% 98%;

  --destructive: 0 62.8% 30.6%;
  --destructive-foreground: 210 40% 98%;

  --border: 217.2 32.6% 17.5%; /* #334155 */
  --input: 217.2 32.6% 17.5%;
  --ring: 212.7 26.8% 83.9%;
}
```

---

## 4. Tipografia

### 4.1 Fontes

- **Primária:** `Inter` (sans-serif)

  - Usada em todo o sistema: títulos, corpo, labels.
  - Carregada via `next/font/google`.

- **Monoespaçada:** `JetBrains Mono`
  - Códigos, IDs, valores numéricos críticos.

### 4.2 Escala Tipográfica

| Classe Tailwind | Tamanho | Peso | Uso                    |
| --------------- | ------- | ---- | ---------------------- |
| `text-xs`       | 12px    | 400  | Captions, badges       |
| `text-sm`       | 14px    | 400  | Labels, texto auxiliar |
| `text-base`     | 16px    | 400  | Texto padrão           |
| `text-lg`       | 18px    | 500  | Subtítulos             |
| `text-xl`       | 20px    | 600  | Títulos de seção       |
| `text-2xl`      | 24px    | 700  | Títulos de card        |
| `text-3xl`      | 30px    | 700  | Títulos de página      |
| `text-4xl`      | 36px    | 800  | Métricas principais    |

---

## 5. Espaçamento

Base de espaçamento de **4px** (0.25rem) do Tailwind:

| Token | Valor | Uso                     |
| ----- | ----- | ----------------------- |
| `1`   | 4px   | Micro espaçamentos      |
| `2`   | 8px   | Entre elementos inline  |
| `3`   | 12px  | Padding interno pequeno |
| `4`   | 16px  | Padding padrão          |
| `5`   | 20px  | Gap em grids            |
| `6`   | 24px  | Padding de cards        |
| `8`   | 32px  | Seções                  |
| `10`  | 40px  | Gaps grandes            |
| `12`  | 48px  | Separação de blocos     |

**Guidelines:**

- Cards: padding `p-6` (24px)
- Formulários: gap `gap-4` (16px) entre campos
- Seções: margin `my-8` (32px)

---

## 6. Bordas e Raios

### Border Radius

| Token          | Valor  | Uso             |
| -------------- | ------ | --------------- |
| `rounded-sm`   | 6px    | Badges, chips   |
| `rounded`      | 8px    | Inputs, buttons |
| `rounded-lg`   | 12px   | Cards           |
| `rounded-xl`   | 16px   | Modais          |
| `rounded-full` | 9999px | Avatares, pills |

### Borders

- Cor padrão: `border-border` (mapeado para `--border`)
- Espessura: `border` (1px)

---

## 7. Sombras

| Token       | Valor | Uso               |
| ----------- | ----- | ----------------- |
| `shadow-sm` | Leve  | Inputs, badges    |
| `shadow`    | Média | Cards             |
| `shadow-lg` | Forte | Dropdowns, modais |

**Uso em Tailwind:**

```tsx
<div className="shadow-lg rounded-xl bg-card">...</div>
```

---

## 8. Responsividade (OBRIGATÓRIO)

> ⚠️ **REGRA INEGOCIÁVEL:** Toda página, componente e layout do NEXO **DEVE** ser responsivo.
> PRs sem responsividade adequada **SERÃO REJEITADOS**.

### 8.1 Filosofia Mobile-First

O NEXO adota a abordagem **Mobile-First**:

1. Escreva o CSS/classes pensando **primeiro no mobile**.
2. Use breakpoints para **adicionar** estilos em telas maiores.
3. **Nunca** faça o contrário (desktop-first com `max-width`).

### 8.2 Breakpoints

| Token | Min-width | Dispositivo       | Uso                                 |
| ----- | --------- | ----------------- | ----------------------------------- |
| Base  | 0px       | Mobile            | Estilos padrão (sem prefixo)        |
| `sm`  | 640px     | Mobile landscape  | Ajustes leves                       |
| `md`  | 768px     | Tablet            | Layouts 2 colunas                   |
| `lg`  | 1024px    | Desktop           | Layouts 3+ colunas, sidebar visível |
| `xl`  | 1280px    | Desktop wide      | Conteúdo mais espaçado              |
| `2xl` | 1536px    | Monitores grandes | Máxima largura de containers        |

### 8.3 Regras por Tipo de Componente

| Componente        | Mobile                                | Tablet (md)                   | Desktop (lg+)           |
| ----------------- | ------------------------------------- | ----------------------------- | ----------------------- |
| **Grid de Cards** | 1 coluna                              | 2 colunas                     | 3-4 colunas             |
| **Formulários**   | 1 coluna, campos full-width           | 2 colunas para campos curtos  | Máx 3 colunas           |
| **Tabelas**       | Scroll horizontal ou cards empilhados | Tabela com colunas essenciais | Tabela completa         |
| **Sidebar**       | Drawer (Sheet)                        | Drawer ou colapsada           | Fixa e expandida        |
| **Modais**        | Full-screen ou quase                  | Centralizado, 80% largura     | Centralizado, max-width |
| **Gráficos**      | Simplificado ou scroll                | Tamanho médio                 | Tamanho completo        |

### 8.4 Padrões Obrigatórios

```tsx
// ✅ CORRETO: Mobile-first
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">

// ❌ ERRADO: Desktop-first (não usar)
<div className="grid grid-cols-3 md:grid-cols-2 sm:grid-cols-1 gap-4">

// ✅ CORRETO: Container responsivo
<div className="container mx-auto px-4 sm:px-6 lg:px-8">

// ✅ CORRETO: Texto responsivo
<h1 className="text-2xl md:text-3xl lg:text-4xl font-bold">

// ✅ CORRETO: Espaçamento responsivo
<section className="py-8 md:py-12 lg:py-16">
```

### 8.5 Checklist de Responsividade

Antes de considerar uma página/componente pronto:

- [ ] Testou em viewport 375px (iPhone SE)?
- [ ] Testou em viewport 768px (iPad)?
- [ ] Testou em viewport 1024px (Desktop)?
- [ ] Textos não quebram de forma estranha?
- [ ] Botões/links têm área de toque mínima de 44px?
- [ ] Tabelas têm scroll horizontal ou alternativa mobile?
- [ ] Modais não cortam conteúdo em telas pequenas?
- [ ] Sidebar funciona como drawer no mobile?

---

## 9. Animações (Framer Motion)

### Transições Padrão

```typescript
// Entrada suave
export const fadeIn = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  transition: { duration: 0.2 },
};

// Slide up
export const slideUp = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  transition: { duration: 0.3, ease: 'easeOut' },
};
```

### Durações

| Tipo   | Duração | Uso               |
| ------ | ------- | ----------------- |
| Rápida | 150ms   | Hover, focus      |
| Normal | 200ms   | Transições gerais |
| Suave  | 300ms   | Modais, drawers   |
