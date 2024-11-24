-- noinspection LanguageDetectionInspectionForFile

CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    "group" TEXT NOT NULL,
    name TEXT NOT NULL,
    link TEXT NOT NULL,
    release_date DATE NOT NULL
);

CREATE TABLE verses (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    verse_number INT NOT NULL,
    verse_text TEXT NOT NULL,
    CONSTRAINT unique_verse_number_per_song UNIQUE (song_id, verse_number)
);


INSERT INTO songs ("group", name, link, release_date) VALUES
    ('Wings', 'Live and Let Die', 'https://www.youtube.com/watch?v=wYQZHNwIUq8','1973-06-01');

INSERT INTO verses (song_id, verse_number, verse_text) VALUES
(1, 0, 'When you were young\n' ||
    'And your heart was an open book\n' ||
    'You used to say, "Live and let live"\n' ||
    '(You know you did\n' ||
    'You know you did\n' ||
    'You know you did)\n' ||
    'But if this ever-changing world in which we''re livin''\n' ||
    'Makes you give in and cry'),
(1, 1, 'Say live and let die\n' || --не куплет конечно, но ведь часть песни...
    '(Live and let die)\n' ||
    'Live and let die\n' ||
    '(Live and let die)'),
(1, 2, 'What does it matter to ya?\n' ||
    'When you got a job to do\n' ||
    'You got to do it well\n' ||
    'You gotta give the other fella hell!'),
(1, 3, 'You used to say live and let live\n' ||
    '(You know you did\n' ||
    'You know you did\n' ||
    'You know you did)\n' ||
    'But if this ever-changing world in which we''re livin''\n' ||
    'Makes you give in and cry'),
(1, 4, 'Say live and let die\n' ||
    '(Live and let die)\n' ||
    'Live and let die\n' ||
    '(Live and let die)');

INSERT INTO songs ("group", name, link, release_date) VALUES
    ('Swans', 'Screen Shot', 'https://www.youtube.com/watch?v=6qDq9eGUmMI', '2014-5-12');

INSERT INTO verses (song_id, verse_number, verse_text) VALUES
(2, 0, 'Love, child, reach, rise; sight, blind, steal, light\n' ||
    'Mind, scar, clear, fire; clean, right, pure, kind\n' ||
    'Sun, come, sky, tar; mouth, sand, teeth, tongue\n' ||
    'Cut, push, reach, inside; feed, breathe, touch, come'),
(2, 1, 'No pain, no death, no fear, no hate\n' ||
    'No time, no now, no suffering\n' ||
    'No touch, no loss, no hand, no sense\n' ||
    'No wound, no waste, no lust, no fear\n' ||
    'No mind, no greed, no suffering\n' ||
    'No thought, no hurt, no hands to reach\n' ||
    'No knife, no words, no lie, no cure\n' ||
    'No need, no hate, no will, no speech'),
(2, 2, 'No dream, no sleep, no suffering\n' ||
    'No dream, no sleep, no suffering\n' ||
    'No dream, no sleep, no suffering\n' ||
    'No dream, no sleep, no suffering'),
(2, 3, 'No pain, no now, no time, no here\n' ||
    'No pain, no now, no time, no here\n' ||
    'No pain, no now, no time, no here\n' ||
    'No pain, no now, no time, no here\n' ||
    'No knife, no mind, no hand, no fear\n' ||
    'No knife, no mind, no hand, no fear\n' ||
    'No knife, no mind, no hand, no fear\n' ||
    'No knife, no mind, no hand, no here'),
(2, 4, 'Love! Now!\n' ||
    'Breathe! Now!\n' ||
    'Love! Now!\n' ||
    'Breathe! Now!\n' ||
    'Love! Now!\n' ||
    'Breathe! Now!\n' ||
    'Love! Now!\n' ||
    'Breathe! Now!\n' ||
    'Here! Now!\n' ||
    'Here! Now!\n' ||
    'Here! Now!\n' ||
    'Here! Now!\n' ||
    'Here! Now!\n' ||
    'Here! Now!\n' ||
    'Here! Now!\n' ||
    'Here! Now!');