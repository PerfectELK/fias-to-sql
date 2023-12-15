CREATE MATERIALIZED VIEW settlements_parents AS
SELECT fias.id,
       fias.settlement_id,
       fias.parent_id
FROM (WITH cities AS (SELECT fias_objects.object_id
                      FROM fias_objects
                               JOIN fias_object_kladr ON fias_objects.object_id = fias_object_kladr.object_id
                      WHERE (level < 6 OR type_name IN
                                          ('г', 'г.', 'пгт', 'пгт.', 'Респ', 'обл', 'обл.', 'Аобл', 'а.обл.', 'а.окр.',
                                           'АО', 'г.ф.з.')))
      SELECT fias_objects_hierarchy.id, fias_objects_hierarchy.object_id AS settlement_id, parent_object_id AS parent_id
      FROM fias_objects_hierarchy
               JOIN cities AS c1 ON c1.object_id = fias_objects_hierarchy.object_id
               JOIN cities AS c2 ON c2.object_id = fias_objects_hierarchy.parent_object_id) AS fias;