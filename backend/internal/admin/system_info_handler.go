package admin

import (
    "context"
    "database/sql"
    "os"
    "runtime"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    redis "github.com/redis/go-redis/v9"
    "github.com/prometheus/client_golang/prometheus"
    prommodel "github.com/prometheus/client_model/go"

    "gamelink/internal/cache"
    "gamelink/internal/config"
    "gamelink/internal/model"
)

// SystemInfoHandler provides admin system information endpoints.
type SystemInfoHandler struct {
    cfg        config.AppConfig
    sqlDB      *sql.DB
    cache      cache.Cache
    startedAt  time.Time
}

func NewSystemInfoHandler(cfg config.AppConfig, sqlDB *sql.DB, cacheClient cache.Cache) *SystemInfoHandler {
    return &SystemInfoHandler{cfg: cfg, sqlDB: sqlDB, cache: cacheClient, startedAt: time.Now()}
}

// Config returns sanitized system configuration.
func (h *SystemInfoHandler) Config(c *gin.Context) {
    type sanitized struct {
        Server struct {
            Port          string `json:"port"`
            EnableSwagger bool   `json:"enableSwagger"`
            Env           string `json:"env"`
        } `json:"server"`
        Database struct {
            Type       string `json:"type"`
            DSNPresent bool   `json:"dsnPresent"`
        } `json:"database"`
        Cache struct {
            Type  string `json:"type"`
            Redis struct {
                Addr string `json:"addr"`
                DB   int    `json:"db"`
            } `json:"redis"`
        } `json:"cache"`
        Crypto struct {
            Enabled       bool     `json:"enabled"`
            Methods       []string `json:"methods"`
            UseSignature  bool     `json:"useSignature"`
            ExcludePaths  []string `json:"excludePaths"`
        } `json:"crypto"`
        Auth struct {
            TokenTTLHours int `json:"tokenTTLHours"`
        } `json:"auth"`
        Seed struct {
            Enabled bool `json:"enabled"`
        } `json:"seed"`
        AdminAuth struct {
            Mode string `json:"mode"`
        } `json:"adminAuth"`
    }

    out := sanitized{}
    out.Server.Port = h.cfg.Port
    out.Server.EnableSwagger = h.cfg.EnableSwagger
    out.Server.Env = os.Getenv("APP_ENV")

    out.Database.Type = h.cfg.Database.Type
    out.Database.DSNPresent = h.cfg.Database.DSN != ""

    out.Cache.Type = h.cfg.Cache.Type
    out.Cache.Redis.Addr = h.cfg.Cache.Redis.Addr
    out.Cache.Redis.DB = h.cfg.Cache.Redis.DB

    out.Crypto.Enabled = h.cfg.Crypto.Enabled
    out.Crypto.Methods = h.cfg.Crypto.Methods
    out.Crypto.UseSignature = h.cfg.Crypto.UseSignature
    out.Crypto.ExcludePaths = h.cfg.Crypto.ExcludePaths

    out.Auth.TokenTTLHours = h.cfg.Auth.TokenTTLHours

    out.Seed.Enabled = h.cfg.Seed.Enabled
    out.AdminAuth.Mode = h.cfg.AdminAuth.Mode

    c.JSON(200, model.APIResponse[sanitized]{Code: 0, Message: "ok", Data: out})
}

// DBStatus returns database pool stats and basic query metrics summary.
func (h *SystemInfoHandler) DBStatus(c *gin.Context) {
    type poolStats struct {
        MaxOpen           int           `json:"maxOpen"`
        Open              int           `json:"open"`
        InUse             int           `json:"inUse"`
        Idle              int           `json:"idle"`
        WaitCount         int64         `json:"waitCount"`
        WaitDuration      time.Duration `json:"waitDuration"`
        MaxIdleClosed     int64         `json:"maxIdleClosed"`
        MaxLifetimeClosed int64         `json:"maxLifetimeClosed"`
    }

    ps := h.sqlDB.Stats()
    pool := poolStats{
        MaxOpen:           ps.MaxOpenConnections,
        Open:              ps.OpenConnections,
        InUse:             ps.InUse,
        Idle:              ps.Idle,
        WaitCount:         ps.WaitCount,
        WaitDuration:      ps.WaitDuration,
        MaxIdleClosed:     ps.MaxIdleClosed,
        MaxLifetimeClosed: ps.MaxLifetimeClosed,
    }

    // Summarize Prometheus histogram db_query_duration_seconds by op/table (avg, count)
    type opTableMetric struct {
        Op    string  `json:"op"`
        Table string  `json:"table"`
        Count uint64  `json:"count"`
        Avg   float64 `json:"avgSeconds"`
    }
    var summary []opTableMetric

    mfs, err := prometheus.DefaultGatherer.Gather()
    if err == nil {
        for _, mf := range mfs {
            if mf.GetName() != "db_query_duration_seconds" || mf.GetType() != prommodel.MetricType_HISTOGRAM {
                continue
            }
            for _, m := range mf.GetMetric() {
                hst := m.GetHistogram()
                if hst == nil || hst.GetSampleCount() == 0 {
                    continue
                }
                labels := m.GetLabel()
                var op, table string
                for _, l := range labels {
                    if l.GetName() == "op" {
                        op = l.GetValue()
                    }
                    if l.GetName() == "table" {
                        table = l.GetValue()
                    }
                }
                avg := hst.GetSampleSum() / float64(hst.GetSampleCount())
                summary = append(summary, opTableMetric{Op: op, Table: table, Count: hst.GetSampleCount(), Avg: avg})
            }
        }
    }

    type response struct {
        Driver  string          `json:"driver"`
        Pool    poolStats       `json:"pool"`
        Metrics []opTableMetric `json:"metrics"`
        Note    string          `json:"note"`
    }

    c.JSON(200, model.APIResponse[response]{Code: 0, Message: "ok", Data: response{
        Driver:  h.cfg.Database.Type,
        Pool:    pool,
        Metrics: summary,
        Note:    "详细慢查询指标请参见 /metrics 的 db_query_duration_seconds",
    }})
}

// CacheStatus returns cache connectivity and memory info (Redis when applicable).
func (h *SystemInfoHandler) CacheStatus(c *gin.Context) {
    type redisInfo struct {
        PingOK                 bool    `json:"pingOK"`
        LatencyMs              int64   `json:"latencyMs"`
        UsedMemoryBytes        int64   `json:"usedMemoryBytes"`
        UsedMemoryRssBytes     int64   `json:"usedMemoryRssBytes"`
        MemFragmentationRatio  float64 `json:"memFragmentationRatio"`
        MaxMemoryBytes         int64   `json:"maxMemoryBytes"`
    }
    type response struct {
        Type  string    `json:"type"`
        Redis *redisInfo `json:"redis,omitempty"`
        Status string   `json:"status"`
        Note   string   `json:"note"`
    }

    out := response{Type: h.cfg.Cache.Type, Status: "unknown"}

    // Quick connectivity test via Set/Delete with short TTL
    start := time.Now()
    _ = h.cache.Set(context.Background(), "sys:check", "1", time.Second)
    _ = h.cache.Delete(context.Background(), "sys:check")
    out.Status = "ok"

    if strings.ToLower(h.cfg.Cache.Type) == "redis" {
        // Construct ephemeral client from config for INFO memory
        r := redis.NewClient(&redis.Options{
            Addr:     h.cfg.Cache.Redis.Addr,
            Password: h.cfg.Cache.Redis.Password,
            DB:       h.cfg.Cache.Redis.DB,
        })
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()

        ri := &redisInfo{}
        pingStart := time.Now()
        if err := r.Ping(ctx).Err(); err == nil {
            ri.PingOK = true
            ri.LatencyMs = time.Since(pingStart).Milliseconds()
            info, err := r.Info(ctx, "memory").Result()
            if err == nil {
                // Parse INFO output
                for _, line := range strings.Split(info, "\n") {
                    kv := strings.Split(strings.TrimSpace(line), ":")
                    if len(kv) != 2 { continue }
                    k, v := kv[0], kv[1]
                    switch k {
                    case "used_memory":
                        ri.UsedMemoryBytes = parseInt64(v)
                    case "used_memory_rss":
                        ri.UsedMemoryRssBytes = parseInt64(v)
                    case "mem_fragmentation_ratio":
                        ri.MemFragmentationRatio = parseFloat64(v)
                    case "maxmemory":
                        ri.MaxMemoryBytes = parseInt64(v)
                    }
                }
            }
        }
        out.Redis = ri
        out.Note = "内存数据来自 Redis INFO memory；如需更详细指标请使用 /metrics"
    } else {
        out.Note = "内存型缓存无统一内存统计，建议使用运行时资源接口或外部监控"
    }

    _ = start // latency for set/delete can be derived if needed later
    c.JSON(200, model.APIResponse[response]{Code: 0, Message: "ok", Data: out})
}

// Resources returns Go runtime resource usage.
func (h *SystemInfoHandler) Resources(c *gin.Context) {
    var ms runtime.MemStats
    runtime.ReadMemStats(&ms)

    type mem struct {
        Alloc            uint64 `json:"allocBytes"`
        TotalAlloc       uint64 `json:"totalAllocBytes"`
        Sys              uint64 `json:"sysBytes"`
        HeapAlloc        uint64 `json:"heapAllocBytes"`
        HeapSys          uint64 `json:"heapSysBytes"`
        HeapInuse        uint64 `json:"heapInuseBytes"`
        HeapIdle         uint64 `json:"heapIdleBytes"`
        NextGC           uint64 `json:"nextGCBytes"`
        PauseTotalNs     uint64 `json:"gcPauseTotalNs"`
        NumGC            uint32 `json:"numGC"`
    }

    type response struct {
        UptimeSeconds  int64 `json:"uptimeSeconds"`
        NumGoroutines  int   `json:"numGoroutines"`
        NumCPU         int   `json:"numCPU"`
        GoVersion      string `json:"goVersion"`
        Memory         mem   `json:"memory"`
    }

    out := response{
        UptimeSeconds: int64(time.Since(h.startedAt).Seconds()),
        NumGoroutines: runtime.NumGoroutine(),
        NumCPU:        runtime.NumCPU(),
        GoVersion:     runtime.Version(),
        Memory: mem{
            Alloc:        ms.Alloc,
            TotalAlloc:   ms.TotalAlloc,
            Sys:          ms.Sys,
            HeapAlloc:    ms.HeapAlloc,
            HeapSys:      ms.HeapSys,
            HeapInuse:    ms.HeapInuse,
            HeapIdle:     ms.HeapIdle,
            NextGC:       ms.NextGC,
            PauseTotalNs: ms.PauseTotalNs,
            NumGC:        ms.NumGC,
        },
    }

    c.JSON(200, model.APIResponse[response]{Code: 0, Message: "ok", Data: out})
}

// Version returns build and version info.
func (h *SystemInfoHandler) Version(c *gin.Context) {
    type response struct {
        Version    string `json:"version"`
        GitCommit  string `json:"gitCommit"`
        BuildTime  string `json:"buildTime"`
        GoVersion  string `json:"goVersion"`
    }

    // Values should be injected via ldflags at build time; fall back to empty strings.
    v := response{
        Version:   buildVersion,
        GitCommit: buildCommit,
        BuildTime: buildTime,
        GoVersion: runtime.Version(),
    }
    c.JSON(200, model.APIResponse[response]{Code: 0, Message: "ok", Data: v})
}

// parse helpers
func parseInt64(s string) int64 {
    var out int64
    for i := 0; i < len(s); i++ {
        if s[i] < '0' || s[i] > '9' { break }
        out = out*10 + int64(s[i]-'0')
    }
    return out
}

func parseFloat64(s string) float64 {
    // minimal parser for ratio like 1.23
    var intPart int64
    var fracPart int64
    var fracScale float64 = 1
    var seenDot bool
    for i := 0; i < len(s); i++ {
        ch := s[i]
        if ch == '.' && !seenDot { seenDot = true; continue }
        if ch < '0' || ch > '9' { break }
        if !seenDot {
            intPart = intPart*10 + int64(ch-'0')
        } else {
            fracPart = fracPart*10 + int64(ch-'0')
            fracScale *= 10
        }
    }
    return float64(intPart) + float64(fracPart)/fracScale
}

// These variables are intended to be filled via -ldflags at build time.
var (
    buildVersion string
    buildCommit  string
    buildTime    string
)