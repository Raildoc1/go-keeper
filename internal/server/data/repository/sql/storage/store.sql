INSERT INTO data (owner, guid, metadata, data)
VALUES ($1, $2, $3, $4)
RETURNING id