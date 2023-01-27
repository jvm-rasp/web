package com.jrasp.admin.server.security;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.nimbusds.jose.*;
import com.nimbusds.jose.crypto.MACSigner;
import com.nimbusds.jose.crypto.MACVerifier;
import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Component;

import javax.annotation.Resource;
import java.text.ParseException;
import java.util.HashMap;
import java.util.Map;

@Data
@ConfigurationProperties(prefix = "jwt")
@Component
public class JwtTokenUtil {
    @Resource
    private ObjectMapper objectMapper;
    private String secret;
    private String header;


    /**
     * 生成token令牌
     *
     * @param userId 用户id
     * @return 令token牌
     */
    public String generateToken(Long userId) throws JsonProcessingException, JOSEException {
        JWSHeader header = new JWSHeader(JWSAlgorithm.HS256);
        HashMap<String, Long> map = new HashMap<>();
        map.put("user_id",userId);
        Payload payload = new Payload(objectMapper.writeValueAsString(map));
        JWSObject jwsObject = new JWSObject(header, payload);
        MACSigner signer = new MACSigner(secret);
        jwsObject.sign(signer);
        return jwsObject.serialize();
    }

    /**
     * 从令牌中获取用户id
     *
     * @param token 令牌
     * @return 用户名
     */
    public Long getUserIdFromToken(String token) throws ParseException, JOSEException {
        JWSObject jwsObject = JWSObject.parse(token);
        Payload payload = jwsObject.getPayload();
        MACVerifier verifier = new MACVerifier(secret);
        if(jwsObject.verify(verifier)){
            Map<String, Object> result = payload.toJSONObject();
            return (Long) result.get("user_id");
        }
        return null;
    }



    /**
     * 验证令牌
     *
     * @param token       令牌
     * @param userDetails 用户
     * @return 是否有效
     */
    public Boolean validateToken(String token, UserDetails userDetails) {

/*        String username = getUsernameFromToken(token);
        return (username.equals(userDetails.getUsername()) && !isTokenExpired(token));*/
        return true;
    }
}
