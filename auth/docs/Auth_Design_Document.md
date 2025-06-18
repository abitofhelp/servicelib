# Authentication and Authorization Design Document

## 1. Introduction

### 1.1 Purpose
This document outlines the design for implementing authentication and authorization in the Family Service using JWT tokens. The solution ensures that all requests have a valid access token and that users are authorized to access specific endpoints based on their roles.

### 1.2 Scope
The authentication and authorization system will:
- Validate access tokens for all requests
- Verify that the source sending the request is authorized to use an endpoint
- Support role-based access control with 'admin' and 'authuser' roles
- Integrate with the existing GraphQL API

### 1.3 References
- OAuth 2.0 and OpenID Connect (OIDC) standards
- JWT (JSON Web Token) specification
- Family Service GraphQL schema

## 2. Requirements

### 2.1 Functional Requirements

#### 2.1.1 Authentication Requirements
1. All requests must include a valid access token in the Authorization header
2. The system must validate the access token's signature, expiration, and issuer
3. The system must extract user information (ID and roles) from the token
4. The system must support both JWT and OIDC token validation

#### 2.1.2 Authorization Requirements
1. The system must enforce role-based access control:
   - 'admin' role: Can perform all operations (queries and mutations)
   - 'authuser' role: Can only perform query operations (no mutations)
2. Unauthorized requests must be rejected with appropriate error messages
3. The system must log authentication and authorization events

### 2.2 Non-Functional Requirements

#### 2.2.1 Security
1. Tokens must be signed with a secure algorithm (HS256 or RS256)
2. Token validation must check for expiration and proper signature
3. Sensitive information must not be exposed in logs

#### 2.2.2 Performance
1. Token validation should add minimal overhead to request processing
2. The system should use caching where appropriate to improve performance

#### 2.2.3 Maintainability
1. The authentication and authorization code should be modular and reusable
2. The system should follow the existing architectural patterns

## 3. System Design

### 3.1 Architecture Overview

The authentication and authorization system will be implemented as middleware that intercepts all requests to the GraphQL API. It will validate the access token, extract user information, and make authorization decisions based on the user's role and the requested operation.

### 3.2 Component Diagram

```
+------------------+     +------------------+     +------------------+
|                  |     |                  |     |                  |
|  Client          |---->|  Auth Middleware |---->|  GraphQL Server  |
|                  |     |                  |     |                  |
+------------------+     +------------------+     +------------------+
                               |
                               v
                         +------------------+
                         |                  |
                         |  Auth Services   |
                         |  - JWT Service   |
                         |  - OIDC Service  |
                         |                  |
                         +------------------+
```

### 3.3 Sequence Diagram

#### 3.3.1 Authentication Flow

```
+--------+                  +----------------+              +-------------+              +-------------+
| Client |                  | Auth Middleware |              | JWT Service |              | GraphQL API |
+--------+                  +----------------+              +-------------+              +-------------+
    |                               |                              |                            |
    | 1. Request with JWT token     |                              |                            |
    |------------------------------>|                              |                            |
    |                               |                              |                            |
    |                               | 2. Extract token             |                            |
    |                               |------------------------      |                            |
    |                               |                       |      |                            |
    |                               |<-----------------------      |                            |
    |                               |                              |                            |
    |                               | 3. Validate token            |                            |
    |                               |----------------------------->|                            |
    |                               |                              |                            |
    |                               | 4. Return claims             |                            |
    |                               |<-----------------------------|                            |
    |                               |                              |                            |
    |                               | 5. Add user info to context  |                            |
    |                               |------------------------      |                            |
    |                               |                       |      |                            |
    |                               |<-----------------------      |                            |
    |                               |                              |                            |
    |                               | 6. Forward request           |                            |
    |                               |---------------------------------------------------------->|
    |                               |                              |                            |
    |                               |                              |                            |
    | 7. Response                   |                              |                            |
    |<-------------------------------------------------------------------------------------------
    |                               |                              |                            |
```

#### 3.3.2 Authorization Flow

```
+--------+                  +----------------+              +-------------------+              +-------------+
| Client |                  | GraphQL Server |              | Auth Service      |              | Resolver    |
+--------+                  +----------------+              +-------------------+              +-------------+
    |                               |                              |                                 |
    | 1. GraphQL mutation request   |                              |                                 |
    |------------------------------>|                              |                                 |
    |                               |                              |                                 |
    |                               | 2. Check authorization       |                                 |
    |                               |----------------------------->|                                 |
    |                               |                              |                                 |
    |                               | 3. Return authorization      |                                 |
    |                               |     decision                 |                                 |
    |                               |<-----------------------------|                                 |
    |                               |                              |                                 |
    |                               | 4. If authorized, execute    |                                 |
    |                               |------------------------------------------>|                   |
    |                               |                              |                                 |
    |                               |                              |                                 |
    |                               | 5. Return result             |                                 |
    |                               |<------------------------------------------|                   |
    |                               |                              |                                 |
    | 6. Response                   |                              |                                 |
    |<------------------------------|                              |                                 |
    |                               |                              |                                 |
```

### 3.4 Class Diagram

```
+-------------------+       +-------------------+       +-------------------+
| AuthMiddleware    |       | JWTService        |       | OIDCService       |
+-------------------+       +-------------------+       +-------------------+
| - jwtService      |------>| - config          |       | - provider        |
| - oidcService     |------>| - logger          |       | - verifier        |
| - logger          |       +-------------------+       | - oauth2Config    |
| - tracer          |       | + GenerateToken() |       | - adminRoleName   |
+-------------------+       | + ValidateToken() |       | - logger          |
| + Middleware()    |       +-------------------+       | - tracer          |
+-------------------+                                   +-------------------+
        |                                               | + ValidateToken() |
        |                                               | + IsAdmin()       |
        v                                               +-------------------+
+-------------------+
| AuthService       |
+-------------------+
| - logger          |
| - tracer          |
+-------------------+
| + IsAuthorized()  |
| + IsAdmin()       |
| + GetUserID()     |
| + GetUserRoles()  |
+-------------------+
```

## 4. Implementation Details

### 4.1 Token Structure

The JWT token will have the following structure:

```json
{
  "sub": "user-id",
  "roles": ["admin", "authuser"],
  "iss": "family-service",
  "exp": 1619712000,
  "iat": 1619625600,
  "nbf": 1619625600
}
```

### 4.2 Authentication Flow

1. Client obtains a JWT token from an authentication server (not in scope)
2. Client includes the token in the Authorization header of requests
3. AuthMiddleware extracts the token from the header
4. AuthMiddleware validates the token using JWTService or OIDCService
5. If valid, user information is added to the request context
6. If invalid, the request is rejected with an appropriate error

### 4.3 Authorization Flow

1. GraphQL resolvers check if the user is authorized for the requested operation
2. For query operations, both 'admin' and 'authuser' roles are allowed
3. For mutation operations, only 'admin' role is allowed
4. If unauthorized, the request is rejected with an appropriate error

### 4.4 Error Handling

1. Missing token: 401 Unauthorized with message "Authorization required"
2. Invalid token: 401 Unauthorized with message "Invalid token"
3. Expired token: 401 Unauthorized with message "Token expired"
4. Unauthorized operation: 403 Forbidden with message "Operation not permitted"

## 5. Deployment

### 5.1 Deployment Diagram

```
+------------------+     +------------------+     +------------------+
|                  |     |                  |     |                  |
|  Client          |---->|  API Gateway     |---->|  Family Service  |
|                  |     |                  |     |                  |
+------------------+     +------------------+     +------------------+
                               |
                               v
                         +------------------+
                         |                  |
                         |  Auth Server     |
                         |  (OIDC Provider) |
                         |                  |
                         +------------------+
```

### 5.2 Configuration

The authentication and authorization system will be configured using environment variables and configuration files:

```yaml
auth:
  jwt:
    secretKey: ${JWT_SECRET_KEY}
    tokenDuration: 24h
    issuer: family-service
  oidc:
    issuerURL: ${OIDC_ISSUER_URL}
    clientID: ${OIDC_CLIENT_ID}
    clientSecret: ${OIDC_CLIENT_SECRET}
    redirectURL: ${OIDC_REDIRECT_URL}
    scopes: [openid, profile, email]
    adminRoleName: admin
```

### 5.3 Security Considerations

1. JWT secret keys must be stored securely and rotated periodically
2. OIDC client secrets must be stored securely
3. All communication must use HTTPS
4. Token lifetimes should be limited to reduce the risk of token theft

## 6. Testing

### 6.1 Unit Testing

1. Test token generation and validation
2. Test role-based authorization logic
3. Test error handling

### 6.2 Integration Testing

1. Test authentication middleware with mock tokens
2. Test GraphQL resolvers with authenticated and unauthenticated requests
3. Test authorization with different user roles

### 6.3 Security Testing

1. Test token expiration handling
2. Test invalid token handling
3. Test authorization bypass attempts

## 7. Conclusion

This design document outlines a comprehensive approach to implementing authentication and authorization in the Family Service using JWT tokens. The solution ensures that all requests have a valid access token and that users are authorized to access specific endpoints based on their roles.

The implementation will leverage the existing auth-related code in the repository and follow the architectural patterns established in the project. The solution is designed to be secure, performant, and maintainable.