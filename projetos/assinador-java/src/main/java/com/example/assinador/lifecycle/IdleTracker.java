
package com.example.assinador.lifecycle;

import jakarta.servlet.*;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.time.Instant;
import java.util.concurrent.atomic.AtomicReference;

@Component
public class IdleTracker implements Filter {

    private final AtomicReference<Instant> lastRequest =
            new AtomicReference<>(Instant.now());

    @Override
    public void doFilter(
            ServletRequest request,
            ServletResponse response,
            FilterChain chain
    ) throws IOException, ServletException {

        String uri = ((HttpServletRequest) request).getRequestURI();

        if(uri != null && uri.startsWith("/api/")){
            lastRequest.set(Instant.now());
        }

        chain.doFilter(request, response);

    }

    public Instant getLastRequest(){
        return lastRequest.get();
    }

}
