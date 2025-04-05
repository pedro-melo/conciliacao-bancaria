# conciliacao-bancaria

História
"Uma empresa é responsável pela gestão de contas de pagamento e emissão de
boletos para seus clientes. Todo dia, milhares de boletos são emitidos e seus
respectivos pagamentos são recebidos via integrações bancárias.
Atualmente, o processo de conciliação entre boletos emitidos e pagamentos recebidos
é feito manualmente por planilhas e scripts não confiáveis. Isso tem causado
problemas como erros no saldo das contas, dificuldade em identificar boletos pagos e
falta de rastreabilidade para auditoria.
Você foi contratado para prototipar um novo serviço de conciliação automatizada — e
sua entrega será usada como base para o novo sistema de conciliação do core
bancário."
Contexto do problema
O fluxo do sistema envolve receber boletos emitidos e pagamentos recebidos, realizar
a conciliação automática entre eles, e gerar dois conjuntos de resultados: boletos
conciliados e boletos não conciliados, conforme o diagrama abaixo:
Requisitos

1. Recebimento e Armazenamento de Dados
a. O sistema deve receber e armazenar uma lista de boletos emitidos
contendo:
i. billet_id (string) - ID do boleto
ii. bank_account (string) - Representa a conta do cliente que deve
realizar o pagamento
iii. amount (float ou decimal) - valor referente ao boleto
iv. issuance_date (ISO 8601) - data de emissão do boleto
v. reference_id (string, opcional) - Identificador que pode
corresponder ao pagamento
b. O sistema deve receber e armazenar uma lista de pagamentos bancários
recebidos contendo:
i. transaction_id (string) - ID da transação
ii. bank_account (string) - Representa a conta de origem do
pagamento (do cliente)
iii. amount (float ou decimal) - valor da transação
iv. payment_date (ISO 8601) - data em que ocorreu a transação
v. reference_id (string, opcional) - Identificador que pode
corresponder ao boleto
c. Considere que todos os boletos são emitidos para a mesma conta de
destino (da empresa), porém com clientes diferentes (contas de origem
diferentes)
d. Todos os pagamentos recebidos devem, idealmente, ser conciliados
com algum boleto emitido
2. Processo de Conciliação
a. O sistema deve analisar e combinar boletos emitidos com pagamentos
recebidos usando as seguintes estratégias em ordem de prioridade:
i. 1ª Estratégia: Conciliação direta por reference_id quando
disponível e idêntico em ambos (boleto e pagamento)
ii. 2ª Estratégia: Quando reference_id não estiver disponível ou não
corresponder, conciliar por correspondência exata de:
1. conta bancária (deve ser a mesma)
2. valor (correspondência exata ou dentro da tolerância
definida)
3. data (proximidade entre issuance_date e payment_date,
priorizando datas mais próximas)
b. Regras para resolver ambiguidades:
i. Se múltiplos boletos tiverem o mesmo valor e conta, priorizar o
boleto com data de emissão mais próxima da data de pagamento
ii. Se ainda houver ambiguidade, priorizar o boleto mais antigo
c. Deve classificar as conciliações como:
i. "conciliado_com_sucesso" quando valores são idênticos
(diferença = 0)
ii. "valor_diferente" quando há discrepância dentro da tolerância
permitida (considerar até 5% de diferença como tolerável)
iii. "nao_conciliado" para boletos sem correspondência de
pagamento
d. A estratégia de conciliação utilizada deve ser registrada no resultado
3. API de Resultados
a. O sistema deve retornar o resultado da conciliação com dois conjuntos:
i. boletos_conciliados: lista de objetos com:
1. billet_id
2. transaction_id
3. bank_account
4. conciliation_status (ex: "conciliado_com_sucesso" ou
"valor_diferente")
5. conciliation_strategy (ex: "reference_id" ou
"account_amount_date")
6. amount_diff (se houver)
7. reference_id (quando utilizado na conciliação)
ii. boletos_nao_conciliados: lista de boletos não encontrados entre
os pagamentos
4. Persistência e Rastreabilidade
a. b. c. Todas as entidades (boletos, pagamentos e conciliações) devem ser
persistidas no banco de dados
O histórico de conciliações deve ser armazenado para fins de auditoria
Deve ser possível consultar o status atual e histórico de qualquer boleto
ou pagamento
Requisitos Técnicos
1. Arquitetura e Implementação
a. Desenvolver uma API RESTful com endpoints para boletos, pagamentos e
conciliações
b. Implementar uma modelagem de domínio que reflita as regras de
negócio
c. A solução deve ser escalável para processar grandes volumes de dados
(milhares de transações)
2. Banco de Dados
a. Utilizar PostgreSQL ou MySQL (containerizado)
b. O esquema mínimo deve ser criado automaticamente na inicialização
c. Modelagem de dados adequada com relacionamentos e índices
apropriados
3. Infraestrutura
a. Toda a solução deve ser orquestrada via Docker Compose
b. Deve incluir a API, banco de dados e quaisquer outras dependências
necessárias
c. Configurações de ambiente devem ser facilmente ajustáveis via
Docker Compose ou env file
4. Testes Automatizados
a. Requisito eliminatório: Implementar testes unitários e de integração
b. Os testes devem cobrir os casos de borda (valores negativos,
duplicidades, etc.)
c. Deve ser possível executar os testes de forma automatizada
5. Qualidade de Código
a. Código limpo, bem estruturado e de fácil manutenção
b. Separação clara de responsabilidades (camadas, módulos)
c. Tratamento adequado de erros e exceções
Dados de Exemplo
Boletos emitidos:
[ { "billet_id": "B001", "bank_account": "C123", "amount": 100.00, "issuance_date":
"2024-03-01", "reference_id": "REF001" }, { "billet_id": "B002", "bank_account": "C345",
"amount": 200.00, "issuance_date": "2024-03-02", "reference_id": null }, { "billet_id":
"B003", "bank_account": "C678", "amount": 300.00, "issuance_date": "2024-03-03",
"reference_id": "REF003" }, { "billet_id": "B004", "bank_account": "C910", "amount":
150.00, "issuance_date": "2024-03-04", "reference_id": "REF004" } ]
Pagamentos recebidos:
[ { "transaction_id": "T100", "bank_account": "C123", "amount": 100.00,
"data_pagamento": "2024-03-05", "reference_id": "REF001" }, { "transaction_id": "T101",
"bank_account": "C345", "amount": 195.00, "data_pagamento": "2024-03-04",
"reference_id": null }, { "transaction_id": "T102", "bank_account": "C678", "amount":
300.00, "data_pagamento": "2024-03-10", "reference_id": "REF999" }, { "transaction_id":
"T103", "bank_account": "C999", "amount": 150.00, "data_pagamento": "2024-03-06",
"reference_id": "REF004" } ]
Resultado esperado:
{ "boletos_conciliados": [ { "billet_id": "B001", "bank_account": "C123",
"transaction_id": "T100", "conciliation_status": "conciliado_com_sucesso",
"conciliation_strategy": "reference_id", "reference_id": "REF001", "amount_diff": 0.00 },
{ "billet_id": "B002", "bank_account": "C345", "transaction_id": "T101",
"conciliation_status": "valor_diferente", "conciliation_strategy": "conta_valor_data",
"reference_id": null, "amount_diff": 5.00 }, { "billet_id": "B003", "bank_account": "C678",
"transaction_id": "T102", "conciliation_status": "conciliado_com_sucesso",
"conciliation_strategy": "conta_valor_data", "reference_id": null, "amount_diff": 0.00 }, {
"billet_id": "B004", "bank_account": "C910", "transaction_id": "T103",
"conciliation_status": "conciliado_com_sucesso", "conciliation_strategy":
"reference_id", "reference_id": "REF004", "amount_diff": 0.00 } ],
"boletos_nao_conciliados": [] }