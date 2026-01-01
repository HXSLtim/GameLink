import { useEffect } from 'react';
import { View, Text, Button } from '@tarojs/components';
import Taro from '@tarojs/taro';
import { useUserStore } from '@/stores/modules/userStore';
import { useAuthStore } from '@/stores/modules/authStore';

/**
 * 示例：Taro 页面使用 Zustand stores
 */
const ProfileStoresExample = () => {
  const { userInfo, loading, fetchUserInfo, updateProfile } = useUserStore();
  const { isLoggedIn } = useAuthStore();

  useEffect(() => {
    if (isLoggedIn()) {
      fetchUserInfo();
    }
  }, []);

  const handleUpdateName = async () => {
    try {
      await updateProfile({ name: '新昵称' });
      Taro.showToast({ title: '更新成功', icon: 'success' });
    } catch (error) {
      Taro.showToast({ title: '更新失败', icon: 'error' });
    }
  };

  if (!isLoggedIn()) {
    return (
      <View>
        <Text>请先登录</Text>
      </View>
    );
  }

  return (
    <View>
      <Text>个人中心</Text>
      {loading ? (
        <Text>加载中...</Text>
      ) : (
        <>
          <Text>昵称: {userInfo?.name}</Text>
          <Text>手机: {userInfo?.phone}</Text>
          <Button onClick={handleUpdateName}>更新昵称</Button>
        </>
      )}
    </View>
  );
};

export default ProfileStoresExample;
