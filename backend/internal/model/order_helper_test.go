package model

import (
    "testing"
)

func TestGenerateOrderNo(t *testing.T) {
    s := GenerateOrderNo("X")
    if len(s) != 1+20 { t.Fatalf("unexpected length: %d", len(s)) }
    if s[:1] != "X" { t.Fatalf("unexpected prefix") }
}

func TestGenerateEscortOrderNo(t *testing.T) {
    s := GenerateEscortOrderNo()
    if s[:3] != "ESC" { t.Fatalf("unexpected prefix") }
}

func TestGenerateGiftOrderNo(t *testing.T) {
    s := GenerateGiftOrderNo()
    if s[:4] != "GIFT" { t.Fatalf("unexpected prefix") }
}