insert into "user" ("name", "password", "age") values
('user1', 'password', 15),
('user2', 'password', 25),
('user3', 'password', 35);

insert into "coupon" ("name", "code") values
('cpn1000', '
    return 1000
'),

('cpnage', '
    age = get_age()
    if age < 20 then
        return 5000
    else
        return 2500
    end
'),

('cpn10times', '
    limit = (get_data("limit") or 10) + 0
    retvalue = limit * 1000
    limit = (limit <= 0) and 0 or (limit - 1)
    set_data("limit", limit)
    return retvalue
');
