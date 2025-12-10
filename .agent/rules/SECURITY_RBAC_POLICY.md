O sistema segue RBAC estrito.

Obrigatório:
- Toda rota deve validar: Autenticação → Tenant → Role → Permissão.
- Roles: Owner, Manager, Employee, Accountant.
- Barbeiro só pode acessar dados dele.
- Handlers validam autenticação e permissão; Use Cases validam regras de autorização.
- Dados sensíveis nunca podem vazar entre tenants.
- Toda violação deve resultar em 403.
- Logs de auditoria devem registrar user_id, tenant_id, ação e timestamp.
