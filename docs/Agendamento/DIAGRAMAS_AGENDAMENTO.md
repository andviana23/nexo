# Diagramas — Módulo de Agendamento | NEXO v1.0

Este documento contém os diagramas técnicos para o módulo de Agendamento, utilizando a sintaxe **Mermaid**.

---

## 1. Fluxo de Criação de Agendamento (Sequence Diagram)

Este diagrama ilustra a interação entre Frontend, API, Banco de Dados e Google Calendar durante a criação de um agendamento.

```mermaid
sequenceDiagram
    participant User as Usuário (Front)
    participant API as Backend API
    participant DB as PostgreSQL
    participant GCal as Google Calendar

    User->>API: POST /appointments (dados)
    API->>API: Validar Token & Tenant
    API->>DB: Verificar Conflitos (SELECT)
    
    alt Horário Ocupado
        DB-->>API: Conflito Encontrado
        API-->>User: 409 Conflict
    else Horário Livre
        API->>DB: Criar Agendamento (INSERT)
        DB-->>API: Agendamento Criado (ID)
        
        par Sync Google Calendar
            API->>GCal: Criar Evento
            GCal-->>API: Evento Criado (gcal_id)
            API->>DB: Atualizar gcal_id
        and Notificar Cliente
            API->>API: Enviar WhatsApp/Email (Async)
        end
        
        API-->>User: 201 Created (JSON)
    end
```

---

## 2. Máquina de Estados do Agendamento (State Diagram)

Estados possíveis para um agendamento e as transições permitidas.

```mermaid
stateDiagram-v2
    [*] --> CREATED: Agendamento Criado
    
    CREATED --> CONFIRMED: Confirmação Automática/Manual
    CREATED --> CANCELED: Cancelado pelo Cliente/Barbeiro
    
    CONFIRMED --> IN_PROGRESS: Cliente Chegou / Serviço Iniciou
    CONFIRMED --> NO_SHOW: Cliente não apareceu
    CONFIRMED --> CANCELED: Cancelado tardiamente
    
    IN_PROGRESS --> COMPLETED: Serviço Finalizado e Pago
    
    NO_SHOW --> [*]
    CANCELED --> [*]
    COMPLETED --> [*]

    note right of CREATED
        Estado inicial.
        Aguardando confirmação se configurado.
    end note

    note right of COMPLETED
        Gera lançamento financeiro
        e comissão.
    end note
```

---

## 3. Modelo de Dados Simplificado (ER Diagram)

Relacionamento entre as principais tabelas do módulo.

```mermaid
erDiagram
    TENANTS ||--o{ APPOINTMENTS : "possui"
    PROFESSIONALS ||--o{ APPOINTMENTS : "realiza"
    CUSTOMERS ||--o{ APPOINTMENTS : "solicita"
    
    APPOINTMENTS ||--|{ APPOINTMENT_SERVICES : "contém"
    SERVICES ||--o{ APPOINTMENT_SERVICES : "define"

    APPOINTMENTS {
        uuid id PK
        uuid tenant_id FK
        uuid professional_id FK
        uuid customer_id FK
        timestamp start_time
        timestamp end_time
        string status
        decimal total_price
        string google_calendar_event_id
    }

    APPOINTMENT_SERVICES {
        uuid appointment_id FK
        uuid service_id FK
        decimal price_at_booking
        int duration_at_booking
    }

    SERVICES {
        uuid id PK
        string name
        decimal default_price
        int default_duration
    }
```

---

## 4. Fluxo de Sincronização Google Calendar (Flowchart)

Lógica de decisão para sincronização bidirecional (futuro) ou unidirecional (MVP).

```mermaid
flowchart TD
    A[Início Sync] --> B{Token Google Válido?}
    B -- Não --> C[Ignorar Sync]
    B -- Sim --> D{Tipo de Operação}
    
    D -- Criar --> E[Insert no GCal]
    D -- Atualizar --> F[Update no GCal]
    D -- Deletar --> G[Delete no GCal]
    
    E --> H[Salvar gcal_event_id no DB]
    F --> I[Logar Sucesso]
    G --> I
    
    H --> J[Fim]
    I --> J
    C --> J
```

---

**Responsável:** Tech Lead  
**Data:** 25/11/2025  
**Status:** ✅ COMPLETO
