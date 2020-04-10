CREATE TABLE games(
  id text primary key,
  winner text,
  turn text not null,
  status text not null,

  created_at timestamp not null,
  updated_at timestamp not null
);

CREATE TABLE moves (
  id text primary key,
  game_id text not null,
  row_index int not null,
  column_index int not null,
  player text not null,

  created_at timestamp not null,
  updated_at timestamp not null,

  foreign key (game_id) references games (id)
)