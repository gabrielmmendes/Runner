
package com.example.assinador.dto;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class ValidateResponse {

    private boolean valid;

    private String message;

}
