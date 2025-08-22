create table horses(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name STRING not null,
  description STRING not null,
  date_of_birth DATE not null,
  gender INT not null,
);

---- create above / drop below ----

drop table horses;
