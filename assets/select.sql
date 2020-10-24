term

 select data from movies where (data->'name')::jsonb ? 'Cabiria';
 select data from movies where (data->'genre')::jsonb ? 'Adventure';
 
 numericterm

 select data from movies where (data->'imdb_score')::jsonb::NUMERIC = 9.1;
 select data from movies where (data->'imdb_score')::jsonb::NUMERIC > 9.0;
 select data from movies where (data->'genre')::jsonb  ?| array['Adventure'];
 select data from movies where (data)::jsonb @> '{"imdb_score":9.1}'::jsonb;
 
 rangequery
 
 select data from movies where ((data->'imdb_score')::jsonb::NUMERIC > 9.0) and ((data->'imdb_score')::jsonb::NUMERIC < 9.2) and ((data->'99popularity')::jsonb::NUMERIC > 99);
terms
 select data from movies where (data)::jsonb @> '{"imdb_score":9.1,"genre":["Adventure"]}'::jsonb;
