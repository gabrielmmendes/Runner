
package com.example.assinador.crypto;

import com.example.assinador.dto.SignRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.security.*;
import java.util.concurrent.ConcurrentHashMap;

@Slf4j
@Service
public class CryptoService {

    @Value("${assinador.pkcs11.library:}")
    private String pkcs11Library;

    private final ConcurrentHashMap<String, Provider> providerCache =
            new ConcurrentHashMap<>();

    private final KeyPair fallbackKeyPair = generateFallbackKeyPair();

    public byte[] sign(
            byte[] data,
            SignRequest.CryptoMaterial material
    ) throws Exception {

        PrivateKey privateKey = resolvePrivateKey(material);

        Signature signature =
                Signature.getInstance("SHA256withRSA");

        signature.initSign(privateKey);
        signature.update(data);

        return signature.sign();

    }

    private PrivateKey resolvePrivateKey(
            SignRequest.CryptoMaterial material
    ) throws Exception {

        if(pkcs11Library == null || pkcs11Library.isBlank()){
            log.debug("PKCS#11 desativado — usando keypair em memoria");
            return fallbackKeyPair.getPrivate();
        }

        if(material == null){
            throw new IllegalArgumentException(
                    "cryptoMaterial obrigatorio quando PKCS#11 ativo"
            );
        }
        if(material.getPin() == null || material.getPin().isBlank()){
            throw new IllegalArgumentException(
                    "cryptoMaterial.pin obrigatorio para PKCS#11"
            );
        }
        if(material.getIdentifier() == null || material.getIdentifier().isBlank()){
            throw new IllegalArgumentException(
                    "cryptoMaterial.identifier obrigatorio para PKCS#11"
            );
        }

        Provider provider = loadProvider(material);
        char[] pin = material.getPin().toCharArray();

        KeyStore ks = KeyStore.getInstance("PKCS11", provider);
        ks.load(null, pin);

        Key key = ks.getKey(material.getIdentifier(), pin);
        if(!(key instanceof PrivateKey)){
            throw new KeyStoreException(
                    "alias '" + material.getIdentifier() +
                            "' nao corresponde a chave privada"
            );
        }
        return (PrivateKey) key;

    }

    private Provider loadProvider(
            SignRequest.CryptoMaterial material
    ){

        String slotKey =
                material.getSlotId() != null
                        ? String.valueOf(material.getSlotId())
                        : (material.getTokenLabel() != null
                                ? material.getTokenLabel()
                                : "default");

        return providerCache.computeIfAbsent(slotKey, k -> {

            // Prefixo "--" sinaliza ao SunPKCS11 que e config inline (nao arquivo).
            StringBuilder cfg = new StringBuilder("--");
            cfg.append("name = assinador-").append(k).append('\n');
            cfg.append("library = ").append(pkcs11Library).append('\n');

            if(material.getSlotId() != null){
                cfg.append("slot = ").append(material.getSlotId()).append('\n');
            }
            if(material.getTokenLabel() != null && !material.getTokenLabel().isBlank()){
                cfg.append("tokenLabel = ").append(material.getTokenLabel()).append('\n');
            }

            Provider base = Security.getProvider("SunPKCS11");
            if(base == null){
                throw new IllegalStateException(
                        "SunPKCS11 indisponivel nesta JVM"
                );
            }
            return base.configure(cfg.toString());

        });

    }

    private static KeyPair generateFallbackKeyPair(){
        try{
            KeyPairGenerator kpg =
                    KeyPairGenerator.getInstance("RSA");
            kpg.initialize(2048);
            return kpg.generateKeyPair();
        }catch(NoSuchAlgorithmException e){
            throw new IllegalStateException(
                    "RSA indisponivel nesta JVM",
                    e
            );
        }
    }

}
