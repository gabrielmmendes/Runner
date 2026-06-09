package com.example.assinador.metrics;

import org.springframework.stereotype.Component;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.atomic.DoubleAdder;

/**
 * In-memory metrics collector exposed in Prometheus text format.
 * Tracks request counts per endpoint, error counts, request-duration
 * histograms (for p50/p95/p99 via histogram_quantile) and uptime.
 */
@Component
public class MetricsRegistry {

    /** Histogram bucket upper bounds in seconds. */
    static final double[] BUCKETS = {
            0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
    };

    private final long startMillis = System.currentTimeMillis();
    private final Map<String, AtomicLong> requests = new ConcurrentHashMap<>();
    private final Map<String, AtomicLong> errors = new ConcurrentHashMap<>();
    private final Map<String, Histogram> latencies = new ConcurrentHashMap<>();

    public void record(String method, String path, int status, long durationNanos) {
        String key = method + "|" + path + "|" + status;
        requests.computeIfAbsent(key, k -> new AtomicLong()).incrementAndGet();
        if (status >= 400) {
            errors.computeIfAbsent(path, k -> new AtomicLong()).incrementAndGet();
        }
        latencies.computeIfAbsent(path, k -> new Histogram())
                .observe(durationNanos / 1_000_000_000.0);
    }

    public double uptimeSeconds() {
        return (System.currentTimeMillis() - startMillis) / 1000.0;
    }

    public String render() {
        StringBuilder sb = new StringBuilder(1024);

        sb.append("# HELP assinador_uptime_seconds Tempo desde o start do servidor.\n");
        sb.append("# TYPE assinador_uptime_seconds gauge\n");
        sb.append("assinador_uptime_seconds ").append(fmt(uptimeSeconds())).append('\n');

        sb.append("# HELP assinador_requests_total Total de requisições por endpoint e status.\n");
        sb.append("# TYPE assinador_requests_total counter\n");
        for (Map.Entry<String, AtomicLong> e : requests.entrySet()) {
            String[] p = e.getKey().split("\\|", 3);
            sb.append("assinador_requests_total{method=\"").append(p[0])
              .append("\",path=\"").append(p[1])
              .append("\",status=\"").append(p[2]).append("\"} ")
              .append(e.getValue().get()).append('\n');
        }

        sb.append("# HELP assinador_request_errors_total Total de respostas com status >= 400.\n");
        sb.append("# TYPE assinador_request_errors_total counter\n");
        for (Map.Entry<String, AtomicLong> e : errors.entrySet()) {
            sb.append("assinador_request_errors_total{path=\"").append(e.getKey())
              .append("\"} ").append(e.getValue().get()).append('\n');
        }

        sb.append("# HELP assinador_request_duration_seconds Latência das requisições.\n");
        sb.append("# TYPE assinador_request_duration_seconds histogram\n");
        for (Map.Entry<String, Histogram> e : latencies.entrySet()) {
            e.getValue().render(sb, e.getKey());
        }
        return sb.toString();
    }

    private static String fmt(double v) {
        return String.valueOf(v);
    }

    /** Fixed-bucket cumulative histogram. */
    static final class Histogram {
        private final AtomicLong[] buckets = new AtomicLong[BUCKETS.length];
        private final AtomicLong count = new AtomicLong();
        private final DoubleAdder sum = new DoubleAdder();

        Histogram() {
            for (int i = 0; i < buckets.length; i++) {
                buckets[i] = new AtomicLong();
            }
        }

        void observe(double seconds) {
            count.incrementAndGet();
            sum.add(seconds);
            for (int i = 0; i < BUCKETS.length; i++) {
                if (seconds <= BUCKETS[i]) {
                    buckets[i].incrementAndGet();
                }
            }
        }

        void render(StringBuilder sb, String path) {
            for (int i = 0; i < BUCKETS.length; i++) {
                sb.append("assinador_request_duration_seconds_bucket{path=\"").append(path)
                  .append("\",le=\"").append(fmt(BUCKETS[i])).append("\"} ")
                  .append(buckets[i].get()).append('\n');
            }
            sb.append("assinador_request_duration_seconds_bucket{path=\"").append(path)
              .append("\",le=\"+Inf\"} ").append(count.get()).append('\n');
            sb.append("assinador_request_duration_seconds_sum{path=\"").append(path)
              .append("\"} ").append(fmt(sum.sum())).append('\n');
            sb.append("assinador_request_duration_seconds_count{path=\"").append(path)
              .append("\"} ").append(count.get()).append('\n');
        }
    }
}
