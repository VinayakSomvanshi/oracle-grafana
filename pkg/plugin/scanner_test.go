package plugin

import (
	"testing"
	"time"
)

func TestTypedColumnBuilderNumber(t *testing.T) {
	b := &TypedColumnBuilder{Name: "val", Kind: KindNumber}
	b.Append(123)
	b.Append(45.67)
	b.Append("89.1")
	b.Append([]byte("99.9"))
	b.Append(nil)

	if len(b.Numbers) != 5 {
		t.Fatalf("expected 5 items, got %d", len(b.Numbers))
	}
	if *b.Numbers[0] != 123.0 {
		t.Errorf("expected 123.0, got %v", *b.Numbers[0])
	}
	if *b.Numbers[1] != 45.67 {
		t.Errorf("expected 45.67, got %v", *b.Numbers[1])
	}
	if *b.Numbers[2] != 89.1 {
		t.Errorf("expected 89.1, got %v", *b.Numbers[2])
	}
	if *b.Numbers[3] != 99.9 {
		t.Errorf("expected 99.9, got %v", *b.Numbers[3])
	}
	if b.Numbers[4] != nil {
		t.Errorf("expected nil for 5th item, got %v", b.Numbers[4])
	}

	f := b.ToField()
	if f.Name != "val" {
		t.Errorf("expected field name 'val', got %s", f.Name)
	}
	if f.Len() != 5 {
		t.Errorf("expected field len 5, got %d", f.Len())
	}
}

func TestTypedColumnBuilderTime(t *testing.T) {
	b := &TypedColumnBuilder{Name: "ts", Kind: KindTime}
	now := time.Now().UTC().Truncate(time.Second)
	b.Append(now)
	b.Append("2026-09-30 14:00:00")
	b.Append(nil)

	if len(b.Times) != 3 {
		t.Fatalf("expected 3 items, got %d", len(b.Times))
	}
	if !b.Times[0].Equal(now) {
		t.Errorf("expected %v, got %v", now, *b.Times[0])
	}
	if b.Times[2] != nil {
		t.Errorf("expected nil for 3rd item, got %v", b.Times[2])
	}

	f := b.ToField()
	if f.Len() != 3 {
		t.Errorf("expected field len 3, got %d", f.Len())
	}
}
