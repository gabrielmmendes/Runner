
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

        byte[] payload = resolvePayload(request);
        SignRequest.CryptoMaterial material = resolveMaterial(request);

        byte[] signature =
                cryptoService.sign(payload, material);

        return SignResponse.builder()
                .success(true)
                .signature(
                        Base64.getEncoder()
                                .encodeToString(signature)
                )
                .algorithm("SHA256withRSA")
                .build();

    }

    private byte[] resolvePayload(SignRequest r){
        if(r.getData() != null){
            return r.getData().getBytes();
        }
        if(r.getBundle() != null){
            return r.getBundle().toString().getBytes();
        }
        return new byte[0];
    }

    private SignRequest.CryptoMaterial resolveMaterial(SignRequest r){

        if(r.getCryptoMaterial() != null){
            return r.getCryptoMaterial();
        }

        SignRequest.CryptoMaterial fallback =
                new SignRequest.CryptoMaterial();
        fallback.setIdentifier(
                r.getAlias() != null ? r.getAlias() : "default"
        );
        fallback.setPin(r.getPin());
        return fallback;

    }

}
