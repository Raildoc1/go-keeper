SELECT metadata, data
FROM gokeeper.public.data
WHERE guid=$1 AND owner=$2