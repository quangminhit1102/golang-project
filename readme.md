- [Get Start](#get-start)
- [Golang Language](#golang-language)
- [JWT Authentication, Authorization](#jwt-authentication-authorization)
- [Gin Framework](#gin-framework)
- [End-Point Routers](#end-point-routers)
    * [Auth Endpoint](#1)
        + [[POST] Login](#hello)
        + [[POST] Register](#hello)
        + [[POST] Forgot Password](#hello)
        + [[POST] Reset Password](#hello)
        + [[POST] Refresh Token](#hello)
    * [Product Endpoint](#2)
        + [[GET] Get list Products](#hello)
        + [[GET] Get Detail Product By ID](#hello)
        + [[POST] Create Product](#hello)
        + [[PUT] Edit Product](#hello)
        + [[DELETE] Delete Product](#hello)
### **Get Start**
REST API, GOLANG, Gin, PostGresDB, Validator-v10, Viber config...
### **Golang Language**
Go is a popular language for good reason. It offers similar performance to other “low-level” programming languages such as Java and C++, but it’s also incredibly simple, which makes the development experience delightful. 
### **Gin Framework**
The Gin framework is lightweight, well-documented, and, of course, extremely fast.
Unlike other Go web frameworks, Gin uses a custom version of HttpRouter, which means it can navigate through your API routes faster than most frameworks out there. The creators also claim it can run 40 times faster than Martini, a relatively similar framework to Gin. You can see a more detailed comparison in this benchmark.
### **JWT Authentication, Authorization**
::: mermaid
sequenceDiagram
---
title: JWT Login
---
    autonumber
    actor User
    participant Auth Server
    User->>Auth Server: POST /auth/login<br/>body {username, password}
    activate  Auth Server
    Auth Server->>Auth Server: Validate user Input
    Auth Server->>Auth Server: Generate Token, Refresh Token
    Auth Server->>Auth Server: Set Cookie
    Auth Server-->>User: Return Access Token, Refresh Token
    deactivate Auth Server
:::

::: mermaid
sequenceDiagram
---
title: Auth
---
    autonumber
    actor User
    participant Auth Server
    User->>Auth Server: Get /product/get-all<br/>Headers: Bear Access Token
    activate  Auth Server
    Auth Server->>Auth Server: Check valid Token
    opt
    Auth Server->>Auth Server: Check User Role
    end
    Auth Server-->>User: Return list products
    deactivate Auth Server
:::
::: mermaid
sequenceDiagram
---
title: JWT Refresh Token
---
    autonumber
    actor Client
    participant Auth Server
    Client ->> Auth Server: POST: /auth/refresh <br/> body{refreshToken}
    Auth Server ->> Auth Server: Check Refresh Token Exists
    Auth Server ->> Auth Server: Check Refresh Token expiration
    opt
    Auth Server ->> Auth Server: Generate Token, Refresh Token
    end
    Auth Server ->> Client: Return New Access Token, New Refresh Token
:::
### **End-Point Routers**
#### **Auth Endpoint**
##### Login : POST
{{BaseAPI}}/{{version}}/auth/login
##### Register
{{BaseAPI}}/{{version}}/auth/register
```
{
    "email":"minh@gmail.com",
    "password":"Abc@123456"
}
```
##### Forgot Password : POST
{{BaseAPI}}/{{version}}/auth/forgot-password
 ```
{
    "email":"minh@gmail.com"
}
```
##### Reset Password : POST
{{BaseAPI}}/{{version}}/auth/?email=abc@gmail.com & token=abcd-bdfe-xyz
```
{
    "newPassword":"Abc@1234567",
    "confirmPassword":"Abc@1234567"
}
```
#### **Product Endpoint** (All require Authentication Token)
##### Get All Product (Pagination) : GET
{{BaseAPI}}/{{version}}/product/get-all/?page=1&search=searchString
##### Get Product Detail : GET
{{BaseAPI}}/{{version}}/product/:id
##### Create New Product : POST
{{BaseAPI}}/{{version}}/product/add
```
{
 "name":"003",
 "price": 12,
 "description":"Love",
 "image":"Image"
}
   ```
##### Update Product : PUT
{{BaseAPI}}/{{version}}/product/:id
Example body:
```
{
    "name":"Update Name",
    "price":111,
    "description":"Updated",
    "image":"Updated"
}
```
##### Delete Product : DELETE
{{BaseAPI}}/{{version}}/product/:id
```
```
```
//                       _oo0oo_
//                      o8888888o
//                      88" . "88
//                      (| -_- |)
//                      0\  =  /0
//                    ___/`---'\___
//                  .' \\|     |// '.
//                 / \\|||  :  |||// \
//                / _||||| -:- |||||- \
//               |   | \\\  -  /// |   |
//               | \_|  ''\---/''  |_/ |
//               \  .-\__  '-'  ___/-. /
//             ___'. .'  /--.--\  `. .'___
//          ."" '<  `.___\_<|>_/___.' >' "".
//         | | :  `- \`.;`\ _ /`;.`/ - ` : | |
//         \  \ `_.   \_ __\ /__ _/   .-` /  /
//     =====`-.____`.___ \_____/___.-`___.-'=====
//                       `=---='
//
//     ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//        Phật phù hộ, Tối nay con được ngủ sớm
//     ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
```
```
