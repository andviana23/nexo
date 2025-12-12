-- ============================================================================
-- MIGRATION 055: Backfill Data Isolation
-- Objective: Fix orphan records by assigning them to a default unit
-- ============================================================================

DO $$
DECLARE
    default_unit_id UUID;
BEGIN
    -- 1. Find a target unit (Prefer Mangabeiras, fallback to any)
    SELECT id INTO default_unit_id FROM units WHERE name ILIKE '%Mangabeiras%' LIMIT 1;
    
    IF default_unit_id IS NULL THEN
        SELECT id INTO default_unit_id FROM units LIMIT 1;
    END IF;

    -- If no unit exists, we cannot proceed with backfill, but we can't fail the migration 
    -- if it's a fresh install (empty tables).
    IF default_unit_id IS NOT NULL THEN
        
        -- 2. Backfill Professionals
        UPDATE profissionais SET unit_id = default_unit_id WHERE unit_id IS NULL;
        
        -- 3. Backfill Services and Categories
        UPDATE servicos SET unit_id = default_unit_id WHERE unit_id IS NULL;
        UPDATE categorias_servicos SET unit_id = default_unit_id WHERE unit_id IS NULL;
        
        -- 4. Backfill Operational Data
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'appointments') THEN
            UPDATE appointments SET unit_id = default_unit_id WHERE unit_id IS NULL;
        END IF;

        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'commands') THEN
            UPDATE commands SET unit_id = default_unit_id WHERE unit_id IS NULL;
        END IF;

        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'subscriptions') THEN
            UPDATE subscriptions SET unit_id = default_unit_id WHERE unit_id IS NULL;
        END IF;

        -- Add other tables if necessary
    END IF;
END $$;
