-- Support the hot list paths with indexes that include their stable ordering.
CREATE INDEX diary_entries_child_order_idx
    ON diary_entries(child_id, observed_at DESC, id DESC);
