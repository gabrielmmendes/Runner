
package com.example.assinador.service;

import com.example.assinador.crypto.CryptoService;
import com.example.assinador.dto.SignRequest;
import com.example.assinador.dto.SignResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.Base64;

@Service
@RequiredArgsConstructor
public class SignatureService {

    private final CryptoService cryptoService;

    public SignResponse sign(SignRequest request) throws Exception{

        byte[] signature =
                cryptoService.sign(
                        request.getData().getBytes(),
                        request.getAlias(),
                        request.getPin()
                );

        return SignResponse.builder()
                .success(true)
                .signature(
                        Base64.getEncoder()
                                .encodeToString(signature)
                )
                .algorithm("SHA256withRSA")
                .build();

    }

}
