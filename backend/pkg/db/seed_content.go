package db

import (
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// seedContentData 创建内容管理模块的种子数据
func seedContentData(tx *gorm.DB, users map[string]*model.User) error {
	// 1. 创建内容分类
	categories, err := seedContentCategories(tx)
	if err != nil {
		return err
	}

	// 2. 创建聊天群组
	groups, err := seedChatGroups(tx, users)
	if err != nil {
		return err
	}

	// 3. 创建动态
	feeds, err := seedFeeds(tx, users, categories)
	if err != nil {
		return err
	}

	// 4. 创建聊天消息
	if err := seedChatMessages(tx, users, groups); err != nil {
		return err
	}

	// 5. 创建动态举报
	if err := seedFeedReports(tx, users, feeds); err != nil {
		return err
	}

	// 6. 创建内容管理权限
	if err := seedContentPermissions(tx); err != nil {
		return err
	}

	return nil
}

// seedContentCategories 创建内容分类种子数据
func seedContentCategories(tx *gorm.DB) (map[string]*model.ContentCategory, error) {
	seeds := []model.ContentCategory{
		{
			Name:        "游戏攻略",
			Description: "分享游戏攻略、技巧和心得",
			SortOrder:   1,
			Status:      model.ContentCategoryStatusActive,
			IconURL:     "https://gamelink.oss.com/icons/strategy.png",
		},
		{
			Name:        "精彩时刻",
			Description: "记录游戏中的精彩瞬间",
			SortOrder:   2,
			Status:      model.ContentCategoryStatusActive,
			IconURL:     "https://gamelink.oss.com/icons/highlight.png",
		},
		{
			Name:        "组队招募",
			Description: "寻找志同道合的队友",
			SortOrder:   3,
			Status:      model.ContentCategoryStatusActive,
			IconURL:     "https://gamelink.oss.com/icons/team.png",
		},
		{
			Name:        "陪玩日记",
			Description: "陪玩师分享工作日常",
			SortOrder:   4,
			Status:      model.ContentCategoryStatusActive,
			IconURL:     "https://gamelink.oss.com/icons/diary.png",
		},
		{
			Name:        "游戏资讯",
			Description: "最新游戏新闻和更新",
			SortOrder:   5,
			Status:      model.ContentCategoryStatusActive,
			IconURL:     "https://gamelink.oss.com/icons/news.png",
		},
		{
			Name:        "装备展示",
			Description: "展示游戏装备和外观",
			SortOrder:   6,
			Status:      model.ContentCategoryStatusActive,
			IconURL:     "https://gamelink.oss.com/icons/equipment.png",
		},
		{
			Name:        "吐槽专区",
			Description: "游戏吐槽和趣事分享",
			SortOrder:   7,
			Status:      model.ContentCategoryStatusInactive,
			IconURL:     "https://gamelink.oss.com/icons/fun.png",
		},
	}

	result := make(map[string]*model.ContentCategory, len(seeds))
	for i := range seeds {
		cat := &seeds[i]
		var existing model.ContentCategory
		if err := tx.Where("name = ?", cat.Name).First(&existing).Error; err == nil {
			result[cat.Name] = &existing
			continue
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if err := tx.Create(cat).Error; err != nil {
			return nil, err
		}
		result[cat.Name] = cat
	}
	return result, nil
}

// seedChatGroups 创建聊天群组种子数据
func seedChatGroups(tx *gorm.DB, users map[string]*model.User) (map[string]*model.ChatGroup, error) {
	adminUser := users["adminA"]
	proA := users["proA"]
	proB := users["proB"]

	seeds := []struct {
		Key       string
		GroupName string
		GroupType model.ChatGroupType
		CreatedBy uint64
		IsActive  bool
		Desc      string
	}{
		{
			Key:       "public_lol",
			GroupName: "英雄联盟交流群",
			GroupType: model.ChatGroupTypePublic,
			CreatedBy: adminUser.ID,
			IsActive:  true,
			Desc:      "英雄联盟玩家交流群，分享攻略和组队",
		},
		{
			Key:       "public_valorant",
			GroupName: "无畏契约玩家群",
			GroupType: model.ChatGroupTypePublic,
			CreatedBy: adminUser.ID,
			IsActive:  true,
			Desc:      "无畏契约玩家交流，战术分享",
		},
		{
			Key:       "public_general",
			GroupName: "GameLink综合交流群",
			GroupType: model.ChatGroupTypePublic,
			CreatedBy: adminUser.ID,
			IsActive:  true,
			Desc:      "平台综合交流群，欢迎所有玩家",
		},
		{
			Key:       "order_group_1",
			GroupName: "订单服务群-001",
			GroupType: model.ChatGroupTypeOrder,
			CreatedBy: proA.ID,
			IsActive:  true,
			Desc:      "陪玩订单专属服务群",
		},
		{
			Key:       "order_group_2",
			GroupName: "订单服务群-002",
			GroupType: model.ChatGroupTypeOrder,
			CreatedBy: proB.ID,
			IsActive:  false,
			Desc:      "已完成订单的服务群",
		},
	}

	result := make(map[string]*model.ChatGroup, len(seeds))
	for _, s := range seeds {
		var existing model.ChatGroup
		if err := tx.Where("group_name = ?", s.GroupName).First(&existing).Error; err == nil {
			result[s.Key] = &existing
			continue
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		group := &model.ChatGroup{
			GroupName:   s.GroupName,
			GroupType:   s.GroupType,
			CreatedBy:   s.CreatedBy,
			MaxMembers:  100,
			IsActive:    s.IsActive,
			Description: s.Desc,
		}
		if err := tx.Create(group).Error; err != nil {
			return nil, err
		}
		result[s.Key] = group
	}
	return result, nil
}

// seedFeeds 创建动态种子数据
func seedFeeds(tx *gorm.DB, users map[string]*model.User, categories map[string]*model.ContentCategory) (map[string]*model.Feed, error) {
	now := time.Now()

	type feedSpec struct {
		Key        string
		AuthorKey  string
		Content    string
		Category   string
		Status     model.FeedModerationStatus
		Visibility model.FeedVisibility
		Images     []string
		LikeCount  uint64
		ViewCount  uint64
	}

	specs := []feedSpec{
		// 已通过的动态
		{
			Key:        "feed_approved_1",
			AuthorKey:  "proA",
			Content:    "今天带老板连胜10把！英雄联盟打野位置真的太重要了，节奏把控好，胜率自然高。分享一下我的打野思路：前期重点关注中路和下路，4级后看情况入侵对面野区。",
			Category:   "游戏攻略",
			Status:     model.FeedModerationApproved,
			Visibility: model.FeedVisibilityPublic,
			Images:     []string{"https://gamelink.oss.com/feeds/lol_win_streak.jpg", "https://gamelink.oss.com/feeds/lol_stats.jpg"},
			LikeCount:  156,
			ViewCount:  1230,
		},
		{
			Key:        "feed_approved_2",
			AuthorKey:  "proB",
			Content:    "无畏契约新赛季开始了！这个版本杰特被削弱了不少，但依然是很强的角色。推荐大家试试新英雄，技能组合很有意思。",
			Category:   "游戏资讯",
			Status:     model.FeedModerationApproved,
			Visibility: model.FeedVisibilityPublic,
			Images:     []string{"https://gamelink.oss.com/feeds/valorant_new_season.jpg"},
			LikeCount:  89,
			ViewCount:  756,
		},
		{
			Key:        "feed_approved_3",
			AuthorKey:  "customerA",
			Content:    "感谢峡谷守护者的陪玩！从黄金一路上到铂金，技术确实提升了很多。推荐给大家，服务态度超好！",
			Category:   "陪玩日记",
			Status:     model.FeedModerationApproved,
			Visibility: model.FeedVisibilityPublic,
			Images:     []string{"https://gamelink.oss.com/feeds/rank_up.jpg"},
			LikeCount:  45,
			ViewCount:  320,
		},
		{
			Key:        "feed_approved_4",
			AuthorKey:  "proC",
			Content:    "CS:GO五杀精彩时刻！最后一个1v5翻盘，心跳加速的感觉太爽了。视频已上传，欢迎观看学习。",
			Category:   "精彩时刻",
			Status:     model.FeedModerationApproved,
			Visibility: model.FeedVisibilityPublic,
			Images:     []string{"https://gamelink.oss.com/feeds/csgo_ace_1.jpg", "https://gamelink.oss.com/feeds/csgo_ace_2.jpg", "https://gamelink.oss.com/feeds/csgo_ace_3.jpg"},
			LikeCount:  234,
			ViewCount:  1890,
		},
		{
			Key:        "feed_approved_5",
			AuthorKey:  "customerB",
			Content:    "周末有人一起打DOTA2吗？段位传奇，主玩辅助位，希望找几个稳定队友一起冲分。",
			Category:   "组队招募",
			Status:     model.FeedModerationApproved,
			Visibility: model.FeedVisibilityPublic,
			LikeCount:  23,
			ViewCount:  189,
		},
		// 待审核的动态
		{
			Key:        "feed_pending_1",
			AuthorKey:  "customerC",
			Content:    "新手求带！刚开始玩英雄联盟，有没有大佬愿意教教我？可以付费陪玩。",
			Category:   "组队招募",
			Status:     model.FeedModerationPending,
			Visibility: model.FeedVisibilityPublic,
			LikeCount:  0,
			ViewCount:  12,
		},
		{
			Key:        "feed_pending_2",
			AuthorKey:  "proD",
			Content:    "魔兽世界怀旧服开荒攻略分享！详细讲解每个BOSS的机制和站位，新手必看。",
			Category:   "游戏攻略",
			Status:     model.FeedModerationPending,
			Visibility: model.FeedVisibilityPublic,
			Images:     []string{"https://gamelink.oss.com/feeds/wow_guide.jpg"},
			LikeCount:  0,
			ViewCount:  5,
		},
		{
			Key:        "feed_pending_3",
			AuthorKey:  "customerD",
			Content:    "今天的陪玩体验一般般，感觉陪玩师不太专业...",
			Category:   "陪玩日记",
			Status:     model.FeedModerationPending,
			Visibility: model.FeedVisibilityPublic,
			LikeCount:  0,
			ViewCount:  3,
		},
		// 被拒绝的动态
		{
			Key:        "feed_rejected_1",
			AuthorKey:  "customerE",
			Content:    "这个游戏太垃圾了，天天遇到挂...",
			Category:   "吐槽专区",
			Status:     model.FeedModerationRejected,
			Visibility: model.FeedVisibilityPublic,
			LikeCount:  0,
			ViewCount:  0,
		},
		{
			Key:        "feed_rejected_2",
			AuthorKey:  "customerF",
			Content:    "加我微信xxx，私下交易更便宜",
			Status:     model.FeedModerationRejected,
			Visibility: model.FeedVisibilityPublic,
			LikeCount:  0,
			ViewCount:  0,
		},
		// 已移除的动态
		{
			Key:        "feed_removed_1",
			AuthorKey:  "customerG",
			Content:    "这条动态因违规已被移除",
			Status:     model.FeedModerationRemoved,
			Visibility: model.FeedVisibilityPublic,
			LikeCount:  5,
			ViewCount:  50,
		},
		// 仅关注者可见
		{
			Key:        "feed_followers_1",
			AuthorKey:  "proE",
			Content:    "给粉丝们的福利！本周末陪玩8折优惠，先到先得~",
			Category:   "陪玩日记",
			Status:     model.FeedModerationApproved,
			Visibility: model.FeedVisibilityFollowers,
			LikeCount:  67,
			ViewCount:  234,
		},
		// 私密动态
		{
			Key:        "feed_private_1",
			AuthorKey:  "proF",
			Content:    "今日收入统计：完成8单，总收入xxx元。继续加油！",
			Category:   "陪玩日记",
			Status:     model.FeedModerationApproved,
			Visibility: model.FeedVisibilityPrivate,
			LikeCount:  0,
			ViewCount:  1,
		},
	}

	result := make(map[string]*model.Feed, len(specs))
	for _, spec := range specs {
		user, ok := users[spec.AuthorKey]
		if !ok {
			continue
		}

		var existing model.Feed
		if err := tx.Where("author_id = ? AND content = ?", user.ID, spec.Content).First(&existing).Error; err == nil {
			result[spec.Key] = &existing
			continue
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		feed := &model.Feed{
			AuthorID:         user.ID,
			Content:          spec.Content,
			Visibility:       spec.Visibility,
			ModerationStatus: spec.Status,
			Metrics: model.FeedMetricFields{
				LikeCount: spec.LikeCount,
				ViewCount: spec.ViewCount,
			},
		}

		if spec.Category != "" {
			if cat, ok := categories[spec.Category]; ok {
				feed.CategoryID = &cat.ID
			}
		}

		if spec.Status == model.FeedModerationApproved {
			t := now.Add(-time.Duration(len(specs)) * time.Hour)
			feed.ManualModeratedAt = &t
		}

		if err := tx.Create(feed).Error; err != nil {
			return nil, err
		}

		// 创建图片
		for i, imgURL := range spec.Images {
			img := &model.FeedImage{
				FeedID: feed.ID,
				URL:    imgURL,
				Order:  i,
				Width:  800,
				Height: 600,
			}
			if err := tx.Create(img).Error; err != nil {
				return nil, err
			}
		}

		result[spec.Key] = feed
	}
	return result, nil
}

// seedChatMessages 创建聊天消息种子数据
func seedChatMessages(tx *gorm.DB, users map[string]*model.User, groups map[string]*model.ChatGroup) error {
	now := time.Now()

	type msgSpec struct {
		GroupKey    string
		SenderKey   string
		Content     string
		MessageType model.ChatMessageType
		AuditStatus model.ChatMessageAuditStatus
		IsDeleted   bool
		TimeOffset  time.Duration
	}

	specs := []msgSpec{
		// 英雄联盟交流群消息
		{GroupKey: "public_lol", SenderKey: "customerA", Content: "大家好，新人报道！", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -5 * time.Hour},
		{GroupKey: "public_lol", SenderKey: "proA", Content: "欢迎欢迎！有什么问题可以问我", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -4*time.Hour - 50*time.Minute},
		{GroupKey: "public_lol", SenderKey: "customerB", Content: "今天有人一起双排吗？", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -4 * time.Hour},
		{GroupKey: "public_lol", SenderKey: "customerC", Content: "我可以！段位是什么？", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -3*time.Hour - 45*time.Minute},
		{GroupKey: "public_lol", SenderKey: "customerB", Content: "铂金2，你呢？", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -3*time.Hour - 40*time.Minute},
		{GroupKey: "public_lol", SenderKey: "proA", Content: "分享一下今天的战绩", MessageType: model.ChatMessageTypeImage, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -2 * time.Hour},
		{GroupKey: "public_lol", SenderKey: "customerD", Content: "太强了！带带我", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -1*time.Hour - 30*time.Minute},

		// 无畏契约玩家群消息
		{GroupKey: "public_valorant", SenderKey: "proB", Content: "新赛季大家冲到什么段位了？", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -6 * time.Hour},
		{GroupKey: "public_valorant", SenderKey: "customerE", Content: "刚上钻石，感觉这赛季好难打", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -5*time.Hour - 30*time.Minute},
		{GroupKey: "public_valorant", SenderKey: "proB", Content: "确实，这赛季匹配机制改了", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -5 * time.Hour},
		{GroupKey: "public_valorant", SenderKey: "customerF", Content: "有没有人教教我怎么玩杰特？", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -3 * time.Hour},

		// 综合交流群消息
		{GroupKey: "public_general", SenderKey: "adminA", Content: "欢迎大家加入GameLink！有任何问题可以在这里提问", MessageType: model.ChatMessageTypeSystem, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -24 * time.Hour},
		{GroupKey: "public_general", SenderKey: "customerG", Content: "平台怎么下单陪玩？", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -12 * time.Hour},
		{GroupKey: "public_general", SenderKey: "adminA", Content: "可以在首页浏览陪玩师，选择喜欢的下单即可", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -11*time.Hour - 50*time.Minute},
		{GroupKey: "public_general", SenderKey: "customerH", Content: "陪玩师怎么申请入驻？", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -8 * time.Hour},

		// 待审核的消息
		{GroupKey: "public_lol", SenderKey: "customerE", Content: "这个队友太菜了吧", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditPending, TimeOffset: -30 * time.Minute},
		{GroupKey: "public_valorant", SenderKey: "customerG", Content: "有没有人私下交易？", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditPending, TimeOffset: -20 * time.Minute},

		// 被拒绝的消息
		{GroupKey: "public_general", SenderKey: "customerF", Content: "加我微信xxx私聊", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditRejected, TimeOffset: -2 * time.Hour},

		// 已删除的消息
		{GroupKey: "public_lol", SenderKey: "customerH", Content: "这条消息已被删除", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, IsDeleted: true, TimeOffset: -4 * time.Hour},

		// 订单服务群消息
		{GroupKey: "order_group_1", SenderKey: "proA", Content: "您好，我是您的陪玩师，请问现在方便开始吗？", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -1 * time.Hour},
		{GroupKey: "order_group_1", SenderKey: "customerA", Content: "可以的，我已经上线了", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -55 * time.Minute},
		{GroupKey: "order_group_1", SenderKey: "proA", Content: "好的，我加您游戏好友", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved, TimeOffset: -50 * time.Minute},
	}

	for _, spec := range specs {
		group, ok := groups[spec.GroupKey]
		if !ok {
			continue
		}
		user, ok := users[spec.SenderKey]
		if !ok {
			continue
		}

		var existing model.ChatMessage
		if err := tx.Where("group_id = ? AND sender_id = ? AND content = ?", group.ID, user.ID, spec.Content).First(&existing).Error; err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		msg := &model.ChatMessage{
			GroupID:     group.ID,
			SenderID:    user.ID,
			Content:     spec.Content,
			MessageType: spec.MessageType,
			AuditStatus: spec.AuditStatus,
			IsDeleted:   spec.IsDeleted,
		}
		msg.CreatedAt = now.Add(spec.TimeOffset)

		if spec.MessageType == model.ChatMessageTypeImage {
			msg.ImageURL = "https://gamelink.oss.com/chat/game_stats.jpg"
		}

		if err := tx.Create(msg).Error; err != nil {
			return err
		}
	}

	return nil
}

// seedFeedReports 创建动态举报种子数据
func seedFeedReports(tx *gorm.DB, users map[string]*model.User, feeds map[string]*model.Feed) error {
	now := time.Now()
	adminUser := users["adminA"]

	type reportSpec struct {
		FeedKey     string
		ReporterKey string
		Reason      string
		Status      string
		Result      string
		IsHandled   bool
		TimeOffset  time.Duration
	}

	specs := []reportSpec{
		// 待处理的举报
		{
			FeedKey:     "feed_approved_1",
			ReporterKey: "customerE",
			Reason:      "内容涉嫌虚假宣传",
			Status:      "pending",
			TimeOffset:  -2 * time.Hour,
		},
		{
			FeedKey:     "feed_approved_2",
			ReporterKey: "customerF",
			Reason:      "疑似广告内容",
			Status:      "pending",
			TimeOffset:  -1 * time.Hour,
		},
		{
			FeedKey:     "feed_approved_3",
			ReporterKey: "customerG",
			Reason:      "内容不实，与实际体验不符",
			Status:      "pending",
			TimeOffset:  -30 * time.Minute,
		},
		// 已处理-通过的举报
		{
			FeedKey:     "feed_approved_4",
			ReporterKey: "customerH",
			Reason:      "怀疑使用外挂",
			Status:      "approved",
			Result:      "经核实，内容属实，已对被举报内容进行处理",
			IsHandled:   true,
			TimeOffset:  -24 * time.Hour,
		},
		// 已处理-驳回的举报
		{
			FeedKey:     "feed_approved_5",
			ReporterKey: "customerA",
			Reason:      "内容无意义",
			Status:      "rejected",
			Result:      "经审核，该内容符合社区规范，举报不成立",
			IsHandled:   true,
			TimeOffset:  -12 * time.Hour,
		},
		{
			FeedKey:     "feed_followers_1",
			ReporterKey: "customerB",
			Reason:      "涉嫌诱导消费",
			Status:      "rejected",
			Result:      "该内容为正常的优惠活动宣传，不违反规定",
			IsHandled:   true,
			TimeOffset:  -6 * time.Hour,
		},
	}

	for _, spec := range specs {
		feed, ok := feeds[spec.FeedKey]
		if !ok {
			continue
		}
		reporter, ok := users[spec.ReporterKey]
		if !ok {
			continue
		}

		var existing model.FeedReport
		if err := tx.Where("feed_id = ? AND reporter_id = ?", feed.ID, reporter.ID).First(&existing).Error; err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		report := &model.FeedReport{
			FeedID:   feed.ID,
			Reporter: reporter.ID,
			Reason:   spec.Reason,
			Status:   spec.Status,
			Result:   spec.Result,
		}
		report.CreatedAt = now.Add(spec.TimeOffset)

		if spec.IsHandled {
			report.HandledBy = &adminUser.ID
			handledAt := now.Add(spec.TimeOffset + time.Hour)
			report.HandledAt = &handledAt
		}

		if err := tx.Create(report).Error; err != nil {
			return err
		}

		// 更新动态的举报计数
		if err := tx.Model(&model.Feed{}).Where("id = ?", feed.ID).
			Update("metrics_report_count", gorm.Expr("metrics_report_count + 1")).Error; err != nil {
			return err
		}
	}

	return nil
}

// seedContentPermissions 创建内容管理权限种子数据
func seedContentPermissions(tx *gorm.DB) error {
	permissions := []model.Permission{
		// 动态管理权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/feeds",
			Code:        "content.feeds.list",
			Group:       "/admin/content",
			Description: "查看动态列表",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/feeds/:id",
			Code:        "content.feeds.get",
			Group:       "/admin/content",
			Description: "查看动态详情",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/content/feeds/:id/approve",
			Code:        "content.feeds.approve",
			Group:       "/admin/content",
			Description: "批准动态",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/content/feeds/:id/reject",
			Code:        "content.feeds.reject",
			Group:       "/admin/content",
			Description: "拒绝动态",
		},
		{
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/admin/content/feeds/:id",
			Code:        "content.feeds.delete",
			Group:       "/admin/content",
			Description: "删除动态",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/content/feeds/batch-approve",
			Code:        "content.feeds.batch_approve",
			Group:       "/admin/content",
			Description: "批量批准动态",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/content/feeds/batch-reject",
			Code:        "content.feeds.batch_reject",
			Group:       "/admin/content",
			Description: "批量拒绝动态",
		},

		// 聊天监控权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/chat/messages",
			Code:        "content.chat.list",
			Group:       "/admin/content",
			Description: "查看聊天消息列表",
		},
		{
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/admin/content/chat/messages/:id",
			Code:        "content.chat.delete",
			Group:       "/admin/content",
			Description: "删除聊天消息",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/content/chat/mute",
			Code:        "content.chat.mute",
			Group:       "/admin/content",
			Description: "禁言用户",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/content/chat/unmute",
			Code:        "content.chat.unmute",
			Group:       "/admin/content",
			Description: "解除禁言",
		},

		// 举报管理权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/reports",
			Code:        "content.reports.list",
			Group:       "/admin/content",
			Description: "查看举报列表",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/reports/:id",
			Code:        "content.reports.get",
			Group:       "/admin/content",
			Description: "查看举报详情",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/content/reports/:id/process",
			Code:        "content.reports.process",
			Group:       "/admin/content",
			Description: "处理举报",
		},

		// 分类管理权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/categories",
			Code:        "content.categories.list",
			Group:       "/admin/content",
			Description: "查看分类列表",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/content/categories",
			Code:        "content.categories.create",
			Group:       "/admin/content",
			Description: "创建分类",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/content/categories/:id",
			Code:        "content.categories.update",
			Group:       "/admin/content",
			Description: "更新分类",
		},
		{
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/admin/content/categories/:id",
			Code:        "content.categories.delete",
			Group:       "/admin/content",
			Description: "删除分类",
		},

		// 统计权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/stats",
			Code:        "content.stats.list",
			Group:       "/admin/content",
			Description: "查看内容统计",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/stats/export",
			Code:        "content.stats.export",
			Group:       "/admin/content",
			Description: "导出内容统计",
		},
	}

	for _, perm := range permissions {
		var existing model.Permission
		// 先按 code 查找
		if err := tx.Where("code = ?", perm.Code).First(&existing).Error; err == nil {
			// 已存在，更新
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"method":      perm.Method,
				"path":        perm.Path,
				"group":       perm.Group,
				"description": perm.Description,
			}).Error; err != nil {
				return err
			}
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// 再按 method+path 查找（处理唯一约束）
		if err := tx.Where("method = ? AND path = ?", perm.Method, perm.Path).First(&existing).Error; err == nil {
			// 已存在，更新 code 和其他字段
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"code":        perm.Code,
				"group":       perm.Group,
				"description": perm.Description,
			}).Error; err != nil {
				return err
			}
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// 都不存在，创建新记录
		if err := tx.Create(&perm).Error; err != nil {
			return err
		}
	}

	return nil
}
