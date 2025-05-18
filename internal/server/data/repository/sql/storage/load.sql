SELECT metadata, data
FROM gokeeper.public.data
WHERE id=$1 AND owner=$2