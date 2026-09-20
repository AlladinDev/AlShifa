```text
github.com/AlladinDev/AlShifa
│
├── cmd/
│   └── server/
│       └── main.go                             # Application entry point
│
├── internal/                                   # Private application code
│   │
│   ├── auth/
│   │   ├── controller/
│   │   │   ├── login/
│   │   │   │   ├── handler.go
│   │   │   │   └── handler_test.go
│   │   │   ├── refresh/
│   │   │   ├── logout/
│   │   │   └── otp/
│   │   │
│   │   ├── service/
│   │   │   ├── login/
│   │   │   ├── refresh/
│   │   │   ├── logout/
│   │   │   └── otp/
│   │   │
│   │   ├── repository/
│   │   │   ├── token/
│   │   │   ├── otp/
│   │   │   └── session/
│   │   │
│   │   ├── dto/
│   │   ├── model/
│   │   ├── validator/
│   │   └── routes/
│   │
│   ├── doctor/
│   │   ├── controller/
│   │   │   ├── register/
│   │   │   │   ├── handler.go
│   │   │   │   └── handler_test.go
│   │   │   ├── login/
│   │   │   ├── profile/
│   │   │   ├── update/
│   │   │   ├── delete/
│   │   │   ├── verify_nmr/
│   │   │   ├── upload_certificate/
│   │   │   └── search/
│   │   │
│   │   ├── service/
│   │   │   ├── register/
│   │   │   ├── login/
│   │   │   ├── profile/
│   │   │   ├── update/
│   │   │   ├── delete/
│   │   │   ├── verify_nmr/
│   │   │   ├── upload_certificate/
│   │   │   └── search/
│   │   │
│   │   ├── repository/
│   │   │   ├── create/
│   │   │   ├── update/
│   │   │   ├── delete/
│   │   │   ├── search/
│   │   │   ├── certificate/
│   │   │   └── indexes/
│   │   │
│   │   ├── model/
│   │   │   └── doctor.go
│   │   │
│   │   ├── dto/
│   │   │   ├── request/
│   │   │   └── response/
│   │   │
│   │   ├── validator/
│   │   │   ├── register/
│   │   │   ├── login/
│   │   │   └── profile/
│   │   │
│   │   └── routes/
│   │       └── routes.go
│   │
│   ├── clinic/
│   │   ├── controller/
│   │   │   ├── register/
│   │   │   ├── update/
│   │   │   ├── profile/
│   │   │   ├── delete/
│   │   │   └── search/
│   │   │
│   │   ├── service/
│   │   ├── repository/
│   │   ├── model/
│   │   ├── dto/
│   │   ├── validator/
│   │   └── routes/
│   │
│   ├── appointment/
│   │   ├── controller/
│   │   │   ├── book/
│   │   │   ├── cancel/
│   │   │   ├── reschedule/
│   │   │   ├── history/
│   │   │   └── availability/
│   │   │
│   │   ├── service/
│   │   ├── repository/
│   │   ├── model/
│   │   ├── dto/
│   │   ├── validator/
│   │   └── routes/
│   │
│   ├── clinicdoctor/
│   │   ├── controller/
│   │   │   ├── assign/
│   │   │   ├── remove/
│   │   │   ├── list/
│   │   │   └── permissions/
│   │   │
│   │   ├── service/
│   │   ├── repository/
│   │   ├── model/
│   │   ├── dto/
│   │   ├── validator/
│   │   └── routes/
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── recover.go
│   │   ├── requestid.go
│   │   └── ratelimit.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   ├── mongodb.go
│   │   ├── indexes.go
│   │   └── migration.go
│   │
│   └── shared/
│       ├── constants/
│       ├── errors/
│       ├── response/
│       ├── validator/
│       └── utils/
│
├── pkg/
│   ├── cloudinary/
│   │   ├── upload.go
│   │   ├── delete.go
│   │   └── cloudinary_test.go
│   │
│   ├── email/
│   │   ├── send.go
│   │   ├── templates/
│   │   └── email_test.go
│   │
│   ├── encryption/
│   │   ├── encrypt.go
│   │   ├── decrypt.go
│   │   └── random.go
│   │
│   └── jwt/
│       ├── generate.go
│       ├── validate.go
│       ├── claims.go
│       └── refresh.go
│
├── .env.example
├── README.md
├── go.mod
├── go.sum
└── LICENSE
```

dependencies between modules:
1->doctor module depends on clinic module 
2->appointment module depends on clinic ,doctor,clinic doctor module
3->clinic doctor mapping module depends on doctor,clinic module
