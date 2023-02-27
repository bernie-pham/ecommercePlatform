# **Ecommerce Platform**
#### Start on Thu 29 Dec 2022
#### by <green>Bernie Pham</green>

<br/><br/>

## **<mark>Project Requirements</mark>**


#### For config management, this project using [Viper](https://github.com/spf13/viper)

> This package is used to load server config from app.env file

#### For log management, this project using [ZeroLog](https://github.com/rs/zerolog)

#### This project use gin for web development framework [Gin](https://github.com/gin-gonic/gin#quick-start)

#### This project use sqlc to generate type-safe code from SQL [SQLC](https://github.com/kyleconroy/sqlc)

#### This project use go migrate to manage database version [Migrate](https://github.com/golang-migrate/migrate)

#### This project use Paseto to generate session token [Paseto](github.com/o1egl/paseto)

#### This project use Asynq to manage backgound jobs [Asynq](https://github.com/hibiken/asynq)
##### Also use asynmon to monitor background job which distributed in Redis [Asynmon](https://github.com/hibiken/asynq)

## Project Building steps
- building the structure of project directory
- Setup configurations/settings handler with Viper
- Setup log monitoring with ZeroLog
- Design databasae scheme
- Setup DB migration
> To create new DB migration file
```
    migrate create -ext sql -dir db/migration -seq <migration name>
```
- Setup Database connection and SQL : ✅ Query with sqlC
Init sql.yml
```
    sqlc init
```
- CRUD for User Table:
    - Login User: ✅
    - Create user: ✅
    - Forgot/Reset Password: ✅
    - Update User (optional parameters): ✅
    - Verify new registered User: ✅
- Setup Session management for login feature: ✅
- Manage session using Redis: ❌
    - Store session in server into Redis ✅
    - Store only session ID on client-side ✅
    - Refresh token automatically when it is expired (implement Refresh Token Rotation) ❌
    - Remove session in redis if user logging out (Token Revocation) ✅
    - Session should be modified if any update on user, such as privilege, name, ... ❌
- Implement Oauth 2.0: ❌
- Implement Caching: ❌
- Implement CRUD for Product Criteria:
    - tags ✅
    - size ✅
    - colour ✅
    - general ✅
    - entry ✅
- Implement Cart
    - Adding ✅
    - Edit quantity ✅
    - Remove ✅
- Implement Order Management:
    - Imlement Order Management for merchant and user workflow: ✅
    - Implement gRPC for OMS (Order Software Management) integration: ✅
- Implement backgrouund job for send mail, notify ✅
- learning ReactJS for implementing the front end part
