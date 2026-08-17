package comment

// GetReactions returns emoji→count map for a post.
func (s *Store) GetReactions(slug string) (map[string]int, error) {
	rows, err := s.db.Query(`SELECT emoji, count FROM reactions WHERE post_slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]int)
	for rows.Next() {
		var emoji string
		var count int
		if err := rows.Scan(&emoji, &count); err != nil {
			return nil, err
		}
		out[emoji] = count
	}
	return out, rows.Err()
}

// IncrementReaction bumps the reaction count for a post+emoji and returns all reactions for that post.
func (s *Store) IncrementReaction(slug, emoji string) (map[string]int, error) {
	_, err := s.db.Exec(`
		INSERT INTO reactions (post_slug, emoji, count) VALUES (?, ?, 1)
		ON CONFLICT(post_slug, emoji) DO UPDATE SET count = count + 1
	`, slug, emoji)
	if err != nil {
		return nil, err
	}
	return s.GetReactions(slug)
}

// DecrementReaction decrements the reaction count for a post+emoji (capping at 0) and returns all reactions for that post.
func (s *Store) DecrementReaction(slug, emoji string) (map[string]int, error) {
	_, err := s.db.Exec(`
		UPDATE reactions SET count = CASE WHEN count > 0 THEN count - 1 ELSE 0 END 
		WHERE post_slug = ? AND emoji = ?
	`, slug, emoji)
	if err != nil {
		return nil, err
	}
	return s.GetReactions(slug)
}
