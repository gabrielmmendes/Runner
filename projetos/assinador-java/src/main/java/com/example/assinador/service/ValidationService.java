
package com.example.assinador.service;

import com.example.assinador.dto.ValidateRequest;
import com.example.assinador.dto.ValidateResponse;
import org.springframework.stereotype.Service;

@Service
public class ValidationService {

    public ValidateResponse validate(
            ValidateRequest request
    ){

        return ValidateResponse.builder()
                .valid(true)
                .message("Validação simulada")
                .build();

    }

}
