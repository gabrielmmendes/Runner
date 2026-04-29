
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
        String alias = resolveAlias(request);
        String pin = resolvePin(request);

        byte[] signature =
                cryptoService.sign(payload, alias, pin);

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

    private String resolveAlias(SignRequest r){
        if(r.getAlias() != null) return r.getAlias();
        if(r.getCryptoMaterial() != null && r.getCryptoMaterial().getIdentifier() != null){
            return r.getCryptoMaterial().getIdentifier();
        }
        return "default";
    }

    private String resolvePin(SignRequest r){
        if(r.getPin() != null) return r.getPin();
        if(r.getCryptoMaterial() != null) return r.getCryptoMaterial().getPin();
        return null;
    }

}
