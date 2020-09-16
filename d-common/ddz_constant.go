package d_common

const (
	GameStop             DdzMessageType = 2000
	GameStart                           = 2001
	GameRestart                         = 2002
	GameCountdown                       = 2003
	GameInvalidOps                      = 2004
	GameOpsTimeout                      = 2005
	GameNextUserOps                     = 2006
	GameExe                             = 2010
	GameDealPoker                       = 2100 // 发送手牌
	GameGrabLandlord                    = 2101 // 抢地主
	GameNGrabLandlord                   = 2102 // 不抢地主
	GameGrabLandlordEnd                 = 2103 // 抢地主阶段结束
	GameGrabHostingOps                  = 2104 // 抢地主托管操作
	GameNoGrabLandlord                  = 2105 // 没人抢地主
	GameNewLandlord                     = 2105 // 地主产生
	GameWaitGrabLandlord                = 2106 // 等待抢地主
	GameDealHolePokers                  = 2107 // 发送底牌
	GameShowHolePokers                  = 2108 // 显示底牌
	GamePlayPoker                       = 2109 // 出牌
)
