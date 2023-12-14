create index settlements_parents_id
    on settlements_parents (id);

create index settlements_parents_settlement_id
    on settlements_parents (settlement_id);

create index settlements_parents_parent_id
    on settlements_parents (parent_id);