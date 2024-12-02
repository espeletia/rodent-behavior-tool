-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS media (
    id UUID NOT NULL,
    mimetype VARCHAR(255) NOT NULL,
    variant VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    entity_type VARCHAR(255) NOT NULL,
    master_id UUID,
    duration int,
    size int NOT NULL,
    width int NOT NULL,
    height int NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT media_pkey PRIMARY KEY (id),
    CONSTRAINT media_master_id_fk FOREIGN KEY (master_id) REFERENCES media (id) ON DELETE CASCADE
);

CREATE TRIGGER update_media_updated_at
BEFORE UPDATE ON media
FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
