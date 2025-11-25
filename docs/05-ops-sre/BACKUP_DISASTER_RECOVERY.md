# 🔄 Backup & Disaster Recovery — Barber Analytics Pro

**Versão:** 1.0
**Data:** 24/11/2025
**Status:** ✅ Implementado

---

## 📋 Visão Geral

Este documento descreve os procedimentos de **backup, restore e disaster recovery** para o banco de dados PostgreSQL (Neon) do Barber Analytics Pro.

**Estratégia de Backup:**

- **Automated Daily Backups:** pg_dump via GitHub Actions (03:00 UTC)
- **Neon PITR (Point-in-Time Recovery):** Retenção de 7-30 dias
- **Storage:** AWS S3 com criptografia AES-256
- **Retenção:** 30 dias para backups automatizados
- **RPO (Recovery Point Objective):** 24 horas (backups diários)
- **RTO (Recovery Time Objective):** 2 horas

---

## 🛡️ Tipos de Backup

### 1. Backup Automatizado Diário (GitHub Actions)

**Frequência:** Diário às 03:00 UTC (00:00 BRT)
**Método:** `pg_dump` compactado com gzip
**Destino:** AWS S3 bucket `s3://barber-analytics-backups/database-backups/`
**Retenção:** 30 dias (cleanup automático)
**Formato:** SQL plain text compactado

**Arquivo Workflow:** `.github/workflows/backup-database.yml`

**Execução Manual:**

```bash
# Via GitHub Actions UI
# https://github.com/<org>/barber-analytics-proV2/actions/workflows/backup-database.yml
# Clicar em "Run workflow"
```

### 2. Neon PITR (Point-in-Time Recovery)

**Frequência:** Contínuo (WAL logs)
**Retenção:** 7 dias (Free Tier) ou 30 dias (Pro Plan)
**Método:** Neon native PITR via console
**Granularidade:** Qualquer ponto no tempo dentro da janela de retenção

**Como acessar:**

```
1. Acessar: https://console.neon.tech
2. Selecionar projeto: barber-analytics-prod
3. Aba "Backups" → "Point-in-Time Recovery"
4. Selecionar timestamp desejado
5. Clicar em "Restore to new branch"
```

### 3. Backup Manual (Ad-hoc)

**Quando usar:**

- Antes de migrations críticas
- Antes de mudanças estruturais no schema
- Testes de restore
- Compliance/auditoria

**Como executar:**

```bash
# Backup manual via pg_dump
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
pg_dump \
  --host=<NEON_HOST> \
  --port=5432 \
  --username=<USER> \
  --dbname=<DB_NAME> \
  --format=plain \
  --no-owner \
  --no-acl \
  --clean \
  --if-exists \
  | gzip > backup_manual_${TIMESTAMP}.sql.gz

# Upload para S3
aws s3 cp backup_manual_${TIMESTAMP}.sql.gz \
  s3://barber-analytics-backups/manual-backups/ \
  --storage-class STANDARD_IA \
  --server-side-encryption AES256
```

---

## 🔧 Procedimentos de Restore

### Restore 1: Banco Completo (Disaster Recovery Total)

**Cenário:** Perda total do banco de dados, corrupção completa, desastre.

**Passos:**

1. **Baixar último backup do S3:**

```bash
# Listar backups disponíveis
aws s3 ls s3://barber-analytics-backups/database-backups/ --recursive | sort

# Baixar backup mais recente
LATEST_BACKUP=$(aws s3 ls s3://barber-analytics-backups/database-backups/ --recursive | sort | tail -1 | awk '{print $4}')
aws s3 cp s3://barber-analytics-backups/$LATEST_BACKUP /tmp/restore.sql.gz
```

2. **Descompactar backup:**

```bash
gunzip /tmp/restore.sql.gz
```

3. **Criar novo branch no Neon (recomendado) ou usar staging:**

```bash
# Via Neon Console:
# 1. Criar new branch: "restore-<timestamp>"
# 2. Copiar connection string do novo branch
```

4. **Executar restore:**

```bash
# Restaurar em novo branch/database
psql \
  --host=<NEON_BRANCH_HOST> \
  --port=5432 \
  --username=<USER> \
  --dbname=<DB_NAME> \
  < /tmp/restore.sql

# Verificar restore
psql -h <HOST> -U <USER> -d <DB> -c "SELECT COUNT(*) FROM tenants;"
psql -h <HOST> -U <USER> -d <DB> -c "SELECT COUNT(*) FROM users;"
```

5. **Validar dados restaurados:**

```bash
# Conferir total de tabelas
psql -h <HOST> -U <USER> -d <DB> -c "\dt" | wc -l

# Conferir integridade de dados críticos
psql -h <HOST> -U <USER> -d <DB> <<EOF
SELECT
  'tenants' as table, COUNT(*) as count FROM tenants
UNION ALL
SELECT 'users', COUNT(*) FROM users
UNION ALL
SELECT 'receitas', COUNT(*) FROM receitas
UNION ALL
SELECT 'despesas', COUNT(*) FROM despesas;
EOF
```

6. **Promover branch para produção (se validado):**

```bash
# Via Neon Console:
# 1. Aba "Branches"
# 2. Selecionar branch "restore-<timestamp>"
# 3. "Set as primary" (se tudo estiver OK)
```

**Tempo Estimado:** 1-2 horas (dependendo do tamanho do banco)

---

### Restore 2: Point-in-Time Recovery (Corrupção Parcial)

**Cenário:** Dados deletados acidentalmente, migration com erro, corrupção parcial.

**Passos:**

1. **Identificar timestamp exato antes do problema:**

```bash
# Exemplo: Erro ocorreu em 2025-11-24 14:30:00 UTC
# Restaurar para: 2025-11-24 14:25:00 UTC (5 min antes)
```

2. **Via Neon Console:**

```
1. Acessar: https://console.neon.tech
2. Projeto: barber-analytics-prod
3. Aba "Backups" → "Point-in-Time Recovery"
4. Selecionar timestamp: 2025-11-24 14:25:00 UTC
5. Clicar "Restore to new branch"
6. Nome do branch: "pitr-recovery-20251124-1425"
```

3. **Conectar ao branch restaurado:**

```bash
# Neon fornece nova connection string
psql <NEW_BRANCH_CONNECTION_STRING>
```

4. **Validar dados restaurados:**

```sql
-- Conferir se dados estão no estado correto
SELECT * FROM tenants WHERE deleted_at IS NULL;
SELECT COUNT(*) FROM users WHERE ativo = true;

-- Verificar última transação registrada
SELECT MAX(created_at) FROM audit_logs;
```

5. **Extrair dados específicos (se necessário):**

```bash
# Se apenas uma tabela foi corrompida, extrair apenas ela
pg_dump \
  --host=<BRANCH_HOST> \
  --dbname=<DB> \
  --table=users \
  --data-only \
  > /tmp/users_recovery.sql

# Aplicar no banco principal (com cuidado!)
psql <MAIN_DB_CONNECTION> < /tmp/users_recovery.sql
```

**Tempo Estimado:** 30 minutos a 1 hora

---

### Restore 3: Tabela/Dados Específicos (Selective Restore)

**Cenário:** Apenas uma tabela ou registros específicos precisam ser restaurados.

**Passos:**

1. **Baixar backup completo:**

```bash
aws s3 cp s3://barber-analytics-backups/database-backups/<BACKUP_FILE> /tmp/backup.sql.gz
gunzip /tmp/backup.sql.gz
```

2. **Extrair apenas a tabela desejada:**

```bash
# Exemplo: Restaurar apenas tabela 'receitas'
grep -A 1000 "CREATE TABLE receitas" /tmp/backup.sql > /tmp/receitas_only.sql

# Ou usar pg_restore se backup estiver em formato custom
pg_restore \
  --table=receitas \
  --data-only \
  /tmp/backup.dump \
  > /tmp/receitas_only.sql
```

3. **Restaurar em banco temporário para validação:**

```bash
# Criar database temporário
createdb -h <HOST> -U <USER> temp_restore_db

# Restaurar tabela
psql -h <HOST> -U <USER> -d temp_restore_db < /tmp/receitas_only.sql
```

4. **Copiar dados específicos para produção:**

```sql
-- Conectar ao banco temporário
\c temp_restore_db

-- Exportar registros específicos
COPY (
  SELECT * FROM receitas
  WHERE created_at BETWEEN '2025-11-01' AND '2025-11-30'
) TO '/tmp/receitas_nov2025.csv' WITH CSV HEADER;

-- Conectar ao banco produção
\c production_db

-- Importar registros
COPY receitas FROM '/tmp/receitas_nov2025.csv' WITH CSV HEADER;
```

**Tempo Estimado:** 15-30 minutos

---

## 🧪 Testes de Restore (Obrigatório Trimestral)

**Frequência:** A cada 3 meses
**Objetivo:** Validar que backups estão funcionais e procedimentos documentados

### Checklist de Teste

```markdown
- [ ] Baixar backup mais recente do S3
- [ ] Criar branch de teste no Neon
- [ ] Executar restore completo
- [ ] Validar integridade de dados:
  - [ ] Total de tabelas (42 tabelas esperadas)
  - [ ] Total de tenants
  - [ ] Total de users
  - [ ] Constraints e indexes intactos
- [ ] Cronometrar tempo de restore
- [ ] Documentar resultados no runbook
- [ ] Deletar branch de teste
```

**Última execução:** [Preencher após primeiro teste]
**Próxima execução:** [Preencher data]

---

## 📊 Monitoramento de Backups

### Métricas Prometheus

```yaml
# /etc/prometheus/prometheus.yml
scrape_configs:
  - job_name: 'backup_monitoring'
    static_configs:
      - targets: ['localhost:9090']
    metrics_path: /metrics
```

**Métricas a monitorar:**

- `backup_last_success_timestamp` — Timestamp do último backup bem-sucedido
- `backup_duration_seconds` — Duração do backup
- `backup_file_size_bytes` — Tamanho do arquivo de backup
- `backup_failures_total` — Total de falhas de backup

### Alertas (Prometheus Alertmanager)

```yaml
# /etc/prometheus/alert_rules.yml
groups:
  - name: backup_alerts
    interval: 5m
    rules:
      - alert: BackupFailed
        expr: time() - backup_last_success_timestamp > 86400
        for: 1h
        labels:
          severity: critical
        annotations:
          summary: 'Backup não executado nas últimas 24h'
          description: 'Último backup bem-sucedido: {{ $value }} segundos atrás'

      - alert: BackupTooSlow
        expr: backup_duration_seconds > 1800
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: 'Backup demorando mais que 30 minutos'
          description: 'Duração atual: {{ $value }} segundos'

      - alert: S3StorageFull
        expr: aws_s3_bucket_size_bytes{bucket="barber-analytics-backups"} > 100000000000
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: 'Bucket S3 de backups está cheio (>100GB)'
          description: 'Tamanho atual: {{ $value }} bytes'
```

---

## 🔐 Segurança dos Backups

### Criptografia

- ✅ **Em trânsito:** TLS 1.3 (pg_dump → S3)
- ✅ **Em repouso:** AES-256 (S3 Server-Side Encryption)
- ✅ **Neon PITR:** Criptografia nativa do Neon

### Controle de Acesso

```
AWS S3 Bucket Policy:
- Apenas GitHub Actions service account tem write
- Apenas DevOps team tem read
- MFA obrigatório para delete
- Versioning habilitado
- Object Lock para backups críticos (opcional)
```

### Auditoria

- ✅ AWS CloudTrail logs habilitados
- ✅ GitHub Actions workflow logs retidos por 90 dias
- ✅ Alertas de acesso suspeito via AWS GuardDuty

---

## 🚨 Runbook de Disaster Recovery

### Situação 1: Banco de dados inacessível (Neon down)

**Ações:**

1. Verificar status do Neon: https://neon.tech/status
2. Se outage confirmado → aguardar restauração automática (Neon SLA: 99.9%)
3. Se outage > 1 hora → executar Restore 1 em branch alternativo
4. Atualizar connection string no backend (.env)

### Situação 2: Dados deletados acidentalmente

**Ações:**

1. Identificar timestamp do erro (via audit_logs)
2. Executar Restore 2 (PITR) para 5 minutos antes do erro
3. Validar dados restaurados
4. Se OK → promover branch para produção

### Situação 3: Migration com erro crítico

**Ações:**

1. **NÃO** executar rollback manual
2. Executar Restore 2 (PITR) para antes da migration
3. Corrigir script de migration
4. Re-testar em branch de teste
5. Re-executar migration corrigida

### Situação 4: Corrupção de índices

**Ações:**

```sql
-- Recriar índices corrompidos
REINDEX DATABASE neondb;

-- Ou recriar índice específico
REINDEX INDEX idx_users_tenant_id;
```

---

## 📚 Referências

- [Neon PITR Documentation](https://neon.tech/docs/manage/backups)
- [PostgreSQL pg_dump Documentation](https://www.postgresql.org/docs/current/app-pgdump.html)
- [AWS S3 Encryption](https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingEncryption.html)
- [GitHub Actions Workflows](https://docs.github.com/en/actions)

---

## 🔄 Changelog

| Versão | Data       | Alteração      |
| ------ | ---------- | -------------- |
| 1.0    | 24/11/2025 | Versão inicial |

**Última Atualização:** 24/11/2025
**Responsável:** DevOps Team
**Revisão:** Trimestral
