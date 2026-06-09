package com.example.assinador.metrics;

import lombok.RequiredArgsConstructor;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * Exposes {@code GET /metrics} in Prometheus text exposition format.
 * Disabled when {@code assinador.metrics.enabled=false}.
 */
@RestController
@RequiredArgsConstructor
@ConditionalOnProperty(name = "assinador.metrics.enabled", havingValue = "true", matchIfMissing = true)
public class MetricsController {

    /** Prometheus text exposition format version. */
    private static final String CONTENT_TYPE = "text/plain; version=0.0.4; charset=utf-8";

    private final MetricsRegistry registry;

    @GetMapping("/metrics")
    public ResponseEntity<String> metrics() {
        return ResponseEntity.ok()
                .contentType(MediaType.parseMediaType(CONTENT_TYPE))
                .body(registry.render());
    }
}
