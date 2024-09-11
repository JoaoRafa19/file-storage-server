## Requisitos Funcionais

| Codigo | Requisito | Descrião |
|--------|-----------|----------|
| RF 1 | Cadastrar Quarto | Deve poder cadastrar um quarto com os campos : `tipo` , `status`, `tipo`, `quant_cama_casal` , `quant_cama_solt`, `quarto`(numero do quarto), `preco`
| RF 2 | Reserva de quartos | Deve armazenar o quarto relacionado ao hóspede e o periodo da hospedagem |
| RF 3 | Cadastro do usuário | O usuário pode fazer cadastro no sistema |
| RF 4 | Autenticação | O usuário deve poder fazer login no sistema com a sua matrícula e senha | 
| RF 5 | Consulta do status | O usuário deve sr capaz de consultar o status de um quarto pelo numero do quarto | 
| RF 6 | Cadastro de hóspede | O usuário deve ser capaz de cadastrar um hóspede |
| RF 7 | Edição de hóspede | O usuário deve ser capaz de editar os dados do hóspede |
| RF 8 | Checkin do hospede | Deve ser possível registrar o checkin do hospede no quarto reservado |
| RF 9 | Consulta de histórico | Deve ser possível consultar o histórico de checkins e checkouts de cada hóspede |
| RF 10 | Fila de espera | O usuário deve ser capaz de colocar um hóspede na fila de espera de um quarto |





## Requisitos Não Funcionais

| Codigo | Requisito | Descrião |
|--------|-----------|----------|
| RNF 1 | Criptografia da senha | A senha deve ser criptografada antes de ser salva no banco de dados |
| RNF 2 | Disponibilidade de dados | Os registros de entrada e saída de hospedes devem ficar armazenados e disponíveis para consulta |



## Regras de Negócios

| Codigo | Requisito | Descrião |
|--------|-----------|----------|
| RN 1 | Restrição de acesso | Somente funcionários podem acessar o sistema|
| RN 2 | Quarto Reservado | Quando um quarto for reservado o status do quarto na base de dados deve ser alterado para `reservado` |
| RN 3 | Data de entrada e saída | Para realizar a reserva é necessário saber a data de checkin e checkout |
| RN 4 | Adiantamento para reserva | O usuário deve fazer o pagamento de 50% do valor para realizar sua reserva | 
| RN 5 | Cadastro de usuário | O usuário deve fornecer nome completo, CPF, endereço e telefone de contato para realizar o cadastro | 
| RN 6 | Cadastro de hóspede | O usuário deve fornecer nome completo, CPF, endereço e telefone de contato do hóspede para realizar o cadastro | 
| RN 7 | CPF Validado | Os CPF\`s fornecidos devem ser validados como CPF\`s legítimos antes de finalizar o cadastro |
| RN 8 | Quarto ocupado | Quando um funcionário registrar o checkin de um hóspede o quarto deve alterar seu status para `ocupado` |
| RN 9 | Limite da Fila | Somente um único hóspede pode ser colocado na lista de espera de cada quarto | 



