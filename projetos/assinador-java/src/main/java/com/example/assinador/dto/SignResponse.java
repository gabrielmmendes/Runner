
package com.example.assinador.dto;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class SignResponse {

    private boolean success;

    private String signature;

    private String algorithm;

}
