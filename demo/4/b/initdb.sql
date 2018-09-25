create sequence "global_seq";

create table "user" (
	"_user" bigint not null default nextval('global_seq'),
	"name" varchar not null,
	"password" varchar not null,
    "age" int not null default 0,
	constraint user__user primary key ("_user"),
	constraint user_name unique ("name")
);

create table "coupon" (
    "_coupon" bigint not null default nextval('global_seq'),
    "name" varchar not null,
    "code" varchar not null,
	constraint coupon__coupon primary key ("_coupon"),
    constraint coupon_name unique ("name")
);

create table "coupon_data" (
    "_coupon_data" bigint not null default nextval('global_seq'),
    "_user" bigint not null,
    "_coupon" bigint not null,
    "data" jsonb not null,
    constraint coupon_data__coupon_data primary key ("_coupon_data"),
    foreign key ("_user") references "user"("_user"),
    foreign key ("_coupon") references "coupon"("_coupon")
);
