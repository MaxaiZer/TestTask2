package entities

type Verse struct {
	ID          int64  `db:"id"`
	SongID      int64  `db:"song_id"`
	VerseNumber int64  `db:"verse_number"`
	Text        string `db:"verse_text"`
}
