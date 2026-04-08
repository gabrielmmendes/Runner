
package com.example.assinador.dto;

import lombok.Builder;
import lombok.Data;

import java.time.LocalDateTime;

@Data
@Builder
public class ErrorResponse {

    private boolean success;

    private String error;

    private LocalDateTime timestamp;

    public static ErrorResponse of(String message){

        return ErrorResponse.builder()
                .success(false)
                .error(message)
                .timestamp(LocalDateTime.now())
                .build();

    }

}
