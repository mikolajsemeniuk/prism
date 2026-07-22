package corpus

import "testing"

func TestParseFrontmatter(t *testing.T) {
	src := "---\nproduct: \"ChatGPT\"\nsource_url: https://example.com/a:b\nnotes: first\n  continued line\n---\n\nBody starts here.\n"
	meta, body, off, err := ParseFrontmatter(src)
	if err != nil {
		t.Fatal(err)
	}
	if meta["product"] != "ChatGPT" {
		t.Fatalf("product = %q", meta["product"])
	}
	if meta["source_url"] != "https://example.com/a:b" {
		t.Fatalf("source_url = %q (first-colon split broken)", meta["source_url"])
	}
	if meta["notes"] != "first continued line" {
		t.Fatalf("notes = %q (continuation broken)", meta["notes"])
	}
	if body != "Body starts here.\n" {
		t.Fatalf("body = %q", body)
	}
	if src[off:] != body {
		t.Fatalf("offset %d does not point at body", off)
	}
}

func TestParseFrontmatterErrors(t *testing.T) {
	if _, _, _, err := ParseFrontmatter("no frontmatter"); err == nil {
		t.Fatal("want error for missing opening fence")
	}
	if _, _, _, err := ParseFrontmatter("---\nkey: value\n"); err == nil {
		t.Fatal("want error for missing closing fence")
	}
}
