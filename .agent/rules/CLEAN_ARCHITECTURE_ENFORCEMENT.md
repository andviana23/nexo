O projeto NEXO segue strict Clean Architecture. As dependências SEMPRE apontam do externo para o interno:

Presentation (HTTP/UI) → Application (Use Cases) → Domain (Regras de negócio) → Infrastructure (DB, APIs externas).

Obrigatório:
- Toda lógica de negócio fica no DOMAIN ou nos USE CASES.
- HANDLERS não podem conter regra de negócio; apenas orquestram.
- REPOSITÓRIOS são acessados via interfaces definidas no DOMAIN.
- O DOMAIN é puro, sem dependências externas.
- VALUE OBJECTS devem ser imutáveis.
- AGGREGATES devem manter invariantes dentro do próprio root.
- DTOs e MAPPERS ficam na camada de Application.
- Queries só podem existir via SQLC, nunca em código direto.
- Nunca quebrar isolamento entre camadas.
