
package com.example.assinador.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.JsonNode;
import lombok.Data;

import java.util.List;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class SignRequest {

    private JsonNode bundle;

    private JsonNode provenance;

    private List<String> certChain;

    private Long timestamp;

    private String strategy;

    private String policyId;

    private CryptoMaterial cryptoMaterial;

    private JsonNode operationalConfig;

    private String data;

    private String alias;

    private String pin;

    @Data
    @JsonIgnoreProperties(ignoreUnknown = true)
    public static class CryptoMaterial {

        private String type;

        private String pin;

        private String identifier;

        private Integer slotId;

        private String tokenLabel;

    }

}
