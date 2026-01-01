/**
 * UserStore Usage Examples
 *
 * This file demonstrates how to use the userStore in your React components.
 */

import { useEffect, useState } from 'react';
import { View, Button, Text } from '@tarojs/components';
import { useUserStore } from './userStore';
import Taro from '@tarojs/taro';

// ============================================
// Example 1: Basic User Info Display
// ============================================
export const UserProfile = () => {
  const { userInfo, loading, fetchUserInfo } = useUserStore();

  useEffect(() => {
    // Fetch user info on component mount
    fetchUserInfo();
  }, []);

  if (loading) {
    return <Text>Loading...</Text>;
  }

  if (!userInfo) {
    return <Text>No user info</Text>;
  }

  return (
    <View>
      <Text>Name: {userInfo.name}</Text>
      <Text>Phone: {userInfo.phone}</Text>
      <Text>Avatar: {userInfo.avatar}</Text>
    </View>
  );
};

// ============================================
// Example 2: Update User Profile
// ============================================
export const UpdateProfile = () => {
  const { userInfo, updateProfile } = useUserStore();
  const [name, setName] = useState(userInfo?.name || '');

  const handleUpdate = async () => {
    try {
      await updateProfile({ name });
      Taro.showToast({
        title: '更新成功',
        icon: 'success',
      });
    } catch (error) {
      console.error('Update failed:', error);
    }
  };

  return (
    <View>
      <Input value={name} onInput={(e) => setName(e.detail.value)} />
      <Button onClick={handleUpdate}>Update Profile</Button>
    </View>
  );
};

// ============================================
// Example 3: Upload Avatar
// ============================================
export const AvatarUpload = () => {
  const { userInfo, uploadAvatar } = useUserStore();

  const handleChooseImage = () => {
    Taro.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: async (res) => {
        const filePath = res.tempFilePaths[0];
        try {
          const avatarUrl = await uploadAvatar(filePath);
          console.log('Avatar uploaded:', avatarUrl);
        } catch (error) {
          console.error('Upload failed:', error);
        }
      },
    });
  };

  return (
    <View>
      <Image src={userInfo?.avatar || '/default-avatar.png'} />
      <Button onClick={handleChooseImage}>Choose Avatar</Button>
    </View>
  );
};

// ============================================
// Example 4: Wallet Balance Display
// ============================================
export const WalletBalance = () => {
  const { wallet, fetchWallet, balanceYuan, frozenBalanceYuan } = useUserStore();

  useEffect(() => {
    fetchWallet();
  }, []);

  return (
    <View>
      <Text>Available Balance: ¥{balanceYuan().toFixed(2)}</Text>
      <Text>Frozen Balance: ¥{frozenBalanceYuan().toFixed(2)}</Text>
      <Button onClick={fetchWallet}>Refresh</Button>
    </View>
  );
};

// ============================================
// Example 5: VIP Status Display
// ============================================
export const VipStatus = () => {
  const { wallet, isVip, vipDaysLeft } = useUserStore();

  if (!isVip()) {
    return <Text>Not a VIP member</Text>;
  }

  return (
    <View>
      <Text>VIP Level: {wallet.vipLevel}</Text>
      <Text>Expires in: {vipDaysLeft()} days</Text>
    </View>
  );
};

// ============================================
// Example 6: Combined User Dashboard
// ============================================
export const UserDashboard = () => {
  const {
    userInfo,
    loading,
    fetchUserInfo,
    fetchWallet,
    updateProfile,
    uploadAvatar,
    balanceYuan,
    isVip,
    vipDaysLeft,
  } = useUserStore();

  useEffect(() => {
    // Fetch all user data on mount
    fetchUserInfo();
    fetchWallet();
  }, []);

  const handleRefresh = () => {
    fetchUserInfo();
    fetchWallet();
  };

  const handleUpdateName = async () => {
    const newName = prompt('Enter new name:');
    if (newName) {
      try {
        await updateProfile({ name: newName });
      } catch (error) {
        console.error('Update failed:', error);
      }
    }
  };

  const handleUploadAvatar = () => {
    Taro.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: async (res) => {
        try {
          await uploadAvatar(res.tempFilePaths[0]);
        } catch (error) {
          console.error('Upload failed:', error);
        }
      },
    });
  };

  if (loading) {
    return <Text>Loading...</Text>;
  }

  return (
    <View>
      {/* User Info Section */}
      <View className="user-info">
        <Image src={userInfo?.avatar || '/default-avatar.png'} />
        <Text>{userInfo?.name}</Text>
        <Text>{userInfo?.phone}</Text>
        <Button onClick={handleUpdateName}>Update Name</Button>
        <Button onClick={handleUploadAvatar}>Upload Avatar</Button>
      </View>

      {/* Wallet Section */}
      <View className="wallet">
        <Text>Balance: ¥{balanceYuan().toFixed(2)}</Text>
        <Button onClick={fetchWallet}>Refresh Balance</Button>
      </View>

      {/* VIP Section */}
      <View className="vip">
        {isVip() ? (
          <View>
            <Text>VIP Level: {wallet.vipLevel}</Text>
            <Text>Expires in: {vipDaysLeft()} days</Text>
          </View>
        ) : (
          <Text>Not a VIP member</Text>
        )}
      </View>

      {/* Refresh Button */}
      <Button onClick={handleRefresh}>Refresh All</Button>
    </View>
  );
};

// ============================================
// Example 7: Using Selectors (Computed Values)
// ============================================
export const ComputedValuesExample = () => {
  const { balanceYuan, frozenBalanceYuan, isVip, vipDaysLeft } = useUserStore();

  const totalBalance = balanceYuan() + frozenBalanceYuan();
  const vipStatusText = isVip() ? `VIP (${vipDaysLeft()} days left)` : 'Regular User';

  return (
    <View>
      <Text>Total Balance: ¥{totalBalance.toFixed(2)}</Text>
      <Text>Status: {vipStatusText}</Text>
    </View>
  );
};

// ============================================
// Example 8: Auto-refresh Wallet on Payment
// ============================================
export const PaymentSuccessHandler = () => {
  const { fetchWallet } = useUserStore();

  const handlePaymentSuccess = () => {
    // After successful payment, refresh wallet
    fetchWallet();
  };

  return (
    <Button onClick={handlePaymentSuccess}>
      Simulate Payment
    </Button>
  );
};

// ============================================
// Example 9: Conditional Rendering Based on VIP
// ============================================
export const VipOnlyFeature = () => {
  const { isVip } = useUserStore();

  if (!isVip()) {
    return (
      <View>
        <Text>This feature is for VIP members only</Text>
        <Button>Upgrade to VIP</Button>
      </View>
    );
  }

  return (
    <View>
      <Text>VIP Exclusive Content</Text>
      <Text>Enjoy your premium benefits!</Text>
    </View>
  );
};

// ============================================
// Example 10: Combining Multiple Stores
// ============================================
import { useAuthStore } from './authStore';

export const CombinedExample = () => {
  const { userInfo, balanceYuan } = useUserStore();
  const { token } = useAuthStore();

  useEffect(() => {
    // Check if user is authenticated
    if (token && !userInfo) {
      // Fetch user info if authenticated but no user data
      useUserStore.getState().fetchUserInfo();
    }
  }, [token, userInfo]);

  return (
    <View>
      <Text>Authenticated: {!!token}</Text>
      <Text>Balance: ¥{balanceYuan().toFixed(2)}</Text>
    </View>
  );
};
