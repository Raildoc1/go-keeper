INSERT INTO data (owner, metadata, data)
VALUES ($1, $2, $3)
RETURNING id