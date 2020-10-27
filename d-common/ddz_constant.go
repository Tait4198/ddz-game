package d_common

const (
	StageWait         GameStage = iota //等待阶段
	StageGrabLandlord                  //抢地主阶段
	StagePlayPoker                     //游戏阶段
	StageSettlement                    //结算阶段
)

const (
	GameStop                DdzMessageType = 2000
	GameStart               DdzMessageType = 2001
	GameRestart             DdzMessageType = 2002
	GameCountdown           DdzMessageType = 2003
	GameInvalidOps          DdzMessageType = 2004
	GameOpsTimeout          DdzMessageType = 2005
	GameNextUserOps         DdzMessageType = 2006
	GameExe                 DdzMessageType = 2010
	GameDealPoker           DdzMessageType = 2100 // 发送手牌
	GameGrabLandlord        DdzMessageType = 2101 // 抢地主
	GameNGrabLandlord       DdzMessageType = 2102 // 不抢地主
	GameGrabLandlordEnd     DdzMessageType = 2103 // 抢地主阶段结束
	GameGrabHostingOps      DdzMessageType = 2104 // 抢地主托管操作
	GameNoGrabLandlord      DdzMessageType = 2105 // 没人抢地主
	GameNewLandlord         DdzMessageType = 2105 // 地主产生
	GameWaitGrabLandlord    DdzMessageType = 2106 // 等待抢地主
	GameDealHolePokers      DdzMessageType = 2107 // 发送底牌
	GameShowHolePokers      DdzMessageType = 2108 // 显示底牌
	GamePlayPoker           DdzMessageType = 2109 // 出牌
	GamePlayPokerUpdate     DdzMessageType = 2110 // 手牌更新
	GamePlayPokerSkip       DdzMessageType = 2111 // 玩家跳过出牌
	GamePlayPokerHostingOps DdzMessageType = 2112 // 玩家出牌托管操作
	GamePlayPokerInvalid    DdzMessageType = 2113 // 出牌无效
	GamePlayPokerRemaining  DdzMessageType = 2114 // 剩余手牌提示
	GameSettlement          DdzMessageType = 2115 // 游戏结算
	GamePokerRemaining      DdzMessageType = 2116 // 显示玩家剩余卡牌数量
)
