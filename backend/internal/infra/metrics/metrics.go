package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// LGPD Endpoints Metrics
	LGPDExportRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lgpd_export_requests_total",
			Help: "Total de solicitações de exportação de dados (LGPD Art. 18, V)",
		},
		[]string{"tenant_id", "status"},
	)

	LGPDExportDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "lgpd_export_duration_seconds",
			Help:    "Duração das exportações de dados em segundos",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"tenant_id"},
	)

	LGPDDeleteAccountTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lgpd_delete_account_total",
			Help: "Total de solicitações de exclusão de conta (LGPD Art. 18, VI)",
		},
		[]string{"tenant_id", "status"},
	)

	LGPDPreferencesUpdatesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lgpd_preferences_updates_total",
			Help: "Total de atualizações de preferências de privacidade",
		},
		[]string{"tenant_id", "consent_type"},
	)

	// Backup Metrics
	BackupLastSuccessTimestamp = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "backup_last_success_timestamp",
			Help: "Timestamp Unix do último backup bem-sucedido",
		},
	)

	BackupDurationSeconds = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "backup_duration_seconds",
			Help:    "Duração do backup em segundos",
			Buckets: []float64{60, 120, 300, 600, 1200, 1800}, // 1min a 30min
		},
	)

	BackupFileSizeBytes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "backup_file_size_bytes",
			Help: "Tamanho do arquivo de backup em bytes",
		},
	)

	BackupFailuresTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "backup_failures_total",
			Help: "Total de falhas de backup",
		},
	)

	// API Metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requisições HTTP",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP em segundos",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "path"},
	)

	// Database Metrics
	DBQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total de queries executadas no banco de dados",
		},
		[]string{"table", "operation"},
	)

	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duração das queries em segundos",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"table", "operation"},
	)

	// Active Users
	ActiveUsersTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_users_total",
			Help: "Número de usuários ativos por tenant",
		},
		[]string{"tenant_id"},
	)
)
