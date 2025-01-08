-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rodent_cages (
    id UUID NOT NULL,
    activation_code TEXT NOT NULL,
    user_id UUID,
    name TEXT NOT NULL,
    description TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    registered_at TIMESTAMP,

    CONSTRAINT rodent_cages_pkey PRIMARY KEY (id),
    CONSTRAINT rodent_cages_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TRIGGER update_rodent_cages_updated_at
BEFORE UPDATE ON rodent_cages
FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
