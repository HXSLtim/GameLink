package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEvidenceURLArray_ValueAndScan_RoundTrip(t *testing.T) {
	original := EvidenceURLArray{"https://a.example.com/1.png", "https://b.example.com/2.png"}

	v, err := original.Value()
	if err != nil {
		t.Fatalf("Value 应该不返回错误，实际: %v", err)
	}

	bytes, ok := v.([]byte)
	if !ok {
		// driver.Value 理论上允许多种类型，这里只要能被 json.Unmarshal 即可
		t.Fatalf("Value 返回了意外类型: %T", v)
	}

	var parsed EvidenceURLArray
	if err := parsed.Scan(bytes); err != nil {
		t.Fatalf("Scan 解析 JSON 失败: %v", err)
	}

	if len(parsed) != len(original) {
		t.Fatalf("Scan 后长度不一致，期望 %d，实际 %d", len(original), len(parsed))
	}
	for i := range original {
		if original[i] != parsed[i] {
			t.Fatalf("元素不一致，索引 %d，期望 %q，实际 %q", i, original[i], parsed[i])
		}
	}
}

func TestEvidenceURLArray_Scan_InvalidType(t *testing.T) {
	var e EvidenceURLArray
	if err := e.Scan("not-bytes"); err == nil {
		t.Fatalf("期望非 []byte 类型触发错误")
	}
}

func TestEvidenceURLArray_Value_JSONShape(t *testing.T) {
	e := EvidenceURLArray{"a", "b"}
	v, err := e.Value()
	if err != nil {
		t.Fatalf("Value 不应报错: %v", err)
	}

	data, ok := v.([]byte)
	if !ok {
		t.Fatalf("Value 返回的不是 []byte: %T", v)
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		t.Fatalf("返回值不是合法 JSON 数组: %v", err)
	}
	if len(arr) != 2 || arr[0] != "a" || arr[1] != "b" {
		t.Fatalf("返回 JSON 内容不符合预期: %v", arr)
	}
}

func TestOrderDispute_IsOverSLA(t *testing.T) {
	d := &OrderDispute{}
	if d.IsOverSLA() {
		t.Fatalf("没有 SLADeadline 时不应判定为超时")
	}

	future := time.Now().Add(1 * time.Hour)
	d.SLADeadline = &future
	if d.IsOverSLA() {
		t.Fatalf("SLADeadline 在未来时不应判定为超时")
	}

	past := time.Now().Add(-1 * time.Hour)
	d.SLADeadline = &past
	if !d.IsOverSLA() {
		t.Fatalf("SLADeadline 在过去时应该判定为超时")
	}
}

func TestOrderDispute_GetSLARemaining(t *testing.T) {
	d := &OrderDispute{}
	if got := d.GetSLARemaining(); got != 0 {
		t.Fatalf("没有 SLADeadline 时剩余时间应为 0，实际: %d", got)
	}

	future := time.Now().Add(2 * time.Second)
	d.SLADeadline = &future
	got := d.GetSLARemaining()
	if got <= 0 {
		t.Fatalf("未来截止时间的剩余秒数应为正，实际: %d", got)
	}

	past := time.Now().Add(-1 * time.Second)
	d.SLADeadline = &past
	if got := d.GetSLARemaining(); got != 0 {
		t.Fatalf("已过期的截止时间剩余秒数应为 0，实际: %d", got)
	}
}

func TestCanInitiateDispute_BeforeCompletion(t *testing.T) {
	// 完成时间为空，但订单处于进行中，允许发起争议
	order := &Order{Status: OrderStatusInProgress}
	if !CanInitiateDispute(order) {
		t.Fatalf("服务进行中时应允许发起争议")
	}

	// 完成时间为空，且状态不是进行中，则不允许
	order.Status = OrderStatusPending
	if CanInitiateDispute(order) {
		t.Fatalf("非进行中状态且未完成的订单不应允许发起争议")
	}
}

func TestCanInitiateDispute_AfterCompletionWindow(t *testing.T) {
	now := time.Now()

	// 在 24 小时窗口内应允许
	completedWithin24h := now.Add(-23 * time.Hour)
	order := &Order{CompletedAt: &completedWithin24h}
	if !CanInitiateDispute(order) {
		t.Fatalf("在完成 24 小时内应允许发起争议")
	}

	// 超过 24 小时则不允许
	completedOver24h := now.Add(-25 * time.Hour)
	order.CompletedAt = &completedOver24h
	if CanInitiateDispute(order) {
		t.Fatalf("超过 24 小时不应允许发起争议")
	}
}
