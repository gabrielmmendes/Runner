
package com.example.assinador.lifecycle;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.context.ApplicationContext;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.time.Duration;
import java.time.Instant;

@Slf4j
@Component
@RequiredArgsConstructor
public class IdleShutdownScheduler {

    private final IdleTracker idleTracker;
    private final ApplicationContext context;

    @Value("${assinador.idle-timeout-min:0}")
    private int idleTimeoutMin;

    @Scheduled(fixedDelay = 30_000L)
    public void check(){

        if(idleTimeoutMin <= 0) return;

        Duration idle =
                Duration.between(
                        idleTracker.getLastRequest(),
                        Instant.now()
                );

        if(idle.toMinutes() >= idleTimeoutMin){

            log.info(
                    "auto-stop: inativo ha {} min (limite={} min)",
                    idle.toMinutes(),
                    idleTimeoutMin
            );

            System.exit(
                    SpringApplication.exit(context, () -> 0)
            );

        }

    }

}
