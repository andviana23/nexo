/**
 * NEXO - Sistema de Gestão para Barbearias
 * Multi-Unit Components Index
 *
 * Exporta todos os componentes relacionados a multi-unidade.
 */

// Seletor de unidade
export {
    ActiveUnitBadge, UnitSelector,
    UnitSelectorCompact
} from './UnitSelector';

// Guarda de unidade (proteção de rotas)
export { UnitGuard, useRequireUnit } from './UnitGuard';

// Banners e indicadores de contexto
export {
    UnitConfirmAlert, UnitContextBanner,
    UnitContextInline
} from './UnitContextBanner';

