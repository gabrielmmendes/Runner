
package com.example.assinador;

import com.example.assinador.dto.SignRequest;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.*;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest(webEnvironment =
SpringBootTest.WebEnvironment.RANDOM_PORT)
class SignatureControllerTest {

    @Autowired
    TestRestTemplate rest;

    @Test
    void shouldSign(){

        SignRequest request =
                new SignRequest();

        request.setData("teste");

        ResponseEntity<String> response =
                rest.postForEntity(
                        "/api/sign",
                        request,
                        String.class
                );

        assertEquals(
                HttpStatus.OK,
                response.getStatusCode()
        );

    }

}
