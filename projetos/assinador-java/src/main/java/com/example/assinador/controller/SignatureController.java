
package com.example.assinador.controller;

import com.example.assinador.dto.*;
import com.example.assinador.service.SignatureService;
import com.example.assinador.service.ValidationService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
@RequiredArgsConstructor
public class SignatureController {

    private final SignatureService signatureService;
    private final ValidationService validationService;

    @PostMapping("/sign")
    public ResponseEntity<?> sign(@RequestBody SignRequest request){

        try{

            return ResponseEntity.ok(
                    signatureService.sign(request)
            );

        }catch(Exception e){

            return ResponseEntity.badRequest()
                    .body(ErrorResponse.of(e.getMessage()));

        }

    }

    @PostMapping("/validate")
    public ResponseEntity<?> validate(@RequestBody ValidateRequest request){

        try{

            return ResponseEntity.ok(
                    validationService.validate(request)
            );

        }catch(Exception e){

            return ResponseEntity.badRequest()
                    .body(ErrorResponse.of(e.getMessage()));

        }

    }

}
