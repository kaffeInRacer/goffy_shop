-- DROP TRIGGERS
DROP TRIGGER IF EXISTS set_updated_at_user ON users;
DROP TRIGGER IF EXISTS set_updated_at_product ON product;
DROP TRIGGER IF EXISTS set_updated_at_order ON orders;
DROP TRIGGER IF EXISTS set_updated_at_payment ON payment;

-- DROP FUNCTION
DROP FUNCTION IF EXISTS update_at_column();