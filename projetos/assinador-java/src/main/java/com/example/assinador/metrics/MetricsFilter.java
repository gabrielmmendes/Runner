package com.example.assinador.metrics;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;

/**
 * Records request count, status and latency for every HTTP request,
 * excluding the {@code /metrics} endpoint itself. Active only when
 * {@code assinador.metrics.enabled=true} (default).
 */
@Component
@RequiredArgsConstructor
@ConditionalOnProperty(name = "assinador.metrics.enabled", havingValue = "true", matchIfMissing = true)
public class MetricsFilter extends OncePerRequestFilter {

    private final MetricsRegistry registry;

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response,
                                    FilterChain chain) throws ServletException, IOException {
        String path = request.getRequestURI();
        if ("/metrics".equals(path)) {
            chain.doFilter(request, response);
            return;
        }
        long start = System.nanoTime();
        try {
            chain.doFilter(request, response);
        } finally {
            registry.record(request.getMethod(), path,
                    response.getStatus(), System.nanoTime() - start);
        }
    }
}
