# Estrutura de Arquivos para o Sistema de Conciliação Bancária

conciliacao-bancaria/
├── cmd/
│   └── api/
│       └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── domain/
│   │   ├── model/
│   │   │   ├── billet.go           # Modelo de domínio para boletos
│   │   │   ├── payment.go          # Modelo de domínio para pagamentos
│   │   │   └── reconciliation.go   # Modelo de domínio para conciliações
│   │   ├── repository/
│   │   │   ├── billet_repository.go    # Interface para repositório de boletos
│   │   │   ├── payment_repository.go   # Interface para repositório de pagamentos
│   │   │   └── reconciliation_repository.go  # Interface para repositório de conciliações
│   │   └── service/
│   │       └── reconciliation_service.go  # Lógica de negócio para conciliação
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── connection.go       # Conexão com o banco de dados
│   │   │   ├── migrations/         # Migrações para o banco de dados
│   │   │   │   └── schema.sql      # Esquema inicial do banco de dados
│   │   │   └── repository/         # Implementações concretas dos repositórios
│   │   │       ├── billet_repository_impl.go
│   │   │       ├── payment_repository_impl.go
│   │   │       └── reconciliation_repository_impl.go
│   │   └── http/
│   │       ├── handler/
│   │       │   ├── billet_handler.go    # Handlers para endpoints de boletos
│   │       │   ├── payment_handler.go   # Handlers para endpoints de pagamentos
│   │       │   └── reconciliation_handler.go  # Handlers para endpoints de conciliação
│   │       ├── middleware/
│   │       │   ├── auth.go         # Middleware de autenticação (se necessário)
│   │       │   └── logging.go      # Middleware de logging
│   │       ├── dto/
│   │       │   ├── request/        # DTOs para requisições
│   │       │   │   ├── billet_request.go
│   │       │   │   ├── payment_request.go
│   │       │   │   └── reconciliation_request.go
│   │       │   └── response/       # DTOs para respostas
│   │       │       ├── billet_response.go
│   │       │       ├── payment_response.go
│   │       │       └── reconciliation_response.go
│   │       └── router.go          # Configuração das rotas da API
│   └── application/
│       └── usecase/
│           ├── billet_usecase.go       # Casos de uso para boletos
│           ├── payment_usecase.go      # Casos de uso para pagamentos
│           └── reconciliation_usecase.go  # Casos de uso para conciliação
├── pkg/
│   ├── errors/
│   │   └── errors.go              # Tratamento de erros customizados
│   └── utils/
│       ├── date_utils.go          # Utilitários para manipulação de datas
│       └── validators.go          # Validadores
├── config/
│   └── config.go                  # Configurações da aplicação
├── test/
│   ├── integration/               # Testes de integração
│   │   ├── billet_test.go
│   │   ├── payment_test.go
│   │   └── reconciliation_test.go
│   ├── unit/                      # Testes unitários
│   │   ├── model_test.go
│   │   ├── service_test.go
│   │   └── usecase_test.go
│   └── mocks/                     # Mocks para testes
│       ├── repository_mock.go
│       └── service_mock.go
├── scripts/                       # Scripts auxiliares
│   └── startup.sh                 # Script de inicialização
├── .env.example                   # Exemplo de variáveis de ambiente
├── Dockerfile                     # Dockerfile para a aplicação
├── docker-compose.yml             # Configuração do Docker Compose
├── go.mod                         # Gerenciamento de dependências
├── go.sum                         # Checksums das dependências
└── README.md                      # Documentação do projeto

## Descrição das Principais Pastas e Arquivos
cmd/
## Contém os pontos de entrada para a aplicação. No caso, temos apenas uma aplicação API, mas esta estrutura permite adicionar mais executáveis facilmente no futuro.
internal/
## Código privado da aplicação que não deve ser importado por outros projetos.
domain/
## Contém os modelos e interfaces que representam o domínio da aplicação.

model/: Define as entidades do domínio (boletos, pagamentos, conciliações).
repository/: Define as interfaces dos repositórios.
service/: Implementa a lógica de negócio, incluindo o algoritmo de conciliação.

infrastructure/
Contém as implementações concretas das interfaces de domínio e o código que interage com componentes externos.

database/: Conexão com o banco de dados e implementações dos repositórios.
http/: Implementação da API REST, incluindo handlers, DTOs e configuração de rotas.

application/
Contém os casos de uso que orquestram o fluxo de trabalho entre os handlers e os serviços de domínio.
pkg/
Código que pode ser reutilizado por outros projetos.

errors/: Definições de erros customizados.
utils/: Funções utilitárias, como manipulação de datas.

config/
Configurações da aplicação, como conexões com banco de dados e parâmetros do servidor.
test/
Testes unitários e de integração.
scripts/
Scripts auxiliares para execução da aplicação.
Arquivos na Raiz

Dockerfile e docker-compose.yml: Configuração para containerização.
.env.example: Exemplo de variáveis de ambiente necessárias.
go.mod e go.sum: Gerenciamento de dependências.
README.md: Documentação do projeto.

Esta estrutura segue os princípios da arquitetura hexagonal (ou arquitetura de portas e adaptadores), que facilita a separação de responsabilidades e a testabilidade do código. O domínio da aplicação está claramente separado da infraestrutura, permitindo que as regras de negócio sejam testadas independentemente das implementações de banco de dados ou API.