package main

type HealthBarValues struct {
	toplife    float32
	oldlife    float32
	midlife    float32
	midlifeMin float32
	mlifetime  int32
	gethit     bool
}

func newHealthBarValues() *HealthBarValues {
	return &HealthBarValues{oldlife: 1, midlife: 1, midlifeMin: 1}
}

type PowerBarValues struct {
	midpower    float32
	midpowerMin float32
	prevLevel   int32
}

func newPowerBarValues() *PowerBarValues {
	return &PowerBarValues{}
}

type GuardBarValues struct {
	midpower    float32
	midpowerMin float32
}

func newGuardBarValues() (gb *GuardBarValues) {
	gb = &GuardBarValues{}
	return
}

type StunBarValues struct {
	midpower    float32
	midpowerMin float32
}

func newStunBarValues() (sb *StunBarValues) {
	sb = &StunBarValues{}
	return
}

type LifeBarFaceValues struct {
	face             *Sprite
	old_spr, old_pal [2]int32
}

func newLifeBarFaceValues() *LifeBarFaceValues {
	return &LifeBarFaceValues{}
}

type LifeBarWinIconValues struct {
	wins          []WinType
	numWins       int
	added, addedP Animation
}

func newLifeBarWinIconValues() *LifeBarWinIconValues {
	return &LifeBarWinIconValues{added: Animation{nilAnim: true},
		addedP: Animation{nilAnim: true}}
}

func (wiv *LifeBarWinIconValues) clone() (result *LifeBarWinIconValues) {
	result = &LifeBarWinIconValues{}
	*result = *wiv

	// Manually copy references that shallow copy poorly, as needed
	// Pointers, slices, maps, functions, channels etc
	result.wins = make([]WinType, len(wiv.wins))
	copy(result.wins, wiv.wins)

	return
}

type LifeBarComboValues struct {
	cur, old   int32
	curd, oldd int32
	curp, oldp float32
	resttime   int32
	counterX   float32
	shaketime  int32
	combo      int32
	fakeCombo  int32
}

func newLifeBarComboValues() *LifeBarComboValues {
	return &LifeBarComboValues{}
}

type LifeBarActionValues struct {
	oldleader int
	messages  []*LbMsg
}

func newLifeBarActionValues() *LifeBarActionValues {
	return &LifeBarActionValues{}
}

func (acv *LifeBarActionValues) clone() (result *LifeBarActionValues) {
	result = &LifeBarActionValues{}
	*result = *acv

	// Manually copy references that shallow copy poorly, as needed
	// Pointers, slices, maps, functions, channels etc
	result.messages = make([]*LbMsg, len(acv.messages))
	for i := range acv.messages {
		result.messages[i] = acv.messages[i].clone()
	}

	return
}

type LifeBarRoundValues struct {
	match_wins         [2]int32
	match_maxdrawgames [2]int32
	cur                int32
	wt, swt, dt        [4]int32
	timerActive        bool
	introState         [2]bool
	firstAttack        [2]bool
}

func newLifeBarRoundValues() *LifeBarRoundValues {
	return &LifeBarRoundValues{match_wins: [...]int32{2, 2},
		match_maxdrawgames: [...]int32{1, 1}}
}

type LifeBarTimerValues struct {
	active bool
}

func newLifeBarTimerValues() *LifeBarTimerValues {
	return &LifeBarTimerValues{}
}

type LifeBarScoreValues struct {
	scorePoints float32
	active      bool
}

func newLifeBarScoreValues() *LifeBarScoreValues {
	return &LifeBarScoreValues{}
}

type LifeBarMatchValues struct {
	active bool
}

func newLifeBarMatchValues() *LifeBarMatchValues {
	return &LifeBarMatchValues{}
}

type LifeBarAiLevelValues struct {
	active bool
}

func newLifeBarAiLevelValues() *LifeBarAiLevelValues {
	return &LifeBarAiLevelValues{}
}

type LifeBarWinCountValues struct {
	wins   int32
	active bool
}

func newLifeBarWinCountValues() *LifeBarWinCountValues {
	return &LifeBarWinCountValues{}
}

type LifebarValues struct {
	order      [2][]int
	hb         [8][]HealthBarValues
	pb         [8][]PowerBarValues
	gb         [8][]GuardBarValues
	sb         [8][]StunBarValues
	fa         [8][]LifeBarFaceValues
	wi         [2]LifeBarWinIconValues
	co         [2]LifeBarComboValues
	ac         [2]LifeBarActionValues
	ro         LifeBarRoundValues
	tr         LifeBarTimerValues
	sc         [2]LifeBarScoreValues
	ma         LifeBarMatchValues
	ai         [2]LifeBarAiLevelValues
	wc         [2]LifeBarWinCountValues
	active     bool
	bars       bool
	mode       bool
	redlifebar bool
	guardbar   bool
	stunbar    bool
	hidebars   bool
}

func (lbv *LifebarValues) clone() (result *LifebarValues) {
	result = &LifebarValues{}
	*result = *lbv

	// Manually copy references that shallow copy poorly, as needed
	// Pointers, slices, maps, functions, channels etc
	for i := range lbv.hb {
		result.hb[i] = make([]HealthBarValues, len(lbv.hb[i]))
		copy(result.hb[i], lbv.hb[i])
	}

	for i := range lbv.pb {
		result.pb[i] = make([]PowerBarValues, len(lbv.pb[i]))
		copy(result.pb[i], lbv.pb[i])
	}

	for i := range lbv.gb {
		result.gb[i] = make([]GuardBarValues, len(lbv.gb[i]))
		copy(result.gb[i], lbv.gb[i])
	}

	for i := range lbv.sb {
		result.sb[i] = make([]StunBarValues, len(lbv.sb[i]))
		copy(result.sb[i], lbv.sb[i])
	}

	for i := range lbv.fa {
		result.fa[i] = make([]LifeBarFaceValues, len(lbv.fa[i]))
		copy(result.fa[i], lbv.fa[i])
	}

	for i := range lbv.wi {
		result.wi[i] = *lbv.wi[i].clone()
	}

	for i := range lbv.ac {
		result.ac[i] = *lbv.ac[i].clone()
	}

	return
}
