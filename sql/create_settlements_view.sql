CREATE MATERIALIZED VIEW settlements AS
SELECT fias.id,
       fias.fias_id,
       fias.kladr_id,
       fias.type,
       fias.type_short,
       fias.name,
       fias.created_at
FROM (SELECT fias_objects.object_id                                                     as id,
             object_guid                                                                as fias_id,
             fias_object_kladr.kladr_id,
             replace(LOWER(fias_object_types.name), '.', '')                            as type,
             replace(LOWER(type_name), '.', '')                                         as type_short,
             fias_objects.name,
             to_char(now(), 'YYYY-MM-DD HH12:MI:SS'::text)::timestamp without time zone AS created_at
      FROM fias_objects
               JOIN fias_object_kladr ON fias_objects.object_id = fias_object_kladr.object_id
               LEFT JOIN fias_object_types ON
                  fias_objects.type_name = fias_object_types.short_name AND fias_objects.level = fias_object_types.level
      WHERE fias_objects.level < 6
         OR type_name IN
            ('г', 'г.', 'пгт', 'пгт.', 'Респ', 'обл', 'обл.', 'Аобл', 'а.обл.', 'а.окр.', 'АО', 'г.ф.з.')) AS fias;