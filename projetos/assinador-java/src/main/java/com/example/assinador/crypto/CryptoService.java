
package com.example.assinador.crypto;

import org.springframework.stereotype.Service;

import java.security.*;

@Service
public class CryptoService {

    public byte[] sign(
            byte[] data,
            String alias,
            String pin
    ) throws Exception{

        KeyPairGenerator kpg =
                KeyPairGenerator.getInstance("RSA");

        kpg.initialize(2048);

        KeyPair kp = kpg.generateKeyPair();

        Signature signature =
                Signature.getInstance("SHA256withRSA");

        signature.initSign(kp.getPrivate());

        signature.update(data);

        return signature.sign();

    }

}
