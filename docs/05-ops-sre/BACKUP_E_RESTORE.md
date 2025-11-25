> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🗄️ Backup & Disaster Recovery — Barber Analytics Pro

**Estratégia de Backup e Recuperação de Desastres**
**Versão:** 1.0.0
**Data:** 15/11/2025
**Status:** 🟡 Em Implementação

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Política de Backup](#política-de-backup)
3. [Backup Automático (Neon)](#backup-automático-neon)
4. [Backup Complementar (pg_dump)](#backup-complementar-pg_dump)
5. [Testes de Restore](#testes-de-restore)
6. [Disaster Recovery Playbook](#disaster-recovery-playbook)
7. [Objetivos RTO/RPO](#objetivos-rtorpo)
8. [Checklist de Validação](#checklist-de-validação)

---

## 🎯 Visão Geral

### Escopo

Este documento descreve a estratégia de **Backup e Disaster Recovery (DR)** para o sistema **Barber Analytics Pro**, incluindo:

- Backups automáticos do banco de dados (Neon PostgreSQL)
- Backups complementares via `pg_dump`
- Procedimentos de restore
- Plano de recuperação de desastres
- Testes periódicos

### Ativos Críticos

| Ativo                                | Criticidade | Backup Necessário        |
| ------------------------------------ | ----------- | ------------------------ |
| **Database (Neon)**                  | 🔴 Crítico  | ✅ SIM                   |
| **Backend Go (código)**              | 🟡 Alto     | ✅ SIM (Git)             |
| **Frontend Next.js 14.2.4 (código)** | 🟡 Alto     | ✅ SIM (Git)             |
| **Chaves JWT (keys/)**               | 🔴 Crítico  | ✅ SIM (secrets manager) |
| **Variáveis de ambiente**            | 🔴 Crítico  | ✅ SIM (secrets manager) |
| **Logs**                             | 🟢 Baixo    | ⏳ Opcional (journald)   |

---

## 📦 Política de Backup

### Retenção

| Tipo de Backup       | Frequência         | Retenção | Responsável         |
| -------------------- | ------------------ | -------- | ------------------- |
| **Neon PITR**        | Contínuo (WAL)     | 7 dias   | Neon (automático)   |
| **pg_dump diário**   | Diário (03:00 UTC) | 30 dias  | GitHub Actions + S3 |
| **Snapshot semanal** | Semanal (domingos) | 90 dias  | GitHub Actions + S3 |
| **Snapshot mensal**  | Mensal (dia 1)     | 1 ano    | GitHub Actions + S3 |
| **Código-fonte**     | Cada push          | Infinito | GitHub              |

### RPO/RTO

| Cenário                         | RPO (Perda Máxima)   | RTO (Tempo de Recuperação) |
| ------------------------------- | -------------------- | -------------------------- |
| **Database corruption**         | < 1 hora (Neon PITR) | < 2 horas                  |
| **Database deletion acidental** | < 24 horas (pg_dump) | < 4 horas                  |
| **Disaster total (AWS outage)** | < 24 horas           | < 8 horas                  |
| **Application bug**             | 0 (rollback código)  | < 30 minutos               |

**Meta:**

- **RPO:** < 24 horas
- **RTO:** < 4 horas

---

## 🚀 Backup Automático (Neon)

### Neon Point-in-Time Recovery (PITR)

**O que é:**

- Neon mantém backups contínuos via Write-Ahead Log (WAL)
- Permite restaurar para qualquer ponto no tempo dentro da janela de retenção

**Configuração atual:**

```yaml
Plano: Pro
Retenção PITR: 7 dias
Snapshots automáticos: Sim (1x/dia)
Região: us-east-2 (AWS)
```

**Como restaurar:**

1. **Via Neon Console:**

   - Acessar: https://console.neon.tech
   - Selecionar projeto: `barber-analytics-prod`
   - Clicar em "Branches" → "Restore to point in time"
   - Escolher timestamp (ex: 2025-11-14 10:30:00 UTC)
   - Criar novo branch com dados restaurados

2. **Via CLI:**

```bash
# Instalar Neon CLI
npm install -g neonctl

# Autenticar
neonctl auth login

# Criar branch de restore
neonctl branches create \
  --project-id ep-winter-leaf-adhqz08p \
  --name "restore-2025-11-14" \
  --point-in-time "2025-11-14T10:30:00Z"

# Obter connection string do novo branch
neonctl connection-string restore-2025-11-14
```

**Vantagens:**

- ✅ Automático (zero configuração)
- ✅ Granularidade de segundos
- ✅ Sem impacto em performance
- ✅ Incluso no plano Pro

**Limitações:**

- ⚠️ Retenção limitada (7 dias no Pro, 30 dias no Business)
- ⚠️ Não protege contra exclusão do projeto Neon

---

## 💾 Backup Complementar (pg_dump)

### Por que pg_dump adicional?

- ✅ Retenção maior (30 dias vs 7 dias Neon)
- ✅ Backup off-site (S3, independente da Neon)
- ✅ Portabilidade (pode restaurar em qualquer PostgreSQL)
- ✅ Proteção contra exclusão acidental do projeto Neon

### Implementação via GitHub Actions

**Arquivo:** `.github/workflows/backup-database.yml`

```yaml
name: Database Backup

on:
  schedule:
    # Diário às 03:00 UTC (00:00 BRT)
    - cron: '0 3 * * *'
  workflow_dispatch: # Permitir trigger manual

jobs:
  backup:
    name: Backup PostgreSQL
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Install PostgreSQL client
        run: |
          sudo apt-get update
          sudo apt-get install -y postgresql-client

      - name: Create backup directory
        run: mkdir -p backups

      - name: Run pg_dump
        env:
          DATABASE_URL: ${{ secrets.DATABASE_URL_PROD }}
        run: |
          TIMESTAMP=$(date +%Y%m%d-%H%M%S)
          BACKUP_FILE="backups/barber-analytics-${TIMESTAMP}.sql"

          echo "Creating backup: $BACKUP_FILE"
          pg_dump "$DATABASE_URL" \
            --clean \
            --if-exists \
            --no-owner \
            --no-acl \
            --format=plain \
            --file="$BACKUP_FILE"

          # Comprimir backup
          gzip "$BACKUP_FILE"
          echo "BACKUP_FILE=${BACKUP_FILE}.gz" >> $GITHUB_ENV

      - name: Upload to S3
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_REGION: us-east-1
          S3_BUCKET: barber-analytics-backups
        run: |
          # Instalar AWS CLI
          pip install awscli

          # Upload com metadata
          aws s3 cp "$BACKUP_FILE" \
            "s3://$S3_BUCKET/daily/$BACKUP_FILE" \
            --metadata "timestamp=$(date -Iseconds)" \
            --storage-class STANDARD_IA

          echo "✅ Backup uploaded to S3"

      - name: Cleanup old backups (30 dias)
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_REGION: us-east-1
          S3_BUCKET: barber-analytics-backups
        run: |
          # Deletar arquivos mais antigos que 30 dias
          aws s3 ls "s3://$S3_BUCKET/daily/" | \
            awk '{print $4}' | \
            while read file; do
              file_date=$(echo $file | grep -oP '\d{8}')
              days_old=$(( ($(date +%s) - $(date -d $file_date +%s)) / 86400 ))

              if [ $days_old -gt 30 ]; then
                echo "Deleting old backup: $file (${days_old} days old)"
                aws s3 rm "s3://$S3_BUCKET/daily/$file"
              fi
            done

      - name: Notify on failure
        if: failure()
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: '❌ Database backup FAILED!'
          webhook_url: ${{ secrets.SLACK_WEBHOOK_URL }}
```

### Secrets Necessários

Configurar no GitHub (Settings → Secrets):

```bash
# Neon connection string
DATABASE_URL_PROD=postgresql://user:pass@ep-xxx.neon.tech/neondb?sslmode=require

# AWS S3 credentials
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=...

# Slack notifications (opcional)
SLACK_WEBHOOK_URL=https://hooks.slack.com/...
```

### Criar S3 Bucket

```bash
# Criar bucket
aws s3 mb s3://barber-analytics-backups \
  --region us-east-1

# Habilitar versionamento
aws s3api put-bucket-versioning \
  --bucket barber-analytics-backups \
  --versioning-configuration Status=Enabled

# Configurar lifecycle (deletar após 30 dias)
cat > lifecycle.json << 'EOF'
{
  "Rules": [
    {
      "Id": "DeleteOldBackups",
      "Status": "Enabled",
      "Prefix": "daily/",
      "Expiration": {
        "Days": 30
      }
    }
  ]
}
EOF

aws s3api put-bucket-lifecycle-configuration \
  --bucket barber-analytics-backups \
  --lifecycle-configuration file://lifecycle.json
```

---

## 🧪 Testes de Restore

### Objetivo

Validar que backups podem ser restaurados corretamente e o sistema funciona.

### Procedimento de Teste (Mensal)

**1. Escolher backup para teste:**

```bash
# Listar backups disponíveis
aws s3 ls s3://barber-analytics-backups/daily/

# Escolher backup recente (ex: de ontem)
BACKUP_FILE=barber-analytics-20251114-030000.sql.gz
```

**2. Criar banco de teste (staging):**

```bash
# Via Neon CLI: Criar branch de teste
neonctl branches create \
  --project-id ep-winter-leaf-adhqz08p \
  --name "restore-test-$(date +%Y%m%d)" \
  --parent main

# Obter connection string
TEST_DB_URL=$(neonctl connection-string restore-test-20251115)
```

**3. Restaurar backup:**

```bash
# Baixar backup do S3
aws s3 cp "s3://barber-analytics-backups/daily/$BACKUP_FILE" .

# Descomprimir
gunzip $BACKUP_FILE

# Restaurar no banco de teste
psql "$TEST_DB_URL" < ${BACKUP_FILE%.gz}
```

**4. Validar dados:**

```bash
# Verificar contagem de registros
psql "$TEST_DB_URL" -c "
SELECT
  (SELECT COUNT(*) FROM tenants) as tenants,
  (SELECT COUNT(*) FROM users) as users,
  (SELECT COUNT(*) FROM receitas) as receitas,
  (SELECT COUNT(*) FROM despesas) as despesas,
  (SELECT COUNT(*) FROM assinaturas) as assinaturas;
"

# Resultado esperado:
#  tenants | users | receitas | despesas | assinaturas
# ---------+-------+----------+----------+-------------
#       15 |    42 |     1250 |      890 |          38
```

**5. Testar aplicação:**

```bash
# Atualizar .env com connection string de teste
export DATABASE_URL="$TEST_DB_URL"

# Iniciar backend
cd backend
go run cmd/api/main.go

# Testar endpoint
curl http://localhost:8080/health
# Deve retornar: {"status":"healthy","database":"connected"}

# Testar login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"qa@barberpro.dev","password":"qa123456"}'
# Deve retornar access_token
```

**6. Medir tempo de restauração:**

```bash
# Anotar duração total do processo:
# - Download: X minutos
# - Descompressão: Y minutos
# - Restore: Z minutos
# Total: < 2 horas (meta RTO)
```

**7. Limpar ambiente de teste:**

```bash
# Deletar branch de teste após validação
neonctl branches delete restore-test-20251115
```

### Registro de Testes

Manter log em `docs/backup-tests.log`:

```
2025-11-15 10:00:00 UTC | Teste de Restore Mensal
Backup: barber-analytics-20251114-030000.sql.gz (450 MB)
Download: 3 min
Restore: 8 min
Validação: OK (15 tenants, 42 users, 1250 receitas)
RTO Real: 15 minutos ✅
Status: SUCESSO ✅
```

---

## 🚨 Disaster Recovery Playbook

### Cenários de Desastre

#### Cenário 1: Corrupção de Dados (Acidental)

**Sintomas:**

- Dados inconsistentes (ex: receitas zeradas, usuários sumindo)
- Erros de integridade referencial
- Aplicação funciona mas dados corrompidos

**Ações:**

1. **Identificar timestamp da corrupção:**

   ```bash
   # Revisar audit_logs
   psql "$DATABASE_URL" -c "
   SELECT * FROM audit_logs
   WHERE criado_em > NOW() - INTERVAL '24 hours'
   ORDER BY criado_em DESC;
   "
   ```

2. **Criar backup da situação atual (por segurança):**

   ```bash
   pg_dump "$DATABASE_URL" > corruption-backup-$(date +%Y%m%d).sql
   ```

3. **Restaurar via Neon PITR:**

   ```bash
   # Criar branch com dados de antes da corrupção
   neonctl branches create \
     --name "restore-before-corruption" \
     --point-in-time "2025-11-14T10:30:00Z"

   # Obter nova connection string
   NEW_DB_URL=$(neonctl connection-string restore-before-corruption)
   ```

4. **Validar dados restaurados:**

   ```bash
   # Testar queries críticas
   psql "$NEW_DB_URL" -c "SELECT COUNT(*) FROM receitas;"
   ```

5. **Promover para produção:**

   ```bash
   # Atualizar DATABASE_URL nos secrets
   # Reiniciar backend com nova connection string
   ssh deploy@vps "sudo systemctl restart barber-api"
   ```

6. **Verificar aplicação:**
   ```bash
   curl https://api.barberpro.dev/health
   ```

**RTO esperado:** < 2 horas

---

#### Cenário 2: Exclusão Acidental de Tabela

**Sintomas:**

- Erro: `relation "users" does not exist`
- Backend crashando ao iniciar

**Ações:**

1. **Parar tráfego para aplicação:**

   ```bash
   # Retornar página de manutenção no NGINX
   ssh deploy@vps "sudo systemctl stop barber-api"
   ```

2. **Baixar último backup pg_dump:**

   ```bash
   LATEST_BACKUP=$(aws s3 ls s3://barber-analytics-backups/daily/ | tail -1 | awk '{print $4}')
   aws s3 cp "s3://barber-analytics-backups/daily/$LATEST_BACKUP" .
   gunzip $LATEST_BACKUP
   ```

3. **Restaurar apenas tabela deletada:**

   ```bash
   # Extrair apenas CREATE + INSERT da tabela users
   grep -A 10000 "CREATE TABLE users" ${LATEST_BACKUP%.gz} > users_restore.sql

   # Aplicar no banco
   psql "$DATABASE_URL" < users_restore.sql
   ```

4. **Recriar índices se necessário:**

   ```bash
   psql "$DATABASE_URL" -c "
   CREATE INDEX IF NOT EXISTS idx_users_tenant_id_email ON users(tenant_id, email);
   "
   ```

5. **Reiniciar aplicação:**
   ```bash
   ssh deploy@vps "sudo systemctl start barber-api"
   ```

**RTO esperado:** < 1 hora

---

#### Cenário 3: Disaster Total (AWS Region Down)

**Sintomas:**

- Neon inacessível
- Toda região us-east-2 fora do ar
- Aplicação completamente offline

**Ações:**

1. **Ativar comunicação de emergência:**

   - Post em status page: "Sistema temporariamente indisponível"
   - Notificar clientes via email/WhatsApp

2. **Provisionar novo banco em região diferente:**

   ```bash
   # Criar projeto Neon em us-west-2
   neonctl projects create \
     --name "barber-analytics-dr" \
     --region us-west-2
   ```

3. **Restaurar último backup:**

   ```bash
   # Baixar backup mais recente
   LATEST_BACKUP=$(aws s3 ls s3://barber-analytics-backups/daily/ | tail -1 | awk '{print $4}')
   aws s3 cp "s3://barber-analytics-backups/daily/$LATEST_BACKUP" .
   gunzip $LATEST_BACKUP

   # Restaurar em novo banco
   DR_DB_URL="postgresql://user:pass@ep-xxx-us-west-2.neon.tech/neondb"
   psql "$DR_DB_URL" < ${LATEST_BACKUP%.gz}
   ```

4. **Atualizar DNS:**

   ```bash
   # Apontar api.barberpro.dev para novo VPS/região
   # (Assumindo VPS multi-região ou novo deploy)
   ```

5. **Atualizar variáveis de ambiente:**

   ```bash
   # GitHub Secrets: DATABASE_URL → novo connection string
   # VPS: /opt/barber-api/.env → DATABASE_URL=$DR_DB_URL
   ```

6. **Deploy em nova região:**

   ```bash
   # Trigger GitHub Actions deploy
   # ou SSH manual
   ssh deploy@vps-dr "sudo systemctl restart barber-api"
   ```

7. **Verificar funcionamento:**
   ```bash
   curl https://api.barberpro.dev/health
   ```

**RTO esperado:** < 8 horas (cenário raro)

---

### Contatos de Emergência

| Papel            | Nome              | Contato              | Responsabilidade    |
| ---------------- | ----------------- | -------------------- | ------------------- |
| **Tech Lead**    | Andrey Viana      | andrey@barberpro.dev | Decisão final em DR |
| **DevOps Lead**  | [TBD]             | devops@barberpro.dev | Execução técnica    |
| **Neon Support** | support@neon.tech | Ticket + Slack       | Suporte Neon        |
| **AWS Support**  | -                 | Console AWS          | Suporte S3/EC2      |

### Checklist de Ativação DR

- [ ] Identificar cenário de desastre
- [ ] Notificar stakeholders (Tech Lead, clientes)
- [ ] Acionar playbook correspondente
- [ ] Documentar cada ação em tempo real
- [ ] Validar restauração com testes
- [ ] Comunicar resolução aos clientes
- [ ] Realizar postmortem (48h após incidente)

---

## 📊 Objetivos RTO/RPO

### Definições

- **RPO (Recovery Point Objective):** Perda máxima de dados aceitável
- **RTO (Recovery Time Objective):** Tempo máximo de indisponibilidade

### Metas Atuais

| Serviço                 | RPO        | RTO       | Implementação                   |
| ----------------------- | ---------- | --------- | ------------------------------- |
| **Database**            | < 1 hora   | < 2 horas | Neon PITR (7 dias)              |
| **Database (disaster)** | < 24 horas | < 4 horas | pg_dump + S3 (30 dias)          |
| **Backend (código)**    | 0 (Git)    | < 30 min  | Git + CI/CD                     |
| **Frontend (código)**   | 0 (Git)    | < 30 min  | Git + Vercel                    |
| **Chaves JWT**          | N/A        | < 1 hora  | Secrets manager + backup manual |

### Medição de Sucesso

**Critérios:**

- ✅ Testes de restore mensais passando
- ✅ RTO real < meta definida
- ✅ RPO real < meta definida
- ✅ Zero perda de dados críticos em 12 meses

**Métricas:**

- Última restauração testada: [Data]
- Tempo de restore médio: [X minutos]
- Taxa de sucesso de backups: [99.x%]

---

## ✅ Checklist de Validação

### Setup Inicial

- [ ] Neon PITR habilitado (7 dias retenção)
- [ ] GitHub Actions workflow criado (backup-database.yml)
- [ ] S3 bucket criado (barber-analytics-backups)
- [ ] Lifecycle policy configurada (30 dias)
- [ ] Secrets configurados (DATABASE_URL, AWS keys)

### Operacional

- [ ] Backups diários rodando com sucesso
- [ ] Alertas configurados (falha de backup → Slack)
- [ ] Teste de restore realizado (mensal)
- [ ] Documentação atualizada (este documento)
- [ ] Equipe treinada em procedimentos DR

### Validação Trimestral

- [ ] Exercício de DR completo (simular disaster)
- [ ] Review de RTO/RPO (ajustar metas se necessário)
- [ ] Atualizar contatos de emergência
- [ ] Audit de backups (verificar integridade de 10 arquivos aleatórios)

---

## 📚 Referências

- [Neon Backup Documentation](https://neon.tech/docs/introduction/point-in-time-restore)
- [PostgreSQL Backup Best Practices](https://www.postgresql.org/docs/current/backup.html)
- [AWS S3 Lifecycle Policies](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lifecycle-mgmt.html)
- [Disaster Recovery Planning (AWS)](https://aws.amazon.com/disaster-recovery/)

---

**Última Atualização:** 15/11/2025
**Versão:** 1.0.0
**Responsável:** Equipe DevOps
**Revisão:** Trimestral
