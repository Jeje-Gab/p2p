BEGIN;

INSERT INTO project.beer_style (style, max_temp, min_temp, average_temp)
VALUES ('Weissbier', 3, -1, 1),
       ('Pilsens', 4, -2, 1),
       ('Weizenbier', 6, -4, 1),
       ('Red ale', 5, -5, 0),
       ('India pale ale', 7, -6, 0),
       ('IPA', 10, -7, 1),
       ('Dunkel', 2, -8, -3),
       ('Imperial Stouts', 13, -10, 1),
       ('Brown ale', 14, 0, 7) ON CONFLICT (style) DO NOTHING;

COMMIT;
