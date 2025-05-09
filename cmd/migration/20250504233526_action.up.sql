
CREATE OR REPLACE FUNCTION update_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER set_updated_at_user
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE FUNCTION update_at_column();


CREATE TRIGGER set_updated_at_product
    BEFORE UPDATE ON product
    FOR EACH ROW
EXECUTE FUNCTION update_at_column();


CREATE TRIGGER set_updated_at_order
    BEFORE UPDATE ON orders
    FOR EACH ROW
EXECUTE FUNCTION update_at_column();


CREATE TRIGGER set_updated_at_payment
    BEFORE UPDATE ON payment
    FOR EACH ROW
EXECUTE FUNCTION update_at_column();
