
package com.example.assinador;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableScheduling;

@SpringBootApplication
@EnableScheduling
public class AssinadorApplication {

    public static void main(String[] args) {
        SpringApplication.run(AssinadorApplication.class, args);
    }

}
