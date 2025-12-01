package admin

import "testing"

// TestValidPassword 测试密码验证逻辑
func TestValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
		reason   string
	}{
		// ✅ 有效密码
		{
			name:     "Valid_AllRequirements",
			password: "Abcd123!",
			want:     true,
			reason:   "包含大写、小写、数字、特殊字符",
		},
		{
			name:     "Valid_LongerPassword",
			password: "MyP@ssw0rd123",
			want:     true,
			reason:   "更长的有效密码",
		},
		{
			name:     "Valid_WithMultipleSpecialChars",
			password: "Test@123#Pass!",
			want:     true,
			reason:   "包含多个特殊字符",
		},

		// ❌ 无效密码 - 长度问题
		{
			name:     "Invalid_TooShort",
			password: "Ab1!",
			want:     false,
			reason:   "密码少于8位",
		},
		{
			name:     "Invalid_SevenChars",
			password: "Abc123!",
			want:     false,
			reason:   "密码只有7位",
		},
		{
			name:     "Invalid_Empty",
			password: "",
			want:     false,
			reason:   "空密码",
		},

		// ❌ 无效密码 - 缺少大写字母
		{
			name:     "Invalid_NoUppercase",
			password: "abcd123!",
			want:     false,
			reason:   "缺少大写字母",
		},

		// ❌ 无效密码 - 缺少小写字母
		{
			name:     "Invalid_NoLowercase",
			password: "ABCD123!",
			want:     false,
			reason:   "缺少小写字母",
		},

		// ❌ 无效密码 - 缺少数字
		{
			name:     "Invalid_NoDigit",
			password: "Abcd!@#$",
			want:     false,
			reason:   "缺少数字",
		},

		// ❌ 无效密码 - 缺少特殊字符
		{
			name:     "Invalid_NoSpecialChar",
			password: "Abcd1234",
			want:     false,
			reason:   "缺少特殊字符",
		},

		// ❌ 无效密码 - 只有小写和数字（旧版本会通过）
		{
			name:     "Invalid_OnlyLowercaseAndDigit",
			password: "abcd1234",
			want:     false,
			reason:   "只有小写和数字（旧版本的弱密码）",
		},

		// ❌ 无效密码 - 常见弱密码
		{
			name:     "Invalid_CommonWeakPassword1",
			password: "password123",
			want:     false,
			reason:   "常见弱密码：缺少大写和特殊字符",
		},
		{
			name:     "Invalid_CommonWeakPassword2",
			password: "123456",
			want:     false,
			reason:   "常见弱密码：过短且缺少字母和特殊字符",
		},

		// 边界情况
		{
			name:     "Valid_ExactlyEightChars",
			password: "Abcd123!",
			want:     true,
			reason:   "正好8位",
		},
		{
			name:     "Valid_AllSpecialCharsTypes",
			password: "Abc123!@#$%^&*()",
			want:     true,
			reason:   "包含多种特殊字符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validPassword(tt.password)
			if got != tt.want {
				t.Errorf("validPassword(%q) = %v, want %v (原因: %s)", tt.password, got, tt.want, tt.reason)
			}
		})
	}
}

// TestValidPassword_SecurityImprovements 测试安全改进
func TestValidPassword_SecurityImprovements(t *testing.T) {
	t.Run("RejectsOldWeakPasswords", func(t *testing.T) {
		// 这些密码在旧版本中会通过，但现在应该被拒绝
		weakPasswords := []string{
			"abc123",     // 6位，只有小写+数字
			"test12",     // 6位，只有小写+数字
			"user99",     // 6位，只有小写+数字
			"password1",  // 只有小写+数字，缺少大写和特殊字符
			"Test1234",   // 缺少特殊字符
			"TEST1234!",  // 缺少小写
			"testtest!",  // 缺少大写和数字
		}

		for _, pw := range weakPasswords {
			if validPassword(pw) {
				t.Errorf("validPassword(%q) = true, 应该拒绝此弱密码", pw)
			}
		}
	})

	t.Run("AcceptsStrongPasswords", func(t *testing.T) {
		strongPasswords := []string{
			"MyP@ssw0rd",
			"Secure123!",
			"C0mpl3x!Pass",
			"Test@123Pass",
			"Abcd!234Efgh",
		}

		for _, pw := range strongPasswords {
			if !validPassword(pw) {
				t.Errorf("validPassword(%q) = false, 应该接受此强密码", pw)
			}
		}
	})

	t.Run("EnforcesMinimumLength", func(t *testing.T) {
		// 测试最小长度要求（8位）
		if validPassword("Ab1!") {
			t.Error("应该拒绝短于8位的密码")
		}
		if validPassword("Abc123!") {
			t.Error("应该拒绝7位密码")
		}
		if !validPassword("Abc123!!") {
			t.Error("应该接受8位强密码")
		}
	})

	t.Run("PreventsDosAttack", func(t *testing.T) {
		// 测试最大长度限制（防止DoS攻击）
		longPassword := make([]byte, 129)
		for i := range longPassword {
			longPassword[i] = 'a'
		}
		if validPassword(string(longPassword)) {
			t.Error("应该拒绝超过128位的密码")
		}
	})
}

// TestContainsRune 测试辅助函数
func TestContainsRune(t *testing.T) {
	tests := []struct {
		name string
		s    string
		r    rune
		want bool
	}{
		{"Found_ExclamationMark", "!@#$", '!', true},
		{"Found_AtSign", "!@#$", '@', true},
		{"NotFound_Letter", "!@#$", 'a', false},
		{"NotFound_Digit", "!@#$", '1', false},
		{"Empty_String", "", 'a', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsRune(tt.s, tt.r)
			if got != tt.want {
				t.Errorf("containsRune(%q, %q) = %v, want %v", tt.s, tt.r, got, tt.want)
			}
		})
	}
}
